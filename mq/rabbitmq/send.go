/*
* @File : send
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/17 18:33
* @Software: GoLand
 */

package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	hlog "git.zhwenxue.com/zhgo/gocontrib/log"
)

// failOnError
// @Description: check error
// @receiver clt
// @param err
// @param msg
func (clt *RMQClient) failOnError(err error, msg string) {
	if err != nil {
		clt.log.Panic(clt.ctx, "rabbit failOnError", hlog.String(msg, err.Error()))
	}
}

// Publish
// @Description: Publish used to send message to topic in message Queue.
// @receiver clt
// @param exchange
// @param routingKey
// @param expire
// @param msg
// @return err
func (clt *RMQClient) Publish(exchange, routingKey string, expire time.Duration, msg []byte) (err error) {
	pMsg := publishMsg{
		exchange:   exchange,
		routingKey: routingKey,
		expire:     expire,
		msg:        msg,
	}

	// 在client中，pubchan是预先建立好的，但是只有在有publish时，才创建publishProc
	// 如果发送过程出现异常导致publishProc退出，此时onPublish被置零，可以再次创建新的publishProc
	// pubChan中可能有残存的msg，如果没有及时新的publishProc启动，则对这些消息的处理是无用的
	// 如果当前client关闭，pubchan不会立刻关闭（等待gc），已经进入发送过程的publish会等待超时
	if atomic.AddInt32(&clt.onPublish, 1) == 1 {
		go clt.publishProc()
	} else {
		atomic.AddInt32(&clt.onPublish, -1)
	}

	timer := time.NewTimer(pubTime)
	defer timer.Stop()

	pMsg.ctx, pMsg.cancel = context.WithCancel(context.Background())
	defer pMsg.cancel()
	select {
	case <-timer.C:
		err = errFull
		return
	case clt.pubChan <- &pMsg:
	}

	timer.Reset(pubTime)
	select {
	case <-pMsg.ctx.Done():
		err = pMsg.ackErr
		break
	case <-timer.C:
		err = errTimeout
		break
	case <-clt.ctx.Done():
		err = errCancel
		break
	}

	return
}

// sendPublish
// @Description: send msg to topic
// @receiver clt
// @param topic
// @param keySuffix
// @param msg
// @param expire
// @return error
func (clt *RMQClient) sendPublish(exchange, routingKey string, msg []byte, expire time.Duration) error {
	if expire <= 0 {
		return errors.New("expiration parameter error")
	}

	//PublishWithDeferredConfirm
	if err := clt.ch.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
			Expiration:  fmt.Sprintf("%d", int64(expire/time.Millisecond)),
		}); err != nil {
		return fmt.Errorf("exchange publish: %s", err)
	}

	return nil
}

func (clt *RMQClient) publishProc() {

	ticker := time.NewTicker(tickTime)
	//发送消息的存储map
	deliveryMap := make(map[uint64]*publishMsg)

	defer func() {
		atomic.AddInt32(&clt.onPublish, -1)
		ticker.Stop()
		for _, msg := range deliveryMap {
			msg.ackErr = errCancel
			msg.cancel()
		}
	}()

	var deliveryTag uint64 = 1
	var ackTag uint64 = 1
	var pMsg *publishMsg

	for {
		select {
		case <-clt.ctx.Done():
			return

		case pMsg = <-clt.pubChan:
			pMsg.startTime = time.Now()
			//msg deliveryTag start  at 1 will auto increase
			err := clt.sendPublish(pMsg.exchange, pMsg.routingKey, pMsg.msg, pMsg.expire)
			if err != nil {
				pMsg.ackErr = err
				pMsg.cancel()
			}
			deliveryMap[deliveryTag] = pMsg
			deliveryTag++

		case c, ok := <-clt.confirm:
			if !ok {
				clt.log.Error(clt.ctx, "rabbit publishProc", hlog.String("err", "client Publish notify channel error"))
				return
			}
			pMsg = deliveryMap[c.DeliveryTag]
			fmt.Println("DeliveryTag:", c.DeliveryTag)
			delete(deliveryMap, c.DeliveryTag)
			if c.Ack {
				pMsg.ackErr = nil
				pMsg.cancel()
			} else {
				pMsg.ackErr = errNack
				pMsg.cancel()
			}

		case <-ticker.C:
			now := time.Now()
			for {
				if len(deliveryMap) == 0 {
					break
				}
				pMsg = deliveryMap[ackTag]
				if pMsg != nil {
					if now.Sub(pMsg.startTime.Add(pubTime)) > 0 {
						pMsg.ackErr = errTimeout
						pMsg.cancel()
						delete(deliveryMap, ackTag)
					} else {
						break
					}
				}
				ackTag++
			}
		}
	}
}
