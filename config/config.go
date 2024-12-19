package config

import "github.com/dnsoftware/mpm-save-get-shares/pkg/kafka_reader"

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App         App
	KafkaReader kafka_reader.Config
}

func New() (Config, error) {
	var config Config

	return config, nil
}
