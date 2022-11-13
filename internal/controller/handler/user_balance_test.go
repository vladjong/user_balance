package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vladjong/user_balance/internal/entities"
	mock_usecase "github.com/vladjong/user_balance/internal/usecase/mocks"
)

func TestHandler_getCustomerBalance(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, id int)
	testTable := []struct {
		name                string
		inputId             int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:    "Ok",
			inputId: 1,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int) {
				s.EXPECT().GetCustomerBalance(id).Return(entities.Customer{
					Id:      1,
					Balance: decimal.NewFromFloat(33.3),
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"balance":"33.3"}`,
		},
		{
			name:    "Service error",
			inputId: 4,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int) {
				s.EXPECT().GetCustomerBalance(id).Return(entities.Customer{}, errors.New("error: id don't exist"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"message\":\"error: id don't exist\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			testCase.mockBehavior(user_balance, testCase.inputId)
			handler := New(user_balance)
			r := gin.New()
			r.GET("/:id", handler.getCustomerBalance)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%d", testCase.inputId), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
