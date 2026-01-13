package v1

import "net/http"

func setRoutes(handler *Handler, mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/test", handler.GetTest)
}
