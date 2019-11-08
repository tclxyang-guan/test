package repo

import (
	"github.com/jinzhu/gorm"
	"test/datasource"
	"test/models"
	"test/utils"
)

type UserRepo interface {
	UserCreate(User *models.User) error
	UserUpdate(m map[string]interface{}) error
	UserDel(ids []interface{}) error
	UserPage(page models.Page, User models.User) (int, []models.User)
	UserByColumn(user models.User) []models.User
	UserRepeat(ID uint, UserName string) []models.User
}

func NewUserRepo() UserRepo {
	return &userRepo{datasource.GetDB()}
}

type userRepo struct {
	db *gorm.DB
}

//用户创建
func (c *userRepo) UserCreate(User *models.User) error {
	return c.db.Create(User).Error
}

//用户修改
func (c *userRepo) UserUpdate(m map[string]interface{}) error {
	return c.db.Model(models.User{}).Updates(m).Error
}

//用户删除
func (c *userRepo) UserDel(ids []interface{}) error {
	tx := c.db.Begin()
	flag := false
	defer utils.Defer(tx, &flag)
	//删除用户
	err := tx.Unscoped().Where("id in (?)", ids).Delete(models.User{}).Error
	if err != nil {
		return err
	}
	//删除中间表user_role
	err = tx.Unscoped().Where("user_id in (?)", ids).Delete(models.UserRole{}).Error
	if err != nil {
		return err
	}
	flag = true
	return nil
}

//用户分页查询
func (c *userRepo) UserPage(page models.Page, user models.User) (count int, rs []models.User) {
	c.db.Where(&user).Order(page.Sort).Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(&rs).
		Offset(-1).Limit(-1).Count(&count)
	return
}

//根据列查询用户
func (c *userRepo) UserByColumn(User models.User) (rs []models.User) {
	c.db.Where(&User).Find(&rs)
	return
}

//根据用户查重
func (c *userRepo) UserRepeat(ID uint, UserName string) (rs []models.User) {
	c.db.Where("id <> ?", ID).Where("User_name=?", UserName).Find(&rs)
	return
}
