package consumers

import (
	"git.zhwenxue.com/zhgo/gocontrib/salog/structs"
)

type NopConsumer struct{}

func InitNopConsumer() *NopConsumer {
	return &NopConsumer{}
}

func (c *NopConsumer) Send(data structs.EventData) error {
	return nil
}

func (c *NopConsumer) Flush() error {
	return nil
}

func (c *NopConsumer) Close() error {
	return nil
}

func (c *NopConsumer) ItemSend(item structs.Item) error {
	return nil
}
