package helper

import (
	"regexp"

	"github.com/labstack/echo/v4"
)

func JSONErrorResponse(ctx echo.Context, status int, message string) error {
	return ctx.JSON(status, map[string]interface{}{
		"success": false,
		"message": message,
	})
}

func IsValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{5,}$`)
	return re.MatchString(username)
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}
