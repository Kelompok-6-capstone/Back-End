package helper

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func JSONSuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}
func GetIDParam(ctx echo.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return 0, err
	}
	return id, nil
}
