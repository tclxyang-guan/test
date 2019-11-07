package route

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"test/web/controllers"
)

func Route(app *iris.Application) {
	mvc.New(app.Party("/user")).Handle(controllers.NewUserController())
	mvc.New(app.Party("/role")).Handle(controllers.NewRoleController())
	mvc.New(app.Party("/menu")).Handle(controllers.NewMenuController())
}
