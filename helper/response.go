package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func JSONSuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}
