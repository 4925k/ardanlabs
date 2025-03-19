package mid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/logger"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(ctx, "request started", "method", r.Method, "path", path, "remoteAddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", path, "remoteAddr", r.RemoteAddr)

			return err
		}

		return h
	}

	return m

}
