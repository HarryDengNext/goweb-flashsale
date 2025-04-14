package main

import (
	"fmt"
	"imooc_product/common"
	"imooc_product/rabbitmq"
	"imooc_product/repositories"
	"imooc_product/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err!=nil {
		fmt.Println(err)
	}
	// 创建product数据库操作实例
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	rabbitmqConsumerSimple := rabbitmq.NewRabbitMQSimple("imoocProduct")
	rabbitmqConsumerSimple.ConsumeSimple(orderService, productService)
}
