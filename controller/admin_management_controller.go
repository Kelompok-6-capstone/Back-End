package controller

import (
	"calmind/helper"
	"calmind/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Tgl_lahir string `json:"tgl_lahir"`
	NoHp      string `json:"no_hp"`
	Pekerjaan string `json:"pekerjaan"`
	Alamat    string `json:"alamat"`
	Gender    string `json:"gender"`
	is_active bool   `json:"is_active"`
}

type AdminManagementController struct {
	AdminUsecase usecase.AdminManagementUsecase
}

func NewAdminManagementController(usecase usecase.AdminManagementUsecase) *AdminManagementController {
	return &AdminManagementController{AdminUsecase: usecase}
}

// user
func (ac *AdminManagementController) GetAllUsers(c echo.Context) error {
	users, err := ac.AdminUsecase.GetAllUsers()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Tgl_lahir: user.Tgl_lahir,
			Pekerjaan: user.Pekerjaan,
			NoHp:      user.NoHp,
			Alamat:    user.Alamat,
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

// dokter
func (ac *AdminManagementController) GetAllDocter(c echo.Context) error {
	doctor, err := ac.AdminUsecase.GetAllDocter()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch doctor")
	}

	var userResponses []UserResponse
	for _, user := range doctor {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Tgl_lahir: user.DateOfBirth,
			NoHp:      user.NoHp,
			Alamat:    user.Address,
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

func (ac *AdminManagementController) GetUserDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	user, err := ac.AdminUsecase.GetUserDetail(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user detail")
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Gender:    user.JenisKelamin,
		Email:     user.Email,
		Tgl_lahir: user.Tgl_lahir,
		Pekerjaan: user.Pekerjaan,
		NoHp:      user.NoHp,
		Alamat:    user.Alamat,
		is_active: user.IsVerified,
	}

	return helper.JSONSuccessResponse(c, response)
}

func (ac *AdminManagementController) GetDocterDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
	}

	doctor, err := ac.AdminUsecase.GetDocterDetail(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch doctor detail")
	}

	response := UserResponse{
		ID:        doctor.ID,
		Username:  doctor.Username,
		Email:     doctor.Email,
		Tgl_lahir: doctor.DateOfBirth,
		NoHp:      doctor.NoHp,
		Alamat:    doctor.Address,
		Gender:    doctor.JenisKelamin,
	}

	return helper.JSONSuccessResponse(c, response)
}
