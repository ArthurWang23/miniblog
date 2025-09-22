package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
)

type kratosAdapter struct {
	s Server
}

func NewKratosTransport(s Server) transport.Server {
	return &kratosAdapter{s: s}
}

func (k *kratosAdapter) Start(ctx context.Context) error {
	go k.s.RunOrDie()
	return nil
}

func (k *kratosAdapter) Stop(ctx context.Context) error {
	k.s.GracefulStop(ctx)
	return nil
}