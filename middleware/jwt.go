package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go_gin/model"
	"go_gin/model/response"
	"strconv"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	UserId      int
	Username    string
	Nickname    string
	AuthorityId string
	BufferTime  int64
	jwt.StandardClaims
}

var (
	TokenExpried   = errors.New("Token is expired")
	TokenNotValid  = errors.New("Token is not active")
	TokenMalformed = errors.New("Not a token")
	TokenInvalid   = errors.New("can not handler this token")
)

func NewJWT() *JWT {
	return &JWT{[]byte("alfa")}
}

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//todo 测试
		//ctx.Next()
		//return
		token := ctx.Request.Header.Get("x-token")
		if token == "" {
			response.JWTFailWithMessage("非法访问", ctx)
			ctx.Abort()
			return
		}
		//todo 拉黑token

		j := NewJWT()
		//解析token
		claims, err := j.ParseToken(token)
		if err != nil {
			response.JWTFailWithMessage("解析token失败", ctx)
			ctx.Abort()
			return
		}
		user := model.User{}
		user.ID = uint(claims.UserId)

		if _, err := user.GetOne(); err != nil {
			//todo 加入黑名单
			response.JWTFailWithMessage("未获取到用户信息", ctx)
			ctx.Abort()
		}
		//验证有效期
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			claims.ExpiresAt = time.Now().Unix() + 60*60*24*7
			newToken, _ := j.CreateToken(*claims)
			newClaims, _ := j.ParseToken(newToken)
			ctx.Header("new-token", newToken)
			ctx.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
			//多点登陆
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

//创建token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.SigningKey)

	return tokenString, err
}

//解析jwt token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, errors.New("token parse fail")
	}
	if token == nil {
		return nil, errors.New("token 不可用")
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("fail to parse")
}
