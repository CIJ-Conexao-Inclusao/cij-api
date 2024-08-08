package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Id        int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Email     string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password  string `gorm:"type:varchar(255);not null" json:"password"`
	ConfigUrl string `gorm:"type:varchar(255);not null" json:"config_url"`
	RoleId    RoleId `gorm:"type:int;not null" json:"role_id"`
	Role      *Role
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id     int    `json:"id"`
	Email  string `json:"email"`
	Config Config `json:"config,omitempty"`
}

func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		Id:    u.Id,
		Email: u.Email,
	}
}
