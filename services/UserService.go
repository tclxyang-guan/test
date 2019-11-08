package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"test/middleware"
	"test/models"
	"test/repo"
	"test/utils"
)

type UserService interface {
	UserLogin(user models.User) *models.Result
	UserCreate(user models.User) *models.Result
	UserUpdate(m map[string]interface{}) *models.Result
	UserDel(ids []interface{}) *models.Result
	UserPage(m map[string]interface{}) *models.Result
}
type userService struct {
	repo         repo.UserRepo
	userRoleRepo repo.UserRoleRepo
	menuRepo     repo.MenuRepo
}

func NewUserService() UserService {
	return &userService{repo.NewUserRepo(), repo.NewUserRoleRepo(), repo.NewMenuRepo()}
}
func (c *userService) UserLogin(user models.User) *models.Result {
	//根据用户名查询用户信息
	user.Password = utils.Md5(user.Password)
	us := c.repo.UserByColumn(user)
	if len(us) == 0 {
		go log.WithFields(map[string]interface{}{"user_name": user.UserName}).Error("用户名或密码错误")
		return models.GetResult("", "用户名或密码错误", errors.New("用户名或密码错误"))
	}
	u := us[0]
	//获取用户的角色
	urs := c.userRoleRepo.UserRoleByColumn(models.UserRole{UserID: u.ID})
	if len(urs) == 0 {
		go log.WithFields(map[string]interface{}{"user_name": user.UserName}).Error("此用户没有角色,请联系管理员")
		return models.GetResult("", "此用户没有角色,请联系管理员", errors.New("此用户没有角色,请联系管理员"))
	}
	var roleIDs []uint
	for i := range urs {
		roleIDs = append(roleIDs, urs[i].RoleID)
	}
	//获取用户的菜单 true只查询菜单不查询功能
	ms := c.menuRepo.MenuByRole(roleIDs, true)
	if len(ms) == 0 {
		go log.WithFields(map[string]interface{}{"user_name": user.UserName}).Error("此用户没有菜单,请联系管理员")
		return models.GetResult("", "此用户没有菜单,请联系管理员", errors.New("此用户没有菜单,请联系管理员"))
	}
	mm := make(map[uint]*models.Menu)
	for i := range ms {
		mm[ms[i].ID] = &ms[i]
	}
	var newMenus []*models.Menu
	for i, v := range ms {
		if v.Superior == 0 {
			newMenus = append(newMenus, mm[ms[i].ID])
		} else {
			parent := mm[ms[i].Superior]
			if parent != nil {
				(*parent).Menu = append((*parent).Menu, &ms[i])
			}
		}
	}
	m := make(map[string]interface{})
	token := middleware.GenerateToken(user)
	m["Token"] = token
	m["Data"] = user
	m["Menu"] = newMenus
	return models.GetResult(m, "", nil)
}
func (c *userService) UserCreate(user models.User) *models.Result {
	ms := c.repo.UserRepeat(0, user.UserName)
	if len(ms) > 0 {
		go log.WithFields(utils.StructToMap(user)).Error("用户名称重复")
		return models.GetResult("", "用户名称重复", errors.New("用户名称重复"))
	}
	user.Password = utils.Md5(user.Password)
	err := c.repo.UserCreate(&user)
	if err != nil {
		go log.WithFields(utils.StructToMap(user)).Error("用户创建失败")
	}
	return models.GetResult(user, "创建失败", err)
}
func (c *userService) UserUpdate(m map[string]interface{}) *models.Result {
	ms := c.repo.UserRepeat(cast.ToUint(m["id"]), cast.ToString(m["user_name"]))
	if len(ms) > 0 {
		go log.WithFields(m).Error("用户名称重复")
		return models.GetResult("", "用户名称重复", errors.New("用户名称重复"))
	}
	if v := cast.ToString(m["password"]); v != "" {
		m["password"] = utils.Md5(v)
	}
	err := c.repo.UserUpdate(m)
	if err != nil {
		go log.WithFields(m).Error("用户修改失败")
	}
	return models.GetResult("修改成功", "修改失败", err)
}
func (c *userService) UserDel(ids []interface{}) *models.Result {
	err := c.repo.UserDel(ids)
	if err != nil {
		go log.WithFields(map[string]interface{}{"ids": ids}).Error("用户删除失败")
	}
	return models.GetResult("删除成功", "删除失败", err)
}
func (c *userService) UserPage(m map[string]interface{}) *models.Result {
	var page models.Page
	if err := utils.DataToAnyData(m["page"], &page); err != nil {
		go log.WithFields(m).Error("参数错误", err.Error())
		return models.GetResult("", "参数错误", err)
	}
	User := models.User{}
	if _, ok := m["user"]; ok {
		if err := utils.DataToAnyData(m["user"], &User); err != nil {
			go log.WithFields(m).Error("参数错误", err.Error())
			return models.GetResult("", "参数错误", err)
		}
	}
	count, ms := c.repo.UserPage(page, User)
	m1 := map[string]interface{}{
		"Count": count,
		"Data":  ms,
	}
	return models.GetResult(m1, "", nil)
}
