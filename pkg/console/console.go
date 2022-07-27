package console

import (
	"fmt"
	"os"

	"github.com/buke/quickjs-go"
)

func consoleFunc(fnType string) (Fn func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value) {
	switch fnType {
	case "trace", "debug", "info", "log", "warn", "error":
		return func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
			fn, _ := ctx.Eval(`(obj) => {
				return JSON.stringify(obj, (key, value) => {
					if(value instanceof Map) {
						return {
						  dataType: 'Map',
						  value: Array.from(value.entries()), // or with spread: value: [...value]
						};
					  } else {
						return value;
					  }
				});
			}`)
			defer fn.Free()

			for _, arg := range args {
				if arg.IsObject() {
					jsonStr := ctx.Invoke(fn, ctx.Null(), arg)
					defer jsonStr.Free()
					fmt.Fprintf(os.Stdout, "%s %s", arg.String(), jsonStr.String())
				} else {
					fmt.Print(arg.String())
				}
				fmt.Print(" ")
			}
			fmt.Println()
			return ctx.Null()
		}
	}
	return nil
}
