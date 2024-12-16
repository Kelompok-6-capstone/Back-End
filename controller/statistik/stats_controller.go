package controller

import (
	"calmind/helper"
	usecase "calmind/usecase/statistik"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StatsControllerImpl struct {
	Usecase usecase.StatsUsecase
}

// DTO untuk statistik
type StatsResponse struct {
	TotalUsers         int64 `json:"totalUsers"`
	TotalDoctors       int64 `json:"totalDoctors"`
	TotalConsultations int64 `json:"totalConsultations"`
	TotalPaid          int64 `json:"totalPaid"`
	TotalPending       int64 `json:"totalPending"`
}

func NewStatsController(usecase usecase.StatsUsecase) *StatsControllerImpl {
	return &StatsControllerImpl{Usecase: usecase}
}

func (sc *StatsControllerImpl) GetStats(ctx echo.Context) error {
	totalUsers, totalDoctors, totalConsultations, totalPaid, totalPending, err := sc.Usecase.GetStats()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch statistics: "+err.Error())
	}

	stats := StatsResponse{
		TotalUsers:         totalUsers,
		TotalDoctors:       totalDoctors,
		TotalConsultations: totalConsultations,
		TotalPaid:          totalPaid,
		TotalPending:       totalPending,
	}

	return helper.JSONSuccessResponse(ctx, stats)
}
