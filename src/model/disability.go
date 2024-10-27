package model

import (
	"gorm.io/gorm"
)

type Disability struct {
	*gorm.Model
	Id          int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Category    string `gorm:"type:varchar(200);not null;index" json:"category"`
	Description string `gorm:"type:varchar(200);not null;index" json:"description"`
	Rate        int    `gorm:"type:int;not null" json:"rate"`
	People      []PersonDisability
}

type DisabilityRequest struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Rate        int    `json:"rate"`
}

type PersonDisabilityRequest struct {
	Id       int  `json:"id"`
	Acquired bool `json:"acquired"`
}

type DisabilityResponse struct {
	Id          int    `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Rate        int    `json:"rate"`
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
