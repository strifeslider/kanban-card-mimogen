package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/auth"
	"github.com/user/kanban-saas/pkg/database"
	appmiddleware "github.com/user/kanban-saas/pkg/middleware"
	"github.com/user/kanban-saas/services/card/internal/handler"
	"github.com/user/kanban-saas/services/card/internal/repository"
	"github.com/user/kanban-saas/services/card/internal/service"
)

func main() {
	env := getEnv("ENV", "local")
	port := getEnv("PORT", "8083")

	logger := setupLogger(env)
	logger.Info("starting card service", "env", env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewPostgresPool(ctx, database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "kanban"),
		Password: getEnv("DB_PASSWORD", "kanban_dev_password"),
		Database: getEnv("DB_NAME", "kanban_card"),
		MaxConns: 10,
		MinConns: 2,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	runMigrations(ctx, db, logger)

	jwtCfg := auth.JWTConfig{
		Secret: getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
	}

	cardRepo := repository.NewCardRepository(db)
	labelRepo := repository.NewLabelRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	cardService := service.NewCardService(cardRepo, labelRepo, commentRepo)

	cardHandler := handler.NewCardHandler(cardService)
	labelHandler := handler.NewLabelHandler(cardService)
	commentHandler := handler.NewCommentHandler(cardService)

	r := chi.NewRouter()

	allowedOrigins := appmiddleware.ParseOrigins(getEnv("CORS_ORIGINS", "http://localhost:3000"))
	r.Use(appmiddleware.CORS(allowedOrigins))
	r.Use(appmiddleware.Logging(logger))
	r.Use(appmiddleware.Recovery(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler.SetupRoutes(r, cardHandler, labelHandler, commentHandler, jwtCfg)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("card service listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down card service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
	logger.Info("card service stopped")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case "local":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(handler)
}

func runMigrations(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS cards (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			column_id UUID NOT NULL,
			board_id UUID NOT NULL,
			title VARCHAR(500) NOT NULL,
			description TEXT,
			position INTEGER NOT NULL DEFAULT 0,
			priority INTEGER DEFAULT 0,
			due_date TIMESTAMPTZ,
			created_by UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);`,
		`CREATE INDEX IF NOT EXISTS idx_cards_column ON cards(column_id) WHERE deleted_at IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_cards_board ON cards(board_id) WHERE deleted_at IS NULL;`,
		`CREATE TABLE IF NOT EXISTS card_assignees (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(card_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS labels (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			workspace_id UUID NOT NULL,
			name VARCHAR(100) NOT NULL,
			color VARCHAR(7) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(workspace_id, name)
		);`,
		`CREATE TABLE IF NOT EXISTS card_labels (
			card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
			label_id UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
			PRIMARY KEY (card_id, label_id)
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_card ON comments(card_id) WHERE deleted_at IS NULL;`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(ctx, m); err != nil {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}
	}
	logger.Info("migrations completed")
}
