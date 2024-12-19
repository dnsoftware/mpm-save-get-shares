package kafka_reader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/internal/controller/kafka_consumer/shares"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/logger"
)

func TestReader(t *testing.T) {

	filePath, err := logger.GetLoggerMainLogPath()
	require.NoError(t, err)
	logger.InitLogger(logger.LogLevelDebug, filePath)
	//logger.InitLogger(logger.LogLevelProduction, filePath)

	cfg := Config{
		Brokers:            []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Group:              constants.KafkaSharesGroup,
		Topic:              constants.KafkaSharesTopic,
		AutoCommitInterval: constants.KafkaSharesAutocommitInterval,
		AutoCommitEnable:   true,
	}

	reader, err := NewKafkaReader(cfg, logger.Log())
	assert.NoError(t, err)

	//logger.Log().Error(fmt.Sprintf("error from consumer group: %s", "errtest"), zap.String("test", "testval"))

	handler := &shares.ShareConsumer{}
	reader.ConsumeMessages(handler)

	time.Sleep(10 * time.Second)
	reader.Close()
}
