package conversion

import (
	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/ArthurWang23/miniblog/pkg/core"
)

// 将不同层之间的数据类型转换都在同一个conversion包中实现  所以要避免出现循环依赖
// 通过将不同层的数据类型定义在一个独立的包中
// 或者将不同层之间的数据类型转换函数都定义在独立的包中 避免循环依赖

// 将模型层的UserM转换为protobuf层的user
func UserModelToUserV1(userModel *model.UserM) *apiv1.User {
	var protoUser apiv1.User
	_ = core.CopyWithConverters(&protoUser, userModel)
	return &protoUser
}

func UserV1ToUserModel(protoUser *apiv1.User) *model.UserM {
	var userModel model.UserM
	_ = core.CopyWithConverters(&userModel, protoUser)
	return &userModel
}
