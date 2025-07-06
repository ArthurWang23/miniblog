package token

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	// 用于签发和解析token的密钥
	key string
	// token中用户身份的键(类似用户id)
	identityKey string
	// 过期时间
	expiration time.Duration
}

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey", 2 * time.Hour}
	once   sync.Once
)

// 设置包级别的配置 config，config会用于本包后面的token签发和解析
func Init(key string, identityKey string, expiration time.Duration) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
		if expiration != 0 {
			config.expiration = expiration
		}
	})
}

// Parse 使用指定的密钥key解析token，解析成功返回token的上下文，否则报错
func Parse(tokenString string, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if key, exists := claims[config.identityKey]; exists {
			if identity, valid := key.(string); valid {
				identityKey = identity
			}
		}
	}
	if identityKey == "" {
		return "", jwt.ErrSignatureInvalid
	}
	return identityKey, nil
}

// 从请求头中获取令牌，并将其传递给parse函数解析
func ParseRequest(ctx context.Context) (string, error) {
	var (
		token string
		err   error
	)
	switch typed := ctx.(type) {
	case *gin.Context:
		header := typed.Request.Header.Get("Authorization")
		if len(header) == 0 {
			return "", errors.New("the length of the `Authorization` header is zero")
		}

		// 从请求头中取出token
		_, _ = fmt.Sscanf(header, "Bearer %s", &token) // 解析 Bearer token
	default:
		token, err = auth.AuthFromMD(typed, "Bearer")
		if err != nil {
			return "", status.Errorf(codes.Unauthenticated, "invalid auth token")
		}
	}
	return Parse(token, config.key)
}

// 使用jwtSecret签发token，token的claims中会存放传入的subject
func Sign(identityKey string) (string, time.Time, error) {
	expireAt := time.Now().Add(config.expiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              expireAt.Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.key))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expireAt, nil
}
