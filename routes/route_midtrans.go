package routes

import (
	controller_notifikasi "calmind/controller/midtrans_notifikasi"

	"github.com/labstack/echo/v4"
)

func WebhookRoutes(e *echo.Echo, notifikasi *controller_notifikasi.MidtransNotificationController) {
	e.POST("/notifications/midtrans", notifikasi.MidtransNotification)
}
