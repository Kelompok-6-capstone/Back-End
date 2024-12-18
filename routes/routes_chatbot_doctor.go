// routes/doctor_chatbot_routes.go
package routes

import (
	controller "calmind/controller/chatbot_ai_doctor"

	"github.com/labstack/echo/v4"
)

func DoctorChatbotRoutes(e *echo.Group, chatbotController *controller.DoctorChatbotController) {
	e.POST("/chatbot", chatbotController.GetDoctorRecommendation)
}
