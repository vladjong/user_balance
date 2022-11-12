package db

import (
	"time"

	"github.com/vladjong/user_balance/internal/entities"
)

type UserBalanse interface {
	GetCustomerBalance(id int) (customer entities.Customer, err error)
	GetHistoryReport(date time.Time) (report []entities.Report, err error)
	GetCustomerReport(id int, date time.Time) (report []entities.CustomerReport, err error)
	PostCustomerBalance(customer entities.Customer, transaction entities.Transaction) error
	PostReserveBalance(transaction entities.Transaction) error
	PostDeReservingBalance(transaction entities.Transaction, history entities.History) error
}
