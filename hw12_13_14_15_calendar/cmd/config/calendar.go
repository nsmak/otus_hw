package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Calendar struct {
	Logger     LoggerConf `json:"logger"`
	RestServer RestConf   `json:"rest_server"`
	GrpcServer GrpcConf   `json:"grpc_server"`
	Database   DBConf     `json:"database"`
}

func NewCalendar(filePath string) (Calendar, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Calendar{}, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var config Calendar
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Calendar{}, fmt.Errorf("can't decode config: %w", err)
	}
	return config, nil
}

type LoggerConf struct {
	Level    int8   `json:"level"`
	FilePath string `json:"file_path"`
}

type RestConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type GrpcConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type DBConf struct {
	InMem    bool   `json:"in_mem"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	DBName   string `json:"db_name"`
}
