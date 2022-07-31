package timer

import "github.com/buke/quickjs-go"

// https://github.com/robertkrimen/natto/blob/5296a7476556988e191b1f0f739bd2e001273d56/natto.go
// https://github.com/dop251/goja_nodejs/blob/master/eventloop/eventloop_test.go

func InjectTo(ctx *quickjs.Context) error {
	ret, _ := ctx.Eval(timerJs)
	defer ret.Free()

	ctx.Globals().Set("__setTimeout", ctx.Function(setTimeoutFunc))
	ctx.Globals().Set("clearTimeout", ctx.Function(clearTimeoutFunc))
	ctx.Globals().Set("__setInterval", ctx.Function(setIntervalFunc))
	ctx.Globals().Set("clearInterval", ctx.Function(clearIntervalFunc))
	return nil
}
