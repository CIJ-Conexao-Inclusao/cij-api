package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"
	"fmt"

	"gorm.io/gorm"
)

type DisabilityRepo interface {
	BaseRepoMethods

	CreateDisability(disability *model.Disability) utils.Error
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

func (d *disabilityRepo) CreateDisability(disability *model.Disability) utils.Error {
	if err := d.db.Create(disability).Error; err != nil {
		return disabilityRepoError("failed to create the disability", "01")
	}

	return utils.Error{}
}

func (d *disabilityRepo) BatchInsertDisabilities(disabilities []*model.Disability) utils.Error {
	sql := "INSERT INTO disabilities (id, category, description, rate) VALUES "

	for i, disability := range disabilities {
		sql += fmt.Sprintf("(NULL, '%s', '%s', '%s')", disability.Category, disability.Description, disability.Rate)

		if i < len(disabilities)-1 {
			sql += ", "
		}
	}

	if err := d.db.Exec(sql).Error; err != nil {
		fmt.Println("Error", err)
		return disabilityRepoError("failed to batch insert the disabilities", "02")
	}

	return utils.Error{}
}
