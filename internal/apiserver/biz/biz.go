package biz

import (
	postv1 "github.com/ArthurWang23/miniblog/internal/apiserver/biz/v1/post"
	userv1 "github.com/ArthurWang23/miniblog/internal/apiserver/biz/v1/user"
	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
)

// 业务逻辑层
type IBiz interface {
	// 获取用户业务接口
	UserV1() userv1.UserBiz
	// 获取博文业务接口
	PostV1() postv1.PostBiz
}

type biz struct {
	store store.IStore
}

var _ IBiz = (*biz)(nil)

func NewBiz(store store.IStore) *biz {
	return &biz{store: store}
}

func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store)
}

func (b *biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}
