package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type NewsRepo interface {
	BaseRepoMethods

	ListNews() ([]model.News, utils.Error)
	CreateNews(news model.News) utils.Error
}

type newsRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewNewsRepo(db *gorm.DB) NewsRepo {
	repo := &newsRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
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

func (r *newsRepo) CreateNews(news model.News) utils.Error {
	err := r.db.Model(model.News{}).Create(&news).Error
	if err != nil {
		return newsRepoError("failed to create the news", "02")
	}

	return utils.Error{}
}
