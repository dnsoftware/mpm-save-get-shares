package kafka_reader

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type Config struct {
	Brokers []string `envconfig:"KAFKA_READER_BROKERS" required:"true"`
	Group   string   `envconfig:"KAFKA_READER_GROUP" required:"true"`
	Topic   string   `envconfig:"KAFKA_READER_TOPIC" required:"true"`
}

type Reader struct {
	cfg           Config
	consumerGroup sarama.ConsumerGroup
	consumer      *Consumer
}

func New(c Config) (*Reader, error) {

	// Настройка конфигурации
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest   // Чтение с самого конца (если нет сохраненного оффсета)
	config.Consumer.Offsets.AutoCommit.Enable = true        // Включаем автоматическое сохранение оффсетов
	config.Consumer.Offsets.AutoCommit.Interval = 10 * 1000 // Интервал для сохранения оффсетов - 10 секунд

	// Создаем клиента для consumer группы
	consumerGroup, err := sarama.NewConsumerGroup(c.Brokers, c.Group, config)
	if err != nil {
		fmt.Println("Ошибка при создании consumer group: ", err)
	}

	return &Reader{
		cfg:           c,
		consumerGroup: consumerGroup,
		consumer:      &Consumer{},
	}, nil
}

func (r *Reader) ConsumerGroupClose() {
	r.consumerGroup.Close()
}

func (r *Reader) Consume(ctx context.Context) error {
	err := r.consumerGroup.Consume(ctx, []string{r.cfg.Topic}, r.consumer)
	return err
}

//func (r *Reader) GetTopic() string {
//	return r.cfg.Topic
//}

/*************************************** Consumer **************************************/

// Consumer реализует интерфейс sarama.ConsumerGroupHandler
type Consumer struct {
}

// Setup вызывается перед началом обработки
func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup вызывается после завершения обработки
func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиций
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Сообщение получено: topic = %s, partition = %d, offset = %d, value = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Value))
		// session.MarkMessage(message, "") // Сообщаем Kafka, что сообщение обработано (не нужно если автокоммит)
	}
	return nil
}
