package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/internal/cleanup"
	"github.com/FacileStudio/Capsule/apps/api/internal/database"
	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/Capsule/apps/api/internal/middleware"
	"github.com/FacileStudio/Capsule/apps/api/migrations"
	"github.com/FacileStudio/Capsule/apps/api/modules/docs"
	"github.com/FacileStudio/Capsule/apps/api/modules/pastes"

	"github.com/FacileStudio/Journal/sdk/journal"
	"github.com/FacileStudio/tronc/health"
	"github.com/FacileStudio/tronc/healthcheck"
	"github.com/FacileStudio/tronc/httpx"
	"github.com/FacileStudio/tronc/logger"
	troncmiddleware "github.com/FacileStudio/tronc/middleware"
	"github.com/FacileStudio/tronc/migrate"
	"github.com/FacileStudio/tronc/spa"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func main() {
	if healthcheck.Handle(os.Args) {
		return
	}

	os.Exit(run())
}

// run returns the process exit code. Every failure below used to return from
// main, which exits 0 — so a failed migration or an unreachable database looked
// to Docker, Dokploy and any supervisor like a clean shutdown.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	appEnv, err := env.Load()
	if err != nil {
		logger.New(logger.Config{}).Error("failed to load config", slog.Any("error", err))
		return 1
	}

	var journalClient *journal.Client
	appLogger := logger.New(logger.Config{
		Level: appEnv.LogLevel,
		Wrap: func(handler slog.Handler) slog.Handler {
			if appEnv.JournalURL == "" || appEnv.JournalToken == "" {
				return handler
			}
			journalClient = journal.New(journal.Config{URL: appEnv.JournalURL, Token: appEnv.JournalToken})
			return journal.NewHandler(journalClient, handler)
		},
	})
	if journalClient != nil {
		defer journalClient.Close()
	}

	db, err := database.Open(appEnv.DatabaseURL)
	if err != nil {
		appLogger.Error("failed to open database", slog.Any("error", err))
		return 1
	}

	// The *sql.DB is taken before migrating rather than after: the migration
	// runner works at that layer, and so does the readiness check.
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("failed to access database handle", slog.Any("error", err))
		return 1
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("failed to close database", slog.Any("error", err))
		}
	}()

	migrateConfig := migrate.Config{DB: sqlDB, FS: migrations.FS, Logger: appLogger}

	// `docker run <image> migrate status` — the image is ENTRYPOINT-only, and
	// distroless has no shell to run anything else through.
	if handled, err := migrate.Command(ctx, os.Args, migrateConfig); handled {
		if err != nil {
			appLogger.Error("migrate", slog.Any("error", err))
			return 1
		}
		return 0
	}

	if err := migrate.Run(ctx, migrateConfig); err != nil {
		appLogger.Error("failed to run migrations", slog.Any("error", err))
		return 1
	}

	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	go cleanup.Start(cleanupCtx, db, appLogger)

	router, err := buildRouter(db, sqlDB, appEnv, appLogger)
	if err != nil {
		appLogger.Error("failed to build the router", slog.Any("error", err))
		return 1
	}

	addr := ":" + strconv.Itoa(appEnv.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- server.ListenAndServe()
	}()

	appLogger.Info("server starting", slog.String("addr", addr))
	select {
	case err := <-serverErrCh:
		if !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("server stopped", slog.Any("error", err))
			return 1
		}
	case <-ctx.Done():
		appLogger.Info("server shutting down")
		cleanupCancel()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			appLogger.Error("server shutdown failed", slog.Any("error", err))
			return 1
		}
		appLogger.Info("server stopped")
	}
	return 0
}

func buildRouter(db *gorm.DB, sqlDB *sql.DB, appEnv env.Config, appLogger *slog.Logger) (chi.Router, error) {
	pasteService := pastes.NewService(db, appEnv.MaxPasteSize)
	createLimiter := middleware.NewRateLimiter(30, time.Minute)

	router := httpx.NewRouter(httpx.Config{
		Logger: appLogger,
		CORS: troncmiddleware.CORSConfig{
			AllowedOrigins: appEnv.CORSAllowedOrigins,
			AllowedHeaders: append(troncmiddleware.DefaultAllowedHeaders, "X-Delete-Token"),
		},
	})
	health.Mount(router, health.DB(sqlDB))

	router.Route("/api", func(r chi.Router) {
		pastes.RegisterRoutes(r, pasteService, createLimiter)
	})
	if err := docs.RegisterRoutes(router); err != nil {
		return nil, err
	}

	clientDir := spa.DirFromEnv()
	if spa.Available(clientDir) {
		router.Handle("/*", spa.Handler(spa.Config{Dir: clientDir}))
		appLogger.Info("serving client", slog.String("dir", clientDir))
	}

	return router, nil
}
