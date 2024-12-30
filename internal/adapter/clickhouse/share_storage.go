package clickhouse

import "github.com/dnsoftware/mpm-save-get-shares/internal/entity"

type ClickHouseShareStorage struct {
}

func NewClickHouseShareStorage() (*ClickHouseShareStorage, error) {

	ss := &ClickHouseShareStorage{}

	return ss, nil
}

// CreateShare добавление шары в ClickHouse базу данных
// если возвращает nil - вставка прошла успешно
func (c *ClickHouseShareStorage) CreateShare(share entity.Share) error {
	return nil
}
