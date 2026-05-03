// Go packages map to directories. Every .go file in here must declade package server

package server

import (
	"fmt"
	"net/http"
)

// Capital P (Port) exports the struct
// Lowecase h (httpServer) keeps the struct private for this package
type Config struct {
	Port int
}

type Server struct {
	httpServer *http.Server
}

// Constructor returning a pointer to shared Server instance
func New(cfg Config) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port), // Sprintf is like printf but returns the string instead of printing it
			Handler: mux,
		},
	}
}

// A method on *Server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
