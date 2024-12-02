package service

import (
	"cij_api/src/enum"
	"cij_api/src/model"
	modelVacancy "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	repoVacancy "cij_api/src/repo/vacancy"
	"cij_api/src/utils"
	"slices"

	"gorm.io/gorm"
)

type vacancyService struct {
	vacancyRepo             repoVacancy.VacancyRepo
	skillsRepo              repoVacancy.SkillsRepo
	requirementsRepo        repoVacancy.RequirementsRepo
	responsabilitiesRepo    repoVacancy.ResponsabilitiesRepo
	vacancyDisabilitiesRepo repoVacancy.VacancyDisabilityRepo
	vacancyAppliesRepo      repoVacancy.VacancyApplyRepo
	personRepo              repo.PersonRepo
	personDisabilitiesRepo  repo.PersonDisabilityRepo
}

type VacancyService interface {
	CreateVacancy(vacancy modelVacancy.VacancyRequest) utils.Error
	ListVacancies(perPage int, companyId int, disabilityId int, candidateId int, area string, contractType enum.VacancyContractType, searchText string) ([]modelVacancy.VacancySimpleResponse, utils.Error)
	GetVacancyById(id int, candidateId int) (modelVacancy.VacancyResponse, utils.Error)
	UpdateVacancy(vacancy modelVacancy.VacancyRequest, id int) utils.Error
	DeleteVacancy(id int) utils.Error

	CandidateApplyVacancy(candidateId int, vacancyId int) utils.Error
	GetVacancyAppliesByVacancyId(vacancyId int) ([]modelVacancy.VacancyApplyResponse, utils.Error)
	UpdateVacancyApplyStatus(vacancyApplyId int, status enum.VacancyApplyStatus) utils.Error
}

func NewVacancyService(
	vacancyRepo repoVacancy.VacancyRepo,
	skillsRepo repoVacancy.SkillsRepo,
	requirementsRepo repoVacancy.RequirementsRepo,
	responsabilitiesRepo repoVacancy.ResponsabilitiesRepo,
	vacancyDisabilitiesRepo repoVacancy.VacancyDisabilityRepo,
	vacancyAppliesRepo repoVacancy.VacancyApplyRepo,
	personRepo repo.PersonRepo,
	personDisabilitiesRepo repo.PersonDisabilityRepo,
) VacancyService {
	return &vacancyService{
		vacancyRepo:             vacancyRepo,
		skillsRepo:              skillsRepo,
		requirementsRepo:        requirementsRepo,
		responsabilitiesRepo:    responsabilitiesRepo,
		vacancyDisabilitiesRepo: vacancyDisabilitiesRepo,
		vacancyAppliesRepo:      vacancyAppliesRepo,
		personRepo:              personRepo,
		personDisabilitiesRepo:  personDisabilitiesRepo,
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

func (v *vacancyService) ListVacancies(perPage int, companyId int, disabilityId int, candidateId int, area string, contractType enum.VacancyContractType, searchText string) ([]modelVacancy.VacancySimpleResponse, utils.Error) {
	var vacanciesResponse []modelVacancy.VacancySimpleResponse

	vacancies, err := v.vacancyRepo.ListVacancies(companyId, area, contractType, searchText)
	if err.Code != "" {
		return []modelVacancy.VacancySimpleResponse{}, vacancyServiceError("failed to list the vacancies", "02")
	}

DisabilityLoop:
	for _, vacancy := range vacancies {
		var disabilities []model.DisabilityResponse

		vacancyDisabilities, err := v.vacancyDisabilitiesRepo.GetVacancyDisabilities(vacancy.Id)
		if err.Code != "" {
			return []modelVacancy.VacancySimpleResponse{}, vacancyServiceError("failed to get the disabilities", "03")
		}

		uniqueDisabilities := map[int]bool{}

		for _, vacancyDisability := range vacancyDisabilities {
			disabilities = append(disabilities, vacancyDisability.Disability.ToResponse())
			uniqueDisabilities[vacancyDisability.Disability.Id] = true
		}

		if disabilityId != 0 && !uniqueDisabilities[disabilityId] {
			continue DisabilityLoop
		}

		if candidateId != 0 {
			vacancyApplies, err := v.vacancyAppliesRepo.ListVacancyAppliesByVacancyId(vacancy.Id)
			if err.Code != "" {
				return []modelVacancy.VacancySimpleResponse{}, vacancyServiceError("failed to get the vacancy applies", "04")
			}

			var candidateIds []int

			for _, vacancyApply := range vacancyApplies {
				candidateIds = append(candidateIds, vacancyApply.CandidateId)
			}

			if !slices.Contains(candidateIds, candidateId) {
				continue DisabilityLoop
			}
		}

		if len(vacanciesResponse) >= perPage {
			break
		}

		vacanciesResponse = append(vacanciesResponse, vacancy.ToSimpleResponse(disabilities))
	}

	return vacanciesResponse, utils.Error{}
}

func (v *vacancyService) GetVacancyById(id int, candidateId int) (modelVacancy.VacancyResponse, utils.Error) {
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

	vacancyResponse := vacancy.ToResponse(
		disabilities,
		skills,
		responsabilities,
		requirements,
	)

	if candidateId != 0 {
		vacancyApplies, err := v.vacancyAppliesRepo.ListVacancyAppliesByVacancyIdAndCandidateId(id, candidateId)
		if err.Code != "" {
			return modelVacancy.VacancyResponse{}, vacancyServiceError("failed to get the vacancy apply", "08")
		}

		vacancyResponse.CandidateAlreadyApplied = len(vacancyApplies) > 0
	}

	return vacancyResponse, utils.Error{}
}

func (v *vacancyService) UpdateVacancy(vacancy modelVacancy.VacancyRequest, id int) utils.Error {
	vacancyModel := vacancy.ToModel()

	_, err := v.vacancyRepo.GetVacancyById(id)
	if err.Code != "" {
		return vacancyServiceError("failed to get the vacancy", "07")
	}

	vacancyModel.Id = id

	errTx := v.vacancyRepo.BeginTransaction(func(tx *gorm.DB) error {
		err := v.vacancyRepo.UpdateVacancy(*vacancyModel, tx)
		if err.Code != "" {
			return err
		}

		err = v.skillsRepo.DeleteSkillsByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.requirementsRepo.DeleteRequirementsByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.responsabilitiesRepo.DeleteResponsabilitiesByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.vacancyDisabilitiesRepo.ClearVacancyDisability(id, tx)
		if err.Code != "" {
			return err
		}

		for _, skill := range vacancy.Skills {
			skillModel := skill.ToModel()
			skillModel.VacancyId = id

			_, err := v.skillsRepo.CreateSkill(*skillModel, tx)
			if err.Code != "" {
				return err
			}
		}

		for _, requirement := range vacancy.Requirements {
			requirementModel := requirement.ToModel()
			requirementModel.VacancyId = id

			_, err := v.requirementsRepo.CreateRequirement(*requirementModel, tx)
			if err.Code != "" {
				return err
			}
		}

		for _, responsability := range vacancy.Responsabilities {
			responsabilityModel := responsability.ToModel()
			responsabilityModel.VacancyId = id

			_, err := v.responsabilitiesRepo.CreateResponsability(*responsabilityModel, tx)
			if err.Code != "" {
				return err
			}
		}

		for _, disability := range vacancy.Disabilities {
			disabilityModel := modelVacancy.VacancyDisability{
				VacancyId:    id,
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
		return vacancyServiceError("failed to update the vacancy", "08")
	}

	return utils.Error{}
}

func (v *vacancyService) DeleteVacancy(id int) utils.Error {
	_, err := v.vacancyRepo.GetVacancyById(id)
	if err.Code != "" {
		return vacancyServiceError("failed to get the vacancy", "07")
	}

	errTx := v.vacancyRepo.BeginTransaction(func(tx *gorm.DB) error {
		err := v.skillsRepo.DeleteSkillsByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.requirementsRepo.DeleteRequirementsByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.responsabilitiesRepo.DeleteResponsabilitiesByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.vacancyDisabilitiesRepo.ClearVacancyDisability(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.vacancyAppliesRepo.DeleteVacancyAppliesByVacancyId(id, tx)
		if err.Code != "" {
			return err
		}

		err = v.vacancyRepo.DeleteVacancy(id)
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

func (v *vacancyService) CandidateApplyVacancy(candidateId int, vacancyId int) utils.Error {
	_, err := v.vacancyRepo.GetVacancyById(vacancyId)
	if err.Code != "" {
		return vacancyServiceError("failed to get the vacancy", "10")
	}

	_, err = v.personRepo.GetPersonById(candidateId, nil)
	if err.Code != "" {
		return vacancyServiceError("failed to get the person", "11")
	}

	vacancyApplyDb, _ := v.vacancyAppliesRepo.GetVacancyApply(vacancyId, candidateId)
	if vacancyApplyDb.Id != 0 {
		return vacancyServiceError("the candidate already applied to the vacancy", "13")
	}

	vacancyApply := modelVacancy.VacancyApply{
		VacancyId:   vacancyId,
		CandidateId: candidateId,
		Status:      enum.VacancyApplyApplied,
	}

	_, err = v.vacancyAppliesRepo.CreateVacancyApply(vacancyApply)
	if err.Code != "" {
		return vacancyServiceError("failed to apply the vacancy", "12")
	}

	return utils.Error{}
}

func (v *vacancyService) GetVacancyAppliesByVacancyId(vacancyId int) ([]modelVacancy.VacancyApplyResponse, utils.Error) {
	vacancyApplies, err := v.vacancyAppliesRepo.ListVacancyAppliesByVacancyId(vacancyId)
	if err.Code != "" {
		return []modelVacancy.VacancyApplyResponse{}, vacancyServiceError("failed to get the vacancy applies", "13")
	}

	var vacancyAppliesResponse []modelVacancy.VacancyApplyResponse
	for _, vacancyApply := range vacancyApplies {
		person, err := v.personRepo.GetPersonById(vacancyApply.CandidateId, nil)
		if err.Code != "" {
			return []modelVacancy.VacancyApplyResponse{}, vacancyServiceError("failed to get the person", "14")
		}

		candidateDisabilities, err := v.personDisabilitiesRepo.GetPersonDisabilities(vacancyApply.CandidateId)
		if err.Code != "" {
			return []modelVacancy.VacancyApplyResponse{}, vacancyServiceError("failed to get the candidate disabilities", "15")
		}

		candidateDisabilitiesResponse := []model.DisabilityResponse{}
		for _, candidateDisability := range candidateDisabilities {
			candidateDisabilitiesResponse = append(candidateDisabilitiesResponse, candidateDisability.Disability.ToResponse())
		}

		vacancyApplyResponse := modelVacancy.VacancyApplyResponse{
			Candidate: vacancyApply.Candidate.ToCandidateResponse(candidateDisabilitiesResponse, *person.Address),
			Status:    vacancyApply.Status,
			Id:        vacancyApply.Id,
		}

		vacancyAppliesResponse = append(vacancyAppliesResponse, vacancyApplyResponse)
	}

	return vacancyAppliesResponse, utils.Error{}
}

func (v *vacancyService) UpdateVacancyApplyStatus(vacancyApplyId int, status enum.VacancyApplyStatus) utils.Error {
	err := v.vacancyAppliesRepo.UpdateVacancyApplyStatus(vacancyApplyId, status)
	if err.Code != "" {
		return vacancyServiceError("failed to update the vacancy apply status", "14")
	}

	return utils.Error{}
}
