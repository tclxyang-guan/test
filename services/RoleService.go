package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"test/models"
	"test/repo"
	"test/utils"
)

type RoleService interface {
	RoleCreate(role models.Role) *models.Result
	RoleUpdate(m map[string]interface{}) *models.Result
	RoleDel(id []interface{}, force bool) *models.Result
	RolePage(m map[string]interface{}) *models.Result
}
type roleService struct {
	repo         repo.RoleRepo
	userRoleRepo repo.UserRoleRepo
}

func NewRoleService() RoleService {
	return &roleService{repo.NewRoleRepo(), repo.NewUserRoleRepo()}
}

func (c *roleService) RoleCreate(role models.Role) *models.Result {
	ms := c.repo.RoleRepeat(0, role.RoleName)
	if len(ms) > 0 {
		go log.WithFields(utils.StructToMap(role)).Error("角色名称重复")
		return models.GetResult("", "角色名称重复", errors.New("角色名称重复"))
	}
	err := c.repo.RoleCreate(&role)
	if err != nil {
		go log.WithFields(utils.StructToMap(role)).Error("角色创建失败", err.Error())
	}
	return models.GetResult(role, "创建失败", err)
}
func (c *roleService) RoleUpdate(m map[string]interface{}) *models.Result {
	ms := c.repo.RoleRepeat(cast.ToUint(m["id"]), cast.ToString(m["role_name"]))
	if len(ms) > 0 {
		go log.WithFields(m).Error("角色名称重复")
		return models.GetResult("", "角色名称重复", errors.New("角色名称重复"))
	}
	err := c.repo.RoleUpdate(m)
	if err != nil {
		go log.WithFields(m).Error("角色修改失败", err.Error())
	}
	return models.GetResult("修改成功", "修改失败", err)
}
func (c *roleService) RoleDel(ids []interface{}, force bool) *models.Result {
	if force { //不管角色是否使用  连带删除
		err := c.repo.RoleDel(ids)
		if err != nil {
			go log.WithFields(map[string]interface{}{"ids": ids, "force": force}).Error("删除失败", err)
			return models.GetResult("", "删除失败", err)
		}
	} else { //询问方式 查询菜单是否使用
		rms := c.userRoleRepo.UserRoleColumn([]interface{}{ids}, "role_id in (?)")
		if len(rms) > 0 {
			go log.WithFields(map[string]interface{}{"ids": ids, "force": force}).Error("该菜单已被使用")
			return models.GetResult("", "", errors.New("该菜单已被使用"))
		}
		err := c.repo.RoleDel(ids)
		if err != nil {
			go log.WithFields(map[string]interface{}{"ids": ids, "force": force}).Error("删除失败", err)
			return models.GetResult("", "删除失败", err)
		}
	}
	return models.GetResult("删除成功", "", nil)
}
func (c *roleService) RolePage(m map[string]interface{}) *models.Result {
	var page models.Page

	if err := utils.DataToAnyData(m["page"], &page); err != nil {
		go log.WithFields(m).Error("菜单分页参数错误page", err)
		return models.GetResult("", "参数错误", err)
	}
	role := models.Role{}
	if _, ok := m["role"]; ok {
		if err := utils.DataToAnyData(m["role"], &role); err != nil {
			go log.WithFields(m).Error("菜单分页参数错误role", err)
			return models.GetResult("", "参数错误", err)
		}
	}
	count, ms := c.repo.RolePage(page, role)
	m1 := map[string]interface{}{
		"Count": count,
		"Data":  ms,
	}
	return models.GetResult(m1, "", nil)
}
