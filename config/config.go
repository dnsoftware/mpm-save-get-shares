package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type App struct {
	Name    string `yaml:"app_name" envconfig:"APP_NAME"    required:"true"`
	Version string `yaml:"app_version" envconfig:"APP_VERSION" required:"true"`
}

type KafkaShareReaderConfig struct {
	Brokers            []string `yaml:"brokers" envconfig:"KAFKA_SHARE_READER_BROKERS" required:"true"`
	Group              string   `yaml:"group" envconfig:"KAFKA_SHARE_READER_GROUP" required:"true"`
	Topic              string   `yaml:"topic" envconfig:"KAFKA_SHARE_READER_TOPIC" required:"true"`
	AutoCommitEnable   bool     `yaml:"auto_commit_enable" envconfig:"KAFKA_SHARE_AUTO_COMMIT_ENABLE" required:"true"`
	AutoCommitInterval int      `yaml:"auto_commit_interval" envconfig:"KAFKA_SHARE_AUTO_COMMIT_INTERVAL" required:"true"` // в секундах
}

type KafkaMetricWriterConfig struct {
	Brokers []string `yaml:"brokers" envconfig:"KAFKA_METRIC_WRITER_BROKERS" required:"true"`
	Topic   string   `yaml:"topic" envconfig:"KAFKA_METRIC_WRITER_TOPIC" required:"true"`
}

type Config struct {
	App               App                     `yaml:"application"`
	KafkaShareReader  KafkaShareReaderConfig  `yaml:"kafka_share_reader"`
	KafkaMetricWriter KafkaMetricWriterConfig `yaml:"kafka_metric_writer"`
}

func New(filePath string, envFile string) (Config, error) {
	var config Config
	var err error

	// 1. Читаем из config.yaml. Самый низкий приоритет
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if decodeErr := decoder.Decode(&config); decodeErr != nil {
			log.Fatalf("Ошибка при чтении config.yaml: %v", decodeErr)
		}
	} else {
		log.Printf("config.yaml не найден, используются значения по умолчанию: %v", err)
	}

	// 2.1 Загрузка переменных окружения из .env
	err = godotenv.Load(envFile)
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	// 2.2 Переопределяем переменные, полученные из конфиг файла
	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	// 3. Чтение параметров командной строки
	// Регистрируем флаги
	kafkaShareBrokers := flag.String("kafka_share_brokers", "", "Хост для подключения")

	// Устанавливаем тестовые аргументы
	flag.CommandLine.Parse([]string{"-kafka_share_brokers=localhost:9092"})

	flag.Parse()

	// для каждого аргумента проверяем не пустой ли он, и если не пустой - переопределяем переменную конфига
	if *kafkaShareBrokers != "" {
		config.KafkaShareReader.Brokers = strings.Split(*kafkaShareBrokers, ",")
	}

	return config, nil
}
