package userservice

import (
	"github.com/ArthurWang23/miniblog/cmd/mb-userservice/app/options"
	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"context"
	"time"

	"github.com/glebarez/sqlite"
	genericstore "github.com/ArthurWang23/miniblog/pkg/store"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"gorm.io/gorm"
	"github.com/segmentio/kafka-go"
)

type ServerConfig struct {
	cfg *options.Config
	retriever *stubRetriever
	authz     *allowAllAuthorizer

	// 新增：DB 与用户通用 Store
	db        *gorm.DB
	userStore *genericstore.Store[model.UserM]

	// 新增：Kafka Writer（可选）
	kafkaWriter *kafka.Writer
}

// 轻量用户查询：直接用 userID 构造一个用户对象，便于中间件注入上下文
type stubRetriever struct{
	// 新增：从 DB 查询用户
	userStore *genericstore.Store[model.UserM]
}

func (s *stubRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	// 改造：从存储查询真实用户，若不存在则返回一个占位用户（避免中间件报错）
	u, err := s.userStore.Get(ctx, where.F("userID", userID))
	if err == nil {
		return u, nil
	}
	return &model.UserM{
		UserID:   userID,
		Username: "stub",
	}, nil
}

// 轻量鉴权：允许所有请求通过（后续接入真实鉴权）
type allowAllAuthorizer struct{}

func (a *allowAllAuthorizer) Authorize(subject, object, action string) (bool, error) {
	return true, nil
}

// 实现一个 DBProvider，满足 pkg/store 的接口
type gormDBProvider struct{ db *gorm.DB }

func (p *gormDBProvider) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := p.db.WithContext(ctx)
	for _, whr := range wheres {
		if whr != nil {
			db = whr.Where(db)
		}
	}
	return db
}

func NewServerConfig(cfg *options.Config) *ServerConfig {
	// 优先使用 MySQL
	var db *gorm.DB
	var err error
	if cfg.MySQLOptions != nil && cfg.MySQLOptions.Addr != "" {
		db, err = cfg.MySQLOptions.NewDB()
		if err != nil {
			panic(err)
		}
	} else {
		// 回退：SQLite 内存库（开发/学习环境）
		db, err = gorm.Open(sqlite.Open("file:userservice?mode=memory&cache=shared"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}

	// 自动迁移 User 表
	if err := db.AutoMigrate(&model.UserM{}); err != nil {
		panic(err)
	}

	// 构建用户 Store
	provider := &gormDBProvider{db: db}
	userStore := genericstore.NewStore[model.UserM](provider, nil)

	// 新增：初始化 Kafka Writer（可选）
	var kw *kafka.Writer
	if cfg.KafkaOptions != nil && len(cfg.KafkaOptions.Brokers) > 0 && cfg.KafkaOptions.Topic != "" {
		w, werr := cfg.KafkaOptions.Writer()
		if werr != nil {
			panic(werr)
		}
		kw = w
	}

	return &ServerConfig{
		cfg:       cfg,
		retriever: &stubRetriever{userStore: userStore},
		authz:     &allowAllAuthorizer{},
		db:        db,
		userStore: userStore,
		kafkaWriter: kw,
	}
}