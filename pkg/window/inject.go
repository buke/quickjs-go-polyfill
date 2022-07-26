package window

import (
	"errors"

	"github.com/buke/quickjs-go"
)

func InjectTo(ctx *quickjs.Context) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	ret, err := ctx.Eval(windowJs)
	defer ret.Free()

	return err
}
