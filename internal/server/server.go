package server

import (
	"net/http"

	"github.com/abdul-hamid-achik/chessdrill/internal/handler"
	"github.com/abdul-hamid-achik/chessdrill/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router          *chi.Mux
	pageHandler     *handler.PageHandler
	authHandler     *handler.AuthHandler
	drillHandler    *handler.DrillHandler
	statsHandler    *handler.StatsHandler
	settingsHandler *handler.SettingsHandler
	authMiddleware  *middleware.AuthMiddleware
}

func New(
	pageHandler *handler.PageHandler,
	authHandler *handler.AuthHandler,
	drillHandler *handler.DrillHandler,
	statsHandler *handler.StatsHandler,
	settingsHandler *handler.SettingsHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	s := &Server{
		router:          chi.NewRouter(),
		pageHandler:     pageHandler,
		authHandler:     authHandler,
		drillHandler:    drillHandler,
		statsHandler:    statsHandler,
		settingsHandler: settingsHandler,
		authMiddleware:  authMiddleware,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(middleware.Logging)
	s.router.Use(middleware.Recoverer)

	fileServer := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	s.router.NotFound(s.pageHandler.NotFound)

	s.router.Group(func(r chi.Router) {
		r.Use(s.authMiddleware.OptionalAuth)
		r.Get("/", s.pageHandler.Home)
		r.Get("/login", s.pageHandler.Login)
		r.Get("/register", s.pageHandler.Register)
	})

	s.router.Post("/auth/register", s.authHandler.Register)
	s.router.Post("/auth/login", s.authHandler.Login)
	s.router.Post("/auth/logout", s.authHandler.Logout)

	s.router.Group(func(r chi.Router) {
		r.Use(s.authMiddleware.RequireAuth)
		r.Get("/dashboard", s.pageHandler.Dashboard)
		r.Get("/drill", s.pageHandler.DrillSelect)
		r.Get("/drill/{type}", s.pageHandler.Drill)
		r.Get("/stats", s.pageHandler.Stats)
		r.Get("/settings", s.pageHandler.Settings)
	})

	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.authMiddleware.RequireAuth)

		r.Post("/drill/start", s.drillHandler.StartDrill)
		r.Post("/drill/check", s.drillHandler.CheckAnswer)
		r.Post("/drill/end", s.drillHandler.EndDrill)
		r.Get("/drill/moves", s.drillHandler.GetLegalMoves)
		r.Get("/drill/question", s.drillHandler.GetNextQuestion)

		r.Get("/stats/heatmap", s.statsHandler.GetHeatmap)
		r.Get("/stats/overall", s.statsHandler.GetOverall)

		r.Patch("/settings", s.settingsHandler.UpdatePreferences)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
