package kafka_reader

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/internal/controller/kafka_consumer/shares"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/logger"
)

// Тестирование получения шар из кафки
// кафка должна быть запущена и в ней в соответствующем топике должно быть несколько шар
func TestReader(t *testing.T) {

	filePath, err := logger.GetLoggerMainLogPath()
	require.NoError(t, err)
	logger.InitLogger(logger.LogLevelDebug, filePath)
	//logger.InitLogger(logger.LogLevelProduction, filePath)

	cfg := Config{
		Brokers:            []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Group:              constants.KafkaSharesGroup,
		Topic:              "test_topic", // constants.KafkaSharesTopic
		AutoCommitInterval: constants.KafkaSharesAutocommitInterval,
		AutoCommitEnable:   true,
	}

	reader, err := NewKafkaReader(cfg, logger.Log())
	assert.NoError(t, err)
	defer reader.Close()

	// сброс в начало
	err = reader.SetGroupOffset(sarama.OffsetOldest)
	assert.NoError(t, err)

	handler := &shares.ShareConsumer{}
	reader.ConsumeMessages(handler)

	// задержка, чтобы сработал автокоммит
	time.Sleep(10 * time.Second)
	reader.Close()
}
