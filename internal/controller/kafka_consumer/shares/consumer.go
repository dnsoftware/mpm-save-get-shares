// Package shares реализует обработчик шар, полученных из топика кафки
package shares

import (
	"fmt"

	"github.com/IBM/sarama"
)

// ShareConsumer реализует интерфейс sarama.ConsumerGroupHandler
type ShareConsumer struct {
}

// Setup вызывается перед началом обработки
func (consumer *ShareConsumer) Setup(session sarama.ConsumerGroupSession) error {

	return nil
}

// Cleanup вызывается после завершения обработки
func (consumer *ShareConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиций
func (consumer *ShareConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Сообщение получено: topic = %s, partition = %d, offset = %d, value = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Value))
		session.MarkMessage(message, "") // Сообщаем Kafka, что сообщение обработано (вызывать всегда, даже если автокоммит включен)
	}
	return nil
}
