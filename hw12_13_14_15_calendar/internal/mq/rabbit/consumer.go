package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/streadway/amqp"
)

type Consumer struct {
	cfg        config.Rabbit
	conn       *amqp.Connection
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery
}

func NewConsumer(cfg config.Rabbit) (*Consumer, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, NewError("can't connect to rmq", err)
	}

	return &Consumer{cfg: cfg, conn: conn}, nil
}

func (c *Consumer) CloseConn() error {
	return c.conn.Close()
}

func (c *Consumer) OpenChannel() error {
	var err error
	c.channel, err = declareChannel(c.cfg, c.conn)
	if err != nil {
		return NewError("can't create channel", err)
	}
	return nil
}

func (c *Consumer) CloseChannel() error {
	return c.channel.Close()
}

func (c *Consumer) BeginConsume() error {
	if c.channel == nil {
		return ErrChannelIsNil
	}
	queue, err := c.channel.QueueDeclare(
		c.cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return NewError("can't declare queue", err)
	}

	err = c.channel.QueueBind(
		queue.Name,
		c.cfg.RoutingKey,
		c.cfg.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return NewError("can't bind queue", err)
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		c.cfg.ConsumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return NewError("can't consume queue", err)
	}
	c.deliveries = deliveries
	return nil
}

func (c *Consumer) Get() <-chan app.MQMessage {
	mChan := make(chan app.MQMessage)

	go func() {
		for d := range c.deliveries {
			var notif app.MQEventNotification
			err := json.Unmarshal(d.Body, &notif)
			mess := app.MQMessage{
				Notif: notif,
				Err:   err,
			}
			if err == nil {
				_ = d.Ack(false)
			}

			mChan <- mess
		}
	}()

	return mChan
}
