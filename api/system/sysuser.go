package system

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mojocn/base64Captcha"
	"go_gin/middleware"
	"go_gin/model"
	"go_gin/model/request"
	"go_gin/model/response"
	"go_gin/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetPostListHandler2 升级版帖子列表接口
// @Summary 用户登陆
// @Tags Base
// @Accept application/json
// @Produce json
// @Param data body request.Login true "用户名,密码,验证码,验证码id"
// @Success 200 {string} json "{"code":200 , "data":"" ,"message":"success" }"
// @Router /base/login [POST]
func Login(ctx *gin.Context) {
	var err error
	var LoginParams request.Login
	var store = base64Captcha.DefaultMemStore
	err = ctx.ShouldBind(&LoginParams)
	if err != nil {
		response.FailWithMessage("请输入用户名,密码及验证码"+err.Error(), ctx)
		return
	}
	//检查验证码
	ok := store.Verify(LoginParams.CodeId, LoginParams.Code, true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"message": "验证码错误"})
		return
	}
	u := model.User{Username: LoginParams.Username, Password: LoginParams.Password}
	err, user := service.Login(&u)
	if err != nil {
		response.FailWithMessage("登陆失败,用户名或密码错误", ctx)
		return
	}
	//登陆成功处理jwt token
	tokenNext(ctx, user)
}

//@desc 登陆后签发jwt
func tokenNext(ctx *gin.Context, user model.User) {
	//获取用户authorityId
	err, userAuth := service.GetUserAuth(int(user.ID))
	log.Println("用户权限id , ", userAuth.AuthorityId)
	j := &middleware.JWT{[]byte("alfa")}
	claims := middleware.CustomClaims{
		UserId:      int(user.ID),
		Username:    user.Username,
		Nickname:    user.NickName,
		AuthorityId: strconv.Itoa(userAuth.AuthorityId),
		BufferTime:  60 * 60 * 24, //缓冲时间1天  缓冲时间内会刷新token令牌 此时一个用户会存在两个令牌但前端只存一个
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, //过去时间7天
			Issuer:    "Alfa",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		response.FailWithMessage("获取token失败", ctx)
		return
	}
	//todo 多点登陆 redis jwt

	response.SuccessWithDetail(gin.H{"user": user, "token": token, "expire_at": claims.StandardClaims.ExpiresAt * 1000}, "登陆成功", ctx)
}

// @Summary 用户注册
// @Tags User
// @Produce json
// @Param data body request.Register true "用户信息"
// @Success 200 {string} json "{code:"200" , "data" :"" , "message":""}"
// @Router /user/register [POST]
func Register(ctx *gin.Context) {
	var registerForm request.Register
	err := ctx.ShouldBind(&registerForm)
	log.Println(registerForm)
	if err != nil {
		response.FailWithMessage("绑定数据失败"+err.Error(), ctx)
		return
	}
	user := &model.User{
		Username: registerForm.Username,
		NickName: registerForm.Nickname,
		Password: registerForm.Password,
		Phone:    registerForm.Phone,
		Email:    registerForm.Email,
		Sex:      registerForm.Sex,
		Avatar:   registerForm.Avatar,
	}
	err, userCreate := service.RegisterUser(*user)
	if err != nil {
		//用户名已存在
		log.Println(err.Error())
		response.FailWithMessage("注册失败"+err.Error(), ctx)

		return
	}
	userCreate.Password = ""
	response.SuccessWithDetail(gin.H{"message": "注册成功", "user": userCreate}, "success", ctx)
}

// @Summary 修改密码
// @Tags User
// @Produce json
// @Param data body request.ChangePassword true "旧密码,新密码"
// @Success 200 string json "{"code":200, "msg":"success"}"
// @Router /user/change-password [POST]
func ChangePassword(ctx *gin.Context) {
	var user request.ChangePassword
	//获取当前登陆用户
	userId := GetUserId(ctx)
	if err := ctx.ShouldBind(&user); err != nil {
		response.FailWithMessage("绑定数据失败", ctx)
		return
	}
	//验证旧密码
	u := model.User{Password: user.Password}
	u.ID = uint(userId)
	log.Println("正在修改用户id为", userId, "的用户密码")
	err, _ := service.ChangePassword(&u, user.NewPassword)
	if err != nil {
		log.Println(err.Error())
		response.FailWithMessage("修改失败", ctx)
		return
	}
	response.SuccessWithMessage("success", ctx)
}

//@Summary 用户列表
//@Tags User
//@Produce json
//@Param data query request.PageInfo false "页码,每页条数"
//@Success 200 string json "{"code":200 ,"data":"" , "msg":"success"}"
//@Router /user/user-list [GET]
func GetUserList(ctx *gin.Context) {
	var pageParams request.PageInfo
	_ = ctx.ShouldBind(&pageParams)
	if pageParams.PageSize == 0 {
		pageParams.PageSize = 10
	}
	if pageParams.Page == 0 {
		pageParams.Page = 1
	}
	//todo 验证数据
	err, list, total := service.GetUserList(pageParams)
	if err != nil {
		response.FailWithMessage("获取信息失败", ctx)
		return
	}
	response.SuccessWithDetail(response.PageResult{list, total, pageParams.Page, pageParams.PageSize}, "success", ctx)
}

//@Summary 删除用户
//@Tags User
//@Produce json
//@Param data body request.GetById true "用户id"
//@Success 200
//@Router /use/delete [POST]
func DeleteUser(ctx *gin.Context) {
	var params request.GetById
	_ = ctx.ShouldBind(&params)
	if params.Id == 0 {
		response.FailWithMessage("缺少参数", ctx)
		return
	}
	loginId := GetUserId(ctx)
	if loginId == int(params.Id) {
		response.FailWithMessage("自杀失败", ctx)
		return
	}
	err := service.DeleteUser(int(params.Id))
	if err != nil {
		response.FailWithMessage("删除失败"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("删除成功", ctx)
}

//从gin Context中获取jwt解析用户id
func GetUserId(ctx *gin.Context) int {
	if claims, exists := ctx.Get("claims"); !exists {
		log.Println("从gin的Context中获取从jwt解析用户id失败，请检查路由是否使用jwt中间件")
		return 0
	} else {
		waitUse := claims.(*middleware.CustomClaims)
		return waitUse.UserId
	}
}

//@Summary 更新用户信息
//@Tags User
//@Produce json
//@Param data body model.User false "用户信息"
//@Success 200
//@Router /user/set-info [POST]
func SetUserInfo(ctx *gin.Context) {
	//todo 获取用户再更新,批量跟新数据id必须存在
	var user model.User
	_ = ctx.ShouldBind(&user)

	err, _ := service.SetUserInfo(user)
	if err != nil {
		response.FailWithMessage("更新失败", ctx)
		return
	}

	response.SuccessWithMessage("success", ctx)
}

func InsertSysUser(ctx *gin.Context) {
	var user model.User
	err := ctx.BindWith(&user, binding.JSON)
	if err != nil {
		ctx.JSON(200, gin.H{"message": "错误的数据格式"})
		return
	}
	id, err := user.Insert()
	if err != nil {
		ctx.JSON(500, gin.H{"message": "添加用户失败" + err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"id": id, "message": "添加成功"})
}

//@Summary 设置用户权限
//@Tags User
//@Produce json
//@Param data body model.UserAuth true "用户角色"
//@Success 200
//@Router /user/user-auth [POST]
//todo 修改为用户与角色一对多
func UserAuth(ctx *gin.Context) {
	var userAuths model.UserAuth
	_ = ctx.ShouldBind(&userAuths)
	err := service.InsertUserAuth(userAuths)
	if err != nil {
		response.FailWithMessage("更新用户角色失败"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("更新用户角色成功", ctx)
}
