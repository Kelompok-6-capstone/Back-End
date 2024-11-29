package controller

import (
	"calmind/helper"
	"calmind/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminManagementController struct {
	AdminUsecase usecase.AdminManagementUsecase
}

func NewAdminManagementController(usecase usecase.AdminManagementUsecase) *AdminManagementController {
	return &AdminManagementController{AdminUsecase: usecase}
}

func (ac *AdminManagementController) GetAllUsers(c echo.Context) error {
	users, err := ac.AdminUsecase.GetAllUsers()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
	}
	return helper.JSONSuccessResponse(c, users)
}

func (ac *AdminManagementController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	user, err := ac.AdminUsecase.DeleteUsers(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
	}

	return helper.JSONSuccessResponse(c, map[string]interface{}{
		"message": "Data users berhasil dihapus",
		"data":    user,
	})
}
