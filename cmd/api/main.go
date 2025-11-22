package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/config"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/http/handlers"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/http/middleware"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/repository/postgres"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/storage"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/pkg/log"
)

var (
	version   = "0.1.0" // Updated by release script
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := log.New(cfg.Log.Level, cfg.Log.Format)
	logger.WithFields(map[string]interface{}{
		"version":    version,
		"build_time": buildTime,
		"git_commit": gitCommit,
	}).Info("Starting Tiny School Hub API")

	// Initialize database
	db, err := postgres.NewDB(cfg.Database.URL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()
	logger.Info("Database connection established")

	// Initialize storage
	storageClient, err := storage.NewClient(&cfg.Storage)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize storage")
	}
	logger.Info("Storage client initialized")

	// Initialize repositories
	userRepo := postgres.NewUserRepo(db)
	profileRepo := postgres.NewProfileRepo(db)
	classRepo := postgres.NewClassRepo(db)
	memberRepo := postgres.NewClassMemberRepo(db)
	photoRepo := postgres.NewPhotoRepo(db)
	_ = postgres.NewAbsenceRepo(db)      // TODO: use in handlers
	_ = postgres.NewMessageRepo(db)      // TODO: use in handlers
	_ = postgres.NewAnnouncementRepo(db) // TODO: use in handlers
	tokenRepo := postgres.NewRefreshTokenRepo(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, profileRepo, tokenRepo, cfg, logger)
	classHandler := handlers.NewClassHandler(classRepo, memberRepo, cfg, logger)
	photoHandler := handlers.NewPhotoHandler(photoRepo, memberRepo, storageClient, cfg, logger)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestID)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)
	r.Use(rateLimiter.RateLimit)

	// Health endpoints (no auth required)
	r.Get("/healthz", healthzHandler(logger))
	r.Get("/readyz", readyzHandler(db, storageClient, logger))
	r.Get("/version", versionHandler(logger))

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// Auth routes (no auth required)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/logout", authHandler.Logout)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg))

			// User routes
			r.Get("/me", handlers.NotImplemented) // TODO: implement

			// Class routes
			r.Post("/classes", middleware.RequireRole("TEACHER", "ADMIN")(http.HandlerFunc(classHandler.Create)).ServeHTTP)
			r.Get("/classes", classHandler.ListMyClasses)
			r.Get("/classes/{id}", classHandler.GetByID)
			r.Get("/classes/{id}/members", classHandler.ListMembers)

			// Photo routes
			r.Post("/classes/{id}/photos", photoHandler.CreateUpload)
			r.Get("/classes/{id}/photos", photoHandler.List)

			// Absence routes - TODO: implement
			// Message routes - TODO: implement
			// Announcement routes - TODO: implement
		})
	})

	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Infof("Server listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed")
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server shutdown failed")
	}

	logger.Info("Server stopped")
}

// healthzHandler returns a simple liveness check
func healthzHandler(_ *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}

// readyzHandler returns a readiness check that validates dependencies
func readyzHandler(db *postgres.DB, storageClient *storage.Client, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check database
		if err := db.PingContext(ctx); err != nil {
			logger.WithError(err).Error("Database health check failed")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"error","message":"database unavailable"}`))
			return
		}

		// Check storage
		if err := storageClient.HealthCheck(ctx); err != nil {
			logger.WithError(err).Error("Storage health check failed")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"error","message":"storage unavailable"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}

// versionHandler returns version information
func versionHandler(_ *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{"version":%q,"build_time":%q,"git_commit":%q}`, version, buildTime, gitCommit)
		_, _ = w.Write([]byte(response))
	}
}
