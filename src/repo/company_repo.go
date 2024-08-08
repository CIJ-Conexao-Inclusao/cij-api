package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type CompanyRepo interface {
	BaseRepoMethods

	CreateCompany(createCompany model.Company, tx *gorm.DB) utils.Error
	ListCompanies() ([]model.Company, utils.Error)
	GetCompanyById(companyId int) (model.Company, utils.Error)
	GetCompanyByUserId(userId int) (model.Company, utils.Error)
	GetCompanyByCnpj(cnpj string) (model.Company, utils.Error)
	UpdateCompany(company model.Company, companyId int) utils.Error
	DeleteCompany(companyId int) utils.Error
}

type companyRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewCompanyRepo(db *gorm.DB) CompanyRepo {
	repo := &companyRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func companyRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.CompanyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (n *companyRepo) CreateCompany(createCompany model.Company, tx *gorm.DB) utils.Error {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&createCompany).Error; err != nil {
		return companyRepoError("failed to create the company", "01")
	}

	return utils.Error{}
}

func (n *companyRepo) ListCompanies() ([]model.Company, utils.Error) {
	var companies []model.Company

	err := n.db.Model(model.Company{}).Preload("User").Preload("Address").Find(&companies).Error
	if err != nil {
		return companies, companyRepoError("failed to list the companies", "02")
	}

	return companies, utils.Error{}
}

func (n *companyRepo) GetCompanyById(companyId int) (model.Company, utils.Error) {
	var company model.Company

	err := n.db.Model(model.Company{}).Preload("User").Where("id = ?", companyId).Find(&company).Error
	if err != nil {
		return company, companyRepoError("failed to get the company", "03")
	}

	return company, utils.Error{}
}

func (n *companyRepo) GetCompanyByUserId(userId int) (model.Company, utils.Error) {
	var company model.Company

	err := n.db.Model(model.Company{}).Preload("User").Where("user_id = ?", userId).Find(&company).Error
	if err != nil {
		return company, companyRepoError("failed to get the company", "04")
	}

	return company, utils.Error{}
}

func (n *companyRepo) UpdateCompany(company model.Company, companyId int) utils.Error {
	if err := n.db.Model(model.Company{}).Where("id = ?", companyId).Updates(company).Error; err != nil {
		return companyRepoError("failed to update the company", "05")
	}

	return utils.Error{}
}

func (n *companyRepo) DeleteCompany(companyId int) utils.Error {
	if err := n.db.Model(model.Company{}).Where("id = ?", companyId).Unscoped().Delete(&model.Company{}).Error; err != nil {
		return companyRepoError("failed to delete the company", "06")
	}

	return utils.Error{}
}

func (n *companyRepo) GetCompanyByCnpj(cnpj string) (model.Company, utils.Error) {
	var company model.Company

	err := n.db.Model(model.Company{}).Preload("User").Where("cnpj = ?", cnpj).Find(&company).Error
	if err != nil {
		return company, companyRepoError("failed to get the company", "07")
	}

	return company, utils.Error{}
}
