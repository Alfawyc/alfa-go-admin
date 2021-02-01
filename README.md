# Alfa-go-admin 定时任务管理项目
基于 `go` , `vue`的后台管理系统
前端项目地址 `https://github.com/Alfawyc/vue-cli-admin.git`

执行 `cp config.toml.example config.toml` 复制配置文件,修改为对应的配置
## TODO LIST
- [x] 用户管理，角色管理
- [x] 权限管理（基于Casbin）
- [x] 任务添加，暂停，重置
- [ ] 任务调度修改为rpc方式调用
- [ ] jwt 修改为密钥和公钥进行加解密
- [ ] 增加前端页面角色菜单

## 程序使用的组件
+ Web框架 [GIN](https://github.com/gin-gonic/gin)
+ ORM [Gorm](https://github.com/go-gorm/gorm)
+ 任务调度[robfig/cron](https://github.com/robfig/cron)
+ 配置文件管理[Viper](https://github.com/spf13/viper)
+ 前端页面 vue-cli