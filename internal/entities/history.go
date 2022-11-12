package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type History struct {
	Id                 int       `json:"id" db:"id"`
	TransactionId      int       `json:"transaction_id" db:"transaction_id"`
	AccountingDatetime time.Time `json:"accounting_datetime" db:"accounting_datetime"`
	StatusTransaction  bool      `json:"status_transaction" db:"status_transaction"`
}

type Report struct {
	Id     int             `json:"id" db:"id"`
	Name   string          `json:"name" db:"name"`
	AllSum decimal.Decimal `json:"all_sum" db:"all_sum"`
}
