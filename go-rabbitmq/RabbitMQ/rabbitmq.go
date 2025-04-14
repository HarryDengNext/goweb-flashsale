package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//url格式  amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://zhangyafei:zhangyafei@182.254.179.186:5672/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	// 连接信息
	Mqurl string
}

//创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	// 创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败！")
	return rabbitmq
}

//断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 简单模式step1： 1.创建简单模式下的rabbitmq实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// 简单模式step2： 2.简单模式下生产
func (r *RabbitMQ) PublishSimple(message string) {
	// 1. 申请队列，如果队列不存在则自动创建，如果存在则跳过创建
	// 保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2. 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true,根据exchange类型和routekey规则，如果无法找到符合条件的队列，则会把发送的消息返回给发送者
		false,
		// 如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息发还给发送者
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(message)},
	)
}

// 简单模式step3： 3.简单模式下消费
func (r *RabbitMQ) ConsumeSimple() {
	// 1. 申请队列，如果队列不存在则自动创建，如果存在则跳过创建
	// 保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2. 接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果为true，表示不能将同一个conn中的消息发送给这个conn中的消费者
		false,
		// 队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 3. 启用协程处理消息
	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] waiting for messages, to exit process CTRL+C")
	<-forever
}

// 订阅模式下创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 订阅模式下生产
func (r *RabbitMQ) PublishPub(message string) {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // 广播类型
		true,     // 持久化
		false,    // 是否删除
		false,    // true表示这个exchange不可以被client用来推送消息的，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a exchange")
	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 订阅模式消费端的代码
func (r *RabbitMQ) ReceiveSub() {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // 广播类型
		true,     // 持久化
		false,    // 是否删除
		false,    // true表示这个exchange不可以被client用来推送消息的，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a exchange")
	// 2. 试探性创建队列
	q, err := r.channel.QueueDeclare(
		"", // 随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a queue")
	// 绑定队列到 exchange中
	err = r.channel.QueueBind(
		q.Name,
		"", // 在订阅模式下，这里的key为空
		r.Exchange,
		false,
		nil)
	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 3. 启用协程处理消息
	go func() {
		for d := range messages {
			// 实现我们要处理的逻辑函数
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] waiting for messages, to exit process CTRL+C")
	<-forever
}

// 路由模式下创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingkey string) *RabbitMQ {
	// 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingkey)
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // 定向类型
		true,     // 持久化
		false,    // 是否删除
		false,    // true表示这个exchange不可以被client用来推送消息的，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a exchange")
	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 路由模式消费端的代码
func (r *RabbitMQ) ReceiveRouting() {
	// 1. 试探性的创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // 广播类型
		true,     // 持久化
		false,    // 是否删除
		false,    // true表示这个exchange不可以被client用来推送消息的，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a exchange")
	// 2. 试探性创建队列
	q, err := r.channel.QueueDeclare(
		"", // 随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a queue")
	// 绑定队列到 exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key, // 在订阅模式下，这里的key为空
		r.Exchange,
		false,
		nil)
	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 3. 启用协程处理消息
	go func() {
		for d := range messages {
			// 实现我们要处理的逻辑函数
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] waiting for messages, to exit process CTRL+C")
	<-forever
}
