package database

func SetUp() {
	//默认连接mysql
	var db = new(Mysql)
	db.SetUp()
}
