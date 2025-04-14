package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"
	"imooc_product/backend/web/controllers"
	"imooc_product/common"
	"imooc_product/repositories"
	"imooc_product/services"
)

func main() {
	// 1. 创建iris实例
	app := iris.New()
	// 2. 设置debug模式
	app.Logger().SetLevel("debug")
	// 3. 设置注册模板
	template := iris.HTML(
		"./backend/web/views", ".html").Layout(
		"shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. 设置模板目标
	app.HandleDir("/assets", "./backend/web/assets")
	// 出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 5. 注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx,orderService)
	order.Handle(new(controllers.OrderController))

	// 6. 启动服务
	app.Run(
		iris.Addr("127.0.0.1:8080"),
		//iris.WithoutVersionChecker,
		// 忽略服务器错误
		iris.WithoutServerError(iris.ErrServerClosed),
		// 尽可能优化
		iris.WithOptimizations,
	)
}

// go mod init
// 编辑文件
// set GO111MODULE=on
// go mod tidy
