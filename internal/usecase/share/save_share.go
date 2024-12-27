package share

import (
	"github.com/dnsoftware/mpm-save-get-shares/internal/dto"
)

// SaveShare сохранение шары в базе данных (ClickHouse)
// Возвращает nil, если запись быда добавлена успешно
func (u *UseCase) SaveShare(shareFound dto.ShareFound) error {
	// Worker
	// Проверяем в кэше по имени воркера, если есть - запоминаем ID воркера

	// Если в кэше нет запрашиваем в удаленной базе, если есть - запоминаем ID воркера

	// если в базе нет добавляем в базу, получаем ID воркера, кэшируем

	// Wallet
	// получаем имя кошелька из имени воркера

	// Проверяем в кэше по имени кошелька, если есть - запоминаем ID кошелька

	// Если в кэше нет запрашиваем в удаленной базе, если есть - запоминаем ID кошелька

	// если в базе нет добавляем в базу, получаем ID кошелька, кэшируем

	// Coin
	// Проверяем в кэше по имени монеты, если есть - запоминаем ID монеты

	// Если в кэше нет запрашиваем в удаленной базе, если есть - запоминаем ID монеты

	// если в базе нет добавляем в базу, получаем ID монеты, кэшируем

	// Save share
	// Сохраняем шару в хранилище

	return nil
}