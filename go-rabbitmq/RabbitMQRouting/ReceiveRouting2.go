package main

import "RabbitMQ/RabbitMQ/RabbitMQ"

func main()  {
	rabbitmq_imooc_two := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	rabbitmq_imooc_two.ReceiveRouting()
}
