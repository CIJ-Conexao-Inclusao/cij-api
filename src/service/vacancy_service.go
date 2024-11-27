package service

import (
	"cij_api/src/enum"
	"cij_api/src/model"
	modelVacancy "cij_api/src/model/vacancy"
	repo "cij_api/src/repo/vacancy"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type vacancyService struct {
	vacancyRepo             repo.VacancyRepo
	skillsRepo              repo.SkillsRepo
	requirementsRepo        repo.RequirementsRepo
	responsabilitiesRepo    repo.ResponsabilitiesRepo
	vacancyDisabilitiesRepo repo.VacancyDisabilityRepo
}

type VacancyService interface {
	CreateVacancy(vacancy modelVacancy.VacancyRequest) utils.Error
	ListVacancies(page int, perPage int, companyId int, disabilityCategory string, area string, contractType enum.VacancyContractType, searchText string) ([]modelVacancy.VacancySimpleResponse, utils.Error)
	GetVacancyById(id int) (modelVacancy.VacancyResponse, utils.Error)
	UpdateVacancy(vacancy modelVacancy.VacancyRequest, id int) utils.Error
	DeleteVacancy(id int) utils.Error
}

func NewVacancyService(
	vacancyRepo repo.VacancyRepo,
	skillsRepo repo.SkillsRepo,
	requirementsRepo repo.RequirementsRepo,
	responsabilitiesRepo repo.ResponsabilitiesRepo,
	vacancyDisabilitiesRepo repo.VacancyDisabilityRepo,
) VacancyService {
	return &vacancyService{
		vacancyRepo:             vacancyRepo,
		skillsRepo:              skillsRepo,
		requirementsRepo:        requirementsRepo,
		responsabilitiesRepo:    responsabilitiesRepo,
		vacancyDisabilitiesRepo: vacancyDisabilitiesRepo,
	}
}

func vacancyServiceError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (v *vacancyService) CreateVacancy(vacancy modelVacancy.VacancyRequest) utils.Error {
	vacancyModel := vacancy.ToModel()

	errTx := v.vacancyRepo.BeginTransaction(func(tx *gorm.DB) error {
		vacancyId, err := v.vacancyRepo.UpsertVacancy(*vacancyModel, tx)
		if err.Code != "" {
			return err
		}

		for _, skill := range vacancy.Skills {
			skillModel := skill.ToModel()
			skillModel.VacancyId = vacancyId

			_, err := v.skillsRepo.CreateSkill(*skillModel, tx)
			if err.Code != "" {
				return err
			}
		}

		for _, requirement := range vacancy.Requirements {
			requirementModel := requirement.ToModel()
			requirementModel.VacancyId = vacancyId

			_, err := v.requirementsRepo.CreateRequirement(*requirementModel, tx)
			if err.Code != "" {
				return err
			}
		}

		for _, responsability := range vacancy.Responsabilities {
			responsabilityModel := responsability.ToModel()
			responsabilityModel.VacancyId = vacancyId

			_, err := v.responsabilitiesRepo.CreateResponsability(*responsabilityModel, tx)
			if err.Code != "" {
				return err
			}
		}

		for _, disability := range vacancy.Disabilities {
			disabilityModel := modelVacancy.VacancyDisability{
				VacancyId:    vacancyId,
				DisabilityId: int(disability),
			}

			err := v.vacancyDisabilitiesRepo.UpsertVacancyDisability(disabilityModel, tx)
			if err.Code != "" {
				return err
			}
		}

		return nil
	})

	if errTx != nil {
		return vacancyServiceError("failed to create the vacancy", "01")
	}

	return utils.Error{}
}

func (v *vacancyService) ListVacancies(page int, perPage int, companyId int, disabilityCategory string, area string, contractType enum.VacancyContractType, searchText string) ([]modelVacancy.VacancySimpleResponse, utils.Error) {
	var vacanciesResponse []modelVacancy.VacancySimpleResponse

	vacancies, err := v.vacancyRepo.ListVacancies(page, perPage, companyId, disabilityCategory, area, contractType, searchText)
	if err.Code != "" {
		return []modelVacancy.VacancySimpleResponse{}, vacancyServiceError("failed to list the vacancies", "02")
	}

	for _, vacancy := range vacancies {
		var disabilities []model.DisabilityResponse

		vacancyDisabilities, err := v.vacancyDisabilitiesRepo.GetVacancyDisabilities(vacancy.Id)
		if err.Code != "" {
			return []modelVacancy.VacancySimpleResponse{}, vacancyServiceError("failed to get the disabilities", "03")
		}

		for _, vacancyDisability := range vacancyDisabilities {
			disabilities = append(disabilities, vacancyDisability.Disability.ToResponse())
		}

		vacanciesResponse = append(vacanciesResponse, vacancy.ToSimpleResponse(disabilities))
	}

	return vacanciesResponse, utils.Error{}
}

func (v *vacancyService) GetVacancyById(id int) (modelVacancy.VacancyResponse, utils.Error) {
	vacancy, err := v.vacancyRepo.GetVacancyById(id)
	if err.Code != "" {
		return modelVacancy.VacancyResponse{}, vacancyServiceError("failed to get the vacancy", "03")
	}

	skills, err := v.skillsRepo.ListSkillsByVacancyId(id)
	if err.Code != "" {
		return modelVacancy.VacancyResponse{}, vacancyServiceError("failed to get the skills", "04")
	}

	requirements, err := v.requirementsRepo.ListRequirementsByVacancyId(id)
	if err.Code != "" {
		return modelVacancy.VacancyResponse{}, vacancyServiceError("failed to get the requirements", "05")
	}

	responsabilities, err := v.responsabilitiesRepo.ListResponsabilitiesByVacancyId(id)
	if err.Code != "" {
		return modelVacancy.VacancyResponse{}, vacancyServiceError("failed to get the responsabilities", "06")
	}

	vacancyDisabilities, err := v.vacancyDisabilitiesRepo.GetVacancyDisabilities(id)
	if err.Code != "" {
		return modelVacancy.VacancyResponse{}, vacancyServiceError("failed to get the disabilities", "07")
	}

	disabilities := []model.DisabilityResponse{}
	for _, vacancyDisability := range vacancyDisabilities {
		disabilities = append(disabilities, vacancyDisability.Disability.ToResponse())
	}

	return vacancy.ToResponse(
		disabilities,
		skills,
		responsabilities,
		requirements,
	), utils.Error{}
}

func (v *vacancyService) UpdateVacancy(vacancy modelVacancy.VacancyRequest, id int) utils.Error {
	// vacancyModel := vacancy.ToModel()

	errTx := v.vacancyRepo.BeginTransaction(func(tx *gorm.DB) error {
		_, err := v.vacancyRepo.GetVacancyById(id)
		if err.Code != "" {
			return err
		}

		return nil
	})

	if errTx != nil {
		return vacancyServiceError("failed to update the vacancy", "08")
	}

	return utils.Error{}
}

func (v *vacancyService) DeleteVacancy(id int) utils.Error {
	errTx := v.vacancyRepo.BeginTransaction(func(tx *gorm.DB) error {
		_, err := v.vacancyRepo.GetVacancyById(id)
		if err.Code != "" {
			return err
		}

		return nil
	})

	if errTx != nil {
		return vacancyServiceError("failed to delete the vacancy", "09")
	}

	return utils.Error{}
}
