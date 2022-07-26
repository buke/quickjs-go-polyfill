package fetch

import (
	"errors"

	"github.com/buke/quickjs-go"
)

func InjectTo(ctx *quickjs.Context) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	ctx.Globals().Set("fetch", ctx.AsyncFunction(fetchFunc))

	return nil
}
