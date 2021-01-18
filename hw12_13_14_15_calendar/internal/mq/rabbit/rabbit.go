package rabbit

import (
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/streadway/amqp"
)

func declareChannel(cfg config.Rabbit, conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, NewError("can't get channel", err)
	}
	err = channel.ExchangeDeclare(
		cfg.ExchangeName,
		cfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, NewError("can't declare exchange", err)
	}

	return channel, nil
}
