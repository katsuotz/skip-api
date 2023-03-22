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
	kelasRepository := repository.NewKelasRepository(database)
	siswaRepository := repository.NewSiswaRepository(database)
	guruRepository := repository.NewGuruRepository(database)
	dataScoreRepository := repository.NewDataScoreRepository(database)
	scoreSiswaRepository := repository.NewScoreSiswaRepository(database)

	authController := controller.NewAuthController(userRepository, jwtService)
	profileController := controller.NewProfileController(profileRepository)
	jurusanController := controller.NewJurusanController(jurusanRepository)
	tahunAjarController := controller.NewTahunAjarController(tahunAjarRepository)
	kelasController := controller.NewKelasController(kelasRepository)
	siswaController := controller.NewSiswaController(siswaRepository)
	guruController := controller.NewGuruController(guruRepository)
	dataScoreController := controller.NewDataScoreController(dataScoreRepository)
	scoreSiswaController := controller.NewScoreSiswaController(scoreSiswaRepository)

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
		kelasController,
		guruController,
		siswaController,
		dataScoreController,
		scoreSiswaController,
		jwtService,
	)

	r.Init()

	err := app.Run(":9100")
	if err != nil {
		return
	}
}
