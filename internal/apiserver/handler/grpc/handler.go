package handler

import (
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
)

type Handler struct {
	apiv1.UnimplementedMiniBlogServer
}

func NewHandler() *Handler {
	return &Handler{}
}
