package model

import (
	"errors"
	"go_gin/common/global"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Username    string `json:"username" form:"username" gorm:"username"`
	Password    string `json:"password" form:"password"  gorm:"password"`
	NickName    string `json:"nickname" form:"nickname" gorm:"nick_name"`
	Phone       string `json:"phone" form:"phone" gorm:"phone"`
	Avatar      string `json:"avatar" form:"avatar" gorm:"avatar"`
	Sex         string `json:"sex" form:"sex" gorm:"sex"`
	Email       string `json:"email" form:"email" gorm:"email"`
	DeptId      string `json:"dept_id" form:"dept_id" gorm:"dept_id"`
	CreatedBy   string `json:"created_by,omitempty" form:"created_by" gorm:"created_by"`
	UpdatedBy   string `json:"updated_by,omitempty" form:"updated_by" gorm:"updated_by"`
	Remark      string `json:"remark" form:"remark" gorm:"remark"`
	Status      string `json:"status" form:"status" form:"status"`
	AuthorityId int    `json:"authority_id" form:"authority_id" gorm:"-"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Insert() (id int, err error) {
	if err = u.Encrypt(); err != nil {
		return
	}
	//check用户名
	var count int64
	global.Db.Table(u.TableName()).Where("username = ? ", u.Username).Count(&count)
	if count > 0 {
		err = errors.New("用户名已存在")
		return
	}
	//添加数据
	if err = global.Db.Table(u.TableName()).Create(&u).Error; err != nil {
		return
	}
	id = int(u.ID)
	return
}

func (u *User) Encrypt() (err error) {
	if u.Password == "" {
		return
	}
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost); err != nil {
		return
	}
	u.Password = string(hash)
	return
}

func (u *User) GetOne() (User, error) {
	data := User{}
	if err := global.Db.First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
