package database

import (
	"fmt"
	"go_gin/common/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Mysql struct {
}

func (e *Mysql) SetUp() {
	global.Source = e.GetConnect()
	db := mysql.New(mysql.Config{
		DSN: global.Source,
	})
	var err error
	global.Db, err = e.Open(db, &gorm.Config{})
	if err != nil {
		log.Fatalln("connect database error")
	} else {
		log.Println("connect database success")
	}
}

//打开数据库连接
func (e *Mysql) Open(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	/*dsn := "root:root@tcp(127.0.0.1:3307)/go_gin?charset=utf8mb4parseTime=True&loc=Local"
	dialector := mysql.Open(dsn)*/
	return gorm.Open(dialector, config)
}

//获取数据库连接
func (e *Mysql) GetConnect() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.Vp.Get("mysql.username"), global.Vp.Get("mysql.password"), global.Vp.Get("mysql.host"), global.Vp.Get("mysql.port"), global.Vp.Get("mysql.database"))
	fmt.Println(dsn)
	return dsn
}

//获取数据库驱动
func (e *Mysql) GetDriver() string {
	return "mysql"
}
