package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Sender struct {
	Logger   LoggerConf `json:"logger"`
	RabbitMQ Rabbit     `json:"rabbit_mq"`
	Database DBConf     `json:"database"`
}

type Rabbit struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ExchangeName string `json:"exchange_name"`
	ExchangeType string `json:"exchange_type"`
	QueueName    string `json:"queue_name"`
	RoutingKey   string `json:"routing_key"`
	ConsumerTag  string `json:"consumer_tag"`
}

func NewSender(filePath string) (Sender, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Sender{}, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var config Sender
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Sender{}, fmt.Errorf("can't decode config: %w", err)
	}
	return config, nil
}
