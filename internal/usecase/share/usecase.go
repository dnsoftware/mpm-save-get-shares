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
	CreateWallet(ctx context.Context, wallet string) (int64, error)
	CreateWorker(ctx context.Context, worker string) (int64, error)
	GetWalletIDByName(ctx context.Context, wallet string) (int64, error)
	GetWorkerIDByName(ctx context.Context, worker string) (int64, error)
}

// CoinStorage получение данных о монете из хранилища  (Postgresql или кэш (ristretto))
type CoinStorage interface {
	GetCoinIDByName(ctx context.Context, coin string) (int64, error) // получение кода монеты в базе по буквенному коду (ALPH, KAS и т.д.)
}

type CacheStorage interface {
	MinerStorage
	CoinStorage
}

type UseCase struct {
	shareStorage ShareStorage // персистентная база (ClickHouse)
	minerStorage MinerStorage // персистентная база (Postgresql)
	coinStorage  CoinStorage  // персистентная база (Postgresql)
	cacheStorage CacheStorage // кэш в оперативной памяти
}

func NewShareUseCase(s ShareStorage, m MinerStorage, c CoinStorage, cache CacheStorage) *UseCase {
	return &UseCase{
		shareStorage: s,
		minerStorage: m,
		coinStorage:  c,
		cacheStorage: cache,
	}
}
