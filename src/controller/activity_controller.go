package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ActivityController struct {
	activityService service.ActivityService
}

func NewActivityController(activityService service.ActivityService) ActivityController {
	return ActivityController{
		activityService: activityService,
	}
}

func (a *ActivityController) CreateActivity(ctx *fiber.Ctx) error {
	var activityRequest model.ActivityRequest
	var response model.Response

	if err := ctx.BodyParser(&activityRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	activity := activityRequest.ToModel()

	if err := a.activityService.CreateActivity(activity); err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "Activity created successfully",
		Data:    activity,
	}

	return ctx.Status(http.StatusCreated).JSON(response)
}

func (a *ActivityController) GetActivitiesByTypeAndPeriod(ctx *fiber.Ctx) error {
	activityType := ctx.Params("type")
	startDate := ctx.Params("start_date")
	endDate := ctx.Params("end_date")

	startDateInt, err := strconv.ParseInt(startDate, 10, 64)

	activities, err := a.activityService.GetActivitiesByTypeAndPeriod(activityType, startDate, endDate)
	if err.Code != "" {
		response := model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response := model.Response{
		Message: "Activities retrieved successfully",
		Data:    activities,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
