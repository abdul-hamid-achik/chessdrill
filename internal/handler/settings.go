package handler

import (
	"net/http"

	"github.com/abdul-hamid-achik/chessdrill/internal/middleware"
	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/service"
)

type SettingsHandler struct {
	userService *service.UserService
}

func NewSettingsHandler(userService *service.UserService) *SettingsHandler {
	return &SettingsHandler{
		userService: userService,
	}
}

func (h *SettingsHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Build preferences from form data
	prefs := model.Preferences{
		Perspective:     r.FormValue("perspective"),
		ShowCoordinates: r.FormValue("show_coordinates") == "on",
		Theme:           r.FormValue("theme"),
	}

	// Set defaults if empty
	if prefs.Perspective == "" {
		prefs.Perspective = "white"
	}
	if prefs.Theme == "" {
		prefs.Theme = "light"
	}

	if err := h.userService.UpdatePreferences(r.Context(), user.ID, prefs); err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
