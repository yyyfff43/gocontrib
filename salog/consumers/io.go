package consumers

import (
	"encoding/json"
	"git.zhwenxue.com/zhgo/gocontrib/salog/structs"
	"io"
)

type IOConsumer struct {
	w io.Writer
}

func InitIOConsumer(w io.Writer) (*IOConsumer, error) {
	return &IOConsumer{w: w}, nil
}

func (c *IOConsumer) Send(data structs.EventData) error {
	itemData, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	_, _ = c.w.Write(itemData)
	return nil
}

func (c *IOConsumer) Flush() error {
	return nil
}

func (c *IOConsumer) Close() error {
	return nil
}

func (c *IOConsumer) ItemSend(item structs.Item) error {
	itemData, err := json.Marshal(item)
	if err != nil {
		return nil
	}
	_, _ = c.w.Write(itemData)
	return nil
}
