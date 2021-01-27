package service

import (
	"errors"
	"go_gin/common/global"
	"go_gin/model"
	"go_gin/model/request"
	"gorm.io/gorm"
)

func GetTaskList(params request.PageInfo) ([]model.TaskList, int64, error) {
	var total int64
	var list []model.TaskList
	global.Db.Model(&model.TaskList{}).Count(&total)
	err := global.Db.Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Find(&list).Error

	return list, total, err
}

func AddTaskList(task model.TaskList) (model.TaskList, error) {
	//相同任务是否存在
	var t model.TaskList
	if !errors.Is(global.Db.Where("command = ?", task.Command).First(&t).Error, gorm.ErrRecordNotFound) {
		return t, errors.New("相同任务已存在")
	}
	err := global.Db.Create(&task).Error
	return task, err
}

func GetAllTask() ([]model.TaskList, error) {
	var task []model.TaskList
	err := global.Db.Find(&task).Error

	return task, err
}

//@desc 任务详情
func TaskDetail(id int) (model.TaskList, error) {
	var detail model.TaskList
	err := global.Db.Where("id = ?", id).First(&detail).Error

	return detail, err
}
