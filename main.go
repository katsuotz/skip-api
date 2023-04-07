package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/config"
	"gitlab.com/katsuotz/skip-api/controller"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/router"
	"gitlab.com/katsuotz/skip-api/service"
	"os"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")

	database := config.SetupDatabaseConnection()
	jwtService := service.NewJWTService(database)
	userRepository := repository.NewUserRepository(database)
	profileRepository := repository.NewProfileRepository(database)
	jurusanRepository := repository.NewJurusanRepository(database)
	tahunAjarRepository := repository.NewTahunAjarRepository(database)
	kelasRepository := repository.NewKelasRepository(database)
	siswaRepository := repository.NewSiswaRepository(database)
	guruRepository := repository.NewGuruRepository(database)
	dataPoinRepository := repository.NewDataPoinRepository(database)
	poinSiswaRepository := repository.NewPoinSiswaRepository(database)
	poinLogRepository := repository.NewPoinLogRepository(database)
	settingRepository := repository.NewSettingRepository(database)

	authController := controller.NewAuthController(userRepository, jwtService)
	profileController := controller.NewProfileController(profileRepository)
	jurusanController := controller.NewJurusanController(jurusanRepository)
	tahunAjarController := controller.NewTahunAjarController(tahunAjarRepository)
	kelasController := controller.NewKelasController(kelasRepository)
	siswaController := controller.NewSiswaController(siswaRepository, poinLogRepository)
	guruController := controller.NewGuruController(guruRepository)
	dataPoinController := controller.NewDataPoinController(dataPoinRepository)
	poinSiswaController := controller.NewPoinSiswaController(poinSiswaRepository, poinLogRepository)
	settingController := controller.NewSettingController(settingRepository)
	infoController := controller.NewInfoController(poinLogRepository, poinSiswaRepository)

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
		dataPoinController,
		poinSiswaController,
		settingController,
		infoController,
		jwtService,
	)

	r.Init()

	err := app.Run(":9200")
	if err != nil {
		return
	}
}
