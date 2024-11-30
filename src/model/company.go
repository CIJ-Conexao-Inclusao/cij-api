package model

import "gorm.io/gorm"

type Company struct {
	*gorm.Model
	Id        int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Name      string `gorm:"type:varchar(200);not null" json:"name"`
	Cnpj      string `gorm:"type:char(14);not null;unique" json:"cnpj"`
	Phone     string `gorm:"type:char(13);not null" json:"phone"`
	UserId    int    `gorm:"type:int;not null;unique" json:"user_id"`
	AddressId *int   `gorm:"type:int;not null;unique" json:"address_id"`
	User      *User
	Address   *Address
}

type CompanyRequest struct {
	Name    string         `json:"name"`
	Cnpj    string         `json:"cnpj"`
	Phone   string         `json:"phone"`
	User    UserRequest    `json:"user"`
	Address AddressRequest `json:"address"`
}

type CompanyResponse struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Cnpj    string          `json:"cnpj"`
	Phone   string          `json:"phone"`
	User    UserResponse    `json:"user"`
	Address AddressResponse `json:"address"`
}

func (c *Company) ToResponse(user User) CompanyResponse {
	return CompanyResponse{
		Id:      c.Id,
		Name:    c.Name,
		Cnpj:    c.Cnpj,
		Phone:   c.Phone,
		User:    user.ToResponse(),
		Address: c.Address.ToResponse(),
	}
}

func (c *CompanyRequest) ToModel(user User) Company {
	return Company{
		Name:   c.Name,
		Cnpj:   c.Cnpj,
		Phone:  c.Phone,
		UserId: user.Id,
	}
}

func (c *CompanyRequest) ToUser() User {
	return User{
		Email:    c.User.Email,
		Password: c.User.Password,
	}
}

func (c *CompanyRequest) ToAddress() Address {
	return c.Address.ToModel()
}
