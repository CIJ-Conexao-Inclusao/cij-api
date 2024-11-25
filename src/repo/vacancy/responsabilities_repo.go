package repo

import (
	model "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type ResponsabilitiesRepo interface {
	repo.BaseRepoMethods

	CreateResponsability(createResponsability model.VacancyResponsability, tx *gorm.DB) (int, utils.Error)
	ListResponsabilitiesByVacancyId(vacancyId int) ([]model.VacancyResponsability, utils.Error)
	UpdateResponsability(responsability model.VacancyResponsability, responsabilityId int, tx *gorm.DB) utils.Error
	DeleteResponsability(responsabilityId int) utils.Error
}

type responsabilitiesRepo struct {
	repo.BaseRepo
	db *gorm.DB
}

func NewResponsabilitiesRepo(db *gorm.DB) ResponsabilitiesRepo {
	repo := &responsabilitiesRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func responsabilitiesRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (r *responsabilitiesRepo) CreateResponsability(createResponsability model.VacancyResponsability, tx *gorm.DB) (int, utils.Error) {
	databaseConn := r.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&createResponsability).Error; err != nil {
		return 0, responsabilitiesRepoError("failed to create the responsability", "01")
	}

	return createResponsability.Id, utils.Error{}
}

func (r *responsabilitiesRepo) ListResponsabilitiesByVacancyId(vacancyId int) ([]model.VacancyResponsability, utils.Error) {
	var responsabilities []model.VacancyResponsability

	if err := r.db.Where("vacancy_id = ?", vacancyId).Find(&responsabilities).Error; err != nil {
		return []model.VacancyResponsability{}, responsabilitiesRepoError("failed to list the responsabilities", "02")
	}

	return responsabilities, utils.Error{}
}

func (r *responsabilitiesRepo) UpdateResponsability(responsability model.VacancyResponsability, responsabilityId int, tx *gorm.DB) utils.Error {
	databaseConn := r.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Model(model.VacancyResponsability{}).Where("id = ?", responsabilityId).Updates(responsability).Error; err != nil {
		return responsabilitiesRepoError("failed to update the responsability", "03")
	}

	return utils.Error{}
}

func (r *responsabilitiesRepo) DeleteResponsability(responsabilityId int) utils.Error {
	if err := r.db.Where("id = ?", responsabilityId).Delete(&model.VacancyResponsability{}).Error; err != nil {
		return responsabilitiesRepoError("failed to delete the responsability", "04")
	}

	return utils.Error{}
}
