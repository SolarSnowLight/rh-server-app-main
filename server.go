package todo

import (
	"context"
	"net/http"
	"time"
)

/*
* Structure for working with an HTTP server
 */
type Server struct {
	httpServer *http.Server
}

/*
* Server startup function
 */
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

/*
* Server shutdown function
 */
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
