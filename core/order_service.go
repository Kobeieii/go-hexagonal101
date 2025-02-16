package core

import (
	"errors"
)

type OrderService interface {
	CreateOrder(order *Order) error
}

type OrderServiceImp struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderService {
	return &OrderServiceImp{repo: repo}
}

func (s *OrderServiceImp) CreateOrder(order *Order) error {
	if order.Total < 0 {
		return errors.New("total must be positive")
	}
	if err := s.repo.Save(order); err != nil {
		return err
	}
	return nil
}