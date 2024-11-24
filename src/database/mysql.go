package database

import (
	"cij_api/src/config"
	"fmt"
	"log"

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

	createFunctionToNormalizeText(client)

	fmt.Print("Database connected\n\n")

	return client
}

func createFunctionToNormalizeText(client *gorm.DB) {
	var count int
	err := client.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = DATABASE() 
		AND ROUTINE_NAME = 'NormalizeText' 
		AND ROUTINE_TYPE = 'FUNCTION';
	`).Scan(&count).Error

	if err != nil {
		log.Fatal("Erro ao verificar existência da função:", err)
	}

	if count > 0 {
		return
	}

	clientNormalizerFunctionSql := `
		CREATE FUNCTION NormalizeText(text VARCHAR(255)) RETURNS VARCHAR(255)
		DETERMINISTIC
		BEGIN
			SET text = REPLACE(text, 'ã', 'a');
			SET text = REPLACE(text, 'á', 'a');
			SET text = REPLACE(text, 'à', 'a');
			SET text = REPLACE(text, 'â', 'a');
			SET text = REPLACE(text, 'é', 'e');
			SET text = REPLACE(text, 'ê', 'e');
			SET text = REPLACE(text, 'í', 'i');
			SET text = REPLACE(text, 'ó', 'o');
			SET text = REPLACE(text, 'õ', 'o');
			SET text = REPLACE(text, 'ô', 'o');
			SET text = REPLACE(text, 'ú', 'u');
			SET text = REPLACE(text, 'ç', 'c');
			RETURN text;
		END;
	`

	err = client.Exec(clientNormalizerFunctionSql).Error
	if err != nil {
		panic("failed to create normalizer function")
	}
}
