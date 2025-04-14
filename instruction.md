# Go Web商品秒杀系统

## 一、项目介绍

### 1. 课程目标

- 应用GoWeb快速构建秒杀系统
- 全流程应用开发及架构化设计思维梳理
- 逐级优化，轻松应对“秒杀”及类似高并发场景

### 2. 知识储备

> - RabbitMQ入门
> - Iris入门

### 3. 基础功能开发

> - 后端商品管理功能表开发
> - 后端订单管理功能开发
> - 前台用户登录
> - 商品展示功能开发

### 4. 性能优化

> - 架构调优
> - 前端优化
> - 服务端优化
> - 安全优化

## 二、需求整理&系统设计

### 2.1 需求分析

- 主要功能点

> - 前台用户登录，商品展示，商品抢购
>
>
> - 后台订单管理

### 2.2 需求原型设计

- 主要设计页面

> - 前台用户登录页面，商品展示页面，商品抢购页面
>
>
> - 后台订单管理页面

### 2.3 系统架构设计

- 系统需求分析

> - 前端页面需要承载大流量
> - 在大并发状态下要解决超卖问题
> - 后端接口需要满足横向扩展

## 三、环境搭建

### 3.1 RabbitMQ介绍

#### 3.1.1 定义和特征

> 1. RbbitMQ是面向消息的中间件，用于组件之间的解耦，主要体现在消息的发送者和消费者之间无强依赖关系
> 2. RabbitMQ特点：高可用，可扩展，多语言客户端，管理界面等；
> 3. 主要使用场景：流量削峰，异步处理，应用解耦等；

#### 3.1.2 安装

> - ubuntu 的参照：    <https://gitee.com/zhangyafeii/rabbitmq>
> - windows的参照：<https://www.cnblogs.com/JustinLau/p/11738511.html>
> - 以下为centos7的安装过程

- 安装erlang

> ```
> # centos7
> wget https://github.com/rabbitmq/rabbitmq-server/releases/download/v3.7.17/rabbitmq-server-3.7.17-1.el7.noarch.rpm
> yum install epel-release
> yum install unixODBC unixODBC-devel wxBase wxGTK SDL wxGTK-gl
> rpm -ivh esl-erlang_22.0.7-1~centos~7_amd64.rpm
> ```

- 安装RabbitMQ

> ```
> wget https://github.com/rabbitmq/rabbitmq-server/releases/download/v3.8.0/rabbitmq-server-3.8.0-1.el7.noarch.rpm
> yum -y install socat
> rpm -ivh rabbitmq-server-3.8.0-1.el7.noarch.rpm
> ```

#### 3.1.3 启动

##### 3.1.3.1 centos

- 启动命令

> ```
> chkconfig rabbitmq-server on                     # 开机启动
> systemctl start rabbitmq-server.service          # 启动
> systemctl stop rabbitmq-server.service    		# 停止
> systemctl restart rabbitmq-server.service		# 重启
> rabbitmqctl status							  # 查看状态
> rabbitmq-plugins enable rabbitmq_management      # 启动Web管理器
> ```

![第一次访问rabbitmq页面](images\第一次访问rabbitmq页面.png)

- 修改配置

> vi /usr/lib/rabbitmq/lib/rabbitmq_server-3.8.0/ebin/rabbit.app
>
> 将：{loopback_users, [<<”guest”>>]}， 改为：{loopback_users, []}， 原因：rabbitmq从3.3.0开始禁止使用guest/guest权限通过   除localhost外的访问
>
> systemctl restart rabbitmq-server.service    # 重启服务

![rabbitmq安装成功](images/rabbitmq安装成功.png)

##### 3.1.3.2 windows

- rabbitmq启动方式有2种
- 1、以应用方式启动

```
rabbitmq-server -detached 或者 rabbitmq-server 直接启动，如果你关闭窗口或者需要在改窗口使用其他命令时应用就会停止
后台启动：rabbitmqctl start_app
关闭：rabbitmqctl stop
```

- 2、以服务方式启动（安装完之后在任务管理器中服务一栏能看到RabbtiMq）

```
rabbitmq-service install 安装服务
rabbitmq-service start 开始服务
Rabbitmq-service stop  停止服务
Rabbitmq-service enable 使服务有效
Rabbitmq-service disable 使服务无效
rabbitmq-service help 帮助
当rabbitmq-service install之后默认服务是enable的，如果这时设置服务为disable的话，rabbitmq-service start就会报错。
当rabbitmq-service start正常启动服务之后，使用disable是没有效果的
关闭:rabbitmqctl stop
```

- 3、Rabbitmq 管理插件启动，可视化界面

```
rabbitmq-plugins enable rabbitmq_management  # 启动
rabbitmq-plugins disable rabbitmq_management # 关闭
rabbitmq-plugins list                        # 展示插件情况
```

4、Rabbitmq节点管理方式

```
Rabbitmqctl
```

#### 3.1.4 核心概念

- VirtualHost
- Connections
- Exchange
- Channel
- Queue
- Binding 

#### 3.1.5 运行模式

##### 3.1.5.1 simple模式

> ![simple模式](images\简单模式.png)

- RabbitMQ/rabbitmq.go

```go
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

// 简单模式step1：1.创建简单模式下的rabbitmq实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

//简单模式step2： 2.简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string)  {
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
		if err != nil{
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
			amqp.Publishing{ContentType:"text/plain", Body:[]byte(message)},
			)
	}

func (r *RabbitMQ)ConsumeSimple() {
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
	if err != nil{
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
	<- forever
}
```

- mainSimplePublish

```go
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
```

- mainSimpleReceive.go

```go
package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
)

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}
```

- 先运行消息接收，再运行发送

![简单模式-生产者](images/简单模式-生产者.png)

![简单模式-生产者](images/简单模式-消费者.png)

##### 3.1.5.2 Work， 工作模式。

-  起到负载均衡的作用

> 一个消息只能被一个消费者获取
>
> ![工作模式](images/工作模式.png)

- mainWorkPublish

```go
package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	for i := 0; i<= 100; i++ {
		rabbitmq.PublishSimple("hello imooc!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
```

- mainWorkReceive1和mainWorkReceive2

```go
package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
)

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}
```

![工作模式-生产者](images/工作模式-生产者.png)

![工作模式-生产者](images/工作模式-消费者1.png)

![工作模式-生产者](images/工作模式-消费者2.png)

##### 3.1.5.3 Publish/Subscribe 订阅模式

> 消息被路由投递给多个队列，一个消息被多个消费者获取,生产者不能指定队列

![订阅模式](images/订阅模式.png)

- 创建RabbitMQ实例

```go
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
```



- 生产者

```go
package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("NewProduct")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		fmt.Println("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		time.Sleep(1 * time.Second)
	}
}
```

- 消费者

```go
package main

import (
	"RabbitMQ/RabbitMQ/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("NewProduct")
	rabbitmq.ReceiveSub()
}
```

![](images/订阅模式-生产者.png)

![](images/订阅模式-消费者1.png)

![](images/订阅模式-消费者2.png)

![](images/订阅模式-消费者3（延迟运行）.png)

##### 3.1.5.4 Routing, 路由模式

> 一个消息被多个消费者获取，并且消息的目标队列可悲生产者指定

![](images/路由模式.png)

- 创建实例

```go
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
```

- 生产者

```go
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
```

- 消费者

```go
package main

import "RabbitMQ/RabbitMQ/RabbitMQ"

func main()  {
	rabbitmq_imooc_one := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	rabbitmq_imooc_one.ReceiveRouting()
}
```

```go
package main

import "RabbitMQ/RabbitMQ/RabbitMQ"

func main()  {
	rabbitmq_imooc_two := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	rabbitmq_imooc_two.ReceiveRouting()
}
```

![](images/路由模式生产者-消费者.png)

### 3.2 Iris

## 四、商品后台管理开发

### 主要内容

> - 商品模型设计开发
> - 商品增删改查功能开发
> - 后台商品管理页面-

开发顺序

> 1. model
> 2. repositories
> 3. services
> 4. controllers
> 5. views

## 五、秒杀订单管理开发

### 主要内容

> - 订单模型设计开发
> - 订单管理功能开发
> - 订单功能页面制作你

作业：

>1. 继续完成订单管理的其他功能展示
>   - 查询单个订单信息
>   - 删除订单
>   - 更新修改订单
>2. 简化controller注册代码
>   - 简单创建controller实现注册
>3. 如何简化model层代码（中小项目）
>   - gorm

## 六、前台核心功能开发

### 用户登录页面开发

### 前端商品展示功能开发

### 秒杀数据控制开发

## 七、秒杀核心功能开发

### 前端商品展示功能开发

- 使用原有商品service查询信息
- 商品控制器开发
- 商品前端展示页面制作

### 秒杀数据控制开发

- 查询商品数量
- 扣除商品
- 生成订单

## 八、秒杀系统优化

### 整体架构分析

![](images/整体架构分析.png)

### 架构优化

![](images/架构优化.png)

### 前端优化

- CDN原理讲解和作用
- 阿里云添加CDN
- 部署前端静态文件

### 后端优化

- 后端优化思路

  系统特征：高并发，大流量

  优化方向：提高网站性能，保护数据库

  具体措施：静态化。分布式，消息队列

  ![](images/基础架构.png)

  ![](images/优化架构.png)

- 突破Session限制

  - Cookie替代Session集群原理

  - 登录代码重构

    Cookie：采用的是自客户端保存状态和数据的方案

    Session：采用的是在服务端保存状态和数据的方案

  - Cookie和Session的区别？

    - CookIe数据保存在客户的浏览器上，session数据放在服务器上
    - session保存在服务器上，当访问增多，会比较占用大量资源
    - 单个Cookie保存的数据不能超过4k

- 分布式接口实现

  - 引入分布式及代码结构调整

  - 一致性Hash原理

    用途：快速定位资源，均匀分布

    场景：分布式存储，分布式缓存，负载均衡

    ![](images/哈希算法分布式.png)


  - 分布式存储实现

  - 什么是分布式？

    分布式系统是相对于集中式系统来说的概念

    集中式系统：所有的应用程序和组件放在童泰机器上运行

    分布式定义：分布式系统是若干个独立计算机得我集合，这些计算机对于用户来说就像是单个相关系统

  - 分布式更直观的感受

    第一个感受：分布式是由多态机器组成

    第二个感受：分布式系统是一个整体，从外部感受不到多态机器的存在

    ![](images/集中式和分布式.png)

- 解决接口超卖问题

  互斥锁

  Wrk压测工具

  - 工具安装

    centos7执行：git clone https://github.com/wg/wrk

    cd wrk

    make

  - Wrk压测命令

    ![](images/Wrk压测命令.png)

    ![](images/Wrk执行结果.png)

  - 只支持类UNIX系统

  - 能用少量的线程测大量的连接

  Redis使用方式

  - 单机或主从方式

    ![](images/redis单机.png)

  - Redis Cluster 集群方式

    ![](images/redis集群.png)

  - Redis分布式方式使用

- RabbitMQ实现秒杀队列

  - RabbitMQ实现生产端代码
  - RabbitMQ实现消费端代码
  - RabbitMQ整体效果演示

### 安全优化

- 限流的作用和意义

  ![](images/为什么限流.png)

- 前端页面限流

  ![](images/前端页面限流.png)

- 服务端防止for循环请求

  > 设置时间间隔策略，判断同一个用户在指定间隔内只能抢购一次


- 其他限流措施
  - 添加黑名单（已实现）；
  - 添加图形验证码；
  - 添加自定义限流，比如漏桶法限流，令牌限流（可在getOne接口中实现）；
  - 添加隐藏秒杀接口地址（秒杀开始前不暴露秒杀接口地址）；
  - 服务端限制开始时间（秒杀开始服务端添加规则，限制访问）；
- 其他安全建议
  - 做好流量估算，预估每次秒杀访问量区间
  - 秒杀商品太多可以分批次秒杀
  - 在资金允许的情况下冗余服务器

接口

后端接口

> - 商品增加： http://127.0.0.1:8080/product/add
> - 商品删除： http://127.0.0.1:8080/product/delete 
> - 商品修改 ： http://127.0.0.1:8080/product/manager
> - 商品更新： http://127.0.0.1:8080/product/update
> - 商品查看：http://127.0.0.1:8080/product/all 
> - 订单查看： http://127.0.0.1:8080/order/all
> - 订单修改：http://127.0.0.1:8080/order/manager
> - 订单更新：http://127.0.0.1:8080/order/update（post)
> - 订单删除： http://127.0.0.1:8080/order/delete (post)

![](images/订单管理.png)

![](images/商品管理.png)

![](images/添加商品.png)

前端接口

> 用户登录 <http://127.0.0.1:8082/user/login>
>
> 用户注册 <http://127.0.0.1:8082/user/register>

![](images/用户登录.png)

![](images/用户注册.png)

分布式验证接口

```
# 网页抢购   点击立即抢购
http://127.0.0.1:8083/html/htmlProduct.html
# api抢购  将订单加入rabbitmq消息队列
http://127.0.0.1:8083/check?productID=2
// 将rabbitmq订单队列导入mysql
consumer.go
```

![](images/立即抢购.png)

![](images/数量控制接口.png)

![](images/rabbitmq将订单消息队列写入mysql.png)

## 九、课程总结

### 网站开发流程总结

![](images/网站开发流程总结.png)

### 秒杀难点总结

- RabbitMQ掌握和使用

- Iris框架的掌握和使用

- 高并发分布式验证

- 去Session权限验证

- 超卖问题解决（自定义接口代替redis）

- 横向扩展设计

- wrk压测工具

  ![](images/wrk压测工具.png)

### 秒杀优化思路总结

![](images/优化架构.png)

- 在秒杀系统每环节尽可能拦截流量
- 静态资源，页面能走CDN都走CDN
- 启用消息队列达到流量消峰作用

### 项目经验分享

- 在大并发系统中单台机器有限
- 从软件结构上设计的满足横向扩展
- Go写高并发确实简单高效并且很合适