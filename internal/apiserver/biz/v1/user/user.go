package user

import (
	"context"
	"sync"
	"time"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/apiserver/pkg/conversion"
	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	"github.com/ArthurWang23/miniblog/internal/pkg/auth"
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUsersRequest) (*apiv1.ListUsersResponse, error)

	UserExpansion
}

type UserExpansion interface {
	Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error)
	RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error)
	ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUsersRequest) (*apiv1.ListUsersResponse, error)
}

type userBiz struct {
	store store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(store store.IStore) *userBiz {
	return &userBiz{
		store: store,
	}
}
func (b *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userM model.UserM
	_ = copier.Copy(&userM, rq)
	if err := b.store.User().Create(ctx, &userM); err != nil {
		return nil, err
	}
	return &apiv1.CreateUserResponse{
		UserID: userM.UserID,
	}, nil
}

func (b *userBiz) Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}
	if rq.Username != nil {
		userM.Username = rq.GetUsername()
	}
	if rq.Email != nil {
		userM.Email = rq.GetEmail()
	}
	if rq.Nickname != nil {
		userM.Nickname = rq.GetNickname()
	}
	if rq.Phone != nil {
		userM.Phone = rq.GetPhone()
	}
	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}
	return &apiv1.UpdateUserResponse{}, nil
}

func (b *userBiz) Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	// 只有 root 可以删除用户
	// 所以这里不用where.T() 因为where.T() 会查询root用户自己
	if err := b.store.User().Delete(ctx, where.F("userID", rq.GetUserID())); err != nil {
		return nil, err
	}
	return &apiv1.DeleteUserResponse{}, nil
}

func (b *userBiz) Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}
	return &apiv1.GetUserResponse{
		User: conversion.UserModelToUserV1(userM),
	}, nil
}

// 查询所有用户列表 统计用户所属的博客数
// 这种方法需要遍历多个列表，且对列表中每个元素都有耗时处理逻辑的代码，性能较差
// 使用errgroup 并发查询每个用户的博客数

// 因为要并发处理userList列表中的每个元素，所以需要一个并发安全的数据类型保存处理后的数据
// 使用sync.Map  直接用sync.Map的store方法添加kv对
// store层返回的数据类型为*model.UserM 需要转换为Biz层使用的数据类型*apiv1.User
// 这种转换在Biz层经常发生，因此统一实现conversion
func (b *userBiz) List(ctx context.Context, rq *apiv1.ListUsersRequest) (*apiv1.ListUsersResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}
	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(known.MaxErrGroupConcurrency)

	for _, user := range userList {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := b.store.Post().List(ctx, where.T(ctx))
				if err != nil {
					return err
				}
				converted := conversion.UserModelToUserV1(user)
				converted.PostCount = count
				m.Store(user.ID, converted)
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		log.W(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}
	users := make([]*apiv1.User, 0, len(userList))
	for _, item := range userList {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*apiv1.User))
	}
	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))
	return &apiv1.ListUsersResponse{
		TotalCount: count,
		Users:      users,
	}, nil
}

func (b *userBiz) Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	whr := where.F("username", rq.GetUsername())
	userM, err := b.store.User().Get(ctx, whr)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	if err := auth.Compare(userM.Password, rq.GetPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}
	return &apiv1.LoginResponse{
		Token:    "<placeholder>",
		ExpireAt: timestamppb.New(time.Now().Add(time.Hour * 2)),
	}, nil
}

func (b *userBiz) RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error) {
	// 还没实现
	return &apiv1.RefreshTokenResponse{Token: "<placeholder>", ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour))}, nil
}

func (b *userBiz) ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if err := auth.Compare(userM.Password, rq.GetOldPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}
	// BeforeCreate钩子，在创建用户时会自动加密，但更新信息时不会调用钩子，要手动加密
	userM.Password, _ = auth.Encrypt(rq.GetNewPassword())
	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}
	return &apiv1.ChangePasswordResponse{}, nil
}

func (b *userBiz) ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUsersRequest) (*apiv1.ListUsersResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}
	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	users := make([]*apiv1.User, 0, len(userList))
	for _, user := range userList {
		count, _, err := b.store.Post().List(ctx, where.T(ctx))
		if err != nil {
			return nil, err
		}
		converted := conversion.UserModelToUserV1(user)
		converted.PostCount = count
		users = append(users, converted)
	}
	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))
	return &apiv1.ListUsersResponse{TotalCount: count, Users: users}, nil
}
