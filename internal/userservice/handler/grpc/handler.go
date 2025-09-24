package grpc

import (
	"context"
	pb "github.com/ArthurWang23/miniblog/pkg/api/userservice/v1"
	"time"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	genericstore "github.com/ArthurWang23/miniblog/pkg/store"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"github.com/ArthurWang23/miniblog/pkg/auth"
	"github.com/ArthurWang23/miniblog/pkg/token"
	"github.com/segmentio/kafka-go"
	"encoding/json"
)

type Handler struct {
	// 预留：后续可注入 biz/store/log 等依赖
	userStore   *genericstore.Store[model.UserM]
	// 新增：Kafka Writer（可选）
	kafkaWriter *kafka.Writer
}

func NewHandler(userStore *genericstore.Store[model.UserM], kafkaWriter *kafka.Writer) *Handler {
	return &Handler{userStore: userStore, kafkaWriter: kafkaWriter}
}

func (h *Handler) Healthz(ctx context.Context, _ *pb.HealthzRequest) (*pb.HealthzReply, error) {
	return &pb.HealthzReply{Status: "ok"}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	// 满足非空约束：nickname/email/phone 简单默认值
	u := &model.UserM{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(), // BeforeCreate 钩子会自动加密
		Nickname:  req.GetUsername(),
		Email:     req.GetUsername() + "@example.com",
		Phone:     "19900000000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.userStore.Create(ctx, u); err != nil {
		return nil, err
	}
	// 新增：发送用户创建事件到 Kafka（可选）
	if h.kafkaWriter != nil {
		type UserCreatedEvent struct {
			UserID    string    `json:"user_id"`
			Username  string    `json:"username"`
			CreatedAt time.Time `json:"created_at"`
		}
		evt := UserCreatedEvent{UserID: u.UserID, Username: u.Username, CreatedAt: u.CreatedAt}
		if payload, err := json.Marshal(evt); err == nil {
			_ = h.kafkaWriter.WriteMessages(ctx, kafka.Message{
				Key:   []byte(u.UserID),
				Value: payload,
			})
		}
	}
	return &pb.CreateUserReply{UserId: u.UserID}, nil
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	u, err := h.userStore.Get(ctx, where.F("username", req.GetUsername()))
	if err != nil {
		return nil, err
	}
	if err := auth.Compare(u.Password, req.GetPassword()); err != nil {
		return nil, err
	}
	tokenStr, expireAt, err := token.Sign(u.UserID)
	if err != nil {
		return nil, err
	}
	expiresIn := int64(time.Until(expireAt).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}
	return &pb.LoginReply{AccessToken: tokenStr, ExpiresIn: expiresIn}, nil
}