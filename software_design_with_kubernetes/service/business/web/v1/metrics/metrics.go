package metrics

import (
	"context"
	"expvar"
	"runtime"
)

var m *metrics

// metrics
// the fields are safe to be accessed concurrently thanks to expvar
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constructs the metrics value that will be used to capture metrics
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

const key ctxKey = 1

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines refreshes the goroutine metric every 100 requests
func AddGoroutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			g := int64(runtime.NumGoroutine())
			v.goroutines.Set(g)
			return g
		}
	}

	return 0
}

// AddRequests increments the request metric by 1
func AddRequests(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
		return v.requests.Value()
	}

	return 0
}

// AddErrors increments the error metric by 1
func AddErrors(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
		return v.errors.Value()
	}

	return 0
}

// AddPanics increments the panic metric by 1
func AddPanics(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
		return v.panics.Value()
	}
	return 0
}
