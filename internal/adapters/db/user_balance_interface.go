package db

import (
	"github.com/vladjong/user_balance/internal/entities"
)

type UserBalanse interface {
	GetCustomerBalance(id int) (customer entities.Customer, err error)
	PostCustomerBalance(customer entities.Customer) error
	PostReserveBalance(transaction entities.Transaction, customer entities.Customer, account entities.Acount) error
	PostDeReservingBalance(customer entities.Customer, account entities.Acount, history entities.History) error
	GetTransactionId(transaction entities.Transaction) (int, error)
}
