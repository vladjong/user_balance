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
	customer, err := u.storage.GetCustomerBalance(customerId)
	if err != nil {
		return err
	}
	if customer.Balance.LessThan(value) {
		return errors.New("error: customer balance less than transaction cost")
	}
	customer.Balance = customer.Balance.Sub(value)
	account := entities.Acount{
		CustomerId: customerId,
		Balance:    value,
	}
	transaction := entities.Transaction{
		CustomeId:           customerId,
		ServiceID:           serviceId,
		OrderID:             orderId,
		Cost:                value,
		TransactionDatiTime: time.Now(),
	}
	return u.storage.PostReserveBalance(transaction, customer, account)
}

func (u *userBalanseUseCase) PostDeReservingBalance(customerId, serviceId, orderId int, value decimal.Decimal, status bool) error {
	if ok := checkValue(value); !ok {
		return errors.New("error: value is negative")
	}
	customer, err := u.storage.GetCustomerBalance(customerId)
	if err != nil {
		return err
	}
	if !status {
		customer.Balance = customer.Balance.Add(value)
	}
	account := entities.Acount{
		CustomerId: customerId,
		Balance:    value,
	}
	id, err := u.storage.GetTransactionId(entities.Transaction{
		CustomeId: customerId,
		ServiceID: serviceId,
		OrderID:   orderId,
		Cost:      value,
	})
	if err != nil {
		return err
	}
	history := entities.History{
		TransactionId:      id,
		AccountingDatetime: time.Now(),
		StatusTransaction:  status,
	}
	return u.storage.PostDeReservingBalance(customer, account, history)
}

func checkValue(value decimal.Decimal) bool {
	return !value.IsNegative()
}
