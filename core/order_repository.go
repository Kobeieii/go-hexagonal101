package core

type OrderRepository interface {
	Save(order *Order) error
}
