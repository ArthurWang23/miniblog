package post

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/apiserver/pkg/conversion"
	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"github.com/jinzhu/copier"
)

type PostBiz interface {
	Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error)
	Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error)
	List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error)

	PostExpansion
}

type PostExpansion interface{}

type postBiz struct {
	store store.IStore
}

var _ PostBiz = (*postBiz)(nil)

func New(store store.IStore) *postBiz {
	return &postBiz{store: store}
}

func (b *postBiz) Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, rq)
	postM.UserID = contextx.UserID(ctx)
	if err := b.store.Post().Create(ctx, &postM); err != nil {
		return nil, err
	}
	return &apiv1.CreatePostResponse{PostID: postM.PostID}, nil
}

func (b *postBiz) Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}
	if rq.Title != nil {
		postM.Title = rq.GetTitle()
	}
	if rq.Content != nil {
		postM.Content = rq.GetContent()
	}
	if err := b.store.Post().Update(ctx, postM); err != nil {
		return nil, err
	}
	return &apiv1.UpdatePostResponse{}, nil
}

func (b *postBiz) Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostIDs())
	if err := b.store.Post().Delete(ctx, whr); err != nil {
		return nil, err
	}
	return &apiv1.DeletePostResponse{}, nil
}

func (b *postBiz) Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}
	return &apiv1.GetPostResponse{Post: conversion.PostModelToPostV1(postM)}, nil
}

func (b *postBiz) List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	whr := where.T(ctx).P(int(rq.GetOffset()), int(rq.GetLimit()))
	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*apiv1.Post, 0, len(postList))
	for _, post := range postList {
		posts = append(posts, conversion.PostModelToPostV1(post))
	}
	return &apiv1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
