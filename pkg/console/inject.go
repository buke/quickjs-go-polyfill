package console

import (
	"github.com/buke/quickjs-go"
)

func InjectTo(ctx *quickjs.Context) error {
	consoleObj := ctx.Object()
	consoleObj.Set("trace", ctx.Function(consoleFunc("trace")))
	consoleObj.Set("debug", ctx.Function(consoleFunc("debug")))
	consoleObj.Set("info", ctx.Function(consoleFunc("info")))
	consoleObj.Set("log", ctx.Function(consoleFunc("log")))
	consoleObj.Set("warn", ctx.Function(consoleFunc("warn")))
	consoleObj.Set("error", ctx.Function(consoleFunc("error")))

	ctx.Globals().Set("console", consoleObj)

	return nil
}
