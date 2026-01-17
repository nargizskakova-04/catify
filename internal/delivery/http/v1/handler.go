package v1

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Handler struct {
	logger      zerolog.Logger
	userService UserService
}

func NewHandler(
	logger zerolog.Logger,
	userService UserService) *Handler {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}

func SetHandler(
	logger zerolog.Logger,
	userService UserService,
	mux *http.ServeMux) {
	handler := NewHandler(logger, userService)
	setRoutes(handler, mux)
}
