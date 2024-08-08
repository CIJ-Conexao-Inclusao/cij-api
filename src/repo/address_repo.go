package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type AddressRepo interface {
	BaseRepoMethods

	GetAddressById(id int) (model.Address, utils.Error)
	UpsertAddress(address model.Address, tx *gorm.DB) (int, utils.Error)
	DeleteAddress(id int) utils.Error
}

type addressRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) AddressRepo {
	repo := &addressRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func addressRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.AddressErrorType, code)

	return utils.NewError(message, errorCode)
}

func (n *addressRepo) GetAddressById(id int) (model.Address, utils.Error) {
	var address model.Address

	err := n.db.Model(model.Address{}).Where("id = ?", id).Find(&address).Error
	if err != nil {
		return address, addressRepoError("failed to get the address", "01")
	}

	return address, utils.Error{}
}

func (n *addressRepo) UpsertAddress(address model.Address, tx *gorm.DB) (int, utils.Error) {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Save(&address).Error; err != nil {
		return 0, addressRepoError("failed to upsert the address", "02")
	}

	return address.Id, utils.Error{}
}

func (n *addressRepo) DeleteAddress(id int) utils.Error {
	if err := n.db.Model(model.Address{}).Where("id = ?", id).Unscoped().Delete(&model.Address{}).Error; err != nil {
		return addressRepoError("failed to delete the address", "03")
	}

	return utils.Error{}
}
