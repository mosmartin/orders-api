package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("error connecting to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			slog.Error("failed to close redis", err)
		}
	}()

	slog.Info("... 🚀 starting server on port 8080")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("error starting server: %w", err)
		}

		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		slog.Info("... 🛑 gracefully shutting down server")

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
