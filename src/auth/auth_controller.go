package auth

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService    AuthService
	personService  service.PersonService
	companyService service.CompanyService
	addressService service.AddressService
	configService  service.ConfigService
}

type TokenRequest struct {
	Token string `json:"token"`
}

func NewAuthController(
	authService AuthService,
	personService service.PersonService,
	companyService service.CompanyService,
	addressService service.AddressService,
	configService service.ConfigService,
) *AuthController {
	return &AuthController{
		authService:    authService,
		personService:  personService,
		companyService: companyService,
		addressService: addressService,
		configService:  configService,
	}
}

// Login
// @Summary Do login.
// @Description do login and returns token.
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param credentials body model.Credentials true "Credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /login [post]
func (c *AuthController) Authenticate(ctx *fiber.Ctx) error {
	var credentials model.Credentials
	var response model.LoginResponse

	if err := ctx.BodyParser(&credentials); err != nil {
		response = model.LoginResponse{
			Message: err.Error(),
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	user, err := c.authService.Authenticate(credentials)
	if err.Code != "" {
		response = model.LoginResponse{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	var userConfig interface{}
	userConfig = model.DefaultConfig

	if user.ConfigUrl != "" {
		userConfig, err = c.configService.GetUserConfig(user.ConfigUrl)
		if err.Code != "" {
			response = model.LoginResponse{
				Message: err.Error(),
				Code:    err.Code,
			}

			return ctx.Status(http.StatusInternalServerError).JSON(response)
		}
	}

	token, err := c.authService.GenerateToken(user)
	if err.Code != "" {
		response = model.LoginResponse{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	userResponse := user.ToResponse()
	userResponse.Config = userConfig

	response = model.LoginResponse{
		Token:    token,
		UserInfo: userResponse,
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetUserData
// @Summary Get user information.
// @Description get user information by token.
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param token body TokenRequest true "Token"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /get-user-data [post]
func (c *AuthController) GetUserData(ctx *fiber.Ctx) error {
	var token TokenRequest
	var response model.LoginResponse

	if err := ctx.BodyParser(&token); err != nil {
		response = model.LoginResponse{
			Message: "token not found",
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	user, err := c.authService.GetUserData(token.Token)
	if err.Code != "" {
		response = model.LoginResponse{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	var userConfig interface{}
	userConfig = model.DefaultConfig

	if user.ConfigUrl != "" {
		userConfig, err = c.configService.GetUserConfig(user.ConfigUrl)
		if err.Code != "" {
			response = model.LoginResponse{
				Message: err.Error(),
				Code:    err.Code,
			}

			return ctx.Status(http.StatusInternalServerError).JSON(response)
		}
	}

	if user.RoleId == 2 {
		company, err := c.companyService.GetCompanyByUserId(user.Id)
		if err.Code != "" {
			response = model.LoginResponse{
				Message: err.Error(),
				Code:    err.Code,
			}

			return ctx.Status(http.StatusInternalServerError).JSON(response)
		}

		companyResponse := company.ToResponse(user)
		companyResponse.User.Config = userConfig

		if company.AddressId != nil {
			address, err := c.addressService.GetAddressById(*company.AddressId)
			if err.Code != "" {
				response = model.LoginResponse{
					Message: err.Error(),
					Code:    err.Code,
				}

				return ctx.Status(http.StatusInternalServerError).JSON(response)
			}

			if address.Id != 0 {
				addressResponse := address.ToResponse()
				companyResponse.Address = addressResponse
			}
		}

		response = model.LoginResponse{
			UserInfo: companyResponse,
		}

		return ctx.Status(http.StatusOK).JSON(response)
	} else {
		person, err := c.personService.GetPersonByUserId(user.Id)
		if err.Code != "" {
			response = model.LoginResponse{
				Message: err.Error(),
				Code:    err.Code,
			}

			return ctx.Status(http.StatusInternalServerError).JSON(response)
		}

		personResponse := person.ToResponse(user)
		personResponse.User.Config = userConfig

		if person.AddressId != nil {
			address, err := c.addressService.GetAddressById(*person.AddressId)
			if err.Code != "" {
				response = model.LoginResponse{
					Message: err.Error(),
					Code:    err.Code,
				}

				return ctx.Status(http.StatusInternalServerError).JSON(response)
			}

			if address.Id != 0 {
				addressResponse := address.ToResponse()
				personResponse.Address = &addressResponse
			}
		}

		response = model.LoginResponse{
			UserInfo: personResponse,
		}

		return ctx.Status(http.StatusOK).JSON(response)
	}
}
