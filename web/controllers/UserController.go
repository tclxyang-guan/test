package controllers

import (
	"errors"
	"github.com/kataras/iris"
	"test/models"
	"test/services"
)

type UserController struct {
	Ctx     iris.Context
	service services.UserService
}

func NewUserController() *UserController {
	return &UserController{service: services.NewUserService()}
}

/*
PostLogin
用户登录
*/
func (c *UserController) PostLogin() (result *models.Result) {
	var user models.User
	c.Ctx.ReadJSON(&user)
	return c.service.UserLogin(user)
}

/*
PostCreate
新增用户
*/
func (c *UserController) PostCreate() (result *models.Result) {
	var user models.User
	c.Ctx.ReadJSON(&user)
	return c.service.UserCreate(user)
}

/*
PostUpdate
修改用户
*/
func (c *UserController) PostUpdate() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	return c.service.UserUpdate(m)
}

/*
PostDel
删除用户
*/
func (c *UserController) PostDel() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	if v, ok := m["ids"].([]interface{}); !ok {
		return models.GetResult("", "参数错误", errors.New("参数错误"))
	} else {
		return c.service.UserDel(v)
	}
}

/*
PostPage
用户分页
*/
func (c *UserController) PostPage() (result *models.Result) {
	var m map[string]interface{}
	c.Ctx.ReadJSON(&m)
	return c.service.UserPage(m)
}
