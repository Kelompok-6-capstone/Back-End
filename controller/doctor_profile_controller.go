package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type DoctorProfileController struct {
	DoctorProfileUsecase usecase.DoctorProfileUseCase
}

func NewDoctorProfileController(DoctorProfileUsecase usecase.DoctorProfileUseCase) *DoctorProfileController {
	return &DoctorProfileController{DoctorProfileUsecase: DoctorProfileUsecase}
}

func (c *DoctorProfileController) GetProfile(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access")
	}

	doctor, err := c.DoctorProfileUsecase.GetDoctorProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil dokter: "+err.Error())
	}

	// Format response lengkap dengan Title dan Tags
	type TagsResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type TitleResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type DoctorResponse struct {
		ID           int            `json:"id"`
		Username     string         `json:"username"`
		NoHp         string         `json:"no_hp"`
		Email        string         `json:"email"`
		Avatar       string         `json:"avatar"`
		DateOfBirth  string         `json:"date_of_birth"`
		Address      string         `json:"address"`
		Schedule     string         `json:"schedule"`
		IsVerified   bool           `json:"is_verified"`
		IsActive     bool           `json:"is_active"`
		Price        float64        `json:"price"`
		Experience   int            `json:"experience"`
		STRNumber    string         `json:"str_number"`
		About        string         `json:"about"`
		JenisKelamin string         `json:"jenis_kelamin"`
		Title        TitleResponse  `json:"title"`
		Tags         []TagsResponse `json:"tags"`
		CreatedAt    string         `json:"created_at"`
		UpdatedAt    string         `json:"updated_at"`
	}

	var tagsResponse []TagsResponse
	for _, tag := range doctor.Tags {
		tagsResponse = append(tagsResponse, TagsResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	doctorProfile := DoctorResponse{
		ID:           doctor.ID,
		Username:     doctor.Username,
		NoHp:         doctor.NoHp,
		Email:        doctor.Email,
		Avatar:       doctor.Avatar,
		DateOfBirth:  doctor.DateOfBirth,
		Address:      doctor.Address,
		Schedule:     doctor.Schedule,
		IsVerified:   doctor.IsVerified,
		IsActive:     doctor.IsActive,
		Price:        doctor.Price,
		Experience:   doctor.Experience,
		STRNumber:    doctor.STRNumber,
		About:        doctor.About,
		JenisKelamin: doctor.JenisKelamin,
		Title: TitleResponse{
			ID:   doctor.Title.ID,
			Name: doctor.Title.Name,
		},
		Tags:      tagsResponse,
		CreatedAt: doctor.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt: doctor.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	return helper.JSONSuccessResponse(ctx, doctorProfile)
}
func (c *DoctorProfileController) UpdateProfile(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access")
	}

	// Struct untuk menerima request body
	var request struct {
		Username     string       `json:"username"`
		NoHp         string       `json:"no_hp"`
		Avatar       string       `json:"avatar"`
		DateOfBirth  string       `json:"date_of_birth"`
		Address      string       `json:"address"`
		Schedule     string       `json:"schedule"`
		Title        string       `json:"title"`
		Experience   int          `json:"experience"`
		STRNumber    string       `json:"str_number"`
		About        string       `json:"about"`
		JenisKelamin string       `json:"jenis_kelamin"`
		Tags         []model.Tags `json:"tags"`
	}

	// Bind JSON ke struct
	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal memproses data: "+err.Error())
	}

	// Validasi nilai JenisKelamin
	if request.JenisKelamin != "" && request.JenisKelamin != "Laki-laki" && request.JenisKelamin != "Perempuan" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Nilai JenisKelamin tidak valid")
	}

	// Map request ke model Doctor
	doctor := &model.Doctor{
		Username:     request.Username,
		NoHp:         request.NoHp,
		Avatar:       request.Avatar,
		DateOfBirth:  request.DateOfBirth,
		Address:      request.Address,
		Schedule:     request.Schedule,
		Title:        model.Title{Name: request.Title},
		Experience:   request.Experience,
		STRNumber:    request.STRNumber,
		About:        request.About,
		JenisKelamin: request.JenisKelamin,
		Tags:         request.Tags,
	}

	// Panggil Usecase untuk update profile
	updatedDoctor, err := c.DoctorProfileUsecase.UpdateDoctorProfile(claims.UserID, doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate profil dokter: "+err.Error())
	}

	// Berikan response sukses
	return helper.JSONSuccessResponse(ctx, updatedDoctor)
}

// SetActiveStatus updates the active status of the logged-in doctor
func (c *DoctorProfileController) SetActiveStatus(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access")
	}

	var statusRequest struct {
		IsActive bool `json:"is_active"`
	}

	if err := ctx.Bind(&statusRequest); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal memproses data: "+err.Error())
	}

	err := c.DoctorProfileUsecase.SetDoctorActiveStatus(claims.UserID, statusRequest.IsActive)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status aktif dokter: "+err.Error())
	}

	message := "Anda telah mengubah status ke tidak aktif"
	if statusRequest.IsActive {
		message = "Anda telah mengubah status ke aktif"
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message": message,
	})
}

func (c *DoctorProfileController) UploadAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	doctorID := claims.UserID

	// Ambil file dari form
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan file: "+err.Error())
	}

	// Validasi ukuran file (maksimal 5 MB)
	if file.Size > 5*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Ukuran file maksimal 5 MB")
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Hanya file dengan format .jpg, .jpeg, atau .png yang diperbolehkan")
	}

	// Simpan file di direktori uploads
	uploadDir := "uploads/avatars"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload")
		}
	}

	filePath := fmt.Sprintf("%s/doctor_%d_%s", uploadDir, doctorID, file.Filename)
	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file")
	}
	defer dst.Close()

	if _, err := helper.CopyFile(src, dst); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file")
	}

	// Update URL avatar di database
	avatarURL := fmt.Sprintf("http://%s/%s", ctx.Request().Host, filePath)
	doctor := model.Doctor{
		Avatar: avatarURL,
	}
	_, err = c.DoctorProfileUsecase.UpdateDoctorProfile(doctorID, &doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": avatarURL,
	})
}

func (c *DoctorProfileController) DeleteAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	doctorID := claims.UserID

	// Ambil data dokter untuk mendapatkan URL avatar
	doctor, err := c.DoctorProfileUsecase.GetDoctorProfile(doctorID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil dokter: "+err.Error())
	}

	// Hapus file avatar
	if doctor.Avatar != "" {
		filePath := "." + doctor.Avatar // Tambahkan "." untuk path relatif
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus file avatar: "+err.Error())
			}
		}
	}

	// Update avatar menjadi kosong di database
	doctor.Avatar = ""
	_, err = c.DoctorProfileUsecase.UpdateDoctorProfile(doctorID, doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar di database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
