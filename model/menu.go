package model

import "go_gin/common/global"

type Menu struct {
	MenuId     int    `json:"menu_id" gorm:"primaryKey;autoIncrement"`
	MenuName   string `json:"menu_name" gorm:"size:128"`
	Title      string `json:"title" gorm:"size:128"`
	Icon       string `json:"icon" gorm:"size:128"`
	Path       string `json:"path" gorm:"size:128"`
	Paths      string `json:"paths" gorm:"size:128"`
	MenuType   string `json:"menu_type" gorm:"size:1"`
	Action     string `json:"action" gorm:"size:16"`
	Permission string `json:"permission" gorm:"size:255"`
	ParentId   int    `json:"parent_id"`
	NoCache    int    `json:"no_cache"`
	BreadCrumb string `json:"bread_crumb" gorm:"size:255"`
	Component  string `json:"component" gorm:"size:255"`
	Sort       int    `json:"sort"`
	Visible    string `json:"visible" gorm:"size:1"`
	CreateBy   string `json:"createBy" gorm:"size:128;"`
	UpdateBy   string `json:"updateBy" gorm:"size:128;"`
	IsFrame    string `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	DataScope  string `json:"dataScope" gorm:"-"`
	Params     string `json:"params" gorm:"-"`
	RoleId     int    `gorm:"-"`
	IsSelect   bool   `json:"is_select" gorm:"-"`
	Children   []Menu `json:"children" gorm:"-"`
	BaseModel
}

func (Menu) TableName() string {
	return "sys_menu"
}

type MenuLabel struct {
	Id       int         `json:"id" gorm:"-"`
	Label    string      `json:"label"`
	Children []MenuLabel `json:"children"`
}

func (e *Menu) SetMenuLabel() ([]MenuLabel, error) {
	//菜单列表
	menuList, err := e.Get()
	menuLabel := make([]MenuLabel, 0)
	for _, val := range menuList {
		if val.ParentId != 0 {
			//获取顶级菜单
			continue
		}
		e := MenuLabel{}
		e.Id = val.MenuId
		e.Label = val.Title
		menusInfo := DiguiMenuLabel(&menuList, e)
		menuLabel = append(menuLabel, menusInfo)
	}

	return menuLabel, err
}

func (e *Menu) Get() ([]Menu, error) {
	table := global.Db.Table(e.TableName())
	if e.MenuName != "" {
		table = table.Where("menu_name = ? ", e.MenuName)
	}
	data := make([]Menu, 0)
	err := table.Order("sort").Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func DiguiMenuLabel(menuList *[]Menu, menu MenuLabel) MenuLabel {
	list := *menuList
	min := make([]MenuLabel, 0)
	for _, val := range list {
		//获取当前顶级菜单下的子菜单
		if menu.Id != val.ParentId {
			continue
		}
		mi := MenuLabel{}
		mi.Id = val.MenuId
		mi.Label = val.Title
		mi.Children = []MenuLabel{}
		if val.MenuType != "F" {
			ms := DiguiMenuLabel(menuList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	menu.Children = min

	return menu
}
