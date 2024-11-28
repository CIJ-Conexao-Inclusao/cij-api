package model

import (
	"cij_api/src/enum"
	"cij_api/src/model"
)

type VacancyApply struct {
	Id          int                     `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	VacancyId   int                     `gorm:"type:int;not null" json:"vacancy_id"`
	CandidateId int                     `gorm:"type:int;not null" json:"candidate_id"`
	Status      enum.VacancyApplyStatus `gorm:"type:varchar(10);not null" json:"status"`
	Vacancy     *Vacancy
	Candidate   *model.Person
}

type VacancyApplyRequest struct {
	VacancyId   int `json:"vacancy_id"`
	CandidateId int `json:"candidate_id"`
}

type VacancyApplyResponse struct {
	Id        int                     `json:"id"`
	Candidate model.CandidateResponse `json:"candidate"`
	Status    enum.VacancyApplyStatus `json:"status"`
}

func (v *VacancyApplyRequest) ToModel() *VacancyApply {
	return &VacancyApply{
		VacancyId: v.VacancyId,
		Status:    enum.VacancyApplyApplied,
	}
}
