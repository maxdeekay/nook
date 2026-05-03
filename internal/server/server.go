// Go packages map to directories. Every .go file in here must declade package server

package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/maxdeekay/nook/internal/middleware"
)

// Capital P (Port) exports the struct
// Lowecase h (httpServer) keeps the struct private for this package
type Config struct {
	Port   int
	Logger *slog.Logger
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// Constructor returning a pointer to shared Server instance
func New(cfg Config) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	indexTmpl := template.Must(template.ParseFiles("web/templates/index.html"))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		// Render template to w, 2nd argument is data available inside template
		if err := indexTmpl.Execute(w, nil); err != nil {
			cfg.Logger.Error("template render failed", "err", err)
		}
	})

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	handler := middleware.Recovery(cfg.Logger)(middleware.Logging(cfg.Logger)(mux))

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port), // Sprintf is like printf but returns the string instead of printing it
			Handler: handler,
		},
		logger: cfg.Logger,
	}
}

// A method on *Server
func (s *Server) Start() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("server shutting down")
	return s.httpServer.Shutdown(ctx)
}
