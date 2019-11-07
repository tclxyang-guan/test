package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"test/models"
	"test/repo"
	"test/utils"
)

type MenuService interface {
	MenuCreate(menu models.Menu) *models.Result
	MenuUpdate(m map[string]interface{}) *models.Result
	MenuDel(ids []interface{}, force bool) *models.Result
	MenuPage(m map[string]interface{}) *models.Result
}
type menuService struct {
	repo         repo.MenuRepo
	roleMenuRepo repo.RoleMenuRepo
}

func NewMenuService() MenuService {
	return &menuService{repo.NewMenuRepo(), repo.NewRoleMenuRepo()}
}

func (c *menuService) MenuCreate(menu models.Menu) *models.Result {
	ms := c.repo.MenuRepeat(0, menu.MenuName)
	if len(ms) > 0 {
		log.Error("菜单名称重复")
		return models.GetResult("", "菜单名称重复", errors.New("菜单名称重复"))
	}
	err := c.repo.MenuCreate(&menu)
	if err != nil {
		log.Error("菜单创建失败", menu)
	}
	return models.GetResult(menu, "创建成功", err)
}
func (c *menuService) MenuUpdate(m map[string]interface{}) *models.Result {
	ms := c.repo.MenuRepeat(cast.ToUint(m["id"]), cast.ToString(m["menu_name"]))
	if len(ms) > 0 {
		log.Println("菜单名称重复")
		return models.GetResult("", "菜单名称重复", errors.New("菜单名称重复"))
	}
	err := c.repo.MenuUpdate(m)
	if err != nil {
		go log.WithFields(m).Error("菜单修改失败")
	}
	return models.GetResult("修改成功", "修改失败", err)
}
func (c *menuService) MenuDel(ids []interface{}, force bool) *models.Result {
	if force { //不管菜单是否使用连带删除
		err := c.repo.MenuDel(ids)
		if err != nil {
			log.Error("菜单删除失败", ids)
			return models.GetResult("", "删除失败", err)
		}
	} else { //询问方式 查询菜单是否使用
		rms := c.roleMenuRepo.RoleMenuColumn([]interface{}{ids}, "menu_id in (?)")
		if len(rms) > 0 {
			log.Error("该菜单已被使用", ids)
			return models.GetResult("", "", errors.New("该菜单已被使用"))
		}
		err := c.repo.MenuDel(ids)
		if err != nil {
			log.Error("删除失败", ids)
			return models.GetResult("", "删除失败", err)
		}
	}
	return models.GetResult("删除成功", "", nil)
}
func (c *menuService) MenuPage(m map[string]interface{}) *models.Result {
	var page models.Page

	if err := utils.DataToAnyData(m["page"], &page); err != nil {
		log.WithFields(m).Error("菜单分页参数错误page")
		return models.GetResult("", "参数错误", err)
	}
	menu := models.Menu{}
	if _, ok := m["menu"]; ok {
		if err := utils.DataToAnyData(m["menu"], &menu); err != nil {
			log.WithFields(m).Error("菜单分页参数错误menu")
			return models.GetResult("", "参数错误", err)
		}
	}
	count, ms := c.repo.MenuPage(page, menu)
	m1 := map[string]interface{}{
		"Count": count,
		"Data":  ms,
	}
	return models.GetResult(m1, "", nil)
}
