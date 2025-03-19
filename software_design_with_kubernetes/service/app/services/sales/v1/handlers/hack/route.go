package hack

import (
	"net/http"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
)

func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
