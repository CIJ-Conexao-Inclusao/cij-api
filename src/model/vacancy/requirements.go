package model

import (
	"gorm.io/gorm"

	"cij_api/src/enum"
)

type VacancyRequirement struct {
	*gorm.Model
	Id          int                         `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Requirement string                      `gorm:"type:text;not null" json:"requirement"`
	Type        enum.VacancyRequirementType `gorm:"type:varchar(200);not null" json:"type"`
	VacancyId   int                         `gorm:"type:int;not null" json:"vacancy_id"`
	Vacancy     *Vacancy
}

type VacancyRequirementResponse struct {
	Requirement string                      `json:"requirement"`
	Type        enum.VacancyRequirementType `json:"type"`
}

type VacancyRequirementRequest struct {
	Requirement string                      `json:"requirement"`
	Type        enum.VacancyRequirementType `json:"type"`
}

func (v *VacancyRequirementRequest) ToModel() *VacancyRequirement {
	return &VacancyRequirement{
		Requirement: v.Requirement,
		Type:        v.Type,
	}
}

func (v *VacancyRequirement) ToResponse() *VacancyRequirementResponse {
	return &VacancyRequirementResponse{
		Requirement: v.Requirement,
		Type:        v.Type,
	}
}
