package usecase

import (
	"github.com/shopspring/decimal"
	"github.com/vladjong/user_balance/internal/entities"
)

type UserBalanse interface {
	GetCustomerBalance(id int) (user entities.Customer, err error)
	PostCustomerBalance(id int, value decimal.Decimal) error
}
