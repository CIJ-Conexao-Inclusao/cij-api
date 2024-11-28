package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
	"fmt"

	"gorm.io/gorm"
)

type CompanyService interface {
	CreateCompany(createCompany model.CompanyRequest) utils.Error
	ListCompanies() ([]model.CompanyResponse, utils.Error)
	GetCompanyByUserId(userId int) (model.Company, utils.Error)
	GetCompanyByCnpj(cnpj string) (model.Company, utils.Error)
	GetCompanyById(companyId int) (model.Company, utils.Error)
	GetUserByEmail(email string) (model.User, utils.Error)
	UpdateCompany(company model.CompanyRequest, companyId int) utils.Error
	DeleteCompany(companyId int) utils.Error
}

type companyService struct {
	companyRepo  repo.CompanyRepo
	userRepo     repo.UserRepo
	addressRepo  repo.AddressRepo
	activityRepo repo.ActivityRepo
}

func NewCompanyService(
	companyRepo repo.CompanyRepo,
	userRepo repo.UserRepo,
	addressRepo repo.AddressRepo,
	activityRepo repo.ActivityRepo,
) CompanyService {
	return &companyService{
		companyRepo:  companyRepo,
		userRepo:     userRepo,
		addressRepo:  addressRepo,
		activityRepo: activityRepo,
	}
}

func companyServiceError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ServiceErrorCode, utils.CompanyErrorType, code)

	return utils.NewError(message, errorCode)
}

func (s *companyService) ListCompanies() ([]model.CompanyResponse, utils.Error) {
	companiesResponse := []model.CompanyResponse{}

	companies, err := s.companyRepo.ListCompanies()
	if err.Code != "" {
		return companiesResponse, err
	}

	for _, company := range companies {
		user, err := s.userRepo.GetUserById(company.User.Id)
		if err.Code != "" {
			return companiesResponse, err
		}

		companyResponse := company.ToResponse(user)

		address, err := s.addressRepo.GetAddressById(*company.AddressId)
		if err.Code != "" {
			return companiesResponse, err
		}

		if address.Id != 0 {
			addressResponse := address.ToResponse()
			companyResponse.Address = addressResponse
		}

		userConfig := model.DefaultConfig
		if user.ConfigUrl != "" {
			configService := NewConfigService(s.userRepo)
			userConfig, err = configService.GetUserConfig(user.ConfigUrl)
			if err.Code != "" {
				return companiesResponse, err
			}
		}

		companyResponse.User.Config = userConfig
		companiesResponse = append(companiesResponse, companyResponse)
	}

	return companiesResponse, utils.Error{}
}

func (n *companyService) CreateCompany(createCompany model.CompanyRequest) utils.Error {
	userInfo := createCompany.ToUser()

	hashedPassword, err := utils.EncryptPassword(userInfo.Password)
	if err != nil {
		return companyServiceError("failed to encrypt the password", "01")
	}

	userInfo.Password = hashedPassword
	userInfo.RoleId = model.CompanyRole

	errTx := n.userRepo.BeginTransaction(func(tx *gorm.DB) error {
		userId, userError := n.userRepo.CreateUser(userInfo, tx)
		if userError.Code != "" {
			fmt.Println("Error: ", userError)
			return userError
		}

		addressInfo := createCompany.ToAddress()

		addressId, addresError := n.addressRepo.UpsertAddress(addressInfo, tx)
		if addresError.Code != "" {
			fmt.Println("Error: ", addresError)
			return addresError
		}

		companyInfo := createCompany.ToModel(userInfo)
		companyInfo.UserId = userId
		companyInfo.AddressId = &addressId

		companyError := n.companyRepo.CreateCompany(companyInfo, tx)
		if companyError.Code != "" {
			fmt.Println("Error: ", companyError)
			return companyError
		}

		return nil
	})

	if errTx != nil {
		return companyServiceError("failed to create the company", "02")
	}

	activityService := NewActivityService(n.activityRepo)
	activity := model.Activity{
		Type:        "register_company",
		Description: "Company " + userInfo.Email + " registered",
		Actor:       userInfo.Email,
	}

	activityError := activityService.CreateActivity(&activity)
	if activityError.Code != "" {
		return activityError
	}

	return utils.Error{}
}

func (n *companyService) GetCompanyByUserId(userId int) (model.Company, utils.Error) {
	company, err := n.companyRepo.GetCompanyByUserId(userId)
	if err.Code != "" {
		return company, err
	}

	return company, utils.Error{}
}

func (n *companyService) GetCompanyByCnpj(cnpj string) (model.Company, utils.Error) {
	company, err := n.companyRepo.GetCompanyByCnpj(cnpj)
	if err.Code != "" {
		return company, err
	}

	return company, utils.Error{}
}

func (n *companyService) GetCompanyById(companyId int) (model.Company, utils.Error) {
	company, err := n.companyRepo.GetCompanyById(companyId)
	if err.Code != "" {
		return company, err
	}

	return company, utils.Error{}
}

func (n *companyService) UpdateCompany(updateCompany model.CompanyRequest, companyId int) utils.Error {
	userInfo := updateCompany.ToUser()

	if userInfo.Password == "" {
		hashedPassword, err := utils.EncryptPassword(userInfo.Password)
		if err != nil {
			return companyServiceError("failed to encrypt the password", "02")
		}

		userInfo.Password = hashedPassword

		userError := n.userRepo.UpdateUser(userInfo, companyId)
		if userError.Code != "" {
			return userError
		}
	}

	addressInfo := updateCompany.ToAddress()

	company, companyError := n.companyRepo.GetCompanyById(companyId)
	if companyError.Code != "" {
		return companyError
	}

	addressInfo.Id = *company.AddressId

	addressId, addresError := n.addressRepo.UpsertAddress(addressInfo, nil)
	if addresError.Code != "" {
		return addresError
	}

	companyInfo := updateCompany.ToModel(userInfo)
	companyInfo.AddressId = &addressId

	companyError = n.companyRepo.UpdateCompany(companyInfo, companyId)
	if companyError.Code != "" {
		return companyError
	}

	return utils.Error{}
}

func (n *companyService) DeleteCompany(companyId int) utils.Error {
	company, err := n.companyRepo.GetCompanyById(companyId)
	if err.Code != "" {
		return err
	}

	err = n.companyRepo.DeleteCompany(companyId)
	if err.Code != "" {
		return err
	}

	err = n.userRepo.DeleteUser(company.UserId)
	if err.Code != "" {
		return err
	}

	err = n.addressRepo.DeleteAddress(*company.AddressId)
	if err.Code != "" {
		return err
	}

	return utils.Error{}
}

func (n *companyService) GetUserByEmail(email string) (model.User, utils.Error) {
	user, err := n.userRepo.GetUserByEmail(email)
	if err.Code != "" {
		return user, err
	}

	return user, utils.Error{}
}
