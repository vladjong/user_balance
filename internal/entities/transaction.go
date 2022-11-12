package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id                  int             `json:"id" db:"id"`
	CustomeId           int             `json:"customer_id" db:"customer_id"`
	ServiceID           int             `json:"service_id" db:"service_id"`
	OrderID             int             `json:"order_id" db:"order_id"`
	Cost                decimal.Decimal `json:"cost" db:"cost"`
	TransactionDatiTime time.Time       `json:"transaction_datetime" db:"transaction_datetime"`
}
