# Polyfill for [quickjs-go](https://github.com/buke/quickjs-go)

[![Test](https://github.com/buke/quickjs-go-polyfill/workflows/Test/badge.svg)](https://github.com/buke/quickjs-go-polyfill/actions?query=workflow%3ATest)
[![codecov](https://codecov.io/gh/buke/quickjs-go-polyfill/branch/main/graph/badge.svg?token=4r8TboEuuJ)](https://codecov.io/gh/buke/quickjs-go-polyfill)
[![Go Report Card](https://goreportcard.com/badge/github.com/buke/quickjs-go-polyfill)](https://goreportcard.com/report/github.com/buke/quickjs-go-polyfill)
[![GoDoc](https://pkg.go.dev/badge/github.com/buke/quickjs-go-polyfill?status.svg)](https://pkg.go.dev/github.com/buke/quickjs-go-polyfill?tab=doc)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbuke%2Fquickjs-go-polyfill.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbuke%2Fquickjs-go-polyfill?ref=badge_shield)

## Features
* fetch: `fetch`
* base64: `atob` and `btoa`
* window: `window`
* console: `console.log` and `console.error` and `console.warn` and `console.info` and `console.debug` and `console.trace`
* timers: `setTimeout` and `setInterval` and `clearTimeout` and `clearInterval`

### Usage
```go
package main
import (
	"time"

	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
)

func main() {
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

	// Wait for the timeout to finish
	time.Sleep(time.Millisecond * 100)

	rt.ExecuteAllPendingJobs()

	// Output:
	// buke
}
```

## Documentation
Go Reference & more examples: https://pkg.go.dev/github.com/buke/quickjs-go-polyfill

## License
[MIT](./LICENSE)


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbuke%2Fquickjs-go-polyfill.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbuke%2Fquickjs-go-polyfill?ref=badge_large)