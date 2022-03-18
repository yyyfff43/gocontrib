/*
* @File : connection_test
* @Describe :
* @Author: gongdenglong@zongheng.com
* @Date : 2022/1/17 18:21
* @Software: GoLand
 */
package rabbitmq

import (
	"os"

	hlog "git.zhwenxue.com/zhgo/gocontrib/log"
)

var testLogger = hlog.New(os.Stdout, hlog.DebugLevel, hlog.AddCallerSkip(1))
var addr = "/rabbitmq_config_sample.yml"

func getRMQOption() RMQOption {
	//opt := RMQOption{
	//	Scheme:   "amqp",
	//	Host:     "10.3.138.105",
	//	Port:     5672,
	//	Username: "xmmq",
	//	Password: "xmdev2021",
	//	Vhost:    "/",
	//	Server: Broker{
	//		QuePrefix: "test-queue-0113",
	//		Topics: []Topic{{
	//			ChanName:  "test-topic-0113",
	//			ChanType:  "fanout", //Exchange type - direct|fanout|topic|x-custom
	//			KeyPrefix: "topic-queue-0113",
	//		}},
	//	},
	//}

	dir, _ := os.Getwd()
	path := dir + addr
	opt, err := RMQConfigWithPath(path)
	if err != nil {
		panic(err)
	}

	return opt
}
