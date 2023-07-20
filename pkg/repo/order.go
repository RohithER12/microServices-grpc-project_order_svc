package repo

import (
	"github.com/RohithER12/order-svc/pkg/db"
	"github.com/RohithER12/order-svc/pkg/models"
)

type OrderRepo struct {
	H db.Handler
}

func (o *OrderRepo) CreateOrder(order models.Order) (int64, error) {
	result := o.H.DB.Create(&order)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.Id, nil
}
