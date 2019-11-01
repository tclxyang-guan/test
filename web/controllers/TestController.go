package controllers

import (
	"github.com/1324204490/tool/tool/model"
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
	"test/datasource"
	"test/models"
)

type TestController struct {
	Ctx iris.Context
}

/*
PostInsert
新增一个广告
*/
func (c *TestController) GetInsert() (result *model.Result) {
	log.Error("插入数据")
	datasource.GetDB().Create(&models.Test{Age: 10, Name: "dfa"})
	datasource.GetDB().Model(&models.Test{}).Where("id=?", 1).Update("age", 10)
	datasource.GetDB().Where("id=?", 1).Find(&models.Test{})
	return
}
