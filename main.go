package main

import (
	"FirstTry/config"
	"FirstTry/internal/cache"
	cacheHandler "FirstTry/internal/http_server/handlers/cache"
	"FirstTry/internal/http_server/handlers/cache/users"
	"FirstTry/internal/http_server/handlers/wordsCache"
	"FirstTry/internal/slogpretty"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustSetupConfig()
	logger := setupLogger(cfg.Env)
	ch := cache.NewCache(cfg.CleanupInterval, logger)
	chWords := cache.NewChForWords(cfg.CleanupInterval)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/cache", func(r chi.Router) {
		r.Post("/users", users.AddInCacheUsersAndProductId(logger, ch))
		r.Post("/words", wordsCache.AddInCacheLetter(logger, chWords))
		r.Get("/", cacheHandler.ShowCache(logger, ch))
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("starting server", slog.String("address", cfg.Port))
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server")
		}
	}()
	<-ctx.Done()
	logger.Info("stopping server")
	time.Sleep(time.Second * 5)
}

func setupLogger(env string) *slog.Logger {
	return setupPrettySlog()
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlersOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
