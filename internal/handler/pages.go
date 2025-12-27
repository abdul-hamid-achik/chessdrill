package handler

import (
	"net/http"

	"github.com/abdul-hamid-achik/chessdrill/internal/middleware"
	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/service"
	"github.com/abdul-hamid-achik/chessdrill/templates/pages"
)

type PageHandler struct {
	statsService *service.StatsService
	drillService *service.DrillService
}

func NewPageHandler(statsService *service.StatsService, drillService *service.DrillService) *PageHandler {
	return &PageHandler{
		statsService: statsService,
		drillService: drillService,
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	pages.Home(user).Render(r.Context(), w)
}

func (h *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	pages.Login("").Render(r.Context(), w)
}

func (h *PageHandler) Register(w http.ResponseWriter, r *http.Request) {
	pages.Register("").Render(r.Context(), w)
}

func (h *PageHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	stats, err := h.statsService.GetOverallStats(r.Context(), user.ID)
	if err != nil {
		stats = &model.OverallStats{}
	}

	pages.Dashboard(user, stats).Render(r.Context(), w)
}

func (h *PageHandler) DrillSelect(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	pages.DrillSelect(user).Render(r.Context(), w)
}

func (h *PageHandler) Drill(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	drillType := r.PathValue("type")
	if drillType == "" {
		drillType = "name_square"
	}

	pages.Drill(user, drillType).Render(r.Context(), w)
}

func (h *PageHandler) Stats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	stats, err := h.statsService.GetOverallStats(r.Context(), user.ID)
	if err != nil {
		stats = &model.OverallStats{}
	}

	heatmap, err := h.statsService.GetHeatmapData(r.Context(), user.ID)
	if err != nil {
		heatmap = &model.HeatmapData{}
	}

	pages.Stats(user, stats, heatmap).Render(r.Context(), w)
}

func (h *PageHandler) Settings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pages.Settings(user).Render(r.Context(), w)
}

func (h *PageHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	w.WriteHeader(http.StatusNotFound)
	pages.NotFound(user).Render(r.Context(), w)
}

func (h *PageHandler) InternalError(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	w.WriteHeader(http.StatusInternalServerError)
	pages.InternalError(user).Render(r.Context(), w)
}
