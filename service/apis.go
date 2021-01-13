package service

import (
	"errors"
	"go_gin/common/global"
	"go_gin/model"
	"go_gin/model/request"
	"gorm.io/gorm"
)

//@description 添加api
func CreateApi(api model.Api) error {
	if !errors.Is(global.Db.Where("path = ? and method = ?", api.Path, api.Method).First(&model.Api{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同路径和请求方式的api")
	}
	err := global.Db.Create(&api).Error

	return err
}

func UpdateApi(api model.Api) error {
	if !errors.Is(global.Db.Where("path = ? and method = ? and id != ?", api.Path, api.Method, api.ID).First(&model.Api{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同路径和请求方式的api")
	}
	err := global.Db.Updates(&api).Error

	return err
}

func DeleteApi(id int64) error {
	var api model.Api
	err := global.Db.Where("id = ?", id).Delete(&api).Error
	if err != nil {
		return err
	}
	//清除casbin权限
	ClearCasbin(1, api.Path, api.Method)
	return nil
}

func GetApiList(pageInfo request.PageInfo) (int64, []model.Api, error) {
	var apiList []model.Api
	var total int64
	global.Db.Model(&model.Api{}).Count(&total)
	err := global.Db.Limit(pageInfo.PageSize).Offset((pageInfo.Page - 1) * pageInfo.PageSize).Find(&apiList).Error

	return total, apiList, err
}

//@description 获取所有api
func GetAllApi() (error, []model.Api, []string) {
	var allApi []model.Api
	err := global.Db.Find(&allApi).Error
	var group []string
	global.Db.Model(model.Api{}).Select("group").Group("group").Find(&group)
	return err, allApi, group
}
