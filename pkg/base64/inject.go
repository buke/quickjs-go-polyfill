package base64

import "github.com/buke/quickjs-go"

func InjectTo(ctx *quickjs.Context) error {
	ctx.Globals().Set("atob", ctx.Function(atobFunc))
	ctx.Globals().Set("btoa", ctx.Function(btoaFunc))
	return nil
}
