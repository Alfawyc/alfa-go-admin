package service

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go_gin/common/global"
	"go_gin/model"
	"go_gin/model/request"
	"log"
)

func Casbin() *casbin.Enforcer {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", global.Vp.Get("mysql.username"), global.Vp.Get("mysql.password"), global.Vp.Get("mysql.host"), global.Vp.Get("mysql.port"), global.Vp.Get("mysql.database"))
	a, err := gormadapter.NewAdapter("mysql", dataSourceName, true)
	if err != nil {
		log.Fatalln("gorm NewAdapter error , ", err.Error())
	}
	e, err := casbin.NewEnforcer("rabc_model.config", a)
	if err != nil {
		log.Fatalln("casbin enforcer error, ", err.Error())
	}
	//todo 引入自定义规则
	_ = e.LoadPolicy()

	return e
}

func UpdateCasbin(authorityId string, casbinInfos []request.CasbinInfo) error {
	ClearCasbin(0, authorityId)
	rules := [][]string{}
	for _, value := range casbinInfos {
		cm := model.CasbinModel{
			Ptype:       "p",
			AuthorityId: authorityId,
			Path:        value.Path,
			Method:      value.Method,
		}
		rules = append(rules, []string{cm.AuthorityId, cm.Path, cm.Method})
	}
	e := Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api添加失败")
	}

	return nil
}

func ClearCasbin(v int, p ...string) bool {
	e := Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//@description: 获取权限列表
func GetAuthorityIdPolicy(authorityId string) []request.CasbinInfo {
	e := Casbin()
	rules := e.GetFilteredPolicy(0, authorityId)
	var paths []request.CasbinInfo
	for _, value := range rules {
		temp := request.CasbinInfo{
			Path:   value[1],
			Method: value[2],
		}
		paths = append(paths, temp)
	}

	return paths
}
