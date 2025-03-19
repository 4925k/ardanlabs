package handlers

import (
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/app/services/sales/v1/handlers/hack"
	v1 "github.com/4925k/ardanlabs/software_design_with_kubernetes/service/business/web/v1"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

type Routes struct{}

func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hack.Routes(app)
}
