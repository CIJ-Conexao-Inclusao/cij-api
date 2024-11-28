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
	companyService service.CompanyService
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

	if err := v.validateVacancy(vacancyRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
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

func (v *VacancyController) CandidateApply(ctx *fiber.Ctx) error {
	var response model.Response
	var vacancyApplyRequest vacancy.VacancyApplyRequest

	if err := ctx.BodyParser(&vacancyApplyRequest); err != nil {
		response = model.Response{
			Message: "failed to parse the request body",
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := v.vacancyService.CandidateApplyVacancy(vacancyApplyRequest.CandidateId, vacancyApplyRequest.VacancyId)
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "candidate applied to the vacancy successfully",
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (v *VacancyController) ListVacancyApplies(ctx *fiber.Ctx) error {
	var response model.Response

	vacancyId, _ := strconv.Atoi(ctx.Params("id"))
	vacancyApplies, err := v.vacancyService.GetVacancyAppliesByVacancyId(vacancyId)
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "vacancy applies listed successfully",
		Data:    vacancyApplies,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (v *VacancyController) UpdateVacancyApplyStatus(ctx *fiber.Ctx) error {
	var response model.Response

	vacancyApplyId, _ := strconv.Atoi(ctx.Params("id"))
	status := ctx.Query("status")

	if !enum.VacancyApplyStatus(status).IsValid() {
		response = model.Response{
			Message: "invalid status. valid values are: 'applied', 'approved', 'rejected'",
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := v.vacancyService.UpdateVacancyApplyStatus(vacancyApplyId, enum.VacancyApplyStatus(status))
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "vacancy apply status updated successfully",
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (v *VacancyController) validateVacancy(vacancyRequest vacancy.VacancyRequest) error {
	if vacancyRequest.Code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "code is required")
	}

	if vacancyRequest.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "title is required")
	}

	if vacancyRequest.Description == "" {
		return fiber.NewError(fiber.StatusBadRequest, "description is required")
	}

	if vacancyRequest.Department == "" {
		return fiber.NewError(fiber.StatusBadRequest, "department is required")
	}

	if vacancyRequest.Section == "" {
		return fiber.NewError(fiber.StatusBadRequest, "section is required")
	}

	if vacancyRequest.Turn == "" {
		return fiber.NewError(fiber.StatusBadRequest, "turn is required")
	}

	if vacancyRequest.PublishDate == "" {
		return fiber.NewError(fiber.StatusBadRequest, "publish date is required")
	}

	if vacancyRequest.RegistrationDate == "" {
		return fiber.NewError(fiber.StatusBadRequest, "registration date is required")
	}

	if vacancyRequest.Area == "" {
		return fiber.NewError(fiber.StatusBadRequest, "area is required")
	}

	if len(vacancyRequest.Disabilities) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one disability is required")
	}

	if len(vacancyRequest.Skills) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one skill is required")
	}

	if len(vacancyRequest.Responsabilities) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one responsability is required")
	}

	if len(vacancyRequest.Requirements) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one requirement is required")
	}

	for _, requirement := range vacancyRequest.Requirements {
		if requirement.Type == "" {
			return fiber.NewError(fiber.StatusBadRequest, "requirement type is required")
		}

		if !requirement.Type.IsValid() {
			return fiber.NewError(fiber.StatusBadRequest, "invalid requirement type. valid values are: 'desirable', 'obligatory'")
		}

		if requirement.Requirement == "" {
			return fiber.NewError(fiber.StatusBadRequest, "requirement description is required")
		}
	}

	if vacancyRequest.ContractType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "contract type is required")
	}

	if !vacancyRequest.ContractType.IsValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contract type. valid values are: 'clt', 'pj', 'trainee'")
	}

	if vacancyRequest.CompanyId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "company ID is required")
	}

	company, err := v.companyService.GetCompanyById(vacancyRequest.CompanyId)
	if err.Code != "" {
		return fiber.NewError(fiber.StatusBadRequest, "failed to get the company")
	}

	if company.Id == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "company not found")
	}

	return nil
}
