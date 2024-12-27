package entity

type Wallet struct {
	CoinID       int64
	Name         string
	IsSolo       bool   // оставлено для совместимости TODO убрать
	RewardMethod string // строковый код метода распределения наград
}
