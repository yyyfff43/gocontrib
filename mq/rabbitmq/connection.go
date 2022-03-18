/*
* @File : connection
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/17 18:15
* @Software: GoLand
 */

package rabbitmq

import (
	"context"
	"fmt"
	"io/ioutil"

	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/yaml.v2"

	hlog "git.zhwenxue.com/zhgo/gocontrib/log"
)

// RMQConfigWithPath
// @Description: parse RMQ config
// @param path
// @return RMQOption
// @return error
func RMQConfigWithPath(path string) (RMQOption, error) {
	rmqConfig := RMQOption{}
	fmt.Println("path", path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return rmqConfig, err
	}

	err = yaml.Unmarshal(yamlFile, &rmqConfig)
	if err != nil {
		return rmqConfig, err
	}

	return rmqConfig, nil

}

// NewPublisherClient
// @Description: init a Publisher of rabbitmq client
// @param ctx
// @param option
// @return *RMQClient
func NewPublisherClient(ctx context.Context, option RMQOption, log hlog.Logger) *RMQClient {

	clt := &RMQClient{}
	clt.ctx, clt.cancel = context.WithCancel(ctx)
	clt.server = option.Server
	clt.log = log

	err := clt.connInit(ctx, option)
	if err != nil {
		log.Error(ctx, "rabbit client", hlog.String("msg", "connInit ERROR"), hlog.String("err", err.Error()))
		return nil
	}

	return clt
}

// connInit
// @Description: rabbit client connection init
// @receiver clt
// @param ctx
// @param option
// @return err
func (clt *RMQClient) connInit(ctx context.Context, option RMQOption) (err error) {

	//将外部struct转换成amqp.URI
	opt := amqp.URI{
		Scheme:   option.Scheme,
		Host:     option.Host,
		Port:     option.Port,
		Username: option.Username,
		Password: option.Password,
		Vhost:    option.Vhost,
	}

	url := opt.String()
	clt.log.Info(ctx, "rabbit connInit", hlog.String("conn url", url))
	//fmt.Println("rabbit connInit", "conn url", url)
	conn, dErr := amqp.Dial(url)
	clt.failOnError(dErr, "dial to rabbit")

	ch, cErr := conn.Channel()
	clt.failOnError(cErr, "create channel")

	if option.Reliable {
		//生产者确认（publisher confirm）的模式
		clt.log.Info(clt.ctx, "enabling publishing confirms")
		nErr := ch.Confirm(false)
		clt.failOnError(nErr, "puts channel into confirm mode error")

		clt.confirm = make(chan amqp.Confirmation, 16)
		clt.confirm = ch.NotifyPublish(clt.confirm)

	}

	for _, topic := range clt.server.Topics {
		err = ch.ExchangeDeclare(
			topic.ChanName, // exchange name
			topic.ChanType, // exchange type

			true,  // durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			//return err
			clt.failOnError(err, "declare exchange failed")
		}
	}

	clt.conn = conn
	clt.ch = ch
	clt.pubChan = make(chan *publishMsg, channelBuffSize)

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel is closed.
//func (clt *RMQClient) confirmOne(confirms <-chan amqp.Confirmation) {
//
//	clt.log.Info(clt.ctx, "waiting for confirmation of one publishing")
//	if confirmed := <-confirms; confirmed.Ack {
//		clt.log.Info(clt.ctx, "confirmOne", hlog.Uint64("confirmed delivery with delivery tag", confirmed.DeliveryTag))
//	} else {
//		clt.log.Info(clt.ctx, "confirmOne", hlog.Uint64("failed delivery of delivery tag: %d", confirmed.DeliveryTag))
//	}
//}

// NewConsumer
// @Description: new consumer client
// @param ctx
// @param msgProcess
// @param server
// @param device
// @param option
// @return *RMQClient
func NewConsumer(ctx context.Context, server Broker, device string, option RMQOption, log hlog.Logger) *RMQClient {
	clt := &RMQClient{}

	clt.ctx, clt.cancel = context.WithCancel(ctx)
	clt.server = server
	clt.device = device
	clt.log = log
	clt.Done = make(chan error)

	err := clt.connInit(ctx, option)
	clt.failOnError(err, "init rabbit connection failed")

	err = clt.queInit(clt.server, false)
	if err != nil {
		clt.Close()
		cErr := clt.connInit(ctx, option)
		if cErr != nil {
			clt.log.Error(clt.ctx, "rabbit client", hlog.String("connInit ERROR:", cErr.Error()))
			return nil
		}

		qErr := clt.queInit(clt.server, true)
		if qErr != nil {
			clt.log.Error(clt.ctx, "rabbit client", hlog.String("connInit queInit:", qErr.Error()))
			return nil
		}
	}

	if err != nil {
		clt.Close()
		clt.log.Error(clt.ctx, "rabbit client", hlog.String("queInit ERROR:", err.Error()))
		return nil
	}

	return clt
}

// Close the client
func (clt *RMQClient) Close() error {
	//通过sync.Once控制资源释放只执行一次
	clt.destructor.Do(func() {
		if clt.conn != nil {
			clt.conn.Close()
			clt.conn = nil
			clt.ch = nil
		}
		clt.cancel()
	})

	// wait for MsgProcess() to exit
	return <-clt.Done
}
