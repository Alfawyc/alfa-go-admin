package request

type UserAuthCreate struct {
	UserId      int `json:"user_id" form:"user_id" binding:"required"`
	AuthorityId int `json:"authority_id" form:"authority_id" binding:"required"`
}
