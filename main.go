package main

import (
	"cij_api/src/config"
	"cij_api/src/database"
	"cij_api/src/model"
	"cij_api/src/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

// @title CIJ Project API
// @version 1.0
// @description This is the API for the CIJ project
// @contact.name API Support
// @contact.email cauakathdev@gmail.com
// @host conexao-inclusao.com
// @BasePath /
func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load enviroment variables", err)
	}

	db := database.ConnectionDB(&loadConfig)

	if err := db.Migrator().DropTable(
		&model.User{},
		&model.Address{},
		&model.Person{},
		&model.Disability{},
		&model.PersonDisability{},
		&model.Company{},
		&model.News{},
		&model.Role{},
		&model.Activity{},
	); err != nil {
		log.Fatal("failed to drop tables", err)
	}

	log.Default().Print("Tables dropped")

	migrateDb(db)

	startServer(db)
}

func migrateDb(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Address{})
	db.AutoMigrate(&model.Person{})
	db.AutoMigrate(&model.Disability{})
	db.AutoMigrate(&model.PersonDisability{})
	db.AutoMigrate(&model.Company{})
	db.AutoMigrate(&model.News{})
	db.AutoMigrate(&model.Role{})
	db.AutoMigrate(&model.Activity{})

	createDefaultRoles(db)
	createDefaultDisabilities(db)
}

func createDefaultRoles(db *gorm.DB) {
	db.Exec("INSERT IGNORE INTO roles (name) VALUES ('person')")
	db.Exec("INSERT IGNORE INTO roles (name) VALUES ('company')")
	db.Exec("INSERT IGNORE INTO roles (name) VALUES ('admin')")
}

func createDefaultDisabilities(db *gorm.DB) error {
	disabilities := []model.Disability{
		{Category: "Visual", Description: "Baixa visão, dificuldade em enxergar a longa distância", Rate: 30},
		{Category: "Visual", Description: "Cegueira completa", Rate: 80},
		{Category: "Visual", Description: "Dificuldade em diferenciar cores", Rate: 20},
		{Category: "Hearing", Description: "Perda auditiva parcial em um ouvido", Rate: 25},
		{Category: "Hearing", Description: "Surdez total", Rate: 90},
		{Category: "Hearing", Description: "Sensibilidade a sons altos", Rate: 15},
		{Category: "Physical", Description: "Paralisia parcial nos membros inferiores", Rate: 60},
		{Category: "Physical", Description: "Dificuldade de mobilidade devido a esclerose", Rate: 75},
		{Category: "Physical", Description: "Limitação no movimento das articulações", Rate: 40},
		{Category: "Intellectual", Description: "Transtorno do espectro autista", Rate: 50},
		{Category: "Intellectual", Description: "Déficit de atenção e hiperatividade", Rate: 30},
		{Category: "Intellectual", Description: "Deficiência intelectual leve", Rate: 45},
		{Category: "Psychosocial", Description: "Transtorno de ansiedade", Rate: 25},
		{Category: "Psychosocial", Description: "Depressão grave", Rate: 70},
		{Category: "Psychosocial", Description: "Transtorno bipolar", Rate: 65},
		{Category: "Psychosocial", Description: "Transtorno de estresse pós-traumático (TEPT)", Rate: 60},
	}

	for _, disability := range disabilities {
		var existing model.Disability
		if err := db.Where("category = ? AND description = ?", disability.Category, disability.Description).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&disability).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func startServer(db *gorm.DB) {
	app := fiber.New()

	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
	}))

	routes := router.NewRouter(app, db)

	err := routes.Listen(":3040")
	if err != nil {
		panic(err)
	}
}
