package entities

import "time"

type History struct {
	Id                 int       `json:"id" db:"id"`
	TransactionId      int       `json:"transaction_id" db:"transaction_id"`
	AccountingDatetime time.Time `json:"accounting_datetime" db:"accounting_datetime"`
	StatusTransaction  bool      `json:"status_transaction" db:"status_transaction"`
}
