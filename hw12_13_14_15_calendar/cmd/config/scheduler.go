package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Scheduler struct {
	Logger        LoggerConf `json:"logger"`
	RabbitMQ      Rabbit     `json:"rabbit_mq"`
	Database      DBConf     `json:"database"`
	IntervalInSec int64      `json:"interval_in_sec"`
}

func NewScheduler(filePath string) (Scheduler, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Scheduler{}, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var config Scheduler
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Scheduler{}, fmt.Errorf("can't decode config: %w", err)
	}
	return config, nil
}
