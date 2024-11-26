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

// CreateVacancy
// @Summary Create a vacancy
// @Description Create a vacancy
// @Tags Vacancies
// @Accept json
// @Produce json
// @Param vacancy body vacancy.VacancyRequest true "Vacancy"
// @Success 201 {object} model.Response
// @Router /vacancies [post]
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

// UpdateVacancy
// @Summary Update a vacancy
// @Description Update a vacancy
// @Tags Vacancies
// @Accept json
// @Produce json
// @Param page query string false "Page"
// @Param per_page query string false "Per Page"
// @Param company_id query string false "Company ID"
// @Param disability query string false "Disability"
// @Param area query string false "Area"
// @Param contract_type query string false "Contract Type"
// @Param search_text query string false "Search Text"
// @Success 200 {object} model.Response
// @Router /vacancies [get]
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

// GetVacancyById
// @Summary Get a vacancy by ID
// @Description Get a vacancy by ID
// @Tags Vacancies
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} model.Response
// @Router /vacancies/{id} [get]
func (v *VacancyController) GetVacancyById(ctx *fiber.Ctx) error {
	var response model.Response

	id, _ := strconv.Atoi(ctx.Params("id"))

	vacancy, err := v.vacancyService.GetVacancyById(id)

	if err.Message == "failed to get the vacancy" {
		response = model.Response{
			Message: err.Message,
		}

		return ctx.Status(fiber.StatusNotFound).JSON(response)
	}

	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "vacancy retrieved successfully",
		Data:    vacancy,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
