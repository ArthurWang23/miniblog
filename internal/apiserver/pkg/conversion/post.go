package conversion

import (
	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/onexstack/onexstack/pkg/core"
)

func PostModelToPostV1(postModel *model.PostM) *apiv1.Post {
	var protoPost apiv1.Post
	_ = core.CopyWithConverters(&protoPost, postModel)
	return &protoPost
}

func PostV1ToPostModel(protoPost *apiv1.Post) *model.PostM {
	var postModel model.PostM
	_ = core.CopyWithConverters(&postModel, protoPost)
	return &postModel
}
