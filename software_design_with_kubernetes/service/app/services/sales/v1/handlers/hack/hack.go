package hack

import (
	"encoding/json"
	"net/http"
)

func Hack(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "ok",
	}

	json.NewEncoder(w).Encode(status)
}
