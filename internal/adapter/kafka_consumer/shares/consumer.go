// Package shares реализует обработчик шар, полученных из топика кафки
package shares

import (
	"context"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/dnsoftware/mpm-save-get-shares/internal/usecase/share"
)

// ShareConsumer реализует интерфейс sarama.ConsumerGroupHandler
type ShareConsumer struct {
	MsgChan chan *sarama.ConsumerMessage
	share.ShareUseCase
}

func NewShareConsumer(msgChan chan *sarama.ConsumerMessage, shareUseCase share.ShareUseCase) (*ShareConsumer, error) {
	return &ShareConsumer{
		MsgChan:      msgChan,
		ShareUseCase: shareUseCase,
	}, nil
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
// если используем канал consumer.MsgChan (или какой-то подобный) - нужно его вычитывать снаружи
func (consumer *ShareConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Создаем carrier для извлечения заголовков
		carrier := propagation.MapCarrier{}
		// Извлекаем контекст трассировки из заголовков
		// Преобразуем заголовки Kafka в формат map
		for _, header := range msg.Headers {
			carrier[string(header.Key)] = string(header.Value)
		}
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

		tracer := otel.Tracer("consume-share")
		ctx, span := tracer.Start(ctx, "process")

		consumer.MsgChan <- msg
		// Тут идет вызов usecase

		// Если сообщение успешно обработано - помечаем, как обработанное TODO
		session.MarkMessage(msg, "")
		// иначе - логируем ошибку и делаем пометку "алерт" TODO
		// ...

		span.End()
	}
	return nil
}

// NormalizeShare Обработчик шары
// Получение кодов майнеров/воркеров из кеша или из Postgres
func (consumer *ShareConsumer) NormalizeShare(shareData []byte) {

}
