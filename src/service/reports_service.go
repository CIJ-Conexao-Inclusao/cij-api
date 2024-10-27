package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
)

type ReportsService interface {
	GetDisabilityTotals() (model.DisabilityTotals, utils.Error)
	GetDisabilityTotalsByNeighborhood(neighborhood string) (model.DisabilityTotalsByNeighborhood, utils.Error)
}

type reportsService struct {
	personDisabilityRepo repo.PersonDisabilityRepo
}

func NewReportsService(personDisabilityRepo repo.PersonDisabilityRepo) ReportsService {
	return &reportsService{
		personDisabilityRepo: personDisabilityRepo,
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
