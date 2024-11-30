package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"cij_api/src/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CompanyController struct {
	companyService service.CompanyService
}

func NewCompanyController(companyService service.CompanyService) *CompanyController {
	return &CompanyController{
		companyService: companyService,
	}
}

func companyControllerError(message string, code string, fields []model.Field) utils.Error {
	errorCode := utils.NewErrorCode(utils.ControllerErrorCode, utils.CompanyErrorType, code)

	return utils.NewErrorWithFields(message, errorCode, fields)
}

// CreateCompany
// @Summary Create a new company.
// @Description create a new company and their user.
// @Tags Companies
// @Accept application/json
// @Produce json
// @Param company body model.CompanyRequest true "Company"
// @Param Authorization header string true "Token"
// @Success 200 {object} string "success"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /companies [post]
func (n *CompanyController) CreateCompany(ctx *fiber.Ctx) error {
	var companyRequest model.CompanyRequest
	var response model.Response

	if err := ctx.BodyParser(&companyRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := validateCompanyRequiredFields(companyRequest); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
			Fields:  err.Fields,
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.validateCompany(companyRequest); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
			Fields:  err.Fields,
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := utils.ValidateUser(companyRequest.User); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
			Fields:  err.Fields,
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := utils.ValidateAddress(companyRequest.Address); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
			Fields:  err.Fields,
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.companyService.CreateCompany(companyRequest); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// ListCompanies
// @Summary List all registered companies.
// @Description list all registered companies and their users.
// @Tags Companies
// @Accept application/json
// @Produce json
// @Success 200 {array} model.CompanyResponse
// @Failure 404 {object} string "not found"
// @Failure 500 {object} string "internal server error"
// @Router /companies [get]
func (n *CompanyController) ListCompanies(ctx *fiber.Ctx) error {
	var response model.Response

	companies, err := n.companyService.ListCompanies()
	if err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
		Data:    companies,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetCompany
// @Summary Get a company by ID.
// @Description get a company by ID and their user.
// @Tags Companies
// @Accept application/json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} model.CompanyResponse
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /companies/:id [get]
func (n *CompanyController) GetCompany(ctx *fiber.Ctx) error {
	var response model.Response

	companyId := ctx.Params("id")

	idInt, errConv := strconv.Atoi(companyId)
	if errConv != nil {
		response = model.Response{
			Message: errConv.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	company, err := n.companyService.GetCompanyById(idInt)
	if err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
		Data:    company.ToResponse(*company.User),
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// UpdateCompany
// @Summary Update a company.
// @Description update an existent company and their user.
// @Tags Companies
// @Accept application/json
// @Produce json
// @Param company body model.CompanyRequest true "Company"
// @Param id path string true "Company ID"
// @Param Authorization header string true "Token"
// @Success 200 {object} string "success"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /companies/:id [put]
func (n *CompanyController) UpdateCompany(ctx *fiber.Ctx) error {
	var companyRequest model.CompanyRequest
	var response model.Response

	if err := ctx.BodyParser(&companyRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	companyId := ctx.Params("id")

	idInt, err := strconv.Atoi(companyId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.companyService.UpdateCompany(companyRequest, idInt); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// DeleteCompany
// @Summary Delete a company.
// @Description delete an existent company and their user.
// @Tags Companies
// @Accept application/json
// @Produce json
// @Param id path string true "Company ID"
// @Param Authorization header string true "Token"
// @Success 200 {object} string "success"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /companies/:id [delete]
func (n *CompanyController) DeleteCompany(ctx *fiber.Ctx) error {
	var response model.Response

	companyId := ctx.Params("id")

	idInt, err := strconv.Atoi(companyId)
	if err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.companyService.DeleteCompany(idInt); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func validateCompanyRequiredFields(company model.CompanyRequest) utils.Error {
	fieldsWithError := []model.Field{}

	if company.Cnpj == "" {
		fieldsWithError = append(fieldsWithError, model.Field{Name: "cnpj"})
	}

	if company.Name == "" {
		fieldsWithError = append(fieldsWithError, model.Field{Name: "name"})
	}

	if company.Phone == "" {
		fieldsWithError = append(fieldsWithError, model.Field{Name: "phone"})
	}

	if len(fieldsWithError) > 0 {
		errorCode := utils.NewErrorCode(utils.ValidationErrorCode, utils.CompanyErrorType, "01")

		return utils.NewErrorWithFields("required fields are missing", errorCode, fieldsWithError)
	}

	return utils.Error{}
}

func (c *CompanyController) validateCompany(companyRequest model.CompanyRequest) utils.Error {
	fieldsWithError := []model.Field{}

	if len(companyRequest.Cnpj) != 14 {
		fieldsWithError = append(fieldsWithError, model.Field{Name: "cnpj", Value: "cnpj must have 14 digits"})
	}

	company, err := c.companyService.GetCompanyByCnpj(companyRequest.Cnpj)
	if err.Code != "" {
		return err
	}

	if company.Id != 0 {
		return companyControllerError("cnpj already registered", "01", nil)
	}

	companyUser, err := c.companyService.GetUserByEmail(companyRequest.User.Email)
	if err.Code != "" {
		return err
	}

	if companyUser.Id != 0 {
		return companyControllerError("email already registered", "02", nil)
	}

	if len(fieldsWithError) > 0 {
		errorCode := utils.NewErrorCode(utils.ValidationErrorCode, utils.CompanyErrorType, "02")

		return utils.NewErrorWithFields("invalid fields", errorCode, fieldsWithError)
	}

	return utils.Error{}
}
