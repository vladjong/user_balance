package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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
