package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"test/datasource"
	"test/models"
	"test/utils"
)

type MenuRepo interface {
	MenuCreate(menu *models.Menu) error
	MenuUpdate(m map[string]interface{}) error
	MenuDel(ids []interface{}) error
	MenuPage(page models.Page, menu models.Menu) (int, []models.Menu)
	MenuByColumn(menu models.Menu) []models.Menu
	MenuRepeat(ID uint, MenuName string) []models.Menu
	MenuByRole(roleIds []uint, Type bool) (ms []models.Menu)
}

func NewMenuRepo() MenuRepo {
	return &menuRepo{datasource.GetDB()}
}

type menuRepo struct {
	db *gorm.DB
}

//菜单创建
func (c *menuRepo) MenuCreate(menu *models.Menu) error {
	tx := c.db.Begin()
	flag := false
	defer utils.Defer(tx, &flag)
	//查询菜单
	var seq int
	tx.Table("menu").Select("max(sort)").Where("type=?", menu.Type).Row().Scan(&seq)
	menu.Sort = seq + 1
	err := tx.Create(menu).Error
	if err != nil {
		return err
	}
	flag = true
	return nil
}

//菜单修改
func (c *menuRepo) MenuUpdate(m map[string]interface{}) error {
	return c.db.Model(models.Menu{}).Updates(m).Error
}

//菜单排序修改
func (c *menuRepo) MenuUpdateSeq(m map[string]interface{}) error {
	tx := c.db.Begin()
	flag := false
	defer utils.Defer(tx, &flag)
	err := tx.Model(&models.Menu{}).Where("id=?", cast.ToUint(m["id"])).Update("sort", cast.ToInt(m["sort1"])).Error
	if err != nil {
		return err
	}
	err = tx.Model(&models.Menu{}).Where("id=?", cast.ToUint(m["id1"])).Update("sort", cast.ToInt(m["sort"])).Error
	if err != nil {
		return err
	}
	flag = true
	return c.db.Model(models.Menu{}).Updates(m).Error
}

//菜单删除
func (c *menuRepo) MenuDel(ids []interface{}) error {
	tx := c.db.Begin()
	flag := false
	defer utils.Defer(tx, &flag)
	//删除菜单
	err := tx.Unscoped().Where("id in (?)", ids).Delete(models.Menu{}).Error
	if err != nil {
		return err
	}
	//删除role_menu中间表
	err = tx.Unscoped().Where("menu_id in (?)", ids).Delete(models.RoleMenu{}).Error
	if err != nil {
		return err
	}
	flag = true
	return nil
}

//菜单分页查询
func (c *menuRepo) MenuPage(page models.Page, menu models.Menu) (count int, ms []models.Menu) {
	c.db.Where(&menu).Order(page.Sort).Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(&ms).
		Offset(-1).Limit(-1).Count(&count)
	return
}

//根据列查询菜单
func (c *menuRepo) MenuByColumn(menu models.Menu) (ms []models.Menu) {
	c.db.Where(&menu).Find(&ms)
	return
}

//根据角色查询菜单以及功能(Type为true只查菜单)
func (c *menuRepo) MenuByRole(roleIds []uint, Type bool) (ms []models.Menu) {
	db := c.db.Where("id in (?)", c.db.Table("role_menu").Select("distinct(menu_id)").Where("role_id in (?)", roleIds).SubQuery()).Order("`level`,Sort")
	if Type {
		db = db.Where("type=?", 1)
	}
	db.Find(&ms)
	return
}

//根据菜单查重
func (c *menuRepo) MenuRepeat(ID uint, MenuName string) (ms []models.Menu) {
	c.db.Where("id <> ?", ID).Where("menu_name=?", MenuName).Find(&ms)
	return
}
