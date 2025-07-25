package store

import (
	"context"
	"errors"
	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	genericstore "github.com/ArthurWang23/miniblog/pkg/store"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"gorm.io/gorm"
)

type PostStore interface {
	Create(ctx context.Context, obj *model.PostM) error
	Update(ctx context.Context, obj *model.PostM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.PostM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.PostM, error)

	PostExpansion
}

type PostExpansion interface{}

type postStore struct {
	*genericstore.Store[model.PostM]
}

var _ PostStore = (*postStore)(nil)

func newPostStore(store *datastore) *postStore {
	return &postStore{Store: genericstore.NewStore[model.PostM](store, NewLogger())}
}
