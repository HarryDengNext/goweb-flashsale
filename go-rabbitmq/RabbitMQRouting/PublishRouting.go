package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main()  {
	rabbit_imooc_one := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	rabbit_imooc_two := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	for i := 0; i <= 100; i++ {
		rabbit_imooc_one.PublishRouting("Hello imooc one!" + strconv.Itoa(i))
		rabbit_imooc_two.PublishRouting("Hello imooc two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
