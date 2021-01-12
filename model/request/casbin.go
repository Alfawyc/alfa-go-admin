package request

type CasbinInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type CasbinReceive struct {
	AuthorityId string       `json:"authority_id" binding:"required"`
	CasbinInfos []CasbinInfo `json:"casbin_infos,omitempty" binding:"required"`
}
