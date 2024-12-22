package main

import (
	"context"
	"log"

	"github.com/dnsoftware/mpm-save-get-shares/config"
	"github.com/dnsoftware/mpm-save-get-shares/internal/app"
	"github.com/dnsoftware/mpm-save-get-shares/internal/constants"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/utils"
)

func main() {
	ctx := context.Background()

	basePath, err := utils.GetProjectRoot(constants.ProjectRootAnchorFile)
	if err != nil {
		log.Fatalf("GetProjectRoot failed: %s", err.Error())
	}
	configFile := basePath + "/config.yaml"
	envFile := basePath + "/.env"

	c, err := config.New(configFile, envFile)
	if err != nil {

	}

	err = app.Run(ctx, c)

}
