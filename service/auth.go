package service

import (
	"errors"
	"go_gin/common/global"
	"go_gin/model"
	"go_gin/model/request"
	"gorm.io/gorm"
	"log"
	"strconv"
)

//添加角色
func CreateAuth(auth model.Auth) (error, model.Auth) {
	var data model.Auth
	if err := global.Db.Where("authority_id = ? ", auth.AuthorityId).First(&data).Error; err == nil {
		return errors.New("存在相同的角色id"), data
	}
	err := global.Db.Create(&auth).Error
	return err, auth
}

//更新角色
func UpdateAuthority(auth model.Auth) (error, model.Auth) {
	err := global.Db.Where("authority_id = ? ", auth.AuthorityId).First(&model.Auth{}).Updates(&auth).Error
	return err, auth
}

//获取角色列表
func GetAuthList(pageInfo request.PageInfo) (error, []model.Auth, int64) {
	var authList []model.Auth
	err := global.Db.Limit(pageInfo.PageSize).Offset((pageInfo.Page - 1) * pageInfo.PageSize).Where("parent_id = 0").Find(&authList).Error
	if len(authList) >= 1 {
		for key := range authList {
			err = FindChildrenAuth(&authList[key])
		}
	}
	var total int64
	global.Db.Model(model.Auth{}).Count(&total)

	return err, authList, total
}

//获取子角色
func FindChildrenAuth(auth *model.Auth) error {
	err := global.Db.Where("parent_id = ?", auth.AuthorityId).Find(&auth.Children).Error
	log.Println(auth.Children)
	if len(auth.Children) >= 1 {
		for key := range auth.Children {
			err = FindChildrenAuth(&auth.Children[key])
		}
	}

	return err
}

func DeleteAuth(auth model.Auth) error {
	if !errors.Is(global.Db.Where("authority_id = ?", auth.AuthorityId).First(&model.UserAuth{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色用户正在使用,禁止删除")
	}
	if !errors.Is(global.Db.Where("parent_id = ?", auth.AuthorityId).First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("改角色存在子角色,禁止删除")
	}
	//删除角色 , 删除成功之后清空casbin权限数据
	err := global.Db.Where("authority_id = ?", auth.AuthorityId).Delete(&auth).Error
	if err == nil {
		ClearCasbin(0, strconv.Itoa(auth.AuthorityId))
	}

	return err
}

func AllAuthList() (error, []model.Auth) {
	var authList []model.Auth
	err := global.Db.Where("parent_id = 0").Find(&authList).Error
	if len(authList) >= 1 {
		for key := range authList {
			err = FindChildrenAuth(&authList[key])
		}
	}
	return err, authList
}
