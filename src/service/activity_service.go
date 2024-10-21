package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
)

type ActivityService interface {
	CreateActivity(activity *model.Activity) utils.Error
	GetActivitiesByTypeAndPeriod(activityType string, startDate int64, endDate int64) ([]model.Activity, utils.Error)
}

type activityService struct {
	activityRepo repo.ActivityRepo
}

func NewActivityService(activityRepo repo.ActivityRepo) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
	}
}

func activityServiceError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ServiceErrorCode, utils.ActivityErrorType, code)

	return utils.NewError(message, errorCode)
}

func (a *activityService) CreateActivity(activity *model.Activity) utils.Error {
	return a.activityRepo.CreateActivity(activity)
}

func (a *activityService) GetActivitiesByTypeAndPeriod(activityType string, startDate int64, endDate int64) ([]model.Activity, utils.Error) {
	startDateStr := utils.GetFormattedDate(startDate)
	endDateStr := utils.GetFormattedDate(endDate)

	return a.activityRepo.GetActivitiesByTypeAndPeriod(activityType, startDateStr, endDateStr)
}
