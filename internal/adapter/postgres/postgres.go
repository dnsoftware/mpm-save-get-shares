package postgres

import (
	"context"
	"database/sql"

	"github.com/dnsoftware/mpm-save-get-shares/internal/entity"
	"github.com/dnsoftware/mpm-save-get-shares/pkg/logger"
)

type PostgresMinerStorage struct {
	db *sql.DB
}

func NewPostgresMinerStorage(dsn string) (*PostgresMinerStorage, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log().Error("NewPostgresMinerStorage: " + err.Error())
		return nil, err
	}

	storage := &PostgresMinerStorage{
		db: db,
	}

	return storage, nil
}

func (p *PostgresMinerStorage) CreateWallet(ctx context.Context, wallet entity.Wallet) (int64, error) {

	return 0, nil
}

func (p *PostgresMinerStorage) CreateWorker(ctx context.Context, worker entity.Worker) (int64, error) {

	return 0, nil
}

func (p *PostgresMinerStorage) GetWalletIDByName(ctx context.Context, wallet string) (int64, error) {

	return 0, nil
}

func (p *PostgresMinerStorage) GetWorkerIDByName(ctx context.Context, worker string) (int64, error) {

	return 0, nil
}
