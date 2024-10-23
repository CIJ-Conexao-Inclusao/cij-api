package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type ActivityRepo interface {
	BaseRepoMethods

	CreateActivity(activity *model.Activity) utils.Error
	GetActivitiesByTypeAndPeriod(activityType string, startDate string, endDate string) ([]model.Activity, utils.Error)
}

type activityRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewActivityRepo(db *gorm.DB) ActivityRepo {
	repo := &activityRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func activityRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.ActivityErrorType, code)

	return utils.NewError(message, errorCode)
}

func (a *activityRepo) CreateActivity(activity *model.Activity) utils.Error {
	if err := a.db.Create(activity).Error; err != nil {
		return activityRepoError("failed to create the activity", "01")
	}

	return utils.Error{}
}

func (a *activityRepo) GetActivitiesByTypeAndPeriod(activityType string, startDate string, endDate string) ([]model.Activity, utils.Error) {
	var activities []model.Activity

	if err := a.db.Where("type = ? AND created_at >= ? AND created_at <= ?", activityType, startDate, endDate).Find(&activities).Error; err != nil {
		return nil, activityRepoError("failed to get the activities", "02")
	}

	return activities, utils.Error{}
}
