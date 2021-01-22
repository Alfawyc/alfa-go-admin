package system

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_gin/model/response"
	"log"
	"os"
	"path"
	"time"
)

type FileResponse struct {
	FullPath string `json:"full_path"`
	Path     string `json:"path"`
}

func UploadFile(ctx *gin.Context) {
	urlPrefix := ctx.Request.Host
	_, header, err := ctx.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("获取图片失败", ctx)
		return
	}
	fileName := uuid.New().String() + path.Ext(header.Filename)
	dir := time.Now().Format("200601")
	err = CheckAndMkDir("static/upload/" + dir)
	if err != nil {
		response.FailWithMessage("创建文件夹失败 ,"+err.Error(), ctx)
		return
	}
	savePath := fmt.Sprintf("static/upload/%s/%s", dir, fileName)
	err = ctx.SaveUploadedFile(header, savePath)
	if err != nil {
		response.FailWithMessage("保存文件失败 ,"+err.Error(), ctx)
		return
	}
	fullPath := fmt.Sprintf("http://%s/%s", urlPrefix, savePath)
	log.Println(fullPath)
	data := FileResponse{FullPath: fullPath, Path: savePath}
	response.SuccessWithDetail(data, "上传成功", ctx)
}

func CheckAndMkDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		//文件夹存在
		return nil
	}
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0666)
		if err != nil {
			log.Println("mkdir fail")
			return errors.New("mkdir fail")
		}
	}
	return nil
}
