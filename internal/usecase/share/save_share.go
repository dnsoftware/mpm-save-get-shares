package share

import (
	"github.com/dnsoftware/mpm-save-get-shares/internal/dto"
)

// SaveShare сохранение шары в базе данных (ClickHouse)
// Возвращает nil, если запись была добавлена успешно
func (u *ShareUseCase) SaveShare(shareFound dto.ShareFound) error {
	//*** Coin
	// Проверяем в кэше по имени монеты, если есть - запоминаем ID монеты

	// Если в кэше нет запрашиваем в удаленной базе, если есть - запоминаем ID монеты

	// если в базе нет добавляем в базу, получаем ID монеты, кэшируем

	//*** Worker
	// Маппим DTO в entity структуру
	//share := shareFound.ToShare()

	//// Проверяем в кэше по имени воркера, если есть - запоминаем ID воркера
	//workerID, err := u.minerCache.GetWorkerIDByName(shareFound.Workerfull, )
	//if err != nil {
	//	return err
	//}
	//
	//// Если в кэше нет запрашиваем в удаленной базе, если есть - запоминаем ID воркера
	//if workerID == 0 {
	//	ctx, cancel1 := context.WithTimeout(context.Background(), constants.QueryDealine*time.Second)
	//	defer cancel1()
	//
	//	workerID, err = u.minerStorage.GetWorkerIDByName(ctx, shareFound.Workerfull)
	//	if err != nil {
	//		return err
	//	}
	//}

	// если в базе нет добавляем в базу, получаем ID воркера, кэшируем

	// Wallet
	// получаем имя кошелька из имени воркера

	// Проверяем в кэше по имени кошелька, если есть - запоминаем ID кошелька

	// Если в кэше нет запрашиваем в удаленной базе, если есть - запоминаем ID кошелька

	// если в базе нет добавляем в базу, получаем ID кошелька, кэшируем

	// Save share
	// Сохраняем шару в хранилище

	return nil
}
