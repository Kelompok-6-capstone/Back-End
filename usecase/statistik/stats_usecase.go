package usecase

import repository "calmind/repository/statistik"

type StatsUsecase interface {
	GetStats() (int64, int64, int64, error)
}

type StatsUsecaseImpl struct {
	Repo repository.StatsRepo
}

func NewStatsUsecase(repo repository.StatsRepo) StatsUsecase {
	return &StatsUsecaseImpl{Repo: repo}
}

// GetStats method to retrieve total visitors, users, and doctors
func (s *StatsUsecaseImpl) GetStats() (int64, int64, int64, error) {

	totalUsers, err := s.Repo.GetTotalUsers()
	if err != nil {
		return 0, 0, 0, err
	}

	totalDoctors, err := s.Repo.GetTotalDoctors()
	if err != nil {
		return 0, 0, 0, err
	}
	totalKonsultasi, err := s.Repo.GetTotalKonsultasi()
	if err != nil {
		return 0, 0, 0, err
	}

	return totalKonsultasi, totalUsers, totalDoctors, nil
}
