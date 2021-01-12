package request

type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type GetById struct {
	Id int64 `json:"id" form:"id"`
}
