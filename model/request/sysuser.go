package request

type Register struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Nickname string `json:"nickname" form:"nickname" binding:"required"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
	Sex      string `json:"sex" form:"sex"`
}

type Login struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Code     string `json:"code" form:"code" binding:"required"`
	CodeId   string `json:"code_id" form:"code_id" binding:"required"`
}

type ChangePassword struct {
	Password    string `json:"password" form:"password" binding:"required" validate:"min=4,max=8"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required" validate:"min=4,max=8"`
}
