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
	JWTService        service.JWTService
}

func NewRouter(server *gin.Engine,
	authController controller.AuthController,
	profileController controller.ProfileController,
	jwtService service.JWTService,
) *Router {
	return &Router{
		server,
		authController,
		profileController,
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
}
