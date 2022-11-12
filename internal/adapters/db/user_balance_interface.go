package db

import (
	"github.com/vladjong/user_balance/internal/entities"
)

type UserBalanse interface {
	GetCustomerBalance(id int) (customer entities.Customer, err error)
	PostCustomerBalance(customer entities.Customer) error
	PostReserveBalance(transaction entities.Transaction) error
	PostDeReservingBalance(transaction entities.Transaction, history entities.History) error
}
