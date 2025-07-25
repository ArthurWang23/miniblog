// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package http

import (
	"time"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(200, &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatue_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
