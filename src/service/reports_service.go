package service

import (
	"cij_api/src/enum"
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
	"time"
)

type ReportsService interface {
	GetDisabilityTotals() (model.DisabilityTotals, utils.Error)
	GetDisabilityTotalsByNeighborhood(neighborhood string) (model.DisabilityTotalsByNeighborhood, utils.Error)
	CountActivitiesByPeriod(activityType string, period enum.PeriodFilterEnum) (model.CountActivitiesByPeriod, utils.Error)
}

type reportsService struct {
	personDisabilityRepo repo.PersonDisabilityRepo
	activityRepo         repo.ActivityRepo
}

func NewReportsService(personDisabilityRepo repo.PersonDisabilityRepo, activityRepo repo.ActivityRepo) ReportsService {
	return &reportsService{
		personDisabilityRepo: personDisabilityRepo,
		activityRepo:         activityRepo,
	}
}

func (s *reportsService) GetDisabilityTotals() (model.DisabilityTotals, utils.Error) {
	disabilityTotals, err := s.personDisabilityRepo.CountDisability()
	if err.Code != "" {
		return model.DisabilityTotals{}, err
	}

	return disabilityTotals, utils.Error{}
}

func (s *reportsService) GetDisabilityTotalsByNeighborhood(neighborhood string) (model.DisabilityTotalsByNeighborhood, utils.Error) {
	disabilityTotals, err := s.personDisabilityRepo.CountDisabilityByNeighborhood(neighborhood)
	if err.Code != "" {
		return model.DisabilityTotalsByNeighborhood{}, err
	}

	return disabilityTotals, utils.Error{}
}

func (s *reportsService) CountActivitiesByPeriod(activityType string, period enum.PeriodFilterEnum) (model.CountActivitiesByPeriod, utils.Error) {
	var countActivitiesByPeriod model.CountActivitiesByPeriod

	startDate := time.Now().AddDate(0, 0, -utils.PeriodToDays(period))
	endDate := time.Now()

	activities, err := s.activityRepo.GetActivitiesByTypeAndPeriod(activityType, utils.GetFormattedDate(startDate.Unix()), utils.GetFormattedDate(endDate.Unix()))
	if err.Code != "" {
		return model.CountActivitiesByPeriod{}, err
	}

	activitiesByMonth := activitiesByMonthInitial(startDate, endDate)

	for _, activity := range activities {
		month := activity.CreatedAt.Format("2006-01")
		activitiesByMonth[month]++
	}

	countActivitiesByPeriod = model.CountActivitiesByPeriod{
		ActivityType: activityType,
		MonthsCount:  activitiesByMonth,
	}

	return countActivitiesByPeriod, utils.Error{}
}

func activitiesByMonthInitial(startDate time.Time, endDate time.Time) map[string]int {
	activitiesByMonth := make(map[string]int)

	for startDate.Before(endDate) {
		month := startDate.Format("2006-01")
		activitiesByMonth[month] = 0

		startDate = startDate.AddDate(0, 1, 0)
	}

	return activitiesByMonth
}
