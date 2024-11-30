package repo

import (
	"cij_api/src/enum"
	model "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type VacancyRepo interface {
	repo.BaseRepoMethods

	GetVacancyById(id int) (model.Vacancy, utils.Error)
	ListVacancies(
		page int,
		perPage int,
		companyId int,
		area string,
		contractType enum.VacancyContractType,
		searchText string,
	) ([]model.Vacancy, utils.Error)
	UpsertVacancy(vacancy model.Vacancy, tx *gorm.DB) (int, utils.Error)
	DeleteVacancy(id int) utils.Error
}

type vacancyRepo struct {
	repo.BaseRepo
	db *gorm.DB
}

func NewVacancyRepo(db *gorm.DB) VacancyRepo {
	repo := &vacancyRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func vacancyRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (v *vacancyRepo) GetVacancyById(id int) (model.Vacancy, utils.Error) {
	var vacancy model.Vacancy

	if err := v.db.Where("id = ?", id).Preload("Company").First(&vacancy).Error; err != nil {
		return model.Vacancy{}, vacancyRepoError("failed to get the vacancy", "01")
	}

	return vacancy, utils.Error{}
}

func (v *vacancyRepo) ListVacancies(
	page int,
	perPage int,
	companyId int,
	area string,
	contractType enum.VacancyContractType,
	searchText string,
) ([]model.Vacancy, utils.Error) {
	var vacancies []model.Vacancy

	query := v.db.Model(&model.Vacancy{}).
		Preload("Disabilities").
		Preload("Company")

	if area != "" {
		query = query.Where("vacancies.area = ?", area)
	}

	if companyId > 0 {
		query = query.Where("vacancies.company_id = ?", companyId)
	}

	if contractType != "" {
		query = query.Where("vacancies.contract_type = ?", contractType)
	}

	if searchText != "" {
		query = query.Where("(vacancies.code LIKE ? OR vacancies.title LIKE ?)", "%"+searchText+"%", "%"+searchText+"%")
	}

	err := query.Limit(perPage).Offset(perPage * page).Find(&vacancies).Error
	if err != nil {
		return vacancies, vacancyRepoError("failed to list the vacancies", "02")
	}

	return vacancies, utils.Error{}
}

func (v *vacancyRepo) UpsertVacancy(vacancy model.Vacancy, tx *gorm.DB) (int, utils.Error) {
	databaseConn := v.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&vacancy).Error; err != nil {
		return 0, vacancyRepoError("failed to create the vacancy", "03")
	}

	return vacancy.Id, utils.Error{}
}

func (v *vacancyRepo) DeleteVacancy(id int) utils.Error {
	if err := v.db.Where("id = ?", id).Delete(&model.Vacancy{}).Error; err != nil {
		return vacancyRepoError("failed to delete the vacancy", "04")
	}

	return utils.Error{}
}
