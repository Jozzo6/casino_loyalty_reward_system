package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/users/config"
)

type server struct {
	Resource   *config.Resource
	httpServer *http.Server
}

func NewServer(ctx context.Context) (*server, error) {
	var (
		s   server
		err error
	)

	s.Resource, err = config.InitResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resources: %w", err)
	}

	if s.Resource.Config.Environment == "dev" {
		config.Welcome(s.Resource.Config)
	}

	return &s, nil
}

func (s *server) ListenAndServe() error {
	s.httpServer = &http.Server{
		Addr:    ":" + s.Resource.Config.Port,
		Handler: s.routes(),
	}

	s.Resource.Log.Debugf("starting server on port: %s", s.Resource.Config.Port)
	return s.httpServer.ListenAndServe()
}

func (s *server) Close(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return errors.Join(
		s.httpServer.Shutdown(ctx),
		s.Resource.Close(),
	)
}
