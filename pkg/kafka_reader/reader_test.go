package kafka_reader

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {

	cfg := Config{
		Brokers: []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Group:   "sharesGroup",
		Topic:   "shares",
	}

	reader, err := New(cfg)
	assert.NoError(t, err)

	ctx := context.Background()
	err = reader.Consume(ctx)
	if err != nil {
		fmt.Println("Error from consumer group:", err)
	}

}
