/*
* @File : queue
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/17 18:11
* @Software: GoLand
 */

package rabbitmq

import (
	"fmt"
	hlog "git.zhwenxue.com/zhgo/gocontrib/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// queInit
// @Description: init channel queue
// @receiver clt
// @param server
// @param ifFresh
// @return err
func (clt *RMQClient) queInit(server Broker, ifFresh bool) (err error) {
	var num int
	ch := clt.ch
	if ifFresh {
		num, err = ch.QueueDelete(
			server.QuePrefix+"."+clt.device, //queue name
			false,
			false,
			false,
		)
		if err != nil {
			return
		}
		clt.log.Info(clt.ctx, "rabbit queInit", hlog.String(clt.device, "queue deleted with"),
			hlog.Int("num", num), hlog.String("msg", "message purged"))
	}

	args := make(amqp.Table)
	args["x-message-ttl"] = messageTTL
	args["x-expires"] = queueExpire
	q, err := ch.QueueDeclare(
		server.QuePrefix+"."+clt.device, // name of the queue
		true,                            // durable
		false,                           // delete when usused
		false,                           // exclusive
		false,                           // no-wait
		args,                            // arguments
	)
	if err != nil {
		return
	}

	//以下代码实现绑定关系说明
	//exchange:routingKey:queue == 1:n:1
	//exchange:routingKey:queue == n:n:1
	//
	for _, topic := range clt.server.Topics {
		err = ch.QueueBind(
			q.Name,                         //name of the queue
			topic.KeyPrefix+"."+clt.device, //RoutingKey/bindingKey
			topic.ChanName,                 //sourceExchange name
			false,                          // noWait
			nil,                            // arguments
		)

		s := fmt.Sprintf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topic", topic)
		clt.log.Info(clt.ctx, "rabbit queInit", hlog.String("QueueBind", s))
		if err != nil {
			return
		}
	}

	clt.que = q
	return
}
