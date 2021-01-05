package app

import (
	"context"
)

type MQConsumer interface {
	OpenChannel() error
	BeginConsume() error
	Get() <-chan MQMessage
	CloseChannel() error
	CloseConn() error
}

type MQProducer interface {
	Publish(ctx context.Context, body []byte) error
	OpenChannel() error
	CloseChannel() error
	CloseConn() error
}
