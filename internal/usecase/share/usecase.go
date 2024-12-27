package share

import (
	"context"

	"github.com/dnsoftware/mpm-save-get-shares/internal/entity"
)

// SharesStorage сохранение шары в хранилище (ClickHouse)
type ShareStorage interface {
	CreateShare(share entity.Share) error // если возвращает nil - вставка прошла учпешно
}

// MinerStorage сохранение/получение данных майнеров (кошельков) и воркеров в справочники в хранилище (Postgresql или кэш (ristretto))
type MinerStorage interface {
	CreateWallet(ctx context.Context, wallet entity.Wallet) (int64, error)
	CreateWorker(ctx context.Context, worker entity.Worker) (int64, error)
	GetWalletIDByName(ctx context.Context, wallet string) (int64, error) // 0 - если не найден
	GetWorkerIDByName(ctx context.Context, worker string) (int64, error) // 0 - если не найден
}

// CoinStorage работа с данными о монете из хранилища  (Postgresql или кэш (ristretto))
type CoinStorage interface {
	// GetCoinIDByName получение кода монеты в базе по буквенному коду (ALPH, KAS и т.д.), если не найдено - возвращаем ошибку
	GetCoinIDByName(ctx context.Context, coin string) (int64, error)
}

// MinerCache работа с кэшированными данными майнеров
type MinerCache interface {
	CreateWallet(wallet string) (int64, error)
	CreateWorker(worker string) (int64, error)
	GetWalletIDByName(wallet string) (int64, error) // 0 - если не найден
	GetWorkerIDByName(worker string) (int64, error) // 0 - если не найден
}

// CoinCache работа с кэшированными данными о монете
type CoinCache interface {
	GetCoinIDByName(coin string) (int64, error) // получение кода монеты из кэша по буквенному коду (ALPH, KAS и т.д.)
}

type ShareUseCase struct {
	shareStorage ShareStorage // персистентная база (ClickHouse)
	minerStorage MinerStorage // персистентная база (Postgresql)
	coinStorage  CoinStorage  // персистентная база (Postgresql)
	minerCache   MinerCache   // кэш в оперативной памяти для майнеров
	coinCache    CoinCache    // кэш в оперативной памяти для монет
}

func NewShareUseCase(s ShareStorage, m MinerStorage, c CoinStorage, mc MinerCache, cc CoinCache) *ShareUseCase {
	return &ShareUseCase{
		shareStorage: s,
		minerStorage: m,
		coinStorage:  c,
		minerCache:   mc,
		coinCache:    cc,
	}
}
