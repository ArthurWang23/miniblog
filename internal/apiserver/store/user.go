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

// 实现了user模块在store层所实现的方法
// store层之需要对数据库记录进行简单的增删改查即可
// 对于业务代码，对插入数据或查询数据的处理可以放在Biz层，对于查询条件的定制，可以通过提供灵活的查询参数来实现
type UserStore interface {
	// 标准CRUD
	Create(ctx context.Context, obj *model.UserM) error
	Update(ctx context.Context, obj *model.UserM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.UserM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.UserM, error)
	// 用户操作的附加方法
	UserExpansion
}

type UserExpansion interface{}

// 实现了UserStore接口
type userStore struct {
	store *datastore
}

var _ UserStore = (*userStore)(nil)

func newUserStore(store *datastore) *userStore {
	return &userStore{store: store}
}

// 插入一条用户记录
func (s *userStore) Create(ctx context.Context, obj *model.UserM) error {
	// 从context中获取事务实例 若没有则返回*gorm.DB类型的实例并调用*gorm.DB的create方法完成插入操作
	if err := s.store.DB(ctx).Create(obj).Error; err != nil {
		log.Errorw("Failed to insert user into database", "err", err, "user", obj)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *userStore) Update(ctx context.Context, obj *model.UserM) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to update user in database", "err", err, "user", obj)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.UserM)).Error
	// 幂等删除
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to delete user from database", "err", err, "conditions", opts)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.UserM, error) {
	var obj model.UserM
	if err := s.store.DB(ctx, opts).Find(&obj).Error; err != nil {
		log.Errorw("Failed to retrieve user from database", "error", err, "condition", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDBRead.WithMessage(err.Error())
	}
	return &obj, nil
}

func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.UserM, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list users from database", "err", err, "conditions", opts)
		err = errno.ErrDBRead.WithMessage(err.Error())
	}
	return count, ret, nil
}
