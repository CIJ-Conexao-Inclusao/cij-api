package router

import (
	"cij_api/src/auth"
	"cij_api/src/controller"
	"cij_api/src/middleware"
	"cij_api/src/repo"
	vacancy "cij_api/src/repo/vacancy"
	"cij_api/src/service"
	"fmt"

	_ "cij_api/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router *fiber.App, db *gorm.DB) *fiber.App {
	userRepo := repo.NewUserRepo(db)
	activityRepo := repo.NewActivityRepo(db)

	addressRepo := repo.NewAddressRepo(db)
	addressService := service.NewAddressService(addressRepo)

	personDisabilityRepo := repo.NewPersonDisabilityRepo(db)

	personRepo := repo.NewPersonRepo(db)
	personService := service.NewPersonService(personRepo, userRepo, addressRepo, personDisabilityRepo, activityRepo)
	personController := controller.NewPersonController(personService)

	companyRepo := repo.NewCompanyRepo(db)
	companyService := service.NewCompanyService(companyRepo, userRepo, addressRepo, activityRepo)
	companyController := controller.NewCompanyController(companyService)

	newsRepo := repo.NewNewsRepo(db)
	newsService := service.NewNewsService(newsRepo)
	newsController := controller.NewNewsController(newsService)

	disabilityRepo := repo.NewDisabilityRepo(db)
	disabilityService := service.NewDisabilityService(disabilityRepo)
	disabilityController := controller.NewDisabilityController(disabilityService)

	configService := service.NewConfigService(userRepo)
	configController := controller.NewConfigController(configService)

	authService := auth.NewAuthService(userRepo, activityRepo)
	authController := auth.NewAuthController(*authService, personService, companyService, addressService, configService)

	activityService := service.NewActivityService(activityRepo)
	activityController := controller.NewActivityController(activityService)

	vacancyRepo := vacancy.NewVacancyRepo(db)
	vacancySkillsRepo := vacancy.NewSkillsRepo(db)
	vacancyRequirementsRepo := vacancy.NewRequirementsRepo(db)
	vacancyResponsabilitiesRepo := vacancy.NewResponsabilitiesRepo(db)
	vacancyDisabilitiesRepo := vacancy.NewVacancyDisabilityRepo(db)
	vacancyApplyRepo := vacancy.NewVacancyApplyRepo(db)

	vacancyService := service.NewVacancyService(
		vacancyRepo, vacancySkillsRepo, vacancyRequirementsRepo,
		vacancyResponsabilitiesRepo, vacancyDisabilitiesRepo, vacancyApplyRepo, personRepo,
		personDisabilityRepo,
	)
	vacancyController := controller.NewVacancyController(vacancyService, companyService)

	reportsService := service.NewReportsService(personDisabilityRepo, activityRepo)
	reportsController := controller.NewReportsController(reportsService)

	router.Get("/health", HealthCheck)

	router.Get("/swagger/*", swagger.HandlerDefault)

	router.Post("/login", authController.Authenticate)
	router.Post("/get-user-data", authController.GetUserData)

	api := router.Group("/people")
	{
		api.Get("/", personController.ListPeople)
		api.Get("/:id", personController.GetPerson)
		api.Post("/", personController.CreatePerson)

		api.Use(middleware.AuthUser)
		api.Put("/:id", personController.UpdatePerson)
		api.Put("/:id/address", personController.UpdatePersonAddress)
		api.Put("/:id/disabilities", personController.UpdatePersonDisabilities)
		api.Delete("/:id", personController.DeletePerson)
		api.Post("/:id/curriculum", personController.UploadCurriculum)
	}

	api = router.Group("/companies")
	{
		api.Get("/", companyController.ListCompanies)
		api.Get("/:id", companyController.GetCompany)

		api.Use(middleware.AuthAdmin)
		api.Post("/", companyController.CreateCompany)
		api.Put("/:id", companyController.UpdateCompany)
		api.Delete("/:id", companyController.DeleteCompany)
	}

	api = router.Group("/news")
	{
		api.Get("/", newsController.ListNews)
		api.Post("/", newsController.CreateNews)
	}

	api = router.Group("/config")
	{
		api.Put("/:email", configController.UpdateUserConfig)
	}

	api = router.Group("/disabilities")
	{
		api.Post("/", disabilityController.CreateDisability)
	}

	api = router.Group("/activities")
	{
		api.Get("/", activityController.GetActivitiesByTypeAndPeriod)

		api.Use(middleware.AuthAdmin)
		api.Post("/", activityController.CreateActivity)
	}

	api = router.Group("/vacancies")
	{
		api.Get("/", vacancyController.ListVacancies)
		api.Get("/:id", vacancyController.GetVacancyById)
		api.Post("/apply", vacancyController.CandidateApply)

		api.Use(middleware.AuthCompany)
		api.Post("/", vacancyController.CreateVacancy)
		api.Put("/:id", vacancyController.UpdateVacancy)
		api.Delete("/:id", vacancyController.DeleteVacancy)

		api.Get("/apply/:id", vacancyController.ListVacancyApplies)
		api.Patch("/apply/:id", vacancyController.UpdateVacancyApplyStatus)
	}

	api = router.Group("/reports")
	{
		api.Get("/disabilities", reportsController.GetDisabilityTotals)
		api.Get("/disabilities/:neighborhood", reportsController.GetDisabilityTotalsByNeighborhood)
		api.Get("/activities/:type/:period", reportsController.CountActivitiesByPeriod)
	}

	basePath := getBasePath()
	fmt.Printf("API Routes:\n")

	for _, r := range router.GetRoutes() {
		if (r.Method == "GET" || r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE") && r.Path != "/health" && r.Path != "/" {
			fullPath := basePath + r.Path
			paintMethod(r.Method)
			paintPath(fullPath)
		}
	}

	return router
}

func getBasePath() string {
	return "http://localhost:3040"
}

func paintMethod(method string) {
	switch method {
	case "GET":
		color.New(color.FgMagenta).Printf("%s ", method)
	case "POST":
		color.New(color.FgGreen).Printf("%s ", method)
	case "PUT":
		color.New(color.FgYellow).Printf("%s ", method)
	case "DELETE":
		color.New(color.FgRed).Printf("%s ", method)
	default:
		color.New(color.FgWhite).Printf("%s ", method)
	}
}

func paintPath(path string) {
	color.White("%s\n", path)
}

// HealthCheck
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
