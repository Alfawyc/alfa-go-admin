package global

import "gorm.io/gorm"

var (
	Source string
	Driver string
	DBName string
)

var Db *gorm.DB
