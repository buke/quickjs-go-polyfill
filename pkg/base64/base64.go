package base64

import (
	stdBase64 "encoding/base64"

	"github.com/buke/quickjs-go"
)

func atobFunc(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	aStr := args[0].String()
	aByts, err := stdBase64.StdEncoding.DecodeString(aStr)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.String(string(aByts))
}

func btoaFunc(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	bStr := args[0].String()
	aByts := stdBase64.StdEncoding.EncodeToString([]byte(bStr))
	return ctx.String(aByts)
}
