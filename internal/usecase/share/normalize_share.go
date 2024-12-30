package share

import (
	"context"
	"fmt"
	"time"

	"github.com/dnsoftware/mpm-save-get-shares/internal/dto"
	"github.com/dnsoftware/mpm-save-get-shares/internal/entity"
)

// NormalizeShare нормализация шары
// получаем коды монеты. майнера, воркера из кеша или из базы, чтобы сформировать структуру шары для вставки в базу данных
func (u *ShareUseCase) NormalizeShare(shareFound dto.ShareFound) (entity.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	//*** получение кода монеты в базе
	coinID, err := u.coinCache.GetCoinIDByName(shareFound.CoinSymbol)
	if err != nil {
		return entity.Share{}, err
	}
	if coinID == 0 { // в кэше нет
		coinID, err = u.coinStorage.GetCoinIDByName(ctx, shareFound.CoinSymbol)
		if err != nil {
			return entity.Share{}, err
		}
		if coinID == 0 { // в базе нет (а должна быть, так как в миграциях заполнили все монеты в таблице)
			return entity.Share{}, fmt.Errorf("coinID must be greater then 0")
		}

		coinID, err = u.coinCache.CreateCoin(shareFound.CoinSymbol, coinID) // кешируем
		if err != nil {
			return entity.Share{}, err
		}
		if coinID == 0 { // в кеше должна быть уже
			return entity.Share{}, fmt.Errorf("coinID in cache must be greater then 0")
		}

	}

	//*** получение кода майнера(кошелька) в базе
	walletName := WalletFromWorkerfull(shareFound.Workerfull)

	// Запрашиваем данные майнера из кеша
	walletID, err := u.minerCache.GetWalletIDByName(walletName, coinID, shareFound.RewardMethod)
	if err != nil {
		return entity.Share{}, err
	}
	if walletID == 0 { // нет в кеше
		walletID, err = u.minerStorage.GetWalletIDByName(ctx, walletName, coinID, shareFound.RewardMethod)
		if err != nil {
			return entity.Share{}, err
		}
		if walletID == 0 { // нет в базе
			walletEntity := entity.Wallet{
				ID:           0,
				CoinID:       coinID,
				Name:         walletName,
				IsSolo:       shareFound.IsSolo,
				RewardMethod: shareFound.RewardMethod,
			}

			walletID, err = u.minerStorage.CreateWallet(ctx, walletEntity)
			if err != nil {
				return entity.Share{}, err
			}
			walletEntity.ID = walletID

			walletID, err = u.minerCache.CreateWallet(walletEntity)
			walletEntity.ID = walletID
		}
	}

	//*** получение кода воркера в базе
	workerName := WorkerFromWorkerfull(shareFound.Workerfull)

	// Запрашиваем данные воркера из кеша
	workerID, err := u.minerCache.GetWorkerIDByName(shareFound.Workerfull, coinID, shareFound.RewardMethod)
	if err != nil {
		return entity.Share{}, err
	}
	if workerID == 0 { // нет в кеше
		workerID, err = u.minerStorage.GetWorkerIDByName(ctx, shareFound.Workerfull, coinID, shareFound.RewardMethod)
		if err != nil {
			return entity.Share{}, err
		}
		if workerID == 0 { // нет в базе
			workerEntity := entity.Worker{
				ID:           0,
				CoinID:       coinID,
				Workerfull:   shareFound.Workerfull,
				Wallet:       walletName,
				Worker:       workerName,
				ServerID:     shareFound.ServerID,
				IP:           shareFound.MinerIp,
				IsSolo:       shareFound.IsSolo,
				RewardMethod: shareFound.RewardMethod,
			}

			workerID, err = u.minerStorage.CreateWorker(ctx, workerEntity)
			if err != nil {
				return entity.Share{}, err
			}
			workerEntity.ID = workerID

			workerID, err = u.minerCache.CreateWorker(workerEntity)
			workerEntity.ID = workerID
		}
	}

	// формируем entity.Share
	share := shareFound.ToShare()
	share.CoinID = coinID
	share.WalletID = walletID
	share.WorkerID = workerID

	return share, nil
}
