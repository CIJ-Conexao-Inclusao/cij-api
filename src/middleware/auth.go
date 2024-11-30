package middleware

import (
	"cij_api/src/auth"
	"cij_api/src/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	jwt.StandardClaims
}

const PERSON_ROLE = "person"
const COMPANY_ROLE = "company"
const ADMIN_ROLE = "admin"

func AuthUser(ctx *fiber.Ctx) error {
	var response model.Response

	token, err := Auth(ctx)
	if err.Message != "" {
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	claims := token.Claims.(jwt.MapClaims)
	tokenRole := claims["role"].(string)

	if tokenRole != PERSON_ROLE && tokenRole != ADMIN_ROLE {
		response = model.Response{
			Message: "role don't have permission",
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	return ctx.Next()
}

func AuthAdmin(ctx *fiber.Ctx) error {
	var response model.Response

	token, err := Auth(ctx)
	if err.Message != "" {
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	claims := token.Claims.(jwt.MapClaims)
	tokenRole := claims["role"].(string)

	if tokenRole != ADMIN_ROLE {
		response = model.Response{
			Message: "role don't have permission",
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	return ctx.Next()
}

func AuthCompany(ctx *fiber.Ctx) error {
	var response model.Response

	token, err := Auth(ctx)
	if err.Message != "" {
		return ctx.Status(http.StatusBadRequest).JSON(err)
	}

	claims := token.Claims.(jwt.MapClaims)
	tokenRole := claims["role"].(string)

	if tokenRole != COMPANY_ROLE && tokenRole != ADMIN_ROLE {
		response = model.Response{
			Message: "role don't have permission",
		}

		return ctx.Status(http.StatusBadRequest).JSON(response)
	}

	return ctx.Next()
}

func Auth(ctx *fiber.Ctx) (*jwt.Token, model.Response) {
	var response model.Response
	tokenParam := ctx.Get("Authorization")

	if tokenParam == "" {
		response = model.Response{
			Message: "token not found",
		}

		return nil, response
	}

	token, err := auth.ValidateToken(tokenParam)
	if err != nil {
		response = model.Response{
			Message: "token invalid or expired",
		}

		return nil, response
	}

	if !token.Valid {
		response = model.Response{
			Message: "token is not valid",
		}

		return nil, response
	}

	return token, model.Response{}
}
