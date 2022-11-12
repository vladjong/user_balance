package usecase

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/vladjong/user_balance/internal/entities"
)

type UserBalanse interface {
	GetCustomerBalance(id int) (user entities.Customer, err error)
	GetHistoryReport(date time.Time) (string, error)
	GetCustomerReport(id int, date time.Time) (report []entities.CustomerReport, err error)
	PostCustomerBalance(id int, value decimal.Decimal) error
	PostReserveBalance(customerId, serviceId, orderId int, value decimal.Decimal) error
	PostDeReservingBalance(customerId, serviceId, orderId int, value decimal.Decimal, status bool) error
}
