package v1

import (
	"catify/internal/dto"
	"encoding/json"
	"net/http"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("failed to decode request")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
