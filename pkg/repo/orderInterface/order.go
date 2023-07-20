package orderinterface

import (
	"github.com/RohithER12/order-svc/pkg/models"
	"github.com/RohithER12/order-svc/pkg/repo"
)

type Order interface {
	CreateOrder(order models.Order) (int64, error)
}

func NewOrderImpl() Order {
	return &repo.OrderRepo{}
}
