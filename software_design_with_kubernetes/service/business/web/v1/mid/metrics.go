package mid

import (
	"context"
	"net/http"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/business/web/v1/metrics"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)

			err := handler(ctx, w, r)

			n := metrics.AddRequests(ctx)
			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
