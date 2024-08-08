package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PersonDisabilityRepo interface {
	BaseRepoMethods

	GetPersonDisabilities(personId int) ([]model.PersonDisability, utils.Error)
	GetDisabilityById(disabilityId int) (model.Disability, utils.Error)
	UpsertPersonDisability(personDisability model.PersonDisability, tx *gorm.DB) utils.Error
	ClearPersonDisability(personId int, tx *gorm.DB) utils.Error
}

type personDisabilityRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewPersonDisabilityRepo(db *gorm.DB) PersonDisabilityRepo {
	repo := &personDisabilityRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func personDisabilityRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.DisabilityErrorType, code)

	return utils.NewError(message, errorCode)
}

func (n *personDisabilityRepo) GetPersonDisabilities(personId int) ([]model.PersonDisability, utils.Error) {
	var disabilities []model.PersonDisability

	err := n.db.Model(model.PersonDisability{}).Preload("Disability").Where("person_id = ?", personId).Find(&disabilities).Error
	if err != nil {
		return disabilities, personDisabilityRepoError("failed to get the person disabilities", "01")
	}

	return disabilities, utils.Error{}
}

func (n *personDisabilityRepo) GetDisabilityById(disabilityId int) (model.Disability, utils.Error) {
	var disability model.Disability

	err := n.db.Model(model.Disability{}).Where("id = ?", disabilityId).Find(&disability).Error
	if err != nil {
		return disability, personDisabilityRepoError("failed to get the disability", "02")
	}

	return disability, utils.Error{}
}

func (n *personDisabilityRepo) UpsertPersonDisability(personDisability model.PersonDisability, tx *gorm.DB) utils.Error {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	err := databaseConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "person_id"}, {Name: "disability_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"acquired"}),
	}).Create(&personDisability).Error

	if err != nil {
		return personDisabilityRepoError("failed to upsert the person disability", "03")
	}

	return utils.Error{}
}

func (n *personDisabilityRepo) ClearPersonDisability(personId int, tx *gorm.DB) utils.Error {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Where("person_id = ?", personId).Delete(&model.PersonDisability{}).Error; err != nil {
		return personDisabilityRepoError("failed to clear the person disability", "04")
	}

	return utils.Error{}
}
