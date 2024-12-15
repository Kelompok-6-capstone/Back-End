package routes

import (
	controller "calmind/controller/customer_service"

	"github.com/labstack/echo/v4"
)

func UserCustServiceRoutes(e *echo.Echo, custServiceController *controller.CustServiceController) {
	e.POST("/customer-service", custServiceController.GetResponse)
	e.GET("/customer-service", custServiceController.GetQuestion)
}
