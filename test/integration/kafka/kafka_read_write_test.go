package kafka

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"

	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/kafka_reader"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/kafka_writer"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/logger"
)

// Тестирование отправки и получения сообщений в Кафку с использованием testcontainers
func TestKafkaWriteRead(t *testing.T) {
	topic := "test-by-testcontainers"
	group := "test-group"

	/********************** Настройка testcontainers ************************/
	ctx := context.Background()
	kafkaContainer, err := kafka.Run(ctx, "confluentinc/confluent-local:7.5.0", kafka.WithClusterID("kraftCluster"))
	testcontainers.CleanupContainer(t, kafkaContainer)
	require.NoError(t, err)

	assertAdvertisedListeners(t, kafkaContainer)

	if !strings.EqualFold(kafkaContainer.ClusterID, "kraftCluster") {
		t.Fatalf("expected clusterID to be %s, got %s", "kraftCluster", kafkaContainer.ClusterID)
	}
	/******************* КОНЕЦ Настройка testcontainers **********************/

	// Создаем издателя и подписчика, тестируем прием/отправку сообщения
	filePath, err := logger.GetLoggerMainLogPath()
	require.NoError(t, err)
	logger.InitLogger(logger.LogLevelDebug, filePath)

	brokers, err := kafkaContainer.Brokers(ctx)
	require.NoError(t, err)

	cfg := kafka_writer.Config{
		Brokers: brokers,
		Topic:   topic,
	}

	writer, err := kafka_writer.NewKafkaWriter(cfg, logger.Log())
	assert.NoError(t, err)
	defer writer.Close()

	// Запуск продюсера
	writer.Start()

	// Отправка сообщений
	var msgSend []string
	msgSend = append(msgSend, fmt.Sprintf("%v", time.Now().Nanosecond()), fmt.Sprintf("%v", time.Now().Nanosecond()))
	for _, val := range msgSend {
		writer.SendMessage("test_write", val)
	}

	//////////////////////////////////////// читаем из топика
	time.Sleep(2 * time.Second)
	cfgReader := kafka_reader.Config{
		Brokers:            brokers,
		Group:              group,
		Topic:              topic, // constants.KafkaSharesTopic
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
	i := 0
	select {
	case msg := <-msgChan:
		assert.Equal(t, msgSend[i], string(msg.Value))
		i++
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