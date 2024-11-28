package repo

import (
	"cij_api/src/enum"
	model "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type VacancyApplyRepo interface {
	repo.BaseRepoMethods

	CreateVacancyApply(createVacancyApply model.VacancyApply) (int, utils.Error)
	GetVacancyApply(vacancyId int, candidateId int) (model.VacancyApply, utils.Error)
	ListVacancyAppliesByVacancyId(vacancyId int) ([]model.VacancyApply, utils.Error)
	UpdateVacancyApplyStatus(vacancyApplyId int, status enum.VacancyApplyStatus) utils.Error
}

type vacancyApplyRepo struct {
	repo.BaseRepo
	db *gorm.DB
}

func NewVacancyApplyRepo(db *gorm.DB) VacancyApplyRepo {
	repo := &vacancyApplyRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func vacancyApplyRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (v *vacancyApplyRepo) CreateVacancyApply(createVacancyApply model.VacancyApply) (int, utils.Error) {
	if err := v.db.Create(&createVacancyApply).Error; err != nil {
		return 0, vacancyApplyRepoError("failed to create the vacancy apply", "01")
	}

	return createVacancyApply.Id, utils.Error{}
}

func (v *vacancyApplyRepo) GetVacancyApply(vacancyId int, candidateId int) (model.VacancyApply, utils.Error) {
	var vacancyApply model.VacancyApply

	if err := v.db.Where("vacancy_id = ? AND candidate_id = ?", vacancyId, candidateId).Preload("Vacancy").Preload("Candidate").First(&vacancyApply).Error; err != nil {
		return model.VacancyApply{}, vacancyApplyRepoError("failed to get the vacancy apply", "02")
	}

	return vacancyApply, utils.Error{}
}

func (v *vacancyApplyRepo) ListVacancyAppliesByVacancyId(vacancyId int) ([]model.VacancyApply, utils.Error) {
	var vacancyApplies []model.VacancyApply

	if err := v.db.Where("vacancy_id = ?", vacancyId).Preload("Vacancy").Preload("Candidate").Find(&vacancyApplies).Error; err != nil {
		return []model.VacancyApply{}, vacancyApplyRepoError("failed to list the vacancy applies", "02")
	}

	return vacancyApplies, utils.Error{}
}

func (v *vacancyApplyRepo) UpdateVacancyApplyStatus(vacancyApplyId int, status enum.VacancyApplyStatus) utils.Error {
	if err := v.db.Model(model.VacancyApply{}).Where("id = ?", vacancyApplyId).Update("status", status).Error; err != nil {
		return vacancyApplyRepoError("failed to update the vacancy apply status", "03")
	}

	return utils.Error{}
}
