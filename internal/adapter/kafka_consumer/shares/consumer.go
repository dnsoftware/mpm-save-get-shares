// Package shares реализует обработчик шар, полученных из топика кафки
package shares

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/dnsoftware/mpm-save-get-shares/internal/dto"
	"github.com/dnsoftware/mpm-save-get-shares/internal/entity"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/kafka_reader"
)

type Config struct {
	BatchSize     int           // Размер буфера для пакетного чтения
	FlushInterval time.Duration // Максимальное время ожидания для заполнения пакета в секундах
}

type Processor interface {
	NormalizeShare(ctxTask context.Context, shareFound dto.ShareFound) (entity.Share, error)
	AddSharesBatch(shares []entity.Share) error
}

// ShareConsumer реализует интерфейс sarama.ConsumerGroupHandler
type ShareConsumer struct {
	kafkaReader *kafka_reader.KafkaReader
	msgChan     chan *sarama.ConsumerMessage
	Processor
}

func NewShareConsumer(cfg Config, kafkaReader *kafka_reader.KafkaReader, processor Processor) (*ShareConsumer, error) {
	return &ShareConsumer{
		kafkaReader: kafkaReader,
		msgChan:     make(chan *sarama.ConsumerMessage),
		Processor:   processor,
	}, nil
}

// StartConsume Стартует чтение из Кафки
func (consumer *ShareConsumer) StartConsume() {

	go func() {
		consumer.kafkaReader.ConsumeMessages(consumer)
	}()

}

// GetConsumeChan возвращает канал куда пишуться считанные сообщения
func (consumer *ShareConsumer) GetConsumeChan() chan *sarama.ConsumerMessage {
	return consumer.msgChan
}

func (consumer *ShareConsumer) Close() {
	consumer.kafkaReader.Close()
}

// Setup вызывается перед началом обработки (интерфейс ConsumerGroupHandler)
func (consumer *ShareConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup вызывается после завершения обработки (интерфейс ConsumerGroupHandler)
func (consumer *ShareConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиций (интерфейс ConsumerGroupHandler)
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

		consumer.msgChan <- msg
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
