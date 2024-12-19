package app

import (
	"context"

	"github.com/dnsoftware/mpm-save-get-shares/config"
)

type Dependencies struct {
}

func Run(ctx context.Context, c config.Config) (err error) {
	var deps Dependencies
	_ = deps

	return nil
}
