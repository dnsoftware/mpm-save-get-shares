package grpc

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dnsoftware/mpm-save-get-shares/config"
	pb "github.com/dnsoftware/mpm-save-get-shares/internal/adapter/grpc"
	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/utils"
)

// Должен быть запущен сторонний микросервис к которому идут запросы от тестируемого клиента
func TestGRPCStorageTest(t *testing.T) {

	basePath, err := utils.GetProjectRoot(constants.ProjectRootAnchorFile)
	if err != nil {
		log.Fatalf("GetProjectRoot failed: %s", err.Error())
	}
	configFile := basePath + "/config.yaml"
	envFile := basePath + "/.env_example"

	cfg, err := config.New(configFile, envFile)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx,
		cfg.GRPC.CoinTarget, // Адрес:порт
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Отключаем TLS
	)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	storage, err := pb.NewCoinStorage(conn)
	require.NoError(t, err)

	ctx = context.Background()
	id, err := storage.GetCoinIDByName(ctx, "ALPH")
	require.NoError(t, err)
	require.Equal(t, int64(4), id)
}
