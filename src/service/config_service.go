package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ConfigService interface {
	UploadUserConfig(email string, config *interface{}) utils.Error
	GetUserConfig(email string) (interface{}, utils.Error)
}

type configService struct {
	userRepo repo.UserRepo
}

func NewConfigService(userRepo repo.UserRepo) ConfigService {
	return &configService{
		userRepo: userRepo,
	}
}

func configServiceError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.ServiceErrorCode, utils.ConfigErrorType, code)

	return utils.NewError(message, errorCode)
}

func (s *configService) UploadUserConfig(email string, config *interface{}) utils.Error {
	userConfig, err := json.Marshal(model.DefaultConfig)
	if err != nil {
		return configServiceError("failed to marshall user config", "01")
	}

	if config != nil {
		userConfig, err = json.Marshal(*config)
		if err != nil {
			return configServiceError("failed to marshall user config", "02")
		}
	}

	err = os.WriteFile(email+".json", userConfig, 0644)
	if err != nil {
		fmt.Print("Error: ", err)
		return configServiceError("failed to write user config", "03")
	}

	userFile, err := os.Open(email + ".json")
	if err != nil {
		return configServiceError("failed to open user config", "04")
	}

	filesService := NewFilesService()
	fileUrl, err := filesService.UploadFile(userFile, "cij/user_config/"+email)
	if err != nil {
		return configServiceError("failed to upload user config", "05")
	}

	userFile.Close()
	err = os.Remove(email + ".json")
	if err != nil {
		return configServiceError("failed to remove user config", "06")
	}

	updateUserErr := s.userRepo.UpdateUserConfig(fileUrl, email)
	if updateUserErr.Code != "" {
		return updateUserErr
	}

	return utils.Error{}
}

func (s *configService) GetUserConfig(url string) (interface{}, utils.Error) {
	userConfig, err := http.Get(url)
	if err != nil {
		return nil, configServiceError("failed to get user config", "07")
	}

	defer userConfig.Body.Close()

	var config interface{}

	err = json.NewDecoder(userConfig.Body).Decode(&config)
	if err != nil {
		return nil, configServiceError("failed to decode user config", "08")
	}

	return config, utils.Error{}
}
