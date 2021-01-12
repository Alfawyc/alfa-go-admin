package model

import (
	"gorm.io/gorm"
	"time"
)

type Auth struct {
	AuthorityId   int            `json:"authority_id" form:"authority_id" gorm:"authority_id"`
	AuthorityName string         `json:"authority_name" form:"authority_name" gorm:"authority_name"`
	ParentID      int            `json:"parent_id" form:"parent_id" gorm:"parent_id"`
	Status        string         `json:"status" form:"status" gorm:"status"`
	CreatedBy     int            `json:"created_by" form:"created_by" gorm:"created_by"`
	Remark        string         `json:"remark" form:"remark" gorm:"remark"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Children      []Auth         `json:"Children" gorm:"-"`
}

func (Auth) TableName() string {
	return "auth"
}
