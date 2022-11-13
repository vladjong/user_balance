package handler

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vladjong/user_balance/config"
)

func TestCheckNegativeDecimal(t *testing.T) {
	testTable := []struct {
		numbers  decimal.Decimal
		expected bool
	}{
		{
			numbers:  decimal.NewFromFloat(136.02),
			expected: false,
		},
		{
			numbers:  decimal.NewFromFloat(-136.02),
			expected: true,
		},
		{
			numbers:  decimal.NewFromFloat(0),
			expected: false,
		},
	}
	for _, testCase := range testTable {
		result := checkNegativeDecimal(testCase.numbers)
		t.Logf("Calling checkNegativeDecimal(%s), result %v\n", testCase.numbers.String(), result)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestCheckIsBalanceServer(t *testing.T) {
	testTable := []struct {
		numbers  int
		expected bool
	}{
		{
			numbers:  config.OrderBalanceId,
			expected: true,
		},
		{
			numbers:  config.ServiceBalanceId,
			expected: true,
		},
		{
			numbers:  1,
			expected: false,
		},
	}
	for _, testCase := range testTable {
		result := checkIsBalanceServer(testCase.numbers)
		t.Logf("Calling checkIsBalanceServer(%d), result %v\n", testCase.numbers, result)
		assert.Equal(t, testCase.expected, result)
	}
}
