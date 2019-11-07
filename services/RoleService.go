package services

import (
	"errors"
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
		return models.GetResult("", "角色名称重复", errors.New("角色名称重复"))
	}
	err := c.repo.RoleCreate(&role)
	return models.GetResult(role, "创建成功", err)
}
func (c *roleService) RoleUpdate(m map[string]interface{}) *models.Result {
	ms := c.repo.RoleRepeat(cast.ToUint(m["id"]), cast.ToString(m["role_name"]))
	if len(ms) > 0 {
		return models.GetResult("", "角色名称重复", errors.New("角色名称重复"))
	}
	err := c.repo.RoleUpdate(m)
	return models.GetResult("修改成功", "修改失败", err)
}
func (c *roleService) RoleDel(id []interface{}, force bool) *models.Result {
	if force { //不管角色是否使用  连带删除
		err := c.repo.RoleDel(id)
		if err != nil {
			return models.GetResult("", "删除失败", err)
		}
	} else { //询问方式 查询菜单是否使用
		rms := c.userRoleRepo.UserRoleColumn([]interface{}{id}, "role_id in (?)")
		if len(rms) > 0 {
			return models.GetResult("", "", errors.New("该菜单已被使用"))
		}
		err := c.repo.RoleDel(id)
		if err != nil {
			return models.GetResult("", "删除失败", err)
		}
	}
	return models.GetResult("删除成功", "", nil)
}
func (c *roleService) RolePage(m map[string]interface{}) *models.Result {
	var page models.Page

	if err := utils.DataToAnyData(m["page"], &page); err != nil {
		return models.GetResult("", "参数错误", err)
	}
	role := models.Role{}
	if _, ok := m["role"]; ok {
		if err := utils.DataToAnyData(m["role"], &role); err != nil {
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
