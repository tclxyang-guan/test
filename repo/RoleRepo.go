package repo

import (
	"github.com/jinzhu/gorm"
	"test/datasource"
	"test/models"
	"test/utils"
)

type RoleRepo interface {
	RoleCreate(Role *models.Role) error
	RoleUpdate(m map[string]interface{}) error
	RoleDel(ids []interface{}) error
	RolePage(page models.Page, Role models.Role) (int, []models.Role)
	//RoleByColumn(Role models.Role)[]models.Role
	RoleRepeat(ID uint, RoleName string) []models.Role
}

func NewRoleRepo() RoleRepo {
	return &roleRepo{datasource.GetDB()}
}

type roleRepo struct {
	db *gorm.DB
}

//角色创建
func (c *roleRepo) RoleCreate(Role *models.Role) error {
	return c.db.Create(Role).Error
}

//角色修改
func (c *roleRepo) RoleUpdate(m map[string]interface{}) error {
	return c.db.Model(models.Role{}).Updates(m).Error
}

//角色删除
func (c *roleRepo) RoleDel(ids []interface{}) error {
	tx := c.db.Begin()
	flag := false
	defer utils.Defer(tx, &flag)
	//删除角色
	err := tx.Unscoped().Where("id in (?)", ids).Delete(models.Role{}).Error
	if err != nil {
		return err
	}
	//删除中间表user_role
	err = tx.Unscoped().Where("role_id in (?)", ids).Delete(models.UserRole{}).Error
	if err != nil {
		return err
	}
	//删除中间表role_menu
	err = tx.Unscoped().Where("role_id in (?)", ids).Delete(models.RoleMenu{}).Error
	if err != nil {
		return err
	}
	flag = true
	return nil
}

//角色分页查询
func (c *roleRepo) RolePage(page models.Page, role models.Role) (count int, rs []models.Role) {
	c.db.Where(&role).Order(page.Sort).Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(&rs).
		Offset(-1).Limit(-1).Count(&count)
	return
}

//根据列查询角色
func (c *roleRepo) RoleByColumn(role models.Role) (rs []models.Role) {
	c.db.Where(&role).Find(&rs)
	return
}

//根据角色查重
func (c *roleRepo) RoleRepeat(ID uint, RoleName string) (rs []models.Role) {
	c.db.Where("id <> ?", ID).Where("role_name=?", RoleName).Find(&rs)
	return
}
