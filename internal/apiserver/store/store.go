package store

import (
	"context"
	"sync"

	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"gorm.io/gorm"
)

var (
	once sync.Once
	// 全局变量 方便其他包调用已经初始化好的datastore实例
	// 包级变量 通过store.S.User().Create()调用Store层接口
	S *datastore
)

// 定义了Store层需要实现的方法
type IStore interface {
	// 返回Store层的*gorm.DB实例
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
	// 将*gorm.DB类型的实例注入到context中
	// DB 和 TX 在Biz层实现事务，在Store层执行事务
	TX(ctx context.Context, fn func(ctx context.Context) error) error
	// User()和Post()分别返回User资源的store层方法和Post资源的store层方法
	User() UserStore
	Post() PostStore
}

// 用于在context.Context中存储事务上下文的键
type transactionKey struct{}

// datastore是IStore的具体实现
type datastore struct {
	// 可以根据需要添加其他数据库实例
	core *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	// 确保只初始化一次
	once.Do(func() {
		S = &datastore{core: db}
	})
	return S
}

// DB 根据传入的条件（wheres）对数据库实例进行筛选
// 如果未传入任何条件，则返回上下文中的数据库实例（事务实例或核心数据库实例）
func (store *datastore) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := store.core
	// 从上下文中提取事务实例
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}
	// 遍历所有传入的条件并逐一叠加到数据库查询对象上
	for _, whr := range wheres {
		db = whr.Where(db)
	}
	return db
}

// TX返回一个新的事务实例
func (store *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return store.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

func (store *datastore) User() UserStore {
	return newUserStore(store)
}

func (store *datastore) Post() PostStore {
	return newPostStore(store)
}
