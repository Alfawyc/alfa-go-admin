package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloWord(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": "Go Gin"})
}
