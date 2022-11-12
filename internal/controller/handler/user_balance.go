package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

const (
	dateFormat = "2006-01"
)

func (h *handler) GetCustomerBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer id param")
		return
	}
	customer, err := h.userBalance.GetCustomerBalance(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *handler) PostCustomerBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer value param")
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

func (h *handler) PostReserveCustomerBalance(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer id param")
		return
	}
	serviceId, err := strconv.Atoi(c.Param("id_ser"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid service id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("id_ord"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid store id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil {
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

func (h *handler) PostDeReservingBalanceAccept(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer id param")
		return
	}
	serviceId, err := strconv.Atoi(c.Param("id_ser"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid service id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("id_ord"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid store id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid value param")
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

func (h *handler) PostDeReservingBalanceReject(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid customer id param")
		return
	}
	serviceId, err := strconv.Atoi(c.Param("id_ser"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid service id param")
		return
	}
	orderId, err := strconv.Atoi(c.Param("id_ord"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid store id param")
		return
	}
	value, err := decimal.NewFromString(c.Param("val"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid value param")
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

func (h *handler) GetHistoryReport(c *gin.Context) {
	date, err := time.Parse(dateFormat, c.Param("date"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	filename, err := h.userBalance.GetHistoryReport(date)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, filename)
}

func (h *handler) GetCustomerReport(c *gin.Context) {
	date, err := time.Parse(dateFormat, c.Param("date"))
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
	c.JSON(http.StatusOK, report)
}
