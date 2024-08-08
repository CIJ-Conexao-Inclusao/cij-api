package model

type Role struct {
	Id   int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Name string `gorm:"type:varchar(200);not null;unique" json:"name"`
}

type RoleId int

const (
	PersonRole  RoleId = 1
	CompanyRole RoleId = 2
	AdminRole   RoleId = 3
)
