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
