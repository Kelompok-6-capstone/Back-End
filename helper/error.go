package helper

import (
	"github.com/labstack/echo/v4"
)

func JSONErrorResponse(ctx echo.Context, status int, message string) error {
	return ctx.JSON(status, map[string]interface{}{
		"success": false,
		"message": message,
	})
}
