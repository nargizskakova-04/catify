package v1

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Handler struct {
	logger      zerolog.Logger
	userService UserService
	goalService GoalService
}

func NewHandler(
	logger zerolog.Logger,
	userService UserService,
	goalService GoalService) *Handler {
	return &Handler{
		logger:      logger,
		userService: userService,
		goalService: goalService,
	}
}

func SetHandler(
	logger zerolog.Logger,
	userService UserService,
	goalService GoalService,
	mux *http.ServeMux) {
	handler := NewHandler(logger, userService, goalService)
	setRoutes(handler, mux)
}
