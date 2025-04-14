package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"web/GoWeb商品秒杀系统/imooc-iris/web/controllers"
)

func main()  {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// 注册视图
	app.RegisterView(iris.HTML("./web/views", ".html"))
	// 注册控制器
	mvc.New(app.Party("/hello")).Handle(new(controllers.MovieController))
	// 启动项目
	app.Run(iris.Addr("127.0.0.1:8080"))
}
