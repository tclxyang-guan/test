package models

import "github.com/jinzhu/gorm"

/*
角色表
*/
type Role struct {
	gorm.Model
	RoleName string `json:"role_name"`              //角色名
	Describe string `json:"describe"`               //描述
	State    bool   `gorm:"default:1" json:"state"` //1==true 正常 0==false 禁用
}

/*
角色关联菜单表
*/
type RoleMenu struct {
	gorm.Model
	RoleID uint `json:"role_id"`
	MenuID uint `json:"menu_id"`
}
