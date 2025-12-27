package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdul-hamid-achik/chessdrill/internal/config"
	"github.com/abdul-hamid-achik/chessdrill/internal/handler"
	"github.com/abdul-hamid-achik/chessdrill/internal/middleware"
	"github.com/abdul-hamid-achik/chessdrill/internal/mongo"
	"github.com/abdul-hamid-achik/chessdrill/internal/repository"
	"github.com/abdul-hamid-achik/chessdrill/internal/server"
	"github.com/abdul-hamid-achik/chessdrill/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := mongo.NewClient(ctx, cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Close(context.Background()); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	if err := mongoClient.CreateIndexes(ctx); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	db := mongoClient.Database()
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	drillSessionRepo := repository.NewDrillSessionRepository(db)
	attemptRepo := repository.NewAttemptRepository(db)

	authService := service.NewAuthService(userRepo, sessionRepo, cfg.SessionMaxAge)
	drillService := service.NewDrillService(drillSessionRepo, attemptRepo)
	statsService := service.NewStatsService(attemptRepo, drillSessionRepo)
	userService := service.NewUserService(userRepo)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	pageHandler := handler.NewPageHandler(statsService, drillService)
	authHandler := handler.NewAuthHandler(authService, cfg.SessionMaxAge)
	drillHandler := handler.NewDrillHandler(drillService)
	statsHandler := handler.NewStatsHandler(statsService)
	settingsHandler := handler.NewSettingsHandler(userService)

	srv := server.New(pageHandler, authHandler, drillHandler, statsHandler, settingsHandler, authMiddleware)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      srv,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on http://localhost:%s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
