package config

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/utils"
)

func TestConfigNew(t *testing.T) {
	basePath, err := utils.GetProjectRoot(constants.ProjectRootAnchorFile)
	if err != nil {
		log.Fatalf("GetProjectRoot failed: %s", err.Error())
	}
	configFile := basePath + "/config_example.yaml"
	envFile := basePath + "/.env"

	cfg, err := New(configFile, envFile)
	require.NoError(t, err)

	assert.Equal(t, constants.KafkaSharesTopic, cfg.KafkaShareReader.Topic)

}
