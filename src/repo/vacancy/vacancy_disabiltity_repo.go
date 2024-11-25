package repo

import (
	model "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	"cij_api/src/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VacancyDisabilityRepo interface {
	repo.BaseRepoMethods

	GetVacancyDisabilities(vacancyId int) ([]model.VacancyDisability, utils.Error)
	UpsertVacancyDisability(disability model.VacancyDisability, tx *gorm.DB) utils.Error
	ClearVacancyDisability(vacancyId int, tx *gorm.DB) utils.Error
}

type vacancyDisabilityRepo struct {
	repo.BaseRepo
	db *gorm.DB
}

func NewVacancyDisabilityRepo(db *gorm.DB) VacancyDisabilityRepo {
	repo := &vacancyDisabilityRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func vacancyDisabilityRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (v *vacancyDisabilityRepo) GetVacancyDisabilities(vacancyId int) ([]model.VacancyDisability, utils.Error) {
	var disabilities []model.VacancyDisability

	err := v.db.Model(model.VacancyDisability{}).Preload("Disability").Where("vacancy_id = ?", vacancyId).Find(&disabilities).Error
	if err != nil {
		return disabilities, vacancyDisabilityRepoError("failed to get the vacancy disabilities", "01")
	}

	return disabilities, utils.Error{}
}

func (v *vacancyDisabilityRepo) UpsertVacancyDisability(disability model.VacancyDisability, tx *gorm.DB) utils.Error {
	databaseConn := v.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "vacancy_id"}, {Name: "disability_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"disability_id", "vacancy_id"}),
	}).Create(&disability).Error; err != nil {
		return vacancyDisabilityRepoError("failed to upsert the vacancy disability", "02")
	}

	return utils.Error{}
}

func (v *vacancyDisabilityRepo) ClearVacancyDisability(vacancyId int, tx *gorm.DB) utils.Error {
	databaseConn := v.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Where("vacancy_id = ?", vacancyId).Delete(&model.VacancyDisability{}).Error; err != nil {
		return vacancyDisabilityRepoError("failed to clear the vacancy disability", "03")
	}

	return utils.Error{}
}
