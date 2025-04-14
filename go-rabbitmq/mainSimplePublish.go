package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
	"fmt"
)

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.PublishSimple("hello imooc!")
	fmt.Println("发送成功")
}
