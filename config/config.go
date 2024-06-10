package config

import (
	"api-save-fact/model"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

// ReadConfig читает конфигурацию из файла и возвращает объект конфигурации
func ReadConfig(configPath string) (model.Config, error) {
	var cfg model.Config

	// Открываем файл конфигурации
	f, err := os.Open(configPath)
	if err != nil {
		return cfg, fmt.Errorf("error opening config file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	// Декодируем конфигурацию из файла
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("error decoding config file: %v", err)
	}

	// Валидируем конфигурацию
	if err = validator.New().Struct(cfg); err != nil {
		return cfg, fmt.Errorf("error validating config: %v", err)
	}

	return cfg, nil
}
