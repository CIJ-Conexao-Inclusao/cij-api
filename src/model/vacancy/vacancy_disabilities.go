package model

import (
	"cij_api/src/model"
)

type VacancyDisability struct {
	VacancyId    int `gorm:"type:int;not null" json:"vacancy_id"`
	DisabilityId int `gorm:"type:int;not null" json:"disability_id"`
	Vacancy      *Vacancy
	Disability   *model.Disability
}

type VacancyDisabilityRequest int
