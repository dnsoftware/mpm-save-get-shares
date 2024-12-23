package app

import (
	"context"

	"github.com/dnsoftware/mpm-save-get-shares/config"
)

type Dependencies struct {
}

func Run(ctx context.Context, cfg config.Config) (err error) {
	var deps Dependencies
	_ = deps

	/********* Инициализация трассировщика **********/
	/********* КОНЕЦ Инициализация трассировщика **********/

	return nil
}
