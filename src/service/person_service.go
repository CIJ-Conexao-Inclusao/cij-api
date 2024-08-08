package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
	"fmt"

	"gorm.io/gorm"
)

type PersonService interface {
	CreatePerson(createPerson model.PersonRequest) utils.Error
	ListPeople() ([]model.PersonResponse, utils.Error)
	GetPersonByUserId(userId int) (model.Person, utils.Error)
	GetPersonById(personId int) (model.PersonResponse, utils.Error)
	GetPersonByCpf(cpf string) (model.Person, utils.Error)
	GetUserByEmail(email string) (model.User, utils.Error)
	GetDisabilityById(disabilityId int) (model.Disability, utils.Error)
	UpdatePerson(person model.PersonRequest, personId int) utils.Error
	UpdatePersonAddress(address model.AddressRequest, personId int, tx *gorm.DB) utils.Error
	UpdatePersonDisabilities(disabilities []model.PersonDisabilityRequest, personId int, tx *gorm.DB) utils.Error
	DeletePerson(personId int) utils.Error
}

type personService struct {
	personRepo           repo.PersonRepo
	userRepo             repo.UserRepo
	addressRepo          repo.AddressRepo
	personDisabilityRepo repo.PersonDisabilityRepo
}

func NewPersonService(
	personRepo repo.PersonRepo,
	userRepo repo.UserRepo,
	addressRepo repo.AddressRepo,
	personDisabilityRepo repo.PersonDisabilityRepo,
) PersonService {
	return &personService{
		personRepo:           personRepo,
		userRepo:             userRepo,
		addressRepo:          addressRepo,
		personDisabilityRepo: personDisabilityRepo,
	}
}

func personServiceError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ServiceErrorCode, utils.PersonErrorType, code)

	return utils.NewError(message, errorCode)
}

func (s *personService) ListPeople() ([]model.PersonResponse, utils.Error) {
	peopleResponse := []model.PersonResponse{}

	people, err := s.personRepo.ListPeople()
	if err.Code != "" {
		return peopleResponse, err
	}

	for _, person := range people {
		personResponse := model.PersonResponse{}

		s.personToResponse(&personResponse, person)

		peopleResponse = append(peopleResponse, personResponse)
	}

	return peopleResponse, utils.Error{}
}

func (n *personService) CreatePerson(createPerson model.PersonRequest) utils.Error {
	userInfo := createPerson.ToUser()

	hashedPassword, err := utils.EncryptPassword(userInfo.Password)
	if err != nil {
		return personServiceError("failed to encrypt the password", "01")
	}

	userInfo.Password = hashedPassword
	userInfo.RoleId = model.PersonRole

	errTx := n.userRepo.BeginTransaction(func(tx *gorm.DB) error {
		userId, userError := n.userRepo.CreateUser(userInfo, tx)
		if userError.Code != "" {
			fmt.Print("Error: ", userError)
			return userError
		}

		userInfo.Id = userId

		personInfo := createPerson.ToModel(userInfo)
		personInfo.UserId = userId

		personId, personError := n.personRepo.CreatePerson(personInfo, tx)
		if personError.Code != "" {
			fmt.Print("Error: ", personError)
			return personError
		}

		addressError := n.UpdatePersonAddress(createPerson.Address, personId, tx)
		if addressError.Code != "" {
			fmt.Println("Error: ", addressError)
			return addressError
		}

		disabilityError := n.UpdatePersonDisabilities(createPerson.Disabilities, personId, tx)
		if disabilityError.Code != "" {
			fmt.Print("Error: ", disabilityError)
			return disabilityError
		}

		return nil
	})

	if errTx != nil {
		return personServiceError("failed to create the person", "02")
	}

	configService := NewConfigService(n.userRepo)

	uploadErr := configService.UploadUserConfig(userInfo.Email, nil)
	if uploadErr.Code != "" {
		fmt.Print("Error: ", uploadErr)
		return uploadErr
	}

	return utils.Error{}
}

func (n *personService) GetPersonByUserId(userId int) (model.Person, utils.Error) {
	person, err := n.personRepo.GetPersonByUserId(userId)
	if err.Code != "" {
		return person, err
	}

	return person, utils.Error{}
}

func (s *personService) GetPersonById(personId int) (model.PersonResponse, utils.Error) {
	personResponse := model.PersonResponse{}

	person, err := s.personRepo.GetPersonById(personId, nil)
	if err.Code != "" {
		return personResponse, err
	}

	if person.Id == 0 {
		return personResponse, personServiceError("person not found", "01")
	}

	s.personToResponse(&personResponse, person)

	return personResponse, utils.Error{}
}

func (n *personService) GetPersonByCpf(cpf string) (model.Person, utils.Error) {
	person, err := n.personRepo.GetPersonByCpf(cpf)
	if err.Code != "" {
		return person, err
	}

	return person, utils.Error{}
}

func (n *personService) GetUserByEmail(email string) (model.User, utils.Error) {
	user, err := n.userRepo.GetUserByEmail(email)
	if err.Code != "" {
		return user, err
	}

	return user, utils.Error{}
}

func (n *personService) UpdatePerson(updatePerson model.PersonRequest, personId int) utils.Error {
	userInfo := updatePerson.ToUser()

	if userInfo.Password != "" {
		hashedPassword, err := utils.EncryptPassword(userInfo.Password)
		if err != nil {
			return personServiceError("failed to encrypt the password", "02")
		}

		userInfo.Password = hashedPassword

		userError := n.userRepo.UpdateUser(userInfo, personId)
		if userError.Code != "" {
			return userError
		}
	}

	personInfo := updatePerson.ToModel(userInfo)

	personError := n.personRepo.UpdatePerson(personInfo, personId, nil)
	if personError.Code != "" {
		return personError
	}

	return utils.Error{}
}

func (n *personService) UpdatePersonAddress(updateAddress model.AddressRequest, personId int, tx *gorm.DB) utils.Error {
	addressInfo := updateAddress.ToModel()

	person, err := n.personRepo.GetPersonById(personId, tx)
	if err.Code != "" {
		return err
	}

	if person.AddressId != nil {
		addressInfo.Id = *person.AddressId
	}

	addressId, err := n.addressRepo.UpsertAddress(addressInfo, tx)
	if err.Code != "" {
		return err
	}

	person.AddressId = &addressId

	err = n.personRepo.UpdatePerson(person, personId, tx)
	if err.Code != "" {
		return err
	}

	return utils.Error{}
}

func (n *personService) GetDisabilityById(id int) (model.Disability, utils.Error) {
	disability, err := n.personDisabilityRepo.GetDisabilityById(id)
	if err.Code != "" {
		return disability, err
	}

	return disability, utils.Error{}
}

func (n *personService) UpdatePersonDisabilities(disabilities []model.PersonDisabilityRequest, personId int, tx *gorm.DB) utils.Error {
	person, err := n.personRepo.GetPersonById(personId, tx)
	if err.Code != "" {
		return err
	}

	err = n.personDisabilityRepo.ClearPersonDisability(personId, tx)
	if err.Code != "" {
		return err
	}

	for _, disability := range disabilities {
		disability := model.PersonDisability{
			DisabilityId: disability.Id,
			PersonId:     personId,
			Acquired:     disability.Acquired,
		}

		err = n.personDisabilityRepo.UpsertPersonDisability(disability, tx)
		if err.Code != "" {
			return err
		}
	}

	err = n.personRepo.UpdatePerson(person, personId, tx)
	if err.Code != "" {
		return err
	}

	return utils.Error{}
}

func (n *personService) DeletePerson(personId int) utils.Error {
	person, err := n.personRepo.GetPersonById(personId, nil)
	if err.Code != "" {
		return err
	}

	err = n.personDisabilityRepo.ClearPersonDisability(personId, nil)
	if err.Code != "" {
		return err
	}

	err = n.personRepo.DeletePerson(personId)
	if err.Code != "" {
		return err
	}

	err = n.userRepo.DeleteUser(person.UserId)
	if err.Code != "" {
		return err
	}

	err = n.addressRepo.DeleteAddress(*person.AddressId)
	if err.Code != "" {
		return err
	}

	return utils.Error{}
}

func (n *personService) personToResponse(personResponse *model.PersonResponse, person model.Person) (model.PersonResponse, utils.Error) {
	user, err := n.userRepo.GetUserById(person.UserId)
	if err.Code != "" {
		return *personResponse, err
	}

	*personResponse = person.ToResponse(user)

	if person.AddressId != nil {
		address, err := n.addressRepo.GetAddressById(*person.AddressId)
		if err.Code != "" {
			return *personResponse, err
		}

		if address.Id != 0 {
			addressResponse := address.ToResponse()
			personResponse.Address = &addressResponse
		}
	}

	disabilities, err := n.personDisabilityRepo.GetPersonDisabilities(person.Id)
	if err.Code != "" {
		return *personResponse, err
	}

	if len(disabilities) > 0 {
		var disabilitiesResponse []model.PersonDisabilityResponse

		for _, disability := range disabilities {
			disabilityResponse := disability.ToResponse()
			disabilitiesResponse = append(disabilitiesResponse, disabilityResponse)
		}

		personResponse.Disabilities = &disabilitiesResponse
	}

	userConfig := model.DefaultConfig

	if user.ConfigUrl != "" {
		configService := NewConfigService(n.userRepo)
		userConfig, err = configService.GetUserConfig(user.ConfigUrl)
		if err.Code != "" {
			return *personResponse, err
		}
	}

	personResponse.User.Config = userConfig

	return *personResponse, utils.Error{}
}
