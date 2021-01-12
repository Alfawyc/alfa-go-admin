package model

type Api struct {
	BaseModel
	Path        string `json:"path" form:"path" gorm:"path"`
	Description string `json:"description" form:"description" gorm:"description"`
	Group       string `json:"group" form:"group" gorm:"group"`
	Method      string `json:"method" form:"method" gorm:"method"`
}

func (api Api) TableName() string {
	return "apis"
}
