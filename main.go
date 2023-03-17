package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/config"
	"gitlab.com/katsuotz/skip-api/controller"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/service"
)

func main() {
	database := config.SetupDatabaseConnection()
	jwtService := service.NewJWTService()
	userRepository := repository.NewUserRepository(database)

	authController := controller.NewAuthController(userRepository, jwtService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		MaxAge:           86400,
		AllowCredentials: true,
	}))

	r.GET("/", authController.Tes)

	api := r.Group("/api")
	{
		api.POST("/login", authController.Login)
	}

	//v1mid := r.Group("/api/v1")
	//{
	//v1mid.PUT("/user", userController.updateUser)
	//v1mid.POST("/user/password", userController.updatePassword)
	//}

	err := r.Run(":9100")
	if err != nil {
		return
	}
}
