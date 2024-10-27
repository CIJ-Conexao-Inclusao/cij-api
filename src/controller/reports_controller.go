package controller

import (
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
