package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
)

type ActivityService interface {
	CreateActivity(activity *model.Activity) utils.Error
	GetActivitiesByTypeAndPeriod(activityType string, startDate int64, endDate int64) ([]model.ActivityResponse, utils.Error)
}

type activityService struct {
	activityRepo repo.ActivityRepo
}

func NewActivityService(activityRepo repo.ActivityRepo) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
	}
}

func (a *activityService) CreateActivity(activity *model.Activity) utils.Error {
	return a.activityRepo.CreateActivity(activity)
}

func (a *activityService) GetActivitiesByTypeAndPeriod(activityType string, startDate int64, endDate int64) ([]model.ActivityResponse, utils.Error) {
	startDateStr := utils.GetFormattedDate(startDate)
	endDateStr := utils.GetFormattedDate(endDate)

	activities, err := a.activityRepo.GetActivitiesByTypeAndPeriod(activityType, startDateStr, endDateStr)
	if err.Code != "" {
		return nil, err
	}

	activitiesResponse := []model.ActivityResponse{}

	for _, activity := range activities {
		activitiesResponse = append(activitiesResponse, *activity.ToResponse())
	}

	return activitiesResponse, utils.Error{}
}
