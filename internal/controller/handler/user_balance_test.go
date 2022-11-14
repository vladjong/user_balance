package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vladjong/user_balance/config"
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
			r.GET("/:id", handler.GetCustomerBalance)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%d", testCase.inputId), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getHistoryReport(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, date time.Time)
	testTable := []struct {
		name                string
		input               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "Ok",
			input: "2022-01",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, date time.Time) {
				s.EXPECT().GetHistoryReport(date).Return("report_2022-01", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"Filename":"report_2022-01"}`,
		},
		{
			name:  "Status bad request",
			input: "qwerty",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, date time.Time) {
				s.EXPECT().GetHistoryReport(date).Return("", nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"message\":\"parsing time \\\"qwerty\\\" as \\\"2006-01\\\": cannot parse \\\"qwerty\\\" as \\\"2006\\\"\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			date, err := time.Parse(config.DateFormat, testCase.input)
			if err == nil {
				testCase.mockBehavior(user_balance, date)
			}
			handler := New(user_balance)
			r := gin.New()
			r.GET("/report/:date", handler.GetHistoryReport)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/report/%s", testCase.input), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getCustomerReport(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, id int, date time.Time)
	testTable := []struct {
		name                string
		inputId             string
		inputDate           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			inputId:   "1",
			inputDate: "2022-01",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int, date time.Time) {
				s.EXPECT().GetCustomerReport(id, date).Return([]entities.CustomerReport{
					{
						Id:                1,
						ServiceName:       "Balance",
						OrderName:         "Balance",
						Sum:               decimal.NewFromFloat(131.1),
						StatusTransaction: true,
						Date:              time.Date(2006, time.January, 1, 0, 0, 0, 0, time.Local),
					}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"id":1,"service_name":"Balance","order_name":"Balance","sum":"131.1","status_transaction":true,"date":"2006-01-01T00:00:00+06:00"}]`,
		},
		{
			name:      "Status bad request",
			inputId:   "1",
			inputDate: "qwerty",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int, date time.Time) {
				s.EXPECT().GetCustomerReport(id, date).Return(nil, nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"message\":\"parsing time \\\"qwerty\\\" as \\\"2006-01\\\": cannot parse \\\"qwerty\\\" as \\\"2006\\\"\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			check := true
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			id, err := strconv.Atoi(testCase.inputId)
			if err != nil {
				check = false
			}
			date, err := time.Parse(config.DateFormat, testCase.inputDate)
			if err != nil {
				check = false
			}
			if check {
				testCase.mockBehavior(user_balance, id, date)
			}
			handler := New(user_balance)
			r := gin.New()
			r.GET("/history/:id/:date", handler.GetCustomerReport)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/history/%s/%s", testCase.inputId, testCase.inputDate), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_postCustomerBalance(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, id int, value decimal.Decimal)
	testTable := []struct {
		name                string
		inputId             string
		inputValue          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "Ok",
			inputId:    "1",
			inputValue: "100",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int, value decimal.Decimal) {
				s.EXPECT().PostCustomerBalance(id, value).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"Status":"ok"}`,
		},
		{
			name:       "Status bad request",
			inputId:    "qwerty",
			inputValue: "100",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int, value decimal.Decimal) {
				s.EXPECT().PostCustomerBalance(id, value).Return(nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid customer id param"}`,
		},
		{
			name:       "Status bad internal request",
			inputId:    "1",
			inputValue: "1",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id int, value decimal.Decimal) {
				s.EXPECT().PostCustomerBalance(id, value).Return(errors.New("error:don't exits id"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"error:don't exits id"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			check := true
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			id, err := strconv.Atoi(testCase.inputId)
			if err != nil {
				check = false
			}
			value, err := decimal.NewFromString(testCase.inputValue)
			if err != nil {
				check = false
			}
			if check {
				testCase.mockBehavior(user_balance, id, value)
			}
			handler := New(user_balance)
			r := gin.New()
			r.POST("/:id/:val", handler.PostCustomerBalance)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", testCase.inputId, testCase.inputValue), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_postReserveCustomerBalance(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal)
	testTable := []struct {
		name                string
		inputId             string
		inputSer            string
		inputOrd            string
		inputValue          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "Ok",
			inputId:    "1",
			inputSer:   "1",
			inputOrd:   "1",
			inputValue: "100",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal) {
				s.EXPECT().PostReserveBalance(id, idSer, idOrd, value).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"Status":"ok"}`,
		},
		{
			name:       "Status bad request",
			inputId:    "qwerty",
			inputSer:   "1",
			inputOrd:   "1",
			inputValue: "100",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal) {
				s.EXPECT().PostReserveBalance(id, idSer, idOrd, value).Return(nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid customer id param"}`,
		},
		{
			name:       "Status bad internal request",
			inputId:    "1",
			inputSer:   "1",
			inputOrd:   "1",
			inputValue: "100",
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal) {
				s.EXPECT().PostReserveBalance(id, idSer, idOrd, value).Return(errors.New("error:don't exits id"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"error:don't exits id"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			check := true
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			id, err := strconv.Atoi(testCase.inputId)
			if err != nil {
				check = false
			}
			serviceId, err := strconv.Atoi(testCase.inputSer)
			if err != nil {
				check = false
			}
			orderId, err := strconv.Atoi(testCase.inputOrd)
			if err != nil {
				check = false
			}
			value, err := decimal.NewFromString(testCase.inputValue)
			if err != nil {
				check = false
			}
			if check {
				testCase.mockBehavior(user_balance, id, serviceId, orderId, value)
			}
			handler := New(user_balance)
			r := gin.New()
			r.POST("/reserv/:id/:id_ser/:id_ord/:val", handler.PostReserveCustomerBalance)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/reserv/%s/%s/%s/%s", testCase.inputId, testCase.inputSer, testCase.inputOrd, testCase.inputValue), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_postDeReservingBalanceAccept(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool)
	testTable := []struct {
		name                string
		inputId             string
		inputSer            string
		inputOrd            string
		inputValue          string
		inputStatus         bool
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "Ok",
			inputId:     "1",
			inputSer:    "1",
			inputOrd:    "1",
			inputValue:  "100",
			inputStatus: true,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool) {
				s.EXPECT().PostDeReservingBalance(id, idSer, idOrd, value, status).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"Status":"ok"}`,
		},
		{
			name:        "Status bad request",
			inputId:     "qwerty",
			inputSer:    "1",
			inputOrd:    "1",
			inputValue:  "100",
			inputStatus: true,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool) {
				s.EXPECT().PostDeReservingBalance(id, idSer, idOrd, value, status).Return(nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid customer id param"}`,
		},
		{
			name:        "Status bad internal request",
			inputId:     "1",
			inputSer:    "1",
			inputOrd:    "1",
			inputValue:  "100",
			inputStatus: true,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool) {
				s.EXPECT().PostDeReservingBalance(id, idSer, idOrd, value, status).Return(errors.New("error:don't exits id"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"error:don't exits id"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			check := true
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			id, err := strconv.Atoi(testCase.inputId)
			if err != nil {
				check = false
			}
			serviceId, err := strconv.Atoi(testCase.inputSer)
			if err != nil {
				check = false
			}
			orderId, err := strconv.Atoi(testCase.inputOrd)
			if err != nil {
				check = false
			}
			value, err := decimal.NewFromString(testCase.inputValue)
			if err != nil {
				check = false
			}
			if check {
				testCase.mockBehavior(user_balance, id, serviceId, orderId, value, testCase.inputStatus)
			}
			handler := New(user_balance)
			r := gin.New()
			r.POST("/accept/:id/:id_ser/:id_ord/:val", handler.PostDeReservingBalanceAccept)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/accept/%s/%s/%s/%s", testCase.inputId, testCase.inputSer, testCase.inputOrd, testCase.inputValue), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_postDeReservingBalanceReject(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool)
	testTable := []struct {
		name                string
		inputId             string
		inputSer            string
		inputOrd            string
		inputValue          string
		inputStatus         bool
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "Ok",
			inputId:     "1",
			inputSer:    "1",
			inputOrd:    "1",
			inputValue:  "100",
			inputStatus: false,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool) {
				s.EXPECT().PostDeReservingBalance(id, idSer, idOrd, value, status).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"Status":"ok"}`,
		},
		{
			name:        "Status bad request",
			inputId:     "qwerty",
			inputSer:    "1",
			inputOrd:    "1",
			inputValue:  "100",
			inputStatus: false,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool) {
				s.EXPECT().PostDeReservingBalance(id, idSer, idOrd, value, status).Return(nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid customer id param"}`,
		},
		{
			name:        "Status bad internal request",
			inputId:     "1",
			inputSer:    "1",
			inputOrd:    "1",
			inputValue:  "100",
			inputStatus: false,
			mockBehavior: func(s *mock_usecase.MockUserBalanse, id, idSer, idOrd int, value decimal.Decimal, status bool) {
				s.EXPECT().PostDeReservingBalance(id, idSer, idOrd, value, status).Return(errors.New("error:don't exits id"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"error:don't exits id"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			check := true
			user_balance := mock_usecase.NewMockUserBalanse(ctr)
			id, err := strconv.Atoi(testCase.inputId)
			if err != nil {
				check = false
			}
			serviceId, err := strconv.Atoi(testCase.inputSer)
			if err != nil {
				check = false
			}
			orderId, err := strconv.Atoi(testCase.inputOrd)
			if err != nil {
				check = false
			}
			value, err := decimal.NewFromString(testCase.inputValue)
			if err != nil {
				check = false
			}
			if check {
				testCase.mockBehavior(user_balance, id, serviceId, orderId, value, testCase.inputStatus)
			}
			handler := New(user_balance)
			r := gin.New()
			r.POST("/reject/:id/:id_ser/:id_ord/:val", handler.PostDeReservingBalanceReject)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/reject/%s/%s/%s/%s", testCase.inputId, testCase.inputSer, testCase.inputOrd, testCase.inputValue), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
