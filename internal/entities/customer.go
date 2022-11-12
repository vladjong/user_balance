package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Customer struct {
	Id      int             `json:"id" db:"id"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}

type Acount struct {
	CustomerId int             `json:"customer_id" db:"customer_id"`
	Balance    decimal.Decimal `json:"balance" db:"balance"`
}

type CustomerReport struct {
	Id          int             `json:"id" db:"id"`
	ServiceName string          `json:"service_name" db:"service_name"`
	OrderName   string          `json:"order_name" db:"order_name"`
	Sum         decimal.Decimal `json:"sum" db:"sum"`
	Date        time.Time       `json:"date" db:"date"`
}
