package models

/*
菜单表
*/
type Menu struct {
	Model
	Menu         []*Menu //下级菜单切片
	Level        int     `json:"level"`         //等级
	Seq          int     `json:"seq"`           //排序 从小到大排序
	MenuName     string  `json:"menu_name"`     //菜单名称
	MenuUrl      string  `json:"menu_url"`      //菜单路径
	InterfaceUrl string  `json:"interface_url"` //(功能)的接口路径
	RouteUrl     string  `json:"route_url"`     //前端路由路径
	Superior     uint    `json:"superior"`      //菜单上级编号
	Icon         string  `json:"icon"`          //菜单图标地址
	Operation    string  `json:"operation"`     //操作名称
	Type         bool    `json:"type"`          //true菜单  false功能
}
