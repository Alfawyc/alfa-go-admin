package monitor

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func ServerInfo(ctx *gin.Context) {
	//系统信息
	goOs := runtime.GOOS
	arch := runtime.GOARCH
	mem := runtime.MemProfileRate
	complier := runtime.Compiler
	version := runtime.Version()
	numGoroutine := runtime.NumGoroutine()
	osDic := gin.H{
		"goOs":         goOs,
		"arch":         arch,
		"mem":          mem,
		"complier":     complier,
		"version":      version,
		"numGoroutine": numGoroutine,
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "os": osDic})
}
