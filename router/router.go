package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/controller"
	"gitlab.com/katsuotz/skip-api/service"
)

type Router struct {
	server            *gin.Engine
	AuthController    controller.AuthController
	ProfileController controller.ProfileController
	JurusanController controller.JurusanController
	JWTService        service.JWTService
}

func NewRouter(server *gin.Engine,
	authController controller.AuthController,
	profileController controller.ProfileController,
	jurusanController controller.JurusanController,
	jwtService service.JWTService,
) *Router {
	return &Router{
		server,
		authController,
		profileController,
		jurusanController,
		jwtService,
	}
}

func (r *Router) Init() {
	basePath := r.server.Group("/api")

	basePath.GET("/", r.AuthController.Tes)
	basePath.POST("/login", r.AuthController.Login)

	profile := basePath.Group("/profile", r.JWTService.GetUser)
	{
		profile.PATCH("/", r.ProfileController.UpdateProfile)
	}

	jurusan := basePath.Group("/jurusan", r.JWTService.GetUser)
	{
		jurusan.GET("/", r.JurusanController.GetJurusan)
		jurusan.POST("/", r.JurusanController.CreateJurusan)
		jurusan.PATCH("/:jurusan_id", r.JurusanController.UpdateJurusan)
		jurusan.DELETE("/:jurusan_id", r.JurusanController.DeleteJurusan)
	}
}
