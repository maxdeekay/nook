// Entry point, only wiring, a Go app only has one main.go

package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

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

	// net/http returns http.ErrServerClosed as a signal to indicate clean shutdown
	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "err", err) // "message", "key", value
		os.Exit(1)
	}
}
