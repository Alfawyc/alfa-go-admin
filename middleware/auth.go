package middleware

import (
	"go_gin/middleware/handle"
	"go_gin/pkg/jwt"
	"time"
)

func AuthInit() (*jwt.GinMiddleware, error) {
	return jwt.New(&jwt.GinMiddleware{
		Realm:         "Test",
		Timeout:       time.Second * 3600,
		Key:           []byte("Alfa"),
		Authenticator: handle.Authenticator,
		PayloadFunc:   handle.PayloadFunc,
		Unauthorized:  handle.Unauthorized,
	})
}
