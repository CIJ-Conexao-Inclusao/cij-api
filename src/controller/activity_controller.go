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

// @Summary Create a new activity
// @Description Create a new activity with the provided details
// @Tags activities
// @Accept json
// @Produce json
// @Param activity body model.ActivityRequest true "Activity Request"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /activities [post]
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
	}

	return ctx.Status(http.StatusCreated).JSON(response)
}

// @Summary Get activities by type and period
// @Description Retrieve activities filtered by type and date range
// @Tags activities
// @Accept json
// @Produce json
// @Param type query string true "Activity Type"
// @Param start_date query string true "Start Date"
// @Param end_date query string true "End Date"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /activities [get]
func (a *ActivityController) GetActivitiesByTypeAndPeriod(ctx *fiber.Ctx) error {
	activityType := ctx.Query("type")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	startDateInt, err := strconv.ParseInt(startDate, 10, 64)
	if err != nil {
		response := model.Response{
			Message: "Invalid start date",
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	endDateInt, err := strconv.ParseInt(endDate, 10, 64)
	if err != nil {
		response := model.Response{
			Message: "Invalid end date",
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	activities, activitiesError := a.activityService.GetActivitiesByTypeAndPeriod(activityType, startDateInt, endDateInt)
	if activitiesError.Code != "" {
		response := model.Response{
			Message: activitiesError.Error(),
			Code:    activitiesError.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response := model.Response{
		Message: "Activities retrieved successfully",
		Data:    activities,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
