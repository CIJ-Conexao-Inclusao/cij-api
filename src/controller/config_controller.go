package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"cij_api/src/utils"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

type ConfigController struct {
	configService service.ConfigService
}

func NewConfigController(configService service.ConfigService) *ConfigController {
	return &ConfigController{
		configService: configService,
	}
}

// GetUserConfig godoc
// @Summary Get user config
// @Description Get user config
// @Tags config
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} interface{}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /config/{email} [get]
func (c *ConfigController) UpdateUserConfig(ctx *fiber.Ctx) error {
	var configRequest interface{}
	var response model.Response

	if err := ctx.BodyParser(&configRequest); err != nil {
		response = model.Response{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	userEmail := ctx.Params("email")

	// err := validateUserConfig(configRequest)
	// if err.Message != "" {
	// 	response = model.Response{
	// 		Message: err.Message,
	// 		Fields:  err.Fields,
	// 	}

	// 	return ctx.Status(http.StatusBadRequest).JSON(response)
	// }

	err := c.configService.UploadUserConfig(userEmail, &configRequest)
	if err.Code != "" {
		response = model.Response{
			Message: err.Message,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	response = model.Response{
		Message: "User config updated successfully",
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func validateUserConfig(config model.Config) utils.Error {
	errorFields := []model.Field{}

	const MAX_FONT_SIZE = 30
	const MIN_FONT_SIZE = 14

	if !config.Theme.IsValid() {
		errorFields = append(errorFields, model.Field{Name: "theme", Value: "theme must be 'light', 'dark' or 'system'"})
	}

	if config.FontSize < MIN_FONT_SIZE || config.FontSize > MAX_FONT_SIZE {
		errorFields = append(errorFields, model.Field{Name: "font_size", Value: "font_size must be between 14 and 30"})
	}

	if !config.ColorBlindness.IsValid() {
		errorFields = append(errorFields, model.Field{Name: "color_blindness", Value: "color_blindness must be 'normal', 'protanopia', 'deuteranopia' or 'tritanopia'"})
	}

	errorColors := validateConfigColors(config)
	fmt.Print(errorColors)
	errorFields = append(errorFields, errorColors...)
	fmt.Print(errorFields)

	if len(errorFields) > 0 {
		return utils.Error{
			Message: "Invalid user config",
			Fields:  errorFields,
		}
	}

	return utils.Error{}
}

func validateConfigColors(config model.Config) []model.Field {
	errorColors := []model.Field{}
	configColors := config.SystemColors

	primaryColors := map[string]string{
		"primary_color":        configColors.PrimaryColors.PrimaryColor,
		"secondary_color":      configColors.PrimaryColors.SecondaryColor,
		"font_color":           configColors.PrimaryColors.FontColor,
		"secondary_font_color": configColors.PrimaryColors.SecondaryFontColor,
		"input_color":          configColors.PrimaryColors.InputColor,
		"background_color":     configColors.PrimaryColors.BackgroundColor,
	}

	for colorCategory, colorValue := range primaryColors {
		if !isHexString(colorValue) {
			errorColors = append(errorColors, model.Field{Name: colorCategory, Value: colorValue + " must be a valid hex color"})
		}
	}

	for chartCategory, chartColor := range configColors.ChartColors {
		if !isHexString(chartColor) {
			errorColors = append(errorColors, model.Field{Name: string(chartCategory), Value: chartColor + " must be a valid hex color"})
		}
	}

	return errorColors
}

func isHexString(s string) bool {
	hexPattern := "^#([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$" // Hex color pattern

	ok, err := regexp.MatchString(hexPattern, s)
	if err != nil {
		return false
	}

	return ok
}
