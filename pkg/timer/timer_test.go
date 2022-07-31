package timer_test

import (
	"testing"
	"time"

	"github.com/buke/quickjs-go"
	"github.com/buke/quickjs-go-polyfill/pkg/console"
	"github.com/buke/quickjs-go-polyfill/pkg/timer"
	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	rt := quickjs.NewRuntime()
	defer rt.Close()

	ctx := rt.NewContext()
	defer ctx.Close()

	console.InjectTo(ctx)
	timer.InjectTo(ctx)

	ret, _ := ctx.Eval(`
	var i = 1;
	var timerId = setTimeout(() => {
		i++;
	}, 50);
	var clearTimerId = setTimeout(() => {
		i++;
	}, 50);
	clearTimeout(clearTimerId);

	var tickerTimers = 1;
	var tickerId = setInterval(() => { 
		if (tickerTimers > 3) {
			clearInterval(tickerId);
		}
		else {
			i++;
			tickerTimers++;
		}
	}, 100);

	var clearTickerId = setInterval(() => { 
		i++;
		tickerTimers++;
	}, 100);
	clearInterval(clearTickerId);

	`)
	defer ret.Free()

	time.Sleep(time.Millisecond * 1000)
	rt.ExecuteAllPendingJobs()

	i, _ := ctx.Eval("i")
	defer i.Free()

	require.EqualValues(t, int64(5), i.Int64())

}
