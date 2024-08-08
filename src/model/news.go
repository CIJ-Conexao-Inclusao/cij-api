package model

import "gorm.io/gorm"

type News struct {
	*gorm.Model
	Id          int    `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Title       string `gorm:"type:varchar(200);not null" json:"title"`
	Description string `gorm:"type:text;not null" json:"description"`
	Banner      string `gorm:"type:text;" json:"banner"`
	Author      string `gorm:"type:varchar(200);not null" json:"author"`
	AuthorImage string `gorm:"type:text;" json:"author_image"`
	Date        string `gorm:"type:date;not null" json:"date"`
}

type NewsRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Date        string `json:"date"`
}

type NewsResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Date        string `json:"date"`
	Banner      string `json:"banner"`
	AuthorImage string `json:"author_image"`
}

func (n *News) ToResponse() NewsResponse {
	return NewsResponse{
		Id:          n.Id,
		Title:       n.Title,
		Description: n.Description,
		Author:      n.Author,
		Date:        n.Date,
		Banner:      n.Banner,
		AuthorImage: n.AuthorImage,
	}
}

func (n *NewsRequest) ToModel() News {
	return News{
		Title:       n.Title,
		Description: n.Description,
		Author:      n.Author,
		Date:        n.Date,
	}
}
