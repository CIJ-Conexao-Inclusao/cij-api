package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
)

type UserRepo interface {
	BaseRepoMethods

	CreateUser(createUser model.User, tx *gorm.DB) (int, utils.Error)
	ListUsers() ([]model.User, utils.Error)
	GetUserByEmail(email string) (model.User, utils.Error)
	GetUserById(id int) (model.User, utils.Error)
	UpdateUser(user model.User, userId int) utils.Error
	UpdateUserConfig(configUrl string, userEmail string) utils.Error
	DeleteUser(id int) utils.Error
}

type userRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	repo := &userRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func userRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.UserErrorType, code)

	return utils.NewError(message, errorCode)
}

func (n *userRepo) CreateUser(createUser model.User, tx *gorm.DB) (int, utils.Error) {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Create(&createUser).Error; err != nil {
		return 0, userRepoError("failed to create the user", "01")
	}

	return createUser.Id, utils.Error{}
}

func (n *userRepo) ListUsers() ([]model.User, utils.Error) {
	var users []model.User

	err := n.db.Model(model.User{}).Find(&users).Error
	if err != nil {
		return users, userRepoError("failed to list the users", "02")
	}

	return users, utils.Error{}
}

func (n *userRepo) GetUserByEmail(email string) (model.User, utils.Error) {
	var user model.User

	err := n.db.Model(model.User{}).Preload("Role").Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, userRepoError("failed to get the user", "03")
	}

	return user, utils.Error{}
}

func (n *userRepo) GetUserById(id int) (model.User, utils.Error) {
	var user model.User

	err := n.db.Model(model.User{}).Preload("Role").Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, userRepoError("failed to get the user", "04")
	}

	return user, utils.Error{}
}

func (n *userRepo) UpdateUser(user model.User, userId int) utils.Error {
	if err := n.db.Model(model.User{}).Where("id = ?", userId).Updates(user).Error; err != nil {
		return userRepoError("failed to update the user", "05")
	}

	return utils.Error{}
}

func (n *userRepo) UpdateUserConfig(configUrl string, userEmail string) utils.Error {
	if err := n.db.Model(model.User{}).Where("email = ?", userEmail).Update("config_url", configUrl).Error; err != nil {
		return userRepoError("failed to update the user config", "07")
	}

	return utils.Error{}
}

func (n *userRepo) DeleteUser(userId int) utils.Error {
	err := n.db.Model(model.User{}).Where("id = ?", userId).Unscoped().Delete(&model.User{}).Error
	if err != nil {
		return userRepoError("failed to delete the user", "06")
	}

	return utils.Error{}
}
