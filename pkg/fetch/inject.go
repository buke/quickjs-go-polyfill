package fetch

import (
	"github.com/buke/quickjs-go"
)

func InjectTo(ctx *quickjs.Context) error {
	ctx.Globals().Set("fetch", ctx.AsyncFunction(fetchFunc))
	return nil
}
