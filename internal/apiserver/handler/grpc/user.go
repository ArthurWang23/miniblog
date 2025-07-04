package grpc

import (
	"context"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
)

func (h *Handler) Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	return h.biz.UserV1().Login(ctx, rq)
}

func (h *Handler) ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error) {
	return h.biz.UserV1().ChangePassword(ctx, rq)
}

func (h *Handler) RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error) {
	return h.biz.UserV1().RefreshToken(ctx, rq)
}

func (h *Handler) CreateUser(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	return h.biz.UserV1().Create(ctx, rq)
}

func (h *Handler) UpdateUser(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	return h.biz.UserV1().Update(ctx, rq)
}

func (h *Handler) DeleteUser(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	return h.biz.UserV1().Delete(ctx, rq)
}

func (h *Handler) GetUser(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	return h.biz.UserV1().Get(ctx, rq)
}

func (h *Handler) ListUser(ctx context.Context, rq *apiv1.ListUsersRequest) (*apiv1.ListUsersResponse, error) {
	return h.biz.UserV1().List(ctx, rq)
}
