package controller

import (
	"calmind/helper"
	"calmind/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	NoHp     string `json:"no_hp"`
}

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

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			NoHp:     user.NoHp,
		})
	}

	return helper.JSONSuccessResponse(c, userResponses)
}

func (ac *AdminManagementController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	_, err = ac.AdminUsecase.DeleteUsers(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
	}

	return helper.JSONSuccessResponse(c, "User data successfully deleted")
}

func (ac *AdminManagementController) GetAllDocter(c echo.Context) error {
	doctor, err := ac.AdminUsecase.GetAllDocter()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch doctor")
	}

	var userResponses []UserResponse
	for _, user := range doctor {
		userResponses = append(userResponses, UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			NoHp:     user.NoHp,
		})
	}

	return helper.JSONSuccessResponse(c, userResponses)
}

func (ac *AdminManagementController) DeleteDocter(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	_, err = ac.AdminUsecase.DeleteDocter(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to delete dokter")
	}

	return helper.JSONSuccessResponse(c, "Docter data successfully deleted")
}
