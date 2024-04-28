package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type NewsRepo interface {
	ListNews() ([]model.News, utils.Error)
}

type newsRepo struct {
	db *gorm.DB
}

func NewNewsRepo(db *gorm.DB) NewsRepo {
	return &newsRepo{
		db: db,
	}
}

func newsRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.NewsErrorType, code)

	return utils.NewError(message, errorCode)
}

func (r *newsRepo) ListNews() ([]model.News, utils.Error) {
	var news []model.News

	err := r.db.Model(model.News{}).Find(&news).Error
	if err != nil {
		return news, newsRepoError("failed to list the news", "01")
	}

	return news, utils.Error{}
}
