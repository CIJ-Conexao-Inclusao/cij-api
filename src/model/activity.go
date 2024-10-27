package model

import "gorm.io/gorm"

type Activity struct {
	*gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Type        string `gorm:"type:varchar(100);not null" json:"type"`
	Description string `gorm:"type:varchar(200);not null" json:"description"`
	Actor       string `gorm:"type:varchar(100);not null" json:"actor"`
}

type ActivityResponse struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Actor       string `json:"actor"`
}

type ActivityRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Actor       string `json:"actor"`
}

func (a *ActivityRequest) ToModel() *Activity {
	return &Activity{
		Type:        a.Type,
		Description: a.Description,
		Actor:       a.Actor,
	}
}

func (a *Activity) ToResponse() *ActivityResponse {
	return &ActivityResponse{
		ID:          a.ID,
		Type:        a.Type,
		Description: a.Description,
		Actor:       a.Actor,
	}
}
