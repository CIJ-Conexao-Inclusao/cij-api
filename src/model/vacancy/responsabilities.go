package model

import "gorm.io/gorm"

type VacancyResponsability struct {
	*gorm.Model
	Id             int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Responsability string `gorm:"type:text;not null" json:"responsability"`
	VacancyId      int    `gorm:"type:int;not null" json:"vacancy_id"`
	Vacancy        *Vacancy
}

type VacancyResponsabilityResponse string

type VacancyResponsabilityRequest string

func (v *VacancyResponsabilityRequest) ToModel() *VacancyResponsability {
	return &VacancyResponsability{
		Responsability: string(*v),
	}
}

func (v *VacancyResponsability) ToResponse() *VacancyResponsabilityResponse {
	return (*VacancyResponsabilityResponse)(&v.Responsability)
}
