/*
* @File : send_test
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/18 17:28
* @Software: GoLand
 */

package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRMQConfigWithPath(t *testing.T) {
	dir, _ := os.Getwd()
	path := dir + addr
	config, err := RMQConfigWithPath(path)
	fmt.Println(config, err)
	fmt.Println(config.Server.Topics[0].ChanName)
}

func TestNewPublisherClient(t *testing.T) {
	opt := getRMQOption()
	opt.Reliable = true

	ctx := context.Background()

	client := NewPublisherClient(ctx, opt, *testLogger)
	for i := 0; i < 3; i++ {
		msg := fmt.Sprintf("测试测试-%d", i)
		err := client.Publish("test-topic-0113", "*.mirror-0113.*", 2*time.Second, []byte(msg))
		//if err != nil {
		//	fmt.Println("push msg failed", err)
		//}
		require.Nil(t, err)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(3 * time.Second)

}
