package models

import "github.com/jinzhu/gorm"

/*
用户表
*/
type User struct {
	gorm.Model
	UserName string `gorm:"type:varcha(50);not null;unique" json:"user_name"` //用户名
	Password string `json:"password"`                                         //密码
	State    bool   `gorm:"default:1" json:"state"`                           //true启用
	OtherID  uint   `json:"other_id"`
	Other    Other  //其它对象例如 部门学校等等
}

/*
用户角色表
*/
type UserRole struct {
	gorm.Model
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}
