/*
* @File : receive_test
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/18 17:33
* @Software: GoLand
 */

package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func MsgProcess(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		err := d.Ack(false)
		if err != nil {
			//ack失败,需要重发ack消息并记录日志
			fmt.Println("send ack failed, please check it")
		}
	}

	log.Println("handle: deliveries channel closed")
	done <- nil
}

func TestNewConsumer(t *testing.T) {
	opt := getRMQOption()
	opt.Reliable = false
	broker := Broker{
		QuePrefix: "test-queue-0113",
		Topics: []Topic{{
			ChanName:  "test-topic-0113",
			ChanType:  "fanout", //Exchange type - direct|fanout|topic|x-custom
			KeyPrefix: "topic-prefix-0113",
		},
		}}
	ctx := context.Background()

	client := NewConsumer(ctx, broker, "cTags-001", opt, *testLogger)

	err := client.Consume(ctx, MsgProcess)
	client.failOnError(err, "consue message failed for cli")
	require.Nil(t, err)

	//client1 := NewConsumer(ctx, broker, "consumerTag1", opt)
	//err1 := client1.Consume(ctx, msgProcess)
	//failOnError(err1, "consue message failed for cli1")

	//实际测试需要打开 长期监听中, 此时为了测试程序通过行暂时注销
	//select {}

}
