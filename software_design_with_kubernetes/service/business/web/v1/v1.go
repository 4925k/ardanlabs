package v1

import (
	"os"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/business/web/v1/mid"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/logger"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

type RouterAddr interface {
	Add(app *web.App, cfg APIMuxConfig)
}

func APIMux(cfg APIMuxConfig, routeAddr RouterAddr) *web.App {

	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	routeAddr.Add(app, cfg)

	return app

}
