package http

import (
	"github.com/go-chi/chi/v5"
)

// Server handles the HTTP frontend of the application.
type Server struct {
	r chi.Router
	n *http.Server
}

// New retruns a ready to serve instance of the HTTP server.
func New() (*Server, error) {
	s := new(Server)
	return s, nil
}

// Serve binds and serves http on the bound socket.  An error will be
// returned if the server cannot initialize.
func (s *Server) Serve(bind string) error {
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.n.Shutdown(ctx)
}
