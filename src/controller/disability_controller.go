package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"

	"github.com/gofiber/fiber/v2"
)

type DisabilityController struct {
	disabilityService service.DisabilityService
}

type DisabilityPostParameters struct {
	Disabilities []model.DisabilityRequest `json:"disabilities"`
}

func NewDisabilityController(disabilityService service.DisabilityService) *DisabilityController {
	return &DisabilityController{
		disabilityService: disabilityService,
	}
}

// @Summary Create Disabilities
// @Description Create new disabilities
// @Tags disabilities
// @Accept json
// @Produce json
// @Param disabilities body DisabilityPostParameters true "List of disabilities"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /disabilities [post]
func (c *DisabilityController) CreateDisability(ctx *fiber.Ctx) error {
	var disabilityRequest DisabilityPostParameters
	var response model.Response

	if err := ctx.BodyParser(&disabilityRequest); err != nil {
		response := model.Response{
			Message: err.Error(),
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	if len(disabilityRequest.Disabilities) == 0 {
		response := model.Response{
			Message: "No disabilities to create",
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := c.disabilityService.CreateDisability(disabilityRequest.Disabilities)
	if err.Code != "" {
		response := model.Response{
			Message: err.Error(),
			Code:    err.GetCode(),
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = model.Response{
		Message: "Disabilities created successfully",
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
