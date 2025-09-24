package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	pb "github.com/ArthurWang23/miniblog/pkg/api/userservice/v1"
	genericstore "github.com/ArthurWang23/miniblog/pkg/store"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"github.com/ArthurWang23/miniblog/pkg/auth"
	"github.com/ArthurWang23/miniblog/pkg/token"
	"github.com/segmentio/kafka-go"
	"encoding/json"
)

type Handler struct {
	userStore   *genericstore.Store[model.UserM]
	// 新增：Kafka Writer（可选）
	kafkaWriter *kafka.Writer
}

func NewHandler(userStore *genericstore.Store[model.UserM], kafkaWriter *kafka.Writer) *Handler {
	return &Handler{userStore: userStore, kafkaWriter: kafkaWriter}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/users", func(c *gin.Context) {
			var req pb.CreateUserRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			u := &model.UserM{
				Username:  req.GetUsername(),
				Password:  req.GetPassword(), // BeforeCreate 自动加密
				Nickname:  req.GetUsername(),
				Email:     req.GetUsername() + "@example.com",
				Phone:     "19900000000",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := h.userStore.Create(c.Request.Context(), u); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
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
					_ = h.kafkaWriter.WriteMessages(c.Request.Context(), kafka.Message{
						Key:   []byte(u.UserID),
						Value: payload,
					})
				}
			}
			c.JSON(http.StatusOK, &pb.CreateUserReply{UserId: u.UserID})
		})

		v1.POST("/login", func(c *gin.Context) {
			var req pb.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			u, err := h.userStore.Get(c.Request.Context(), where.F("username", req.GetUsername()))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			if err := auth.Compare(u.Password, req.GetPassword()); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
				return
			}
			tokenStr, expireAt, err := token.Sign(u.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "sign token failed"})
				return
			}
			expiresIn := int64(time.Until(expireAt).Seconds())
			if expiresIn < 0 {
				expiresIn = 0
			}
			c.JSON(http.StatusOK, &pb.LoginReply{AccessToken: tokenStr, ExpiresIn: expiresIn})
		})
	}
}