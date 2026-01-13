package v1

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Handler struct {
	logger zerolog.Logger
}

func NewHandler(logger zerolog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func SetHandler(logger zerolog.Logger, mux *http.ServeMux) {
	handler := NewHandler(logger)
	setRoutes(handler, mux)
}
