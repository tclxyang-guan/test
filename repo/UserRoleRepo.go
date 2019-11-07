package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"test/datasource"
	"test/models"
)

type UserRoleRepo interface {
	UserRoleByColumn(userRole models.UserRole) []models.UserRole
	UserRoleColumn(Value []interface{}, Column ...string) []models.UserRole
}

func NewUserRoleRepo() UserRoleRepo {
	return &userRoleRepo{datasource.GetDB()}
}

type userRoleRepo struct {
	db *gorm.DB
}

//用户角色查询
func (c *userRoleRepo) UserRoleByColumn(UserRole models.UserRole) (urs []models.UserRole) {
	c.db.Where(&UserRole).Find(&urs)
	return
}

//用户角色查询
func (c *userRoleRepo) UserRoleColumn(Value []interface{}, Column ...string) (urs []models.UserRole) {
	db := c.db
	for i, v := range Column {
		if cast.ToString(Value[i]) != "" && cast.ToString(Value[i]) != "0" {
			db = db.Where(v, Value[i])
		}
	}
	db.Find(&urs)
	return
}
