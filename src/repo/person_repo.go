package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type PersonRepo interface {
	BaseRepoMethods

	CreatePerson(createPerson model.Person, tx *gorm.DB) (int, utils.Error)
	ListPeople() ([]model.Person, utils.Error)
	GetPersonById(personId int, tx *gorm.DB) (model.Person, utils.Error)
	GetPersonByUserId(userId int) (model.Person, utils.Error)
	GetPersonByCpf(cpf string) (model.Person, utils.Error)
	UpdatePerson(person model.Person, personId int, tx *gorm.DB) utils.Error
	DeletePerson(personId int) utils.Error
	UploadCurriculum(personId int, fileUrl string) utils.Error
}

type personRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewPersonRepo(db *gorm.DB) PersonRepo {
	repo := &personRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func personRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.PersonErrorType, code)

	return utils.NewError(message, errorCode)
}

func (n *personRepo) CreatePerson(createPerson model.Person, tx *gorm.DB) (int, utils.Error) {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&createPerson).Error; err != nil {
		return 0, personRepoError("failed to create the person", "01")
	}

	return createPerson.Id, utils.Error{}
}

func (n *personRepo) ListPeople() ([]model.Person, utils.Error) {
	var people []model.Person

	err := n.db.Model(model.Person{}).Preload("User").Find(&people).Error
	if err != nil {
		return people, personRepoError("failed to list the people", "02")
	}

	return people, utils.Error{}
}

func (n *personRepo) GetPersonById(personId int, tx *gorm.DB) (model.Person, utils.Error) {
	var person model.Person
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	err := databaseConn.Model(model.Person{}).Preload("User").Preload("Address").Where("id = ?", personId).Find(&person).Error
	if err != nil {
		return person, personRepoError("failed to get the person", "03")
	}

	return person, utils.Error{}
}

func (n *personRepo) GetPersonByUserId(userId int) (model.Person, utils.Error) {
	var person model.Person

	err := n.db.Model(model.Person{}).Preload("User").Where("user_id = ?", userId).Find(&person).Error
	if err != nil {
		return person, personRepoError("failed to get the person", "04")
	}

	return person, utils.Error{}
}

func (n *personRepo) GetPersonByCpf(cpf string) (model.Person, utils.Error) {
	var person model.Person

	err := n.db.Model(model.Person{}).Preload("User").Where("cpf = ?", cpf).Find(&person).Error
	if err != nil {
		return person, personRepoError("failed to get the person", "05")
	}

	return person, utils.Error{}
}

func (n *personRepo) UpdatePerson(person model.Person, personId int, tx *gorm.DB) utils.Error {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Model(model.Person{}).Where("id = ?", personId).Updates(person).Error; err != nil {
		return personRepoError("failed to update the person", "06")
	}

	return utils.Error{}
}

func (n *personRepo) DeletePerson(personId int) utils.Error {
	if err := n.db.Model(model.Person{}).Where("id = ?", personId).Unscoped().Delete(&model.Person{}).Error; err != nil {
		return personRepoError("failed to delete the person", "07")
	}

	return utils.Error{}
}

func (n *personRepo) UploadCurriculum(personId int, fileUrl string) utils.Error {
	if err := n.db.Model(model.Person{}).Where("id = ?", personId).Update("curriculum", fileUrl).Error; err != nil {
		return personRepoError("failed to upload the curriculum", "08")
	}

	return utils.Error{}
}
