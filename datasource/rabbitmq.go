package datasource

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"test/config"
)

var rabbitmq *amqp.Channel

func InitRabbitMQ() {
	conn, err := amqp.Dial("amqp://" + config.Sysconfig.RabbitMQ.UserName + ":" + config.Sysconfig.RabbitMQ.Password + "@" + config.Sysconfig.RabbitMQ.Addr + "/")
	//defer conn.Close()
	if err != nil {
		log.Panic(err, "Failed to connect to RabbitMQ")
	}
	ch, err := conn.Channel()
	//defer ch.Close()
	if err != nil {
		log.Panic(err, "Failed to open a channel")
	}
	//创建一个交换机
	/*err = ch.ExchangeDeclare("orderExchange", "fanout", true, false, false, true, nil)
	if err != nil {
		log.Panic(err, "交换机创建失败")
	}
	//创建一个队列Queue
	q, err := ch.QueueDeclare(
		"orderQueue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Panic(err, "队列创建失败")
	}
	//当路由键中无*#时只接收该路由键的消息 相当于单播模式(direct)
	err=ch.QueueBind(q.Name,"orderQueue","orderExchange",false,nil)
	if err!=nil{
		log.Panic(err,"绑定失败")
	}*/
	rabbitmq = ch
	log.Println("rabbitmq初始化成功")
	go recevice()
}
func GetMQ() *amqp.Channel {
	return rabbitmq
}
func recevice() {
	//消费消息
	msgs, err := rabbitmq.Consume(
		"orderQueue", // queue
		"",           // consumer
		false,        // auto-ack需要确认接收
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Panic(err, "Failed to register a consumer")
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Println(string(d.Body))
			rabbitmq.Ack(d.DeliveryTag, false) //确认接收到了rabbitmq 否则退回
		}
	}()
	<-forever
}
