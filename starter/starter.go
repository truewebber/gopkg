package starter

import (
	"context"
	"fmt"
	"sync"
)

type Starter struct {
	recorder *errorRecorder
	servers  []Server
}

func NewStarter() *Starter {
	return &Starter{
		recorder: newErrorRecorder(),
	}
}

func (s *Starter) RegisterServer(server Server) {
	s.servers = append(s.servers, server)
}

func (s *Starter) StartServers(ctx context.Context) error {
	closableContext, cancel := context.WithCancel(ctx)
	s.startServers(closableContext, cancel)

	if err := s.recorder.buildFromRecorded(); err != nil {
		return fmt.Errorf("start servers: %w", err)
	}

	return nil
}

func (s *Starter) startServers(ctx context.Context, cancel context.CancelFunc) {
	var wg sync.WaitGroup

	const waitDelta = 2

	for _, server := range s.servers {
		wg.Add(waitDelta)

		serveWithWgRelease := func(cancel context.CancelFunc, server Server) {
			defer wg.Done()
			s.serveWithCancelOnFinish(cancel, server)
		}
		go serveWithWgRelease(cancel, server)

		shutdownWithWgRelease := func(ctx context.Context, server Server) {
			defer wg.Done()
			s.shutdownOnClosedContext(ctx, server)
		}
		go shutdownWithWgRelease(ctx, server)
	}

	wg.Wait()
}

func (s *Starter) shutdownOnClosedContext(ctx context.Context, server Server) {
	<-ctx.Done()

	if err := server.Shutdown(); err != nil {
		s.recorder.record(err)
	}
}

func (s *Starter) serveWithCancelOnFinish(cancel context.CancelFunc, server Server) {
	defer cancel()

	if err := server.Serve(); err != nil {
		s.recorder.record(err)
	}
}
