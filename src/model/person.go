package model

import (
	"cij_api/src/enum"

	"gorm.io/gorm"
)

type Person struct {
	*gorm.Model
	Id        int             `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Name      string          `gorm:"type:varchar(200);not null" json:"name"`
	Cpf       string          `gorm:"type:char(11);not null;unique" json:"cpf"`
	Phone     string          `gorm:"type:char(13);not null" json:"phone"`
	Gender    enum.GenderEnum `gorm:"type:char(6);not null" json:"gender"`
	UserId    int             `gorm:"type:int;not null;unique" json:"user_id"`
	AddressId *int            `gorm:"type:int;unique" json:"address_id"`
	Address   *Address
	User      *User
}

type PersonRequest struct {
	Name    string          `json:"name"`
	Cpf     string          `json:"cpf"`
	Phone   string          `json:"phone"`
	Gender  enum.GenderEnum `json:"gender"`
	User    UserRequest     `json:"user"`
	Address AddressRequest  `json:"address"`
}

type PersonResponse struct {
	Id      int              `json:"id"`
	Name    string           `json:"name"`
	Cpf     string           `json:"cpf"`
	Phone   string           `json:"phone"`
	Gender  enum.GenderEnum  `json:"gender"`
	User    UserResponse     `json:"user"`
	Address *AddressResponse `json:"address,omitempty"`
}

func (p *Person) ToResponse(user User) PersonResponse {
	return PersonResponse{
		Id:     p.Id,
		Name:   p.Name,
		Cpf:    p.Cpf,
		Phone:  p.Phone,
		Gender: p.Gender,
		User:   user.ToResponse(),
	}
}

func (p *PersonRequest) ToModel(user User) Person {
	return Person{
		Name:   p.Name,
		Cpf:    p.Cpf,
		Phone:  p.Phone,
		Gender: p.Gender,
		UserId: user.Id,
	}
}

func (p *PersonRequest) ToUser() User {
	return User{
		Email:    p.User.Email,
		Password: p.User.Password,
	}
}

func (p *PersonRequest) ToAddress() Address {
	return Address{
		Street:       p.Address.Street,
		Number:       p.Address.Number,
		Neighborhood: p.Address.Neighborhood,
		City:         p.Address.City,
		State:        p.Address.State,
		ZipCode:      p.Address.ZipCode,
		Complement:   p.Address.Complement,
	}
}
