/*
* @File : receive
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/17 18:55
* @Software: GoLand
 */

package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"

	hlog "git.zhwenxue.com/zhgo/gocontrib/log"
)

// Consume
// @Description: consume message to call MsgProcess function
// @receiver clt
// @param ctx
// @param MsgProcess
// @return error
func (clt *RMQClient) Consume(ctx context.Context, MsgProcess func(deliveries <-chan amqp.Delivery, done chan error)) error {
	msgs, err := clt.ch.Consume(
		clt.que.Name, // queue
		clt.device,   // consumer
		false,        // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)

	if err != nil {
		clt.Close()
		clt.log.Error(clt.ctx, "rabbit consume", hlog.String("Start consume ERROR:", err.Error()))
		return nil
	}

	clt.msgs = msgs

	go func() {
		cc := make(chan *amqp.Error)
		e := <-clt.ch.NotifyClose(cc)
		clt.log.Error(clt.ctx, "rabbit consume", hlog.String("channel close error:", e.Error()))
		clt.cancel()
	}()

	//调用业务的callback函数
	go MsgProcess(msgs, clt.Done)

	return nil
}
