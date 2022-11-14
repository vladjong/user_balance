package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/vladjong/user_balance/config"
)

// @Summary Get Customer balance
// @Tags customer
// @Description get by INT id
// @Accept  json
// @Produce  json
// @Param        id   path      int  true  "Customer ID"
// @Success 200 {object} entities.Customer
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /{id} [get]
func (h *handler) GetCustomerBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid customer id param")
		return
	}
	customer, err := h.userBalance.GetCustomerBalance(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, customer)
}

// @Summary Post Customer balance
// @Tags customer
// @Description post by INT id
// @Accept  json
// @Produce  json
// @Param        id   path      int  true  "Customer ID"
// @Param        val   path      string  true  "Value"
// @Success 200 {string} string "Status"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /{id}/{val} [post]
func (h *handler) PostCustomerBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid customer id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil || checkNegativeDecimal(value) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid customer value param")
		return
	}
	err = h.userBalance.PostCustomerBalance(id, value)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "ok",
	})
}

// @Summary Post Reserving balance
// @Tags customer
// @Description post by INT id, id_service, id_order and Decimal value
// @Accept  json
// @Produce  json
// @Param        id   path      int  true  "Customer ID"
// @Param        id_ser   path      int  true  "Service ID"
// @Param        id_ord   path      int  true  "Order ID"
// @Param        val   path      string  true  "Value"
// @Success 200 {string} string "Status"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /reserv/{id}/{id_ser}/{id_ord}/{val} [post]
func (h *handler) PostReserveCustomerBalance(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid customer id param")
		return
	}
	serviceId, err := strconv.Atoi(c.Param("id_ser"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid service id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("id_ord"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid store id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil || checkNegativeDecimal(value) {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid value param")
		return
	}
	err = h.userBalance.PostReserveBalance(customerId, serviceId, orderId, value)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "ok",
	})
}

// @Summary Post Dereserving balance ACCEPT
// @Tags customer
// @Description post by INT id, id_service, id_order and Decimal value
// @Accept  json
// @Produce  json
// @Param        id   path      int  true  "Customer ID"
// @Param        id_ser   path      int  true  "Service ID"
// @Param        id_ord   path      int  true  "Order ID"
// @Param        val   path      string  true  "Value"
// @Success 200 {string} string "Status"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /accept/{id}/{id_ser}/{id_ord}/{val} [post]
func (h *handler) PostDeReservingBalanceAccept(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid customer id param")
		return
	}
	serviceId, err := strconv.Atoi(c.Param("id_ser"))
	if err != nil || checkIsBalanceServer(serviceId) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid service id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("id_ord"))
	if err != nil || checkIsBalanceServer(orderId) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid store id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil || checkNegativeDecimal(value) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid value param")
		return
	}
	err = h.userBalance.PostDeReservingBalance(customerId, serviceId, orderId, value, true)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "ok",
	})
}

// @Summary Post Dereserving balance REJECT
// @Tags customer
// @Description post by INT id, id_service, id_order and Decimal value
// @Accept  json
// @Produce  json
// @Param        id   path      int  true  "Customer ID"
// @Param        id_ser   path      int  true  "Service ID"
// @Param        id_ord   path      int  true  "Order ID"
// @Param        val   path      string  true  "Value"
// @Success 200 {string} string "Status"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /reject/{id}/{id_ser}/{id_ord}/{val} [post]
func (h *handler) PostDeReservingBalanceReject(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid customer id param")
		return
	}
	serviceId, err := strconv.Atoi(c.Param("id_ser"))
	if err != nil || checkIsBalanceServer(serviceId) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid service id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("id_ord"))
	if err != nil || checkIsBalanceServer(orderId) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid store id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil || checkNegativeDecimal(value) {
		NewErrorResponse(c, http.StatusBadRequest, "invalid value param")
		return
	}
	err = h.userBalance.PostDeReservingBalance(customerId, serviceId, orderId, value, false)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "ok",
	})
}

// @Summary Get History report
// @Tags accounting
// @Description get by DATE (YYYY-MM)
// @Accept  json
// @Produce  json
// @Param        date   path      string  true  "Date"
// @Success 200 {string} string "Filename"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /report/{date} [get]
func (h *handler) GetHistoryReport(c *gin.Context) {
	date, err := time.Parse(config.DateFormat, c.Param("date"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	filename, err := h.userBalance.GetHistoryReport(date)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Filename": filename,
	})
}

// @Summary Get Customer report
// @Tags customer
// @Description get INT by ID and DATE (YYYY-MM)
// @Accept  json
// @Produce  json
// @Param        id   path      int  true  "Customer ID"
// @Param        date   path      string  true  "Date"
// @Success 200 {object} []entities.CustomerReport
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /history/{id}/{date} [get]
func (h *handler) GetCustomerReport(c *gin.Context) {
	date, err := time.Parse(config.DateFormat, c.Param("date"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer id param")
		return
	}
	report, err := h.userBalance.GetCustomerReport(id, date)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if report == nil {
		empty := fmt.Sprintf("user id:%d don't have customer history in %s", id, date.String())
		c.String(http.StatusOK, empty)
		return
	}
	c.JSON(http.StatusOK, report)
}
