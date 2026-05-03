// Entry point, only wiring, a Go app only has one main.go

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxdeekay/nook/internal/server"
)

// stdlib (standard library), default libs
// := is a short for inferring type from the value
// fmt.Println writes to stdout, no prefix, followed by a newline. General-purpose "print this value" function. Space-separated args
// log.Println writes to stderr with an automatic timestamp prefix, intended for diagnostic/logging output
// Print (no newline)
// Printf (formatted with %s, %d etc)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	srv := server.New(server.Config{
		Port:   8080,
		Logger: logger,
	})

	serverErrors := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("server error", "err", err)
		os.Exit(1)

	case sig := <-shutdown:
		logger.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("graceful shutdown failed", "err", err)
			os.Exit(1)
		}

		logger.Info("server stopped")
	}
}
