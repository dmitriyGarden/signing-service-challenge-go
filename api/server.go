package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// DeviceHandler mounts HTTP routes on a chi router.
type DeviceHandler interface {
	Register(r chi.Router)
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	handlers      map[string]DeviceHandler
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, handlers map[string]DeviceHandler) *Server {
	return &Server{
		listenAddress: listenAddress,
		handlers:      handlers,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	for prefix, handlers := range s.handlers {
		r.Route(prefix, handlers.Register)
	}
	return http.ListenAndServe(s.listenAddress, r)
}
