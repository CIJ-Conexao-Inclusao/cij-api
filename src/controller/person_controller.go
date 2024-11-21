package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"cij_api/src/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type PersonController struct {
	personService service.PersonService
}

func NewPersonController(personService service.PersonService) *PersonController {
	return &PersonController{
		personService: personService,
	}
}

func personControllerError(message string, code string, fields []model.Field) utils.Error {
	errorCode := utils.NewErrorCode(utils.ControllerErrorCode, utils.PersonErrorType, code)

	return utils.NewErrorWithFields(message, errorCode, fields)
}

// CreatePerson
// @Summary Create a new person.
// @Description create a new person and their user.
// @Tags People
// @Accept application/json
// @Produce json
// @Param person body model.PersonRequest true "Person"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} MessageResponse
// @Router /people [post]
func (n *PersonController) CreatePerson(ctx *fiber.Ctx) error {
	var personRequest model.PersonRequest
	var response model.Response

	if err := ctx.BodyParser(&personRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := validatePersonRequiredFields(personRequest); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
			Fields:  err.GetFields(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := utils.ValidateUser(personRequest.User); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
			Fields:  err.GetFields(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.validatePerson(personRequest); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
			Fields:  err.GetFields(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := utils.ValidateAddress(personRequest.Address); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
			Fields:  err.GetFields(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.validatePersonDisabilities(personRequest.Disabilities); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.personService.CreatePerson(personRequest); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// ListPeople
// @Summary List all registered people.
// @Description list all registered people and their users.
// @Tags People
// @Accept application/json
// @Produce json
// @Success 200 {array} model.PersonResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /people [get]
func (n *PersonController) ListPeople(ctx *fiber.Ctx) error {
	var response model.Response

	people, err := n.personService.ListPeople()
	if err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
		Data:    people,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetPerson
// @Summary Get a person by ID.
// @Description get a person by their ID.
// @Tags People
// @Accept application/json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} model.PersonResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /people/:id [get]
func (n *PersonController) GetPerson(ctx *fiber.Ctx) error {
	var response model.Response

	personId := ctx.Params("id")

	idInt, err := strconv.Atoi(personId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	person, errPerson := n.personService.GetPersonById(idInt)
	if errPerson.Code != "" {
		response = model.Response{
			Message: errPerson.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	if person.Id == 0 {
		response = model.Response{
			Message: "person not found",
		}

		return ctx.Status(http.StatusNotFound).JSON(response)
	}

	response = model.Response{
		Message: "success",
		Data:    person,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// UpdatePerson
// @Summary Update a person.
// @Description update an existent person and their user.
// @Tags People
// @Accept application/json
// @Produce json
// @Param person body model.PersonRequest true "Person"
// @Param id path string true "Person ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} MessageResponse
// @Router /people/:id [put]
func (n *PersonController) UpdatePerson(ctx *fiber.Ctx) error {
	var personRequest model.PersonRequest
	var response model.Response

	if err := ctx.BodyParser(&personRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	personId := ctx.Params("id")

	idInt, err := strconv.Atoi(personId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	person, errPerson := n.personService.GetPersonById(idInt)
	if errPerson.Code != "" {
		response = model.Response{
			Message: errPerson.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	if person.Id == 0 {
		response = model.Response{
			Message: "person not found",
		}

		return ctx.Status(http.StatusNotFound).JSON(response)
	}

	// TODO: Validate only the passed fields
	// if err := n.validatePerson(personRequest); err.Code != "" {
	// 	response = model.Response{
	// 		Message: err.Error(),
	// 		Code:    err.GetCode(),
	// 		Fields:  err.GetFields(),
	// 	}

	// 	return ctx.Status(http.StatusBadRequest).JSON(response)
	// }

	if err := n.personService.UpdatePerson(personRequest, idInt); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// UpdatePersonAddress
// @Summary Update a person address.
// @Description update an existent person address.
// @Tags People
// @Accept application/json
// @Produce json
// @Param address body model.AddressRequest true "Address"
// @Param id path string true "Person ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} MessageResponse
// @Router /people/:id/address [put]
func (n *PersonController) UpdatePersonAddress(ctx *fiber.Ctx) error {
	var addressRequest model.AddressRequest
	var response model.Response

	if err := ctx.BodyParser(&addressRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	personId := ctx.Params("id")

	idInt, err := strconv.Atoi(personId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	person, errPerson := n.personService.GetPersonById(idInt)
	if errPerson.Code != "" {
		response = model.Response{
			Message: errPerson.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	if person.Id == 0 {
		response = model.Response{
			Message: "person not found",
		}

		return ctx.Status(http.StatusNotFound).JSON(response)
	}

	if err := n.personService.UpdatePersonAddress(addressRequest, idInt, nil); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// UpdatePersonDisabilities
// @Summary Update a person disabilities.
// @Description update an existent person disabilities.
// @Tags People
// @Accept application/json
// @Produce json
// @Param disabilities body []model.DisabilityRequest true "Disabilities"
// @Param id path string true "Person ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} MessageResponse
// @Router /people/:id/disabilities [put]
func (n *PersonController) UpdatePersonDisabilities(ctx *fiber.Ctx) error {
	var disabilities []model.PersonDisabilityRequest
	var response model.Response

	if err := ctx.BodyParser(&disabilities); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	personId := ctx.Params("id")

	idInt, err := strconv.Atoi(personId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	person, errPerson := n.personService.GetPersonById(idInt)
	if errPerson.Code != "" {
		response = model.Response{
			Message: errPerson.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	if person.Id == 0 {
		response = model.Response{
			Message: "person not found",
		}

		return ctx.Status(http.StatusNotFound).JSON(response)
	}

	if err := n.validatePersonDisabilities(disabilities); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.personService.UpdatePersonDisabilities(disabilities, idInt, nil); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// DeletePerson
// @Summary Delete a person.
// @Description delete an existent person and their user.
// @Tags People
// @Accept application/json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} MessageResponse
// @Router /people/:id [delete]
func (n *PersonController) DeletePerson(ctx *fiber.Ctx) error {
	var response model.Response

	personId := ctx.Params("id")

	idInt, err := strconv.Atoi(personId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	person, errPerson := n.personService.GetPersonById(idInt)
	if errPerson.Code != "" {
		response = model.Response{
			Message: errPerson.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	if person.Id == 0 {
		response = model.Response{
			Message: "person not found",
		}

		return ctx.Status(http.StatusNotFound).JSON(response)
	}

	if err := n.personService.DeletePerson(idInt); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// UploadCurriculum
// @Summary Upload a person curriculum.
// @Description upload a curriculum for a person.
// @Tags People
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Person ID"
// @Param file formData file true "Curriculum"
// @Success 200 {object} model.Response
// @Failure 400 {object} utils.Error
// @Failure 500 {object} model.Response
// @Router /people/:id/curriculum [post]
func (n *PersonController) UploadCurriculum(ctx *fiber.Ctx) error {
	var response model.Response

	file, err := ctx.FormFile("file")
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	personId := ctx.Params("id")

	idInt, err := strconv.Atoi(personId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	person, errPerson := n.personService.GetPersonById(idInt)
	if errPerson.Code != "" {
		response = model.Response{
			Message: errPerson.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	if person.Id == 0 {
		response = model.Response{
			Message: "person not found",
		}

		return ctx.Status(http.StatusNotFound).JSON(response)
	}

	if err := n.personService.UploadCurriculum(*file, idInt); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func validatePersonRequiredFields(personRequest model.PersonRequest) utils.Error {
	fieldsWithErrors := []model.Field{}

	if personRequest.Name == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "name"})
	}

	if personRequest.Cpf == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "cpf"})
	}

	if personRequest.Phone == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "phone"})
	}

	if personRequest.Gender == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "gender"})
	}

	if len(fieldsWithErrors) > 0 {
		errorCode := utils.NewErrorCode(utils.ValidationErrorCode, utils.PersonErrorType, "01")

		return utils.NewErrorWithFields("required fields are missing", errorCode, fieldsWithErrors)
	}

	return utils.Error{}
}

func (c *PersonController) validatePerson(personRequest model.PersonRequest) utils.Error {
	fieldsWithErrors := []model.Field{}

	if len(personRequest.Cpf) != 11 {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "cpf", Value: "cpf must have 11 digits"})
	}

	person, err := c.personService.GetPersonByCpf(personRequest.Cpf)
	if err.Code != "" {
		return err
	}

	if person.Id != 0 {
		return personControllerError("cpf already registered", "02", nil)
	}

	if len(personRequest.Phone) != 13 {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "phone", Value: "phone must have 13 digits"})
	}

	if !personRequest.Gender.IsValid() {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "gender", Value: "gender is not valid"})
	}

	user, err := c.personService.GetUserByEmail(personRequest.User.Email)
	if err.Code != "" {
		return err
	}

	if user.Id != 0 {
		return personControllerError("email already registered", "03", nil)
	}

	if len(fieldsWithErrors) > 0 {
		errorCode := utils.NewErrorCode(utils.ValidationErrorCode, utils.PersonErrorType, "02")

		return utils.NewErrorWithFields("invalid fields", errorCode, fieldsWithErrors)
	}

	return utils.Error{}
}

func (n *PersonController) validatePersonDisabilities(disabiliesRequest []model.PersonDisabilityRequest) utils.Error {
	for _, disabilityRequest := range disabiliesRequest {
		disability, err := n.personService.GetDisabilityById(disabilityRequest.Id)
		if err.Code != "" {
			return err
		}

		if disability.Id == 0 {
			return personControllerError(fmt.Sprintf("disability with id %d not found", disabilityRequest.Id), "06", nil)
		}
	}

	return utils.Error{}
}
