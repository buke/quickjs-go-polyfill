package console_test

import (
	"testing"

	"github.com/buke/quickjs-go"
	"github.com/buke/quickjs-go-polyfill/pkg/console"
	"github.com/stretchr/testify/require"
)

func TestConsole(t *testing.T) {
	rt := quickjs.NewRuntime()
	defer rt.Close()

	ctx := rt.NewContext()
	defer ctx.Close()

	err := console.InjectTo(ctx)
	require.NoError(t, err)

	ret, err := ctx.Eval(`
	const obj = {"a": 1, "b": 2};
	console.error(obj);

	const map = new Map(Object.entries({foo: 'bar'}))
	console.log(map);

	console.log('hello', 'world');

	`)
	defer ret.Free()

	require.NoError(t, err)

}
