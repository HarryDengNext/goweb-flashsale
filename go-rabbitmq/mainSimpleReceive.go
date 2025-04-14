package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
)

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}