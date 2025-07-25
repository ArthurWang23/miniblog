package store

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	genericstore "github.com/ArthurWang23/miniblog/pkg/store"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
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
	*genericstore.Store[model.UserM]
}

var _ UserStore = (*userStore)(nil)

func newUserStore(store *datastore) *userStore {
	return &userStore{Store: genericstore.NewStore[model.UserM](store, NewLogger())}
}
