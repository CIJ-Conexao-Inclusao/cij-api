package auth

import (
	"cij_api/src/config"
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	userRepo repo.UserRepo
}

func NewAuthService(userRepo repo.UserRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func authServiceError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ServiceErrorCode, utils.UserErrorType, code)

	return utils.NewError(message, errorCode)
}

func getSecretKey() ([]byte, utils.Error) {
	loadConfig, err := config.LoadConfig("../")
	if err != nil {
		return nil, authServiceError("failed to load config", "01")
	}

	return []byte(loadConfig.SecretKey), utils.Error{}
}

func (s *AuthService) GenerateToken(user model.User) (string, utils.Error) {
	secretKey, err := getSecretKey()
	if err.Code != "" {
		return "", err
	}

	claims := &jwt.MapClaims{
		"exp":   jwt.TimeFunc().Add(time.Hour * 24).Unix(),
		"role":  user.Role.Name,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, tokenError := token.SignedString(secretKey)
	if tokenError != nil {
		return "", authServiceError("failed to generate token", "02")
	}

	return tokenString, utils.Error{}
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	secret, err := getSecretKey()
	if err.Code != "" {
		return nil, err
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func (s *AuthService) Authenticate(credentials model.Credentials) (model.User, utils.Error) {
	var user model.User

	user, err := s.userRepo.GetUserByEmail(credentials.Email)
	if err.Code != "" {
		return user, err
	}

	if user.Email == "" {
		return user, authServiceError("user with this email not found", "03")
	}

	if !user.ValidatePassword(credentials.Password) {
		return user, authServiceError("invalid password", "04")
	}

	return user, utils.Error{}
}

func (s *AuthService) GetUserData(token string) (model.User, utils.Error) {
	var user model.User

	tokenData, err := ValidateToken(token)
	if err != nil {
		return user, authServiceError("failed to validate token", "05")
	}

	claims := tokenData.Claims.(jwt.MapClaims)
	tokenEmail := claims["email"].(string)

	user, userError := s.userRepo.GetUserByEmail(tokenEmail)
	if userError.Code != "" {
		return user, userError
	}

	if user.Email == "" {
		return user, authServiceError("user with this email not found", "06")
	}

	return user, utils.Error{}
}
