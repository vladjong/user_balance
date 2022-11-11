package usecase

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/vladjong/user_balance/internal/adapters/db"
	"github.com/vladjong/user_balance/internal/entities"
)

type userBalanseUseCase struct {
	storage db.UserBalanse
}

func New(storage db.UserBalanse) *userBalanseUseCase {
	return &userBalanseUseCase{
		storage: storage,
	}
}

func (u *userBalanseUseCase) GetCustomerBalance(id int) (user entities.Customer, err error) {
	return u.storage.GetCustomerBalance(id)
}

func (u *userBalanseUseCase) PostCustomerBalance(id int, value decimal.Decimal) error {
	if ok := checkValue(value); !ok {
		return errors.New("error: value is negative")
	}
	customer := entities.Customer{
		Id:      id,
		Balance: value,
	}
	return u.storage.PostCustomerBalance(customer)
}

func checkValue(value decimal.Decimal) bool {
	return !value.IsNegative()
}
