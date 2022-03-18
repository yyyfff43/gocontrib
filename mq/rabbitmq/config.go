/*
* @File : config
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/17 18:11
* @Software: GoLand
 */

package rabbitmq

import (
	"context"
	"errors"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	hlog "git.zhwenxue.com/zhgo/gocontrib/log"
)

// RMQClient client basic info
type RMQClient struct {
	server Broker //exchange info and queue prefix
	device string //consume tag
	conn   *amqp.Connection
	ch     *amqp.Channel //channel
	que    amqp.Queue
	msgs   <-chan amqp.Delivery //delivery channel

	destructor sync.Once
	confirm    chan amqp.Confirmation
	pubChan    chan *publishMsg
	ctx        context.Context
	cancel     context.CancelFunc
	onPublish  int32
	log        hlog.Logger
	Done       chan error
}

// RMQOption RabbitMQOption RabbitMQ 配置
//represents a parsed AMQP URI string.
type RMQOption struct {
	Scheme   string `yaml:"scheme"` //protocol amqp or amqps
	Host     string `yaml:"host"`   //mq server host
	Port     int    `yaml:"port"`   //server port
	Username string `yaml:"userName"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"Vhost"`    //虚拟主机，表示一批交换器、消息队列和相关对象
	Server   Broker `yaml:"server"`   //交换器配置信息
	Reliable bool   `yaml:"reliable"` //Reliable publisher confirms require confirm.select support from the connection
}

var (
	pubTime  = time.Second * 16
	tickTime = time.Second * 8

	messageTTL  = int64(time.Hour / time.Millisecond)          // TTL for message in queue
	queueExpire = int64(time.Hour * 24 * 7 / time.Millisecond) // expire time for unused queue

	//errAck     = errors.New("ack")
	errNack    = errors.New("nack")
	errFull    = errors.New("full")
	errCancel  = errors.New("cancel")
	errTimeout = errors.New("timeout")

	channelBuffSize = 4
)

// publishMsg the message info for publish
type publishMsg struct {
	exchange   string
	routingKey string
	msg        []byte
	expire     time.Duration
	startTime  time.Time
	ctx        context.Context
	cancel     context.CancelFunc
	ackErr     error
}

// Broker broker config
type Broker struct {
	QuePrefix string `yaml:"quePrefix"` //the prefix to mark queue

	//此处topic理解为标题, 代表exchange和bingding相关信息而非exchange的会话类型
	//一个exchange管理多个routingKey绑定到一个queue上, 要想创建多个queue 需要多次调用
	//NewPublisher和NewConsumer
	Topics []Topic `yaml:"topics"`
}

// Topic config equal to exchang info
//exchange 和binding 信息
type Topic struct {
	ChanName  string `yaml:"exchangeName"` //exchange
	ChanType  string `yaml:"exchangeType"` //exchange type
	KeyPrefix string `yaml:"keyPrefix"`    //RoutingKey or BindingKey
}
