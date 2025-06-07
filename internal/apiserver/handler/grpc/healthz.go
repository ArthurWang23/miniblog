package handler

import (
	"context"
	"time"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatue_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
