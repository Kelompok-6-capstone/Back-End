package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

func UserChatbotRoutes(e *echo.Group, chatbotController *controller.ChatbotController) {
	e.POST("/chatbot", chatbotController.GetChatResponse)
}
