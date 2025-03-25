package enum

import (
	"database/sql/driver"
	"fmt"
)

type OrderStatus string

const (
	OrderPlaced   OrderStatus = "Placed"
	OrderPaid     OrderStatus = "Paid"
	OrderShipped  OrderStatus = "Shipped out"
	OrderReceived OrderStatus = "Received"
)

func (r *OrderStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}

	*r = OrderStatus(str)
	return nil
}

func (r *OrderStatus) Value() (driver.Value, error) {
	return string(*r), nil
}
