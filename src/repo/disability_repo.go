package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"
	"fmt"

	"gorm.io/gorm"
)

type DisabilityRepo interface {
	BaseRepoMethods

	BatchInsertDisabilities(disabilities []*model.Disability) utils.Error
}

type disabilityRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewDisabilityRepo(db *gorm.DB) DisabilityRepo {
	repo := &disabilityRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func disabilityRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.DisabilityErrorType, code)

	return utils.NewError(message, errorCode)
}

func (d *disabilityRepo) BatchInsertDisabilities(disabilities []*model.Disability) utils.Error {
	if err := d.db.Create(&disabilities).Error; err != nil {
		fmt.Println("Error", err)
		return disabilityRepoError("failed to batch insert the disabilities", "02")
	}

	return utils.Error{}
}
