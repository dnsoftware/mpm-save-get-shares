package kafka_writer

import (
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/kafka_reader"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/logger"
)

// Тестирование отправки шар в кафку
// кафка должна быть запущена
func TestWriter(t *testing.T) {

	filePath, err := logger.GetLoggerMainLogPath()
	require.NoError(t, err)
	logger.InitLogger(logger.LogLevelDebug, filePath)
	//logger.InitLogger(logger.LogLevelProduction, filePath)

	cfg := Config{
		Brokers: []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:   "test_topic",
	}

	writer, err := NewKafkaWriter(cfg, logger.Log())
	assert.NoError(t, err)
	defer writer.Close()

	err = writer.DeleteTopic(writer.topic)
	assert.NoError(t, err)

	// Запуск продюсера
	writer.Start()

	// Отправка сообщений
	msgSend := fmt.Sprintf("%v", time.Now().Nanosecond())
	writer.SendMessage("test_write", msgSend)

	writer.Close()

	//////////////////////////////////////// читаем из топика
	time.Sleep(2 * time.Second)
	cfgReader := kafka_reader.Config{
		Brokers:            []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Group:              constants.KafkaSharesGroup,
		Topic:              "test_topic", // constants.KafkaSharesTopic
		AutoCommitInterval: constants.KafkaSharesAutocommitInterval,
		AutoCommitEnable:   true,
	}

	reader, err := kafka_reader.NewKafkaReader(cfgReader, logger.Log())
	assert.NoError(t, err)
	defer reader.Close()

	err = reader.SetGroupOffset(sarama.OffsetNewest)
	assert.NoError(t, err)

	// Читаем сообщение
	msgChan := make(chan *sarama.ConsumerMessage)
	handler := &testConsumerGroupHandler{msgChan: msgChan}

	go func() {
		reader.ConsumeMessages(handler)
	}()

	// Получаем сообщение
	select {
	case msg := <-msgChan:
		assert.Equal(t, msgSend, string(msg.Value))
	case <-time.After(6 * time.Second):
		t.Fatal("Таймаут при получении сообщения")
	}

}

// Тестовый обработчик ConsumerGroup
type testConsumerGroupHandler struct {
	msgChan chan *sarama.ConsumerMessage
}

func (h *testConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *testConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *testConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.msgChan <- msg
		sess.MarkMessage(msg, "")
		return nil
	}
	return nil
}
