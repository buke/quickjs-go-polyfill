package base64_test

import (
	"testing"

	"github.com/buke/quickjs-go"
	"github.com/buke/quickjs-go-polyfill/pkg/base64"
	"github.com/stretchr/testify/require"
)

func TestBase64(t *testing.T) {
	rt := quickjs.NewRuntime()
	defer rt.Close()

	ctx := rt.NewContext()
	defer ctx.Close()

	err := base64.InjectTo(ctx)
	require.NoError(t, err)

	retBtoa, err := ctx.Eval(`btoa("Hello World")`)
	defer retBtoa.Free()
	require.Equal(t, "SGVsbG8gV29ybGQ=", retBtoa.String())

	retAtob, err := ctx.Eval(`atob("SGVsbG8gV29ybGQ=")`)
	defer retAtob.Free()
	require.Equal(t, "Hello World", retAtob.String())

	uBtoa, err := ctx.Eval(`btoa("你好，世界")`)
	defer uBtoa.Free()
	require.Equal(t, "5L2g5aW977yM5LiW55WM", uBtoa.String())

	uAtob, err := ctx.Eval(`atob("5L2g5aW977yM5LiW55WM")`)
	defer uAtob.Free()
	require.Equal(t, "你好，世界", uAtob.String())

}
