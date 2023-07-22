package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetTotalTransaction() (int, error)
}