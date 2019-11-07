package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"test/datasource"
	"test/models"
)

type RoleMenuRepo interface {
	RoleMenuByColumn(roleMenu models.RoleMenu) []models.RoleMenu
	RoleMenuColumn(Value []interface{}, Column ...string) []models.RoleMenu
}

func NewRoleMenuRepo() RoleMenuRepo {
	return &roleMenuRepo{datasource.GetDB()}
}

type roleMenuRepo struct {
	db *gorm.DB
}

//角色菜单查询
func (c *roleMenuRepo) RoleMenuByColumn(roleMenu models.RoleMenu) (rms []models.RoleMenu) {
	c.db.Where(&roleMenu).Find(&rms)
	return
}

//角色菜单查询
func (c *roleMenuRepo) RoleMenuColumn(Value []interface{}, Column ...string) (rms []models.RoleMenu) {
	db := c.db
	for i, v := range Column {
		if cast.ToString(Value[i]) != "" && cast.ToString(Value[i]) != "0" {
			db = db.Where(v, Value[i])
		}
	}
	db.Find(&rms)
	return
}
