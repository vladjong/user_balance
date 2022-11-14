package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vladjong/user_balance/internal/usecase"
)

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

	//swager
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
