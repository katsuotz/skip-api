package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/controller"
	"gitlab.com/katsuotz/skip-api/service"
)

type Router struct {
	server              *gin.Engine
	AuthController      controller.AuthController
	ProfileController   controller.ProfileController
	JurusanController   controller.JurusanController
	TahunAjarController controller.TahunAjarController
	KelasController     controller.KelasController
	JWTService          service.JWTService
}

func NewRouter(server *gin.Engine,
	authController controller.AuthController,
	profileController controller.ProfileController,
	jurusanController controller.JurusanController,
	tahunAjarController controller.TahunAjarController,
	kelasController controller.KelasController,
	jwtService service.JWTService,
) *Router {
	return &Router{
		server,
		authController,
		profileController,
		jurusanController,
		tahunAjarController,
		kelasController,
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

	tahunAjar := basePath.Group("/tahun-ajar", r.JWTService.GetUser)
	{
		tahunAjar.GET("/", r.TahunAjarController.GetTahunAjar)
		tahunAjar.POST("/", r.TahunAjarController.CreateTahunAjar)
		tahunAjar.PATCH("/:tahun_ajar_id", r.TahunAjarController.UpdateTahunAjar)
		tahunAjar.DELETE("/:tahun_ajar_id", r.TahunAjarController.DeleteTahunAjar)
		tahunAjar.PATCH("/:tahun_ajar_id/set-active", r.TahunAjarController.SetActiveTahunAjar)
	}

	kelas := basePath.Group("/kelas", r.JWTService.GetUser)
	{
		kelas.GET("/", r.KelasController.GetKelas)
		kelas.POST("/", r.KelasController.CreateKelas)
		kelas.PATCH("/:jurusan_id", r.KelasController.UpdateKelas)
		kelas.DELETE("/:jurusan_id", r.KelasController.DeleteKelas)
	}
}
