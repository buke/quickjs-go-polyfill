package timer

import (
	_ "embed"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buke/quickjs-go"
)

//go:embed js/timer.js
var timerJs string

type timerJob struct {
	timer  *time.Timer
	ticker *time.Ticker
}

var timerPtrLen int64
var timerLock sync.Mutex
var timerStore = make(map[int64]*timerJob)

func storeTimer(t *timerJob) int64 {
	id := atomic.AddInt64(&timerPtrLen, 1) - 1
	timerLock.Lock()
	defer timerLock.Unlock()
	timerStore[id] = t
	return id
}

func restoreTimer(ptr int64) *timerJob {
	timerLock.Lock()
	defer timerLock.Unlock()
	return timerStore[ptr]
}

func setTimeoutFunc(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	fnName := args[0].String()
	delay := args[1].Int64()

	job := &timerJob{time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {
		ctx.ScheduleJob(func() {
			ret, _ := ctx.Eval(fmt.Sprintf("%s();", fnName))
			defer ret.Free()
		})
	}), nil}
	jobId := storeTimer(job)

	return ctx.Int64(jobId)
}

func clearTimeoutFunc(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	jobId := args[0].Int64()
	t := restoreTimer(jobId)
	t.timer.Stop()
	return ctx.Null()
}

func setIntervalFunc(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	fnName := args[0].String()
	delay := args[1].Int64()

	t := time.NewTicker(time.Duration(delay) * time.Millisecond)
	go func() {
		for range t.C {
			ctx.ScheduleJob(func() {
				ret, _ := ctx.Eval(fmt.Sprintf("%s();", fnName))
				defer ret.Free()
			})
		}
	}()
	job := &timerJob{nil, t}
	jobId := storeTimer(job)

	return ctx.Int64(jobId)
}

func clearIntervalFunc(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	jobId := args[0].Int64()
	job := restoreTimer(jobId)
	job.ticker.Stop()
	return ctx.Null()
}
