package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/buke/quickjs-go"
)

type QJSRequest struct {
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
	Redirect string            `json:"redirect"`
}

type QJSResponse struct {
	Status     int32  `json:"status"`
	StatusText string `json:"statusText"`
	OK         bool   `json:"ok"`
	Redirected bool   `json:"redirected"`
	URL        string `json:"url"`
}

func prepareReq(ctx *quickjs.Context, args []quickjs.Value) (*http.Request, error) {
	if len(args) <= 0 {
		return nil, errors.New("at lease 1 argument required")
	}
	rawURL := args[0].String()

	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("url '%s' is not valid", rawURL)
	}

	var jsReq QJSRequest
	if len(args) > 1 {
		if !args[1].IsObject() {
			return nil, errors.New("2nd argument must be an object")
		}
		reader := strings.NewReader(args[1].JSONStringify())
		if err := json.NewDecoder(reader).Decode(&jsReq); err != nil {
			return nil, err
		}
	}

	if jsReq.Method == "" {
		jsReq.Method = "GET"
	}

	var body io.Reader
	if jsReq.Method != "GET" {
		body = strings.NewReader(jsReq.Body)
	}

	req, err := http.NewRequest(jsReq.Method, url.String(), body)
	if err != nil {
		return nil, err
	}
	for k, v := range jsReq.Headers {
		headerName := http.CanonicalHeaderKey(k)
		req.Header.Set(headerName, v)
	}

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "*/*")
	}

	if req.Header.Get("Connection") == "" {
		req.Header.Set("Connection", "close")
	}

	req.Header.Set("Redirect", jsReq.Redirect)

	return req, nil
}

func fetch(req *http.Request) (*http.Response, error) {
	redirected := false
	client := &http.Client{
		Transport: http.DefaultTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			switch req.Header.Get("Redirect") {
			case "error":
				return errors.New("redirects are not allowed")
			default:
				if len(via) >= 10 {
					return errors.New("stopped after 10 redirects")
				}
			}

			redirected = true
			return nil
		},
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	res.Header.Set("Redirected", fmt.Sprintf("%v", redirected))

	return res, nil
}

func prepareResp(ctx *quickjs.Context, resp *http.Response) (quickjs.Value, error) {
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ctx.Null(), err
	}

	// prepare response obj
	var jsResp QJSResponse
	jsResp.Status = int32(resp.StatusCode)
	jsResp.StatusText = resp.Status
	jsResp.OK = resp.StatusCode >= 200 && resp.StatusCode < 300
	jsResp.Redirected = resp.Header.Get("Redirected") == "true"
	jsResp.URL = resp.Request.URL.String()

	//  marshal jsResp to json
	b, err := json.Marshal(jsResp)
	if err != nil {
		return ctx.Null(), err
	}

	// parse json to js object
	respObj := ctx.ParseJSON(string(b))

	// header object
	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = strings.Join(v, ",")
	}
	h, err := json.Marshal(headers)
	if err != nil {
		return ctx.Null(), err
	}
	// use js map to set header object
	fn, _ := ctx.Eval(`(str) => {
		return new Map(Object.entries(JSON.parse(str)));
	}`)
	defer fn.Free()

	headerObj := ctx.Invoke(fn, ctx.Null(), ctx.String(string(h)))
	if headerObj.IsError() {
		return ctx.Null(), headerObj.Error()
	}
	respObj.Set("headers", headerObj)

	respObj.Set("text", ctx.AsyncFunction(func(ctx *quickjs.Context, this quickjs.Value, promise quickjs.Value, args []quickjs.Value) quickjs.Value {
		return promise.Call("resolve", ctx.String(string(respBody)))
	}))

	respObj.Set("json", ctx.AsyncFunction(func(ctx *quickjs.Context, this quickjs.Value, promise quickjs.Value, args []quickjs.Value) quickjs.Value {
		retObj := ctx.ParseJSON(string(respBody))
		defer retObj.Free()
		if retObj.IsError() {
			return promise.Call("reject", retObj)
		}

		return promise.Call("resolve", retObj)
	}))

	return respObj, nil

}

func fetchFunc(ctx *quickjs.Context, this quickjs.Value, promise quickjs.Value, args []quickjs.Value) quickjs.Value {
	// prepare request
	req, err := prepareReq(ctx, args)
	if err != nil {
		return promise.Call("reject", ctx.ThrowError(err))
	}

	// fetch http request
	resp, err := fetch(req)
	if err != nil {
		return promise.Call("reject", ctx.ThrowError(err))
	}

	// prepare response obj
	respObj, err := prepareResp(ctx, resp)
	if respObj.IsError() {
		fmt.Println(ctx.Exception())
	}
	if err != nil {
		return promise.Call("reject", ctx.ThrowError(err))
	}
	defer respObj.Free()

	return promise.Call("resolve", respObj)

}
