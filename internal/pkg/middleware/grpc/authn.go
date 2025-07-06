package grpc

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/pkg/token"
	"google.golang.org/grpc"
)

// UserRetriever 用于根据用户名获取用户信息的接口
type UserRetriever interface {
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

// 进行认证
func AuthnInterceptor(retriever UserRetriever) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		//解析 JWT token
		userID, err := token.ParseRequest(ctx)
		if err != nil {
			log.Errorw("Failed to parse request", "err", err)
			return nil, errno.ErrTokenInvalid.WithMessage(err.Error())
		}
		log.Debugw("Token parsing successful", "userID", userID)
		user, err := retriever.GetUser(ctx, userID)
		if err != nil {
			return nil, errno.ErrUnauthenticated.WithMessage(err.Error())
		}
		// 将用户信息存入上下文
		ctx = context.WithValue(ctx, known.XUsername, user.Username)
		ctx = context.WithValue(ctx, known.XUserID, userID)

		ctx = contextx.WithUserID(ctx, user.UserID)
		ctx = contextx.WithUsername(ctx, user.Username)
		return handler(ctx, req)
	}
}
