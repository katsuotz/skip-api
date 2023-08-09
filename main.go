package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	sitiDatabase := config.SitiDatabaseConnection()

	jwtService := service.NewJWTService(database)
	userRepository := repository.NewUserRepository(database)
	profileRepository := repository.NewProfileRepository(database)
	jurusanRepository := repository.NewJurusanRepository(database)
	tahunAjarRepository := repository.NewTahunAjarRepository(database)
	kelasRepository := repository.NewKelasRepository(database)
	siswaRepository := repository.NewSiswaRepository(database)
	pegawaiRepository := repository.NewPegawaiRepository(database)
	dataPoinRepository := repository.NewDataPoinRepository(database)
	poinSiswaRepository := repository.NewPoinSiswaRepository(database)
	poinLogRepository := repository.NewPoinLogRepository(database)
	loginLogRepository := repository.NewLoginLogRepository(database)
	settingRepository := repository.NewSettingRepository(database)
	syncRepository := repository.NewSyncRepository(database, sitiDatabase)

	authController := controller.NewAuthController(userRepository, loginLogRepository, jwtService)
	profileController := controller.NewProfileController(profileRepository)
	jurusanController := controller.NewJurusanController(jurusanRepository)
	tahunAjarController := controller.NewTahunAjarController(tahunAjarRepository)
	kelasController := controller.NewKelasController(kelasRepository)
	siswaController := controller.NewSiswaController(siswaRepository, poinLogRepository)
	pegawaiController := controller.NewPegawaiController(pegawaiRepository)
	dataPoinController := controller.NewDataPoinController(dataPoinRepository)
	poinSiswaController := controller.NewPoinSiswaController(poinSiswaRepository, poinLogRepository)
	poinLogController := controller.NewPoinLogController(poinLogRepository)
	settingController := controller.NewSettingController(settingRepository)
	infoController := controller.NewInfoController(poinLogRepository, poinSiswaRepository)
	fileController := controller.NewFileController()
	syncController := controller.NewSyncController(syncRepository)

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
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
		pegawaiController,
		siswaController,
		dataPoinController,
		poinSiswaController,
		poinLogController,
		settingController,
		infoController,
		fileController,
		syncController,
		jwtService,
	)

	r.Init()

	err := godotenv.Load()
	if err != nil {
		panic("File env not found")
	}
	port := os.Getenv("PORT")

	err = app.Run(":" + port)
	if err != nil {
		return
	}
}
