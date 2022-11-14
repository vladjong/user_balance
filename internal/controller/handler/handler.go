package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/vladjong/user_balance/docs"
	"github.com/vladjong/user_balance/internal/usecase"
)

// @title UserBalance API
// @version 1.0
// @description This is a service that work with the user's balance

// @host      localhost:8080
// @BasePath  /api

type handler struct {
	userBalance usecase.UserBalanse
}

func New(userBalance usecase.UserBalanse) *handler {
	return &handler{
		userBalance: userBalance,
	}
}

func (h *handler) NewRouter() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.GET("/:id", h.GetCustomerBalance)
		api.GET("/report/:date", h.GetHistoryReport)
		api.GET("/history/:id/:date", h.GetCustomerReport)
		api.POST("/:id/:val", h.PostCustomerBalance)
		api.POST("/reserv/:id/:id_ser/:id_ord/:val", h.PostReserveCustomerBalance)
		api.POST("/accept/:id/:id_ser/:id_ord/:val", h.PostDeReservingBalanceAccept)
		api.POST("/reject/:id/:id_ser/:id_ord/:val", h.PostDeReservingBalanceReject)
	}
	return router
}
