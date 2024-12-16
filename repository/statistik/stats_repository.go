package usecase

import repository "calmind/repository/statistik"

type StatsUsecase interface {
	GetStats() (int64, int64, int64, int64, int64, error)
}

type StatsUsecaseImpl struct {
	Repo repository.StatsRepository
}

func NewStatsUsecase(repo repository.StatsRepository) StatsUsecase {
	return &StatsUsecaseImpl{Repo: repo}
}

func (uc *StatsUsecaseImpl) GetStats() (int64, int64, int64, int64, int64, error) {
	totalUsers, err := uc.Repo.GetTotalUsers()
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	totalDoctors, err := uc.Repo.GetTotalDoctors()
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	totalConsultations, err := uc.Repo.GetTotalConsultations()
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	totalPaid, err := uc.Repo.GetTotalConsultationsByPaymentStatus("paid")
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	totalPending, err := uc.Repo.GetTotalConsultationsByPaymentStatus("pending")
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	return totalUsers, totalDoctors, totalConsultations, totalPaid, totalPending, nil
}
