package service

import (
	"go_gin/common/global"
	"go_gin/model"
	"go_gin/model/request"
)

func TaskLogList(params request.PageInfo) (int64, []model.TaskLog, error) {
	var total int64
	var list []model.TaskLog
	global.Db.Model(model.TaskLog{}).Count(&total)
	err := global.Db.Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Order("id desc").Find(&list).Error

	return total, list, err
}
