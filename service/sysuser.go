package service

import (
	"errors"
	"go_gin/common/global"
	"go_gin/model"
	"go_gin/model/request"
	"go_gin/tool"
	"gorm.io/gorm"
	"log"
)

//@description 用户登陆
func Login(u *model.User) (error, model.User) {
	var user model.User
	err := global.Db.Where("username = ?", u.Username).First(&user).Error
	if err != nil {
		return err, user
	}
	ok, _ := tool.CompareHashAndPassword(user.Password, u.Password)
	if !ok {
		return errors.New("密码错误"), user
	}
	//登陆成功,处理用户密码
	user.Password = ""
	return nil, user
}

//@desc 用户注册
//@params u model.User
//@return error model.User
func RegisterUser(u model.User) (error, model.User) {
	var user model.User
	//判断用户名是否已存在
	if !errors.Is(global.Db.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名已存在"), user
	}
	//生成密码
	err := u.Encrypt()
	log.Println("password", u.Password)
	if err != nil {
		return errors.New("生成密码错误"), user
	}
	err = global.Db.Create(&u).Error
	return err, u
}

//@desc 修改用户密码
//@params u *model.User newPassword
//@return error model.User
func ChangePassword(u *model.User, newPassword string) (error, model.User) {
	var user model.User
	//验证密码
	if err := global.Db.Where("id = ? ", u.ID).First(&user).Error; err != nil {
		return errors.New("用户不存在"), user
	}
	_, err := tool.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		return errors.New("原密码不正确"), user
	}
	user.Password = newPassword
	if err := user.Encrypt(); err != nil {
		return errors.New("生产新密码错误"), user
	}
	//更新原密码
	//err = global.Db.Save(&user).Error
	err = global.Db.Model(&user).Update("password", user.Password).Error
	return err, user
}

//@desc 分业获取用户数据
func GetUserList(params request.PageInfo) (err error, list interface{}, total int64) {
	limit := params.PageSize
	offset := params.PageSize * (params.Page - 1)
	dbQuery := global.Db.Model(&model.User{})
	var userList []model.User
	err = dbQuery.Count(&total).Error
	err = dbQuery.Limit(limit).Offset(offset).Find(&userList).Error
	if len(userList) >= 1 {
		for key, value := range userList {
			err, userAuth := GetUserAuth(int(value.ID))
			if err != nil {
				continue
			}
			userList[key].AuthorityId = userAuth.AuthorityId
		}
	}

	return err, userList, total
}

//@desc 删除用户
func DeleteUser(id int) error {
	var user model.User
	err := global.Db.Where("id = ?", id).Delete(&user).Error

	return err
}

//@desc 更新用户信息
func SetUserInfo(user model.User) (error, model.User) {
	err := global.Db.Updates(&user).Error

	return err, user
}

func GetUserAuth(userId int) (error, model.UserAuth) {
	var userAuth model.UserAuth
	err := global.Db.Where("user_id = ?", userId).First(&userAuth).Error

	return err, userAuth
}
