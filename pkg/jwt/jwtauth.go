package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const PayLoad = "JWT_PAYLOAD"

type MapClaims map[string]interface{}

type GinMiddleware struct {
	Realm            string
	SigningAlgorithm string
	Key              []byte
	//有效时间
	Timeout         time.Duration
	MaxRefresh      time.Duration
	Authenticator   func(ctx *gin.Context) (interface{}, error)
	Authorizator    func(data interface{}, ctx *gin.Context) bool
	PayloadFunc     func(data interface{}) MapClaims
	Unauthorized    func(*gin.Context, int, string)
	LoginResponse   func(*gin.Context, int, string, time.Time)
	RefreshResponse func(*gin.Context, int, string, time.Time)
}

var (
	IdentityKey = "identity"

	NiceKey = "nice"
)

func New(m *GinMiddleware) (*GinMiddleware, error) {
	return m, nil
}

//初始化配置
func (m *GinMiddleware) MiddlewareInit() error {
	m.SigningAlgorithm = "HS256"
	m.Timeout = time.Duration(3600) * time.Second
	m.Authorizator = func(data interface{}, ctx *gin.Context) bool {
		return true
	}
	//未授权返回信息
	if m.Unauthorized == nil {
		m.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": message,
			})
		}
	}
	m.LoginResponse = func(ctx *gin.Context, code int, token string, expire time.Time) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":   code,
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
	m.RefreshResponse = func(ctx *gin.Context, i int, token string, expire time.Time) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}

	return nil
}

//生产jwt token
func (m *GinMiddleware) TokenGenerator(data interface{}) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(m.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
	if m.PayloadFunc != nil {
		for key, val := range m.PayloadFunc(data) {
			claims[key] = val
		}
	}
	expire := time.Now().Add(m.Timeout)
	claims["exp"] = expire.Unix()
	claims["original_at"] = time.Now().Unix()
	tokenString, err := m.signedString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

//提取jwt 声明
func ExtraClaims(c *gin.Context) MapClaims {
	claims, exists := c.Get(PayLoad)
	if !exists {
		return make(MapClaims)
	}

	return claims.(MapClaims)
}

//从token提取jwt声明
func ExtraClaimsFromToken(token *jwt.Token) MapClaims {
	if token == nil {
		return make(MapClaims)
	}
	claims := MapClaims{}
	for key, val := range token.Claims.(jwt.MapClaims) {
		claims[key] = val
	}

	return claims
}

func (m *GinMiddleware) signedString(token *jwt.Token) (string, error) {
	//todo 根据配置文件判断key
	tokenString, err := token.SignedString([]byte("jwt"))

	return tokenString, err
}

func (m *GinMiddleware) LoginHandler(ctx *gin.Context) {
	if m.Authenticator == nil {
		m.Unauthorized(ctx, 400, "jwt Authenticator is undefined")
		return
	}
	data, err := m.Authenticator(ctx)
	if err != nil {
		m.Unauthorized(ctx, 400, err.Error())
		return
	}
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)

	if m.PayloadFunc != nil {
		//保存用户信息
		for key, value := range m.PayloadFunc(data) {
			claims[key] = value
		}
	}
	log.Println(claims)
	//过期时间
	expire := time.Now().Add(m.Timeout)
	claims["exp"] = expire
	claims["orig_iat"] = time.Now().Unix()
	tokenString, err := m.signedString(token)
	if err != nil {
		log.Fatalln("jwt token fail", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}
