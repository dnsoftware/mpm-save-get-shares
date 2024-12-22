// Package shares реализует обработчик шар, полученных из топика кафки
package shares

import (
	"github.com/IBM/sarama"
)

// ShareConsumer реализует интерфейс sarama.ConsumerGroupHandler
type ShareConsumer struct {
	MsgChan chan *sarama.ConsumerMessage
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
	for msg := range claim.Messages() {
		consumer.MsgChan <- msg
		session.MarkMessage(msg, "")
		return nil
	}
	return nil
}
