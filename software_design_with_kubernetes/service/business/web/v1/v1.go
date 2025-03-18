package v1

import (
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/app/services/sales/v1/handlers/hack"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/logger"
)

type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

func APIMux(cfg APIMuxConfig) http.Handler {

	mux := httptreemux.NewContextMux()

	mux.Handle(http.MethodGet, "/hack", hack.Hack)

	return mux

}
