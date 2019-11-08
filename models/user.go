package models

/*
用户表
*/
type User struct {
	Model
	UserName string `gorm:"type:varchar(50);not null;unique" json:"user_name"` //用户名
	Password string `json:"password"`                                          //密码
	State    bool   `json:"state"`                                             //true启用
	OtherID  uint   `json:"other_id"`
	Other    Other  `gorm:"-"` //其它对象例如 部门学校等等
}

/*
用户角色表
*/
type UserRole struct {
	Model
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}
