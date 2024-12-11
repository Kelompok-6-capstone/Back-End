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

// Get All Users
func (ac *AdminManagementController) GetAllUsers(c echo.Context) error {
	users, err := ac.AdminUsecase.GetAllUsers()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
	}

	return helper.JSONSuccessResponse(c, users)
}

// Get All Doctors
func (ac *AdminManagementController) GetAllDoctors(c echo.Context) error {
	doctors, err := ac.AdminUsecase.GetAllDoctors()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch doctors")
	}

	return helper.JSONSuccessResponse(c, doctors)
}

// Delete User
func (ac *AdminManagementController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	user, err := ac.AdminUsecase.DeleteUser(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
	}

	return helper.JSONSuccessResponse(c, map[string]interface{}{
		"message": "User deleted successfully",
		"user_id": user.ID,
	})
}

// Delete Doctor
func (ac *AdminManagementController) DeleteDoctor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
	}

	doctor, err := ac.AdminUsecase.DeleteDoctor(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to delete doctor")
	}

	return helper.JSONSuccessResponse(c, map[string]interface{}{
		"message":   "Doctor deleted successfully",
		"doctor_id": doctor.ID,
	})
}
