package kafka_reader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/dnsoftware/mpm-save-get-shares/pkg/logger"
)

type Config struct {
	Brokers            []string `envconfig:"KAFKA_READER_BROKERS" required:"true"`
	Group              string   `envconfig:"KAFKA_READER_GROUP" required:"true"`
	Topic              string   `envconfig:"KAFKA_READER_TOPIC" required:"true"`
	AutoCommitInterval int      `envconfig:"KAFKA_AUTO_COMMIT_INTERVAL" required:"true"` // в секундах
	AutoCommitEnable   bool     `envconfig:"KAFKA_AUTO_COMMIT_ENABLE" required:"true"`
}

type KafkaReader struct {
	group         string
	topic         string
	consumerGroup sarama.ConsumerGroup
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	logger        logger.MPMLogger
}

func NewKafkaReader(c Config, logger logger.MPMLogger) (*KafkaReader, error) {
	// Настройка конфигурации
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                                           // Чтение с самого начала (если нет сохраненного оффсета)
	config.Consumer.Offsets.AutoCommit.Enable = c.AutoCommitEnable                                  // Включаем автоматическое сохранение оффсетов
	config.Consumer.Offsets.AutoCommit.Interval = time.Duration(c.AutoCommitInterval) * time.Second // Интервал для сохранения оффсетов - 10 секунд

	// Создаем клиента для consumer группы
	consumerGroup, err := sarama.NewConsumerGroup(c.Brokers, c.Group, config)
	if err != nil {
		return nil, err
	}

	// Создание контекста для управления остановкой
	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaReader{
		group:         c.Group,
		topic:         c.Topic,
		consumerGroup: consumerGroup,
		ctx:           ctx,
		cancel:        cancel,
		logger:        logger,
	}, nil
}

// ConsumeMessages - запуск чтения сообщений из топика
func (r *KafkaReader) ConsumeMessages(handler sarama.ConsumerGroupHandler) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		// Бесконечный цикл чтения сообщений
		for {
			if err := r.consumerGroup.Consume(r.ctx, []string{r.topic}, handler); err != nil {
				r.logger.Error(fmt.Sprintf("Ошибка при чтении сообщений: %v", err))
			}

			// Если контекст завершён, выходим из цикла
			if r.ctx.Err() != nil {
				return
			}
		}
	}()
}

func (r *KafkaReader) Close() {
	// Отмена контекста
	r.cancel()
	// Ожидание завершения горутин
	r.wg.Wait()
	// Закрытие Consumer Group
	if err := r.consumerGroup.Close(); err != nil {
		r.logger.Error(fmt.Sprintf("Ошибка закрытия Consumer Group: %v", err))
	}

	r.logger.Info(fmt.Sprintf("KafkaConsumer завершён, Group: %s, Topic: %s", r.group, r.topic))
}
