package entities

import "github.com/shopspring/decimal"

type Customer struct {
	Id      int             `json:"id" db:"id"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}
