package database

import (
	"cij_api/src/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB(config *config.Config) *gorm.DB {
	dsn := config.DbConnection
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = client.Exec("CREATE DATABASE IF NOT EXISTS cij").Error
	if err != nil {
		panic("failed to create database cij")
	}

	err = client.Exec("USE cij").Error
	if err != nil {
		panic("failed to enter database cij")
	}

	fmt.Print("Database connected\n\n")

	return client
}
