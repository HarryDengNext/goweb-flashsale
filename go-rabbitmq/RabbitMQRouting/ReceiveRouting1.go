package main

import "RabbitMQ/RabbitMQ/RabbitMQ"

func main()  {
	rabbitmq_imooc_one := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	rabbitmq_imooc_one.ReceiveRouting()
}
