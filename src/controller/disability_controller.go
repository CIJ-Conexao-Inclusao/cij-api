package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"cij_api/src/utils"
	"fmt"

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

func disabilityControllerError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ControllerErrorCode, utils.DisabilityErrorType, code)

	return utils.NewError(message, errorCode)
}

func (c *DisabilityController) CreateDisability(ctx *fiber.Ctx) error {
	fmt.Print("aaaaaaaaa")
	var disabilityRequest DisabilityPostParameters
	fmt.Print("declarou variavel")
	if err := ctx.BodyParser(&disabilityRequest); err != nil {
		fmt.Print("Error", err)
		return disabilityControllerError("failed to parse the request body", "01")
	}
	fmt.Println("passou do body parser", disabilityRequest.Disabilities)

	err := c.disabilityService.CreateDisability(disabilityRequest.Disabilities)
	if err.Code != "" {
		fmt.Print("aaa", err.Message)
		return err
	}

	return utils.Error{}
}
