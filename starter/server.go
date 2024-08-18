package starter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Server interface {
	Serve() error
	Shutdown() error
}

var errServerIsNil = errors.New("server is nil")

type HTTPWrapper struct {
	server *http.Server
}

func WrapHTTP(server *http.Server) Server {
	return &HTTPWrapper{server: server}
}

func (s *HTTPWrapper) Serve() error {
	if s.server == nil {
		return errServerIsNil
	}

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http wrapper listen and serve: %w", err)
	}

	return nil
}

func (s *HTTPWrapper) Shutdown() error {
	if s.server == nil {
		return errServerIsNil
	}

	if err := s.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("http wrapper shutdown: %w", err)
	}

	return nil
}
