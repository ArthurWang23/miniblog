package grpc

import (
	"context"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
)

func (h *Handler) CreatePost(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	return h.biz.PostV1().Create(ctx, rq)
}

func (h *Handler) UpdatePost(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	return h.biz.PostV1().Update(ctx, rq)
}

func (h *Handler) DeletePost(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	return h.biz.PostV1().Delete(ctx, rq)
}

func (h *Handler) GetPost(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	return h.biz.PostV1().Get(ctx, rq)
}

func (h *Handler) ListPost(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	return h.biz.PostV1().List(ctx, rq)
}
