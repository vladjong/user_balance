package usecase

import (
	"errors"
	"time"

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

func (u *userBalanseUseCase) PostReserveBalance(customerId, serviceId, orderId int, value decimal.Decimal) error {
	if ok := checkValue(value); !ok {
		return errors.New("error: value is negative")
	}
	transaction := entities.Transaction{
		CustomeId:           customerId,
		ServiceID:           serviceId,
		OrderID:             orderId,
		Cost:                value,
		TransactionDatiTime: time.Now(),
	}
	return u.storage.PostReserveBalance(transaction)
}

func (u *userBalanseUseCase) PostDeReservingBalance(customerId, serviceId, orderId int, value decimal.Decimal, status bool) error {
	if ok := checkValue(value); !ok {
		return errors.New("error: value is negative")
	}
	transaction := entities.Transaction{
		CustomeId: customerId,
		ServiceID: serviceId,
		OrderID:   orderId,
		Cost:      value,
	}
	history := entities.History{
		TransactionId:      customerId,
		StatusTransaction:  status,
		AccountingDatetime: time.Now(),
	}
	return u.storage.PostDeReservingBalance(transaction, history)
}

func checkValue(value decimal.Decimal) bool {
	return !value.IsNegative()
}
