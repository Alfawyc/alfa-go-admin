package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ResponseBodyWrite struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

//@desc 操作记录中间件
func OperationRecord() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body []byte
		var userId int
		//var err error
		if ctx.Request.Method != http.MethodGet {
			body, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Fatalln("read request body error", err.Error())
			}
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
		if claims, ok := ctx.Get("claims"); ok {
			waitUse := claims.(*CustomClaims)
			userId = waitUse.UserId
		} else {
			id, err := strconv.Atoi(ctx.Request.Header.Get("x-user-id"))
			if err != nil {
				userId = 0
			}
			userId = id
		}
		writer := ResponseBodyWrite{
			ctx.Writer,
			&bytes.Buffer{},
		}
		ctx.Writer = writer
		log.Println(ctx.Writer.Status())
		log.Println(writer.body.String())
		log.Println(string(body))
		log.Printf("status : %s , resp %s , userId %d", ctx.Writer.Status(), writer.body.String(), userId)
		ctx.Next()
	}
}
