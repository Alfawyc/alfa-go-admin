package model

import (
	"go_gin/common/global"
	"go_gin/tool"
)

type Login struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Code     string `json:"code" form:"code"`
	CodeId   string `json:"code_id" form:"code_id"`
}

func (u *Login) GetUser() (user User, e error) {
	e = global.Db.Table("sys_user").Where("username = ? ", u.Username).First(&user).Error
	if e != nil {
		return
	}
	_, e = tool.CompareHashAndPassword(user.Password, u.Password)
	if e != nil {
		return
	}
	return
}
