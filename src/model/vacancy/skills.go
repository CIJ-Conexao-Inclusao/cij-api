package model

import "gorm.io/gorm"

type VacancySkill struct {
	*gorm.Model
	Id        int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Skill     string `gorm:"type:varchar(200);not null" json:"skill"`
	VacancyId int    `gorm:"type:int;not null" json:"vacancy_id"`
	Vacancy   *Vacancy
}

type VacancySkillResponse struct {
	Skill string `json:"skill"`
}

type VacancySkillRequest string

func (v *VacancySkillRequest) ToModel() *VacancySkill {
	return &VacancySkill{
		Skill: string(*v),
	}
}

func (v *VacancySkill) ToResponse() *VacancySkillResponse {
	return &VacancySkillResponse{
		Skill: v.Skill,
	}
}
