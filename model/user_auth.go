package model

import (
	"gorm.io/gorm"
)

type UserAuth struct {
	UserId      int `json:"user_id" form:"user_id" gorm:"user_id"`
	AuthorityId int `json:"authority_id" form:"authority_id" gorm:"authority_id"`
}

func (u UserAuth) TableName() string {
	return "user_auth"
}

func (u *UserAuth) BeforeCreate(db *gorm.DB) error {
	err := db.Where("user_id = ? ", u.UserId).Delete(&UserAuth{}).Error
	return err
}
