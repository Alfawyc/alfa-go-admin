package handle

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go_gin/model"
	"go_gin/pkg/jwt"
	"log"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func Authenticator(ctx *gin.Context) (interface{}, error) {
	var loginVal model.Login
	//自动选择合适的绑定
	err := ctx.ShouldBind(&loginVal)
	log.Println(loginVal)
	if err != nil {
		return nil, errors.New("缺少用户名或密码")
	}
	//验证码
	log.Printf("code_id %s code %s", loginVal.CodeId, loginVal.Code)
	isOk := store.Verify(loginVal.CodeId, loginVal.Code, true)
	if !isOk {
		log.Println("验证码错误")
		return nil, errors.New("验证码错误")
	}
	//获取用户信息
	user, err := loginVal.GetUser()
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	return map[string]interface{}{"user": user}, nil
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(model.User)
		return jwt.MapClaims{
			jwt.IdentityKey: u.ID,
			jwt.NiceKey:     u.Username,
		}
	}

	return jwt.MapClaims{}
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
