package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

func UserCustServiceRoutes(e *echo.Group, custServiceController *controller.CustServiceController) {
	e.POST("/csrespons", custServiceController.GetResponse)
	e.GET("/csquestion", custServiceController.GetQuestion)

}

