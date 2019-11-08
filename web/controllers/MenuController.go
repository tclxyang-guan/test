package controllers

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/spf13/cast"
	"test/models"
	"test/services"
)

type MenuController struct {
	Ctx     iris.Context
	Service services.MenuService
}

func NewMenuController() *MenuController {
	return &MenuController{Service: services.NewMenuService()}
}

/*
PostCreate
新增菜单
*/
func (c *MenuController) PostCreate() (result *models.Result) {
	var menu models.Menu
	err := c.Ctx.ReadJSON(&menu)
	if err != nil {
		return models.GetResult("", "参数错误", err)
	}
	return c.Service.MenuCreate(menu)
}

/*
PostUpdate
修改菜单
*/
func (c *MenuController) PostUpdate() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	return c.Service.MenuUpdate(m)
}

/*
PostDel
删除菜单
*/
func (c *MenuController) PostDel() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	if v, ok := m["ids"].([]interface{}); !ok {
		return models.GetResult("", "参数错误", errors.New("参数错误"))
	} else {
		return c.Service.MenuDel(v, cast.ToBool(m["force"]))
	}
}

/*
PostPage
菜单分页
*/
func (c *MenuController) PostPage() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	return c.Service.MenuPage(m)
}
