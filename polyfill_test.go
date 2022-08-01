package polyfill_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
)

func Example() {
	// Create a new runtime
	rt := quickjs.NewRuntime()
	defer rt.Close()

	// Create a new context
	ctx := rt.NewContext()
	defer ctx.Close()

	// Inject polyfills to the context
	polyfill.InjectAll(ctx)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; utf-8")
		_, _ = w.Write([]byte(`{"status": true}`))
	}))

	ret, _ := ctx.Eval(fmt.Sprintf(`
	setTimeout(() => {
		fetch('%s', {Method: 'GET'}).then(response => response.json()).then(data => {
			console.log(data.status);
		});
	}, 50);
	`, srv.URL))

	defer ret.Free()

	time.Sleep(time.Millisecond * 100)

	rt.ExecuteAllPendingJobs()

	// Output:
	// true
}
