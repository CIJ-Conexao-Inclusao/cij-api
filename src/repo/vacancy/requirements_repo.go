package repo

import (
	model "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type RequirementsRepo interface {
	repo.BaseRepoMethods

	CreateRequirement(createRequirement model.VacancyRequirement, tx *gorm.DB) (int, utils.Error)
	ListRequirementsByVacancyId(vacancyId int) ([]model.VacancyRequirement, utils.Error)
	UpdateRequirement(requirement model.VacancyRequirement, requirementId int, tx *gorm.DB) utils.Error
	DeleteRequirement(requirementId int) utils.Error
}

type requirementsRepo struct {
	repo.BaseRepo
	db *gorm.DB
}

func NewRequirementsRepo(db *gorm.DB) RequirementsRepo {
	repo := &requirementsRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func requirementsRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (r *requirementsRepo) CreateRequirement(createRequirement model.VacancyRequirement, tx *gorm.DB) (int, utils.Error) {
	databaseConn := r.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&createRequirement).Error; err != nil {
		return 0, requirementsRepoError("failed to create the requirement", "01")
	}

	return createRequirement.Id, utils.Error{}
}

func (r *requirementsRepo) ListRequirementsByVacancyId(vacancyId int) ([]model.VacancyRequirement, utils.Error) {
	var requirements []model.VacancyRequirement

	if err := r.db.Where("vacancy_id = ?", vacancyId).Find(&requirements).Error; err != nil {
		return []model.VacancyRequirement{}, requirementsRepoError("failed to list the requirements", "02")
	}

	return requirements, utils.Error{}
}

func (r *requirementsRepo) UpdateRequirement(requirement model.VacancyRequirement, requirementId int, tx *gorm.DB) utils.Error {
	databaseConn := r.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Model(model.VacancyRequirement{}).Where("id = ?", requirementId).Updates(requirement).Error; err != nil {
		return requirementsRepoError("failed to update the requirement", "03")
	}

	return utils.Error{}
}

func (r *requirementsRepo) DeleteRequirement(requirementId int) utils.Error {
	if err := r.db.Where("id = ?", requirementId).Delete(&model.VacancyRequirement{}).Error; err != nil {
		return requirementsRepoError("failed to delete the requirement", "04")
	}

	return utils.Error{}
}
