// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package grpc

import (
	"context"
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	log.W(ctx).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatue_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
