package controllers

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/spf13/cast"
	"test/models"
	"test/services"
)

type RoleController struct {
	Ctx     iris.Context
	service services.RoleService
}

func NewRoleController() *RoleController {
	return &RoleController{service: services.NewRoleService()}
}

/*
PostCreate
新增角色
*/
func (c *RoleController) PostCreate() (result *models.Result) {
	var role models.Role
	err := c.Ctx.ReadJSON(&role)
	if err != nil {
		return models.GetResult("", "参数错误", err)
	}
	return c.service.RoleCreate(role)
}

/*
PostUpdate
修改角色
*/
func (c *RoleController) PostUpdate() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	return c.service.RoleUpdate(m)
}

/*
PostDel
删除角色
*/
func (c *RoleController) PostDel() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	if v, ok := m["ids"].([]interface{}); !ok {
		return models.GetResult("", "参数错误", errors.New("参数错误"))
	} else {
		return c.service.RoleDel(v, cast.ToBool(m["force"]))
	}
}

/*
PostPage
角色分页
*/
func (c *RoleController) PostPage() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	return c.service.RolePage(m)
}
