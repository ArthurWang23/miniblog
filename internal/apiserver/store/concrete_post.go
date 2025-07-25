package store

import (
	"context"
	"errors"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"gorm.io/gorm"
)

// 代码架构升级
// UserStore和PostStore除了返回的GORM Model结构体不一致外 参数定义和方法实现基本是一样的
// 通过Go泛型标准化Store层

// 本代码用于记录优化前的store层 post 代码
type ConcretePostStore interface {
	Create(ctx context.Context, obj *model.PostM) error
	Update(ctx context.Context, obj *model.PostM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.PostM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.PostM, error)

	ConcretePostExpansion
}

type ConcretePostExpansion interface{}

type concretePostStore struct {
	store *datastore
}

var _ ConcretePostStore = (*concretePostStore)(nil)

func newConcretePostStore(store *datastore) *concretePostStore {
	return &concretePostStore{store: store}
}

func (s *concretePostStore) Create(ctx context.Context, obj *model.PostM) error {
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		log.Errorw("Failed to insert post into database", "err", err, "post", obj)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *concretePostStore) Update(ctx context.Context, obj *model.PostM) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to update post in database", "err", err, "post", obj)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *concretePostStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.PostM)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to delete post from database", "err", err, "conditions", opts)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *concretePostStore) Get(ctx context.Context, opts *where.Options) (*model.PostM, error) {
	var obj model.PostM
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		log.Errorw("Failed to retrieve post from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPostNotFound
		}
		return nil, errno.ErrDBRead.WithMessage(err.Error())
	}
	return &obj, nil
}

func (s *concretePostStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.PostM, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list posts from database", "err", err, "conditions", opts)
		err = errno.ErrDBRead.WithMessage(err.Error())
	}
	return count, ret, nil
}
