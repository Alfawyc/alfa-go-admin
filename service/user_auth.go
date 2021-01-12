package service

import (
	"go_gin/common/global"
	"go_gin/model"
	"log"
)

func InsertUserAuth(userAuths model.UserAuth) error {
	//var userAuth []model.UserAuth
	log.Println(userAuths)
	/*for _, value := range userAuths.AuthorityIds {
		temp := model.UserAuth{UserId: userAuths.UserId, AuthorityId: value}
		userAuth = append(userAuth, temp)
	}*/
	err := global.Db.Create(&userAuths).Error

	return err
}
