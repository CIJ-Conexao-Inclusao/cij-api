package repo

import (
	model "cij_api/src/model/vacancy"
	"cij_api/src/repo"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type SkillsRepo interface {
	repo.BaseRepoMethods

	CreateSkill(createSkill model.VacancySkill, tx *gorm.DB) (int, utils.Error)
	ListSkillsByVacancyId(vacancyId int) ([]model.VacancySkill, utils.Error)
	UpdateSkill(skill model.VacancySkill, skillId int, tx *gorm.DB) utils.Error
	DeleteSkill(skillId int) utils.Error
}

type skillsRepo struct {
	repo.BaseRepo
	db *gorm.DB
}

func NewSkillsRepo(db *gorm.DB) SkillsRepo {
	repo := &skillsRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func skillsRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.VacancyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (s *skillsRepo) CreateSkill(createSkill model.VacancySkill, tx *gorm.DB) (int, utils.Error) {
	databaseConn := s.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&createSkill).Error; err != nil {
		return 0, skillsRepoError("failed to create the skill", "01")
	}

	return createSkill.Id, utils.Error{}
}

func (s *skillsRepo) ListSkillsByVacancyId(vacancyId int) ([]model.VacancySkill, utils.Error) {
	var skills []model.VacancySkill

	if err := s.db.Where("vacancy_id = ?", vacancyId).Find(&skills).Error; err != nil {
		return []model.VacancySkill{}, skillsRepoError("failed to list the skills", "02")
	}

	return skills, utils.Error{}
}

func (s *skillsRepo) UpdateSkill(skill model.VacancySkill, skillId int, tx *gorm.DB) utils.Error {
	databaseConn := s.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Model(&model.VacancySkill{}).Where("id = ?", skillId).Updates(&skill).Error; err != nil {
		return skillsRepoError("failed to update the skill", "03")
	}

	return utils.Error{}
}

func (s *skillsRepo) DeleteSkill(skillId int) utils.Error {
	if err := s.db.Where("id = ?", skillId).Delete(&model.VacancySkill{}).Error; err != nil {
		return skillsRepoError("failed to delete the skill", "04")
	}

	return utils.Error{}
}
