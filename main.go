package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/config"
	"gitlab.com/katsuotz/skip-api/controller"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/router"
	"gitlab.com/katsuotz/skip-api/service"
)

func main() {
	database := config.SetupDatabaseConnection()
	jwtService := service.NewJWTService()
	userRepository := repository.NewUserRepository(database)
	profileRepository := repository.NewProfileRepository(database)
	jurusanRepository := repository.NewJurusanRepository(database)
	tahunAjarRepository := repository.NewTahunAjarRepository(database)

	authController := controller.NewAuthController(userRepository, jwtService)
	profileController := controller.NewProfileController(profileRepository, jwtService)
	jurusanController := controller.NewJurusanController(jurusanRepository, jwtService)
	tahunAjarController := controller.NewTahunAjarController(tahunAjarRepository, jwtService)

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		MaxAge:           86400,
		AllowCredentials: true,
	}))

	r := router.NewRouter(
		app,
		authController,
		profileController,
		jurusanController,
		tahunAjarController,
		jwtService,
	)

	r.Init()

	err := app.Run(":9100")
	if err != nil {
		return
	}
}
