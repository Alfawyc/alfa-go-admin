package middleware

import (
	"go_gin/common/global"
	"go_gin/middleware/handle"
	"go_gin/pkg/jwt"
	"time"
)

func AuthInit() (*jwt.GinMiddleware, error) {
	return jwt.New(&jwt.GinMiddleware{
		Realm:         "Test",
		Timeout:       time.Second * 3600,
		Key:           []byte(global.Vp.GetString("jwt.secret")),
		Authenticator: handle.Authenticator,
		PayloadFunc:   handle.PayloadFunc,
		Unauthorized:  handle.Unauthorized,
	})
}
