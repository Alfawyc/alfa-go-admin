package global

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Source string
	Driver string
	DBName string
)

var (
	Db *gorm.DB
	Vp *viper.Viper
)
