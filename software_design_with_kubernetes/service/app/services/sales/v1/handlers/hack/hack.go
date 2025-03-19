package hack

import (
	"context"
	"net/http"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
