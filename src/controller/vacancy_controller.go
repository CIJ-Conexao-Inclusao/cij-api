package controller

import (
	"cij_api/src/enum"
	"cij_api/src/model"
	vacancy "cij_api/src/model/vacancy"
	"cij_api/src/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VacancyController struct {
	vacancyService service.VacancyService
}

func NewVacancyController(vacancyService service.VacancyService) VacancyController {
	return VacancyController{
		vacancyService: vacancyService,
	}
}

func (v *VacancyController) CreateVacancy(ctx *fiber.Ctx) error {
	var vacancyRequest vacancy.VacancyRequest
	var response model.Response

	if err := ctx.BodyParser(&vacancyRequest); err != nil {
		response = model.Response{
			Message: "failed to parse the request body",
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := v.vacancyService.CreateVacancy(vacancyRequest)
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "vacancy created successfully",
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (v *VacancyController) ListVacancies(ctx *fiber.Ctx) error {
	var response model.Response

	page, perPage, companyId, disability := ctx.Query("page"), ctx.Query("per_page"), ctx.Query("company_id"), ctx.Query("disability")
	area, contractType, searchText := ctx.Query("area"), ctx.Query("contract_type"), ctx.Query("search_text")

	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	if perPageInt == 0 {
		perPageInt = 10
	}

	companyIdInt, _ := strconv.Atoi(companyId)

	vacancies, err := v.vacancyService.ListVacancies(pageInt, perPageInt, companyIdInt, disability, area, enum.VacancyContractType(contractType), searchText)
	if err.Code != "" {
		response := model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "vacancies listed successfully",
		Data:    vacancies,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
