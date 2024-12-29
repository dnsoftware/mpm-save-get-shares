package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dnsoftware/mpm-save-get-shares/internal/entity"
)

type PostgresMinerStorage struct {
	pool *pgxpool.Pool
	//db   *sql.DB
}

func NewPostgresMinerStorage(pool *pgxpool.Pool) (*PostgresMinerStorage, error) {

	storage := &PostgresMinerStorage{
		pool: pool,
	}

	return storage, nil
}

func (p *PostgresMinerStorage) CreateWallet(ctx context.Context, wallet entity.Wallet) (int64, error) {
	var newID int64
	err := p.pool.QueryRow(ctx, `INSERT INTO wallets (coin_id, name, is_solo, reward_method) 
			VALUES ($1, $2, $3, $4) RETURNING id`,
		wallet.CoinID, wallet.Name, wallet.IsSolo, wallet.RewardMethod).Scan(&newID)

	return newID, err
}

func (p *PostgresMinerStorage) CreateWorker(ctx context.Context, worker entity.Worker) (int64, error) {
	created_at := time.Now().Format("2006-01-02 15:04:05.000")
	updated_at := time.Now().Format("2006-01-02 15:04:05.000")

	var newID int64
	err := p.pool.QueryRow(ctx, `INSERT INTO workers (coin_id, workerfull, wallet, worker, server_id, created_at, updated_at, reward_method) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		worker.CoinID, worker.Workerfull, worker.Wallet, worker.Worker, worker.ServerID, created_at, updated_at, worker.RewardMethod).Scan(&newID)

	return newID, err
}

func (p *PostgresMinerStorage) GetWalletIDByName(ctx context.Context, wallet string, coinID int64, rewardMethod string) (int64, error) {
	var id int64
	err := p.pool.QueryRow(ctx, `SELECT id FROM wallets WHERE name = $1 AND coin_id = $2 AND reward_method = $3`,
		wallet, coinID, rewardMethod).Scan(&id)

	return id, err
}

func (p *PostgresMinerStorage) GetWorkerIDByName(ctx context.Context, workerFull string, coinID int64, rewardMethod string) (int64, error) {

	var id int64
	err := p.pool.QueryRow(ctx, `SELECT id FROM workers WHERE workerfull = $1 AND coin_id = $2 AND reward_method = $3`,
		workerFull, coinID, rewardMethod).Scan(&id)

	return id, err
}