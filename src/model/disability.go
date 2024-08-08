package model

import (
	"gorm.io/gorm"
)

type Disability struct {
	*gorm.Model
	Id          int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Category    string `gorm:"type:varchar(200);not null" json:"category"`
	Description string `gorm:"type:varchar(200);not null" json:"description"`
	Rate        string `gorm:"type:varchar(200);not null" json:"rate"`
	People      []PersonDisability
}

type DisabilityRequest struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Rate        string `json:"rate"`
}

type PersonDisabilityRequest struct {
	Id       int  `json:"id"`
	Acquired bool `json:"acquired"`
}

type DisabilityResponse struct {
	Id          int    `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Rate        string `json:"rate"`
	Acquired    bool   `json:"acquired"`
}

func (d *Disability) ToResponse() DisabilityResponse {
	return DisabilityResponse{
		Id:          d.Id,
		Category:    d.Category,
		Description: d.Description,
		Rate:        d.Rate,
	}
}

func (dr *DisabilityRequest) ToModel() Disability {
	return Disability{
		Category:    dr.Category,
		Description: dr.Description,
		Rate:        dr.Rate,
	}
}
