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

func NewStatsController(usecase usecase.StatsUsecase) *StatsControllerImpl {
	return &StatsControllerImpl{Usecase: usecase}
}

func (sc *StatsControllerImpl) GetStats(ctx echo.Context) error {
	totalUsers, totalDoctors, totalkonsultasi, err := sc.Usecase.GetStats()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "err :"+err.Error())
	}
	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"totalUsers":      totalUsers,
		"totalDoctors":    totalDoctors,
		"totalKonsultasi": totalkonsultasi,
	})
}
