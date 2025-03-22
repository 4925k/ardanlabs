package hack

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/business/web/v1/response"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100) % 2; n == 0 {
		return response.NewError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
