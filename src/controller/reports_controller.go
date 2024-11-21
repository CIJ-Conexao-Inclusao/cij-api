package controller

import (
	"cij_api/src/enum"
	"cij_api/src/model"
	"cij_api/src/service"
	"cij_api/src/utils"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type ReportsController struct {
	reportsService service.ReportsService
}

func NewReportsController(reportsService service.ReportsService) *ReportsController {
	return &ReportsController{
		reportsService: reportsService,
	}
}

func reportsControllerError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ControllerErrorCode, utils.ReportsErrorType, code)

	return utils.NewError(message, errorCode)
}

// GetDisabilityTotals
// @Summary Get disability totals
// @Description Get disability totals
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} model.Response
// @Router /disabilities [get]
func (c *ReportsController) GetDisabilityTotals(ctx *fiber.Ctx) error {
	var response model.Response

	disabilityTotals, err := c.reportsService.GetDisabilityTotals()
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "Disability totals",
		Data:    disabilityTotals,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetDisabilityTotalsByNeighborhood
// @Summary Get disability totals by neighborhood
// @Description Get disability totals by neighborhood
// @Tags Reports
// @Accept json
// @Produce json
// @Param neighborhood path string true "Neighborhood"
// @Success 200 {object} model.Response
// @Router /disabilities/{neighborhood} [get]
func (c *ReportsController) GetDisabilityTotalsByNeighborhood(ctx *fiber.Ctx) error {
	var response model.Response
	encodedNeighborhood := ctx.Params("neighborhood")

	neighborhood, err := url.QueryUnescape(encodedNeighborhood)
	if err != nil {
		response = model.Response{
			Message: "Failed to decode neighborhood parameter",
			Code:    reportsControllerError("failed to decode neighborhood parameter", "01").GetCode(),
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	if neighborhood == "" {
		response = model.Response{
			Message: "neighborhood is required",
			Code:    reportsControllerError("neighborhood is required", "02").GetCode(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	disabilityTotals, reportErr := c.reportsService.GetDisabilityTotalsByNeighborhood(neighborhood)
	if reportErr.Code != "" {
		response = model.Response{
			Message: reportErr.Message,
			Code:    reportErr.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "Disability totals by neighborhood: " + neighborhood,
		Data:    disabilityTotals,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// CountActivitiesByPeriod
// @Summary Count activities by period
// @Description Count activities by period
// @Tags Reports
// @Accept json
// @Produce json
// @Param type path string true "Type"
// @Param period path string true "Period"
// @Success 200 {object} model.Response
// @Router /activities/{type}/{period} [get]
func (c *ReportsController) CountActivitiesByPeriod(ctx *fiber.Ctx) error {
	var response model.Response
	activityType := ctx.Params("type")
	periodString := ctx.Params("period")

	period := enum.GetPeriodFilterEnum(periodString)
	if period == "" {
		response = model.Response{
			Message: "Invalid period",
			Code:    reportsControllerError("invalid period", "04").GetCode(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	countActivitiesByPeriod, err := c.reportsService.CountActivitiesByPeriod(activityType, period)
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "Activities counted by period",
		Data:    countActivitiesByPeriod,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
