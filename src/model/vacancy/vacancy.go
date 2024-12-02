package model

import (
	"cij_api/src/enum"
	"cij_api/src/model"

	"gorm.io/gorm"
)

type Vacancy struct {
	*gorm.Model
	Id               int                      `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Code             string                   `gorm:"type:varchar(200);not null" json:"code"`
	Title            string                   `gorm:"type:varchar(200);not null" json:"title"`
	Description      string                   `gorm:"type:text;not null" json:"description"`
	Department       string                   `gorm:"type:varchar(200);not null" json:"department"`
	Section          string                   `gorm:"type:varchar(200);not null" json:"section"`
	Turn             string                   `gorm:"type:varchar(200);not null" json:"turn"`
	PublishDate      string                   `gorm:"type:date;not null" json:"publish_date"`
	RegistrationDate string                   `gorm:"type:date;not null" json:"registration_date"`
	Area             string                   `gorm:"type:varchar(200);not null" json:"area"`
	CompanyId        int                      `gorm:"type:int;not null" json:"company_id"`
	ContractType     enum.VacancyContractType `gorm:"type:varchar(200);not null" json:"contract_type"`
	Disabilities     []model.Disability       `gorm:"many2many:vacancy_disabilities" json:"disabilities"`
	Company          model.Company
}

type VacancyResponse struct {
	Id                      int                             `json:"id"`
	Code                    string                          `json:"code"`
	Title                   string                          `json:"title"`
	Description             string                          `json:"description"`
	Department              string                          `json:"department"`
	Section                 string                          `json:"section"`
	Turn                    string                          `json:"turn"`
	PublishDate             string                          `json:"publish_date"`
	RegistrationDate        string                          `json:"registration_date"`
	Area                    string                          `json:"area"`
	CandidateAlreadyApplied bool                            `json:"candidate_already_applied,omitempty"`
	ContractType            enum.VacancyContractType        `json:"contract_type"`
	Company                 string                          `json:"company"`
	Disabilities            []model.DisabilityResponse      `json:"disabilities"`
	Skills                  []VacancySkillResponse          `json:"skills"`
	Responsabilities        []VacancyResponsabilityResponse `json:"responsabilities"`
	Requirements            []VacancyRequirementResponse    `json:"requirements"`
}

type VacancySimpleResponse struct {
	Id           int                        `json:"id"`
	Code         string                     `json:"code"`
	Title        string                     `json:"title"`
	Area         string                     `json:"area"`
	Company      string                     `json:"company"`
	ContractType enum.VacancyContractType   `json:"contract_type"`
	Disabilities []model.DisabilityResponse `json:"disabilities"`
}

type VacancyRequest struct {
	Code             string                         `json:"code"`
	Title            string                         `json:"title"`
	Description      string                         `json:"description"`
	Department       string                         `json:"department"`
	Section          string                         `json:"section"`
	Turn             string                         `json:"turn"`
	PublishDate      string                         `json:"publish_date"`
	RegistrationDate string                         `json:"registration_date"`
	Area             string                         `json:"area"`
	CompanyId        int                            `json:"company_id"`
	ContractType     enum.VacancyContractType       `json:"contract_type"`
	Disabilities     []VacancyDisabilityRequest     `json:"disabilities"`
	Skills           []VacancySkillRequest          `json:"skills"`
	Responsabilities []VacancyResponsabilityRequest `json:"responsabilities"`
	Requirements     []VacancyRequirementRequest    `json:"requirements"`
}

func (v *VacancyRequest) ToModel() *Vacancy {
	return &Vacancy{
		Code:             v.Code,
		Title:            v.Title,
		Description:      v.Description,
		Department:       v.Department,
		Section:          v.Section,
		Turn:             v.Turn,
		PublishDate:      v.PublishDate,
		RegistrationDate: v.RegistrationDate,
		Area:             v.Area,
		ContractType:     v.ContractType,
		CompanyId:        v.CompanyId,
	}
}

func (v *Vacancy) ToResponse(
	disabilities []model.DisabilityResponse,
	skills []VacancySkill,
	responsabilities []VacancyResponsability,
	requirements []VacancyRequirement,
) VacancyResponse {
	var skillsResponse []VacancySkillResponse
	var responsabilitiesResponse []VacancyResponsabilityResponse
	var requirementsResponse []VacancyRequirementResponse

	for _, s := range skills {
		skillsResponse = append(skillsResponse, *s.ToResponse())
	}

	for _, r := range responsabilities {
		responsabilitiesResponse = append(responsabilitiesResponse, *r.ToResponse())
	}

	for _, r := range requirements {
		requirementsResponse = append(requirementsResponse, *r.ToResponse())
	}

	return VacancyResponse{
		Id:               v.Id,
		Code:             v.Code,
		Title:            v.Title,
		Description:      v.Description,
		Department:       v.Department,
		Section:          v.Section,
		Turn:             v.Turn,
		PublishDate:      v.PublishDate,
		RegistrationDate: v.RegistrationDate,
		Area:             v.Area,
		ContractType:     v.ContractType,
		Company:          v.Company.Name,
		Disabilities:     disabilities,
		Skills:           skillsResponse,
		Responsabilities: responsabilitiesResponse,
		Requirements:     requirementsResponse,
	}
}

func (v *Vacancy) ToSimpleResponse(disabilities []model.DisabilityResponse) VacancySimpleResponse {
	return VacancySimpleResponse{
		Id:           v.Id,
		Code:         v.Code,
		Title:        v.Title,
		Area:         v.Area,
		Company:      v.Company.Name,
		ContractType: v.ContractType,
		Disabilities: disabilities,
	}
}
