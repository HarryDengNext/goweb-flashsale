package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("NewProduct")
	rabbitmq.ReceiveSub()
}