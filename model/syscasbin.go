package model

type CasbinModel struct {
	Ptype       string `json:"ptype" gorm:"p_type"`
	AuthorityId string `json:"authority_id" gorm:"v0"`
	Path        string `json:"path" gorm:"v1"`
	Method      string `json:"method" gorm:"v2"`
}
