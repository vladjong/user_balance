package handler

import (
	"github.com/shopspring/decimal"
	"github.com/vladjong/user_balance/config"
)

func checkNegativeDecimal(value decimal.Decimal) bool {
	return value.IsNegative()
}

func checkIsBalanceServer(id int) bool {
	return id == config.ServiceBalanceId
}
