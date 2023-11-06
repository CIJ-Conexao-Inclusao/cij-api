package controller

import (
	"cij_api/src/domain"
	"cij_api/src/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type CompanyController struct {
	companyService domain.CompanyService
}

func NewCompanyController(companyService domain.CompanyService) *CompanyController {
	return &CompanyController{
		companyService: companyService,
	}
}

func (n *CompanyController) CreateCompany(ctx *fiber.Ctx) error {
	var companyRequest model.Company
	var response model.Response

	if err := ctx.BodyParser(&companyRequest); err != nil {
		response = model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	if err := n.companyService.CreateCompany(companyRequest); err != nil {
		response = model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		StatusCode: http.StatusCreated,
		Message:    "success",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func (n *CompanyController) ListCompanies(ctx *fiber.Ctx) error {
	var response model.Response

	companies, err := n.companyService.ListCompanies()
	if err != nil {
		response = model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	companiesResponse := []model.CompanyResponse{}

	for _, company := range companies {
		companiesResponse = append(companiesResponse, company.ToResponse())
	}

	response = model.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       companiesResponse,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}