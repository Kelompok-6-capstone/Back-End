package controller

import (
	"calmind/helper"
	"calmind/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	TglLahir   string `json:"tgl_lahir,omitempty"` // `omitempty` agar tidak ditampilkan jika kosong
	NoHp       string `json:"no_hp,omitempty"`
	Pekerjaan  string `json:"pekerjaan,omitempty"`
	Alamat     string `json:"alamat,omitempty"`
	Gender     string `json:"gender,omitempty"`
	IsVerified bool   `json:"is_verified"`
}

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

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			TglLahir:  user.Tgl_lahir,
			Pekerjaan: user.Pekerjaan,
			NoHp:      user.NoHp,
			Alamat:    user.Alamat,
		})
	}

	return helper.JSONSuccessResponse(c, userResponses)
}

// Get User Detail
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
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		TglLahir:   user.Tgl_lahir,
		NoHp:       user.NoHp,
		Pekerjaan:  user.Pekerjaan,
		Alamat:     user.Alamat,
		Gender:     user.JenisKelamin,
		IsVerified: user.IsVerified,
	}

	return helper.JSONSuccessResponse(c, response)
}

// Get All Doctors
func (ac *AdminManagementController) GetAllDoctors(c echo.Context) error {
	doctors, err := ac.AdminUsecase.GetAllDocter()
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch doctors")
	}

	var doctorResponses []UserResponse
	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, UserResponse{
			ID:       doctor.ID,
			Username: doctor.Username,
			Email:    doctor.Email,
			TglLahir: doctor.DateOfBirth,
			NoHp:     doctor.NoHp,
			Alamat:   doctor.Address,
		})
	}
	return helper.JSONSuccessResponse(c, doctorResponses)
}

// Get Doctor Detail
func (ac *AdminManagementController) GetDoctorDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
	}

	doctor, err := ac.AdminUsecase.GetDocterDetail(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to fetch doctor detail")
	}

	response := UserResponse{
		ID:         doctor.ID,
		Username:   doctor.Username,
		Email:      doctor.Email,
		TglLahir:   doctor.DateOfBirth,
		NoHp:       doctor.NoHp,
		Alamat:     doctor.Address,
		Gender:     doctor.JenisKelamin,
		IsVerified: doctor.IsVerified,
	}

	return helper.JSONSuccessResponse(c, response)
}

// Delete User
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

// Delete Doctor
func (ac *AdminManagementController) DeleteDoctor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
	}

	_, err = ac.AdminUsecase.DeleteDocter(id)
	if err != nil {
		return helper.JSONErrorResponse(c, http.StatusInternalServerError, "Failed to delete doctor")
	}

	return helper.JSONSuccessResponse(c, "Doctor data successfully deleted")
}
