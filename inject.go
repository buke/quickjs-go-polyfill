package polyfill

import (
	"github.com/buke/quickjs-go"
	"github.com/buke/quickjs-go-polyfill/pkg/fetch"
	"github.com/buke/quickjs-go-polyfill/pkg/window"
)

func InjectAll(ctx *quickjs.Context) {
	window.InjectTo(ctx)
	fetch.InjectTo(ctx)
}
