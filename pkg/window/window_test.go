package window_test

import (
	"testing"

	"github.com/buke/quickjs-go"
	"github.com/buke/quickjs-go-polyfill/pkg/window"
	"github.com/stretchr/testify/require"
)

func TestWindow(t *testing.T) {
	rt := quickjs.NewRuntime()
	defer rt.Close()

	ctx := rt.NewContext()
	defer ctx.Close()

	err := window.InjectTo(ctx)
	require.NoError(t, err)

	ret, _ := ctx.Eval("Object.is(globalThis,globalThis.window)")
	defer ret.Free()

	require.NoError(t, err)
	require.EqualValues(t, true, ret.Bool())

}
