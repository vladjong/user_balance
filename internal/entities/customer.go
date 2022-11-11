package entities

import "github.com/shopspring/decimal"

type Customer struct {
	Id      int             `json:"id" db:"id"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}

type Acount struct {
	CustomerId int             `json:"customer_id" db:"customer_id"`
	Balance    decimal.Decimal `json:"balance" db:"balance"`
}
