package polyfill_test

import (
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

	ret, _ := ctx.Eval(`
	setTimeout(() => {
		fetch('https://api.github.com/users/buke', {Method: 'GET'}).then(response => response.json()).then(data => {
			console.log(data.login);
		});
	}, 50);
	`)
	defer ret.Free()

	time.Sleep(time.Millisecond * 100)

	rt.ExecuteAllPendingJobs()

	// Output:
	// buke
}
