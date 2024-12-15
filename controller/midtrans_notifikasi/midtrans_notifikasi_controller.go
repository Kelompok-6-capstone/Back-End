package controller

import (
	usecase "calmind/usecase/konsultasi"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MidtransNotificationController struct {
	ConsultationUsecase *usecase.ConsultationUsecaseImpl
}

func NewMidtransNotificationController(konsultasi *usecase.ConsultationUsecaseImpl) *MidtransNotificationController {
	return &MidtransNotificationController{ConsultationUsecase: konsultasi}
}

func (c *MidtransNotificationController) MidtransNotification(ctx echo.Context) error {
	var notification map[string]interface{}
	if err := ctx.Bind(&notification); err != nil {
		log.Printf("Failed to bind notification: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid payload"})
	}

	// Log payload yang diterima
	log.Printf("Received Midtrans Notification: %+v\n", notification)

	transactionStatus, ok := notification["transaction_status"].(string)
	orderID, ok := notification["order_id"].(string)

	if !ok || transactionStatus == "" || orderID == "" {
		log.Println("Invalid notification data")
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid notification data"})
	}

	// Update payment status di database
	err := c.ConsultationUsecase.UpdatePaymentStatus(orderID, transactionStatus)
	if err != nil {
		log.Printf("Failed to update payment status for order_id=%s: %v", orderID, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update payment status"})
	}

	log.Printf("Payment status updated successfully for order_id=%s", orderID)
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Notification processed successfully"})
}
