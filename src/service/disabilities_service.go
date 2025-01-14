package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
	"fmt"
)

type DisabilityService interface {
	CreateDisability(disability []model.DisabilityRequest) utils.Error
}

type disabilityService struct {
	disabilityRepo repo.DisabilityRepo
}

func NewDisabilityService(disabilityRepo repo.DisabilityRepo) DisabilityService {
	return &disabilityService{
		disabilityRepo: disabilityRepo,
	}
}

func (s *disabilityService) CreateDisability(disabilities []model.DisabilityRequest) utils.Error {
	disabilitiesToInsert := []*model.Disability{}

	for _, disability := range disabilities {
		disabilityModel := disability.ToModel()

		disabilitiesToInsert = append(disabilitiesToInsert, &disabilityModel)
	}

	err := s.disabilityRepo.BatchInsertDisabilities(disabilitiesToInsert)
	if err.Code != "" {
		fmt.Println("Error", err.Message)
		return err
	}

	return utils.Error{}
}
