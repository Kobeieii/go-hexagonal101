package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockOrderRepo struct {
	saveFunc func(order *Order) error
}

func (r *mockOrderRepo) Save(order *Order) error {
	return r.saveFunc(order)
}

func TestCreateOrder(t *testing.T) {
	t.Run("save success", func(t *testing.T) {
		repo := &mockOrderRepo{
			saveFunc: func(order *Order) error {
				return nil
			},
		}

		service := NewOrderService(repo)
		err := service.CreateOrder(&Order{Total: 100})
		assert.NoError(t, err)
	})

	t.Run("total less than 0", func(t *testing.T) {
		repo := &mockOrderRepo{
			saveFunc: func(order *Order) error {
				return nil
			},
		}
		service := NewOrderService(repo)

		err := service.CreateOrder(&Order{Total: -10})
		assert.Error(t, err)
		assert.Equal(t, "total must be positive", err.Error())
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &mockOrderRepo{
			saveFunc: func(order *Order) error {
				return errors.New("save failed")
			},
		}
		service := NewOrderService(repo)

		err := service.CreateOrder(&Order{Total: 10})
		assert.Error(t, err)
		assert.Equal(t, "save failed", err.Error())
	})
}
