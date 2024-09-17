package starter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type Server interface {
	Serve() error
	Shutdown() error
}

var errServerIsNil = errors.New("server is nil")
var errListenerIsNil = errors.New("listener is nil")

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

type GRPCWrapper struct {
	server   *grpc.Server
	listener net.Listener
}

func WrapGRPC(server *grpc.Server, listener net.Listener) Server {
	return &GRPCWrapper{server: server, listener: listener}
}

func (s *GRPCWrapper) Serve() error {
	if s.server == nil {
		return errServerIsNil
	}

	if s.listener == nil {
		return errListenerIsNil
	}

	if err := s.server.Serve(s.listener); err != nil {
		return fmt.Errorf("grpc server serve: %w", err)
	}

	return nil
}

func (s *GRPCWrapper) Shutdown() error {
	s.server.GracefulStop()

	return nil
}
