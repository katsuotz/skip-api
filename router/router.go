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
	GuruController      controller.GuruController
	SiswaController     controller.SiswaController
	DataPoinController  controller.DataPoinController
	PoinSiswaController controller.PoinSiswaController
	SettingController   controller.SettingController
	JWTService          service.JWTService
}

func NewRouter(server *gin.Engine,
	authController controller.AuthController,
	profileController controller.ProfileController,
	jurusanController controller.JurusanController,
	tahunAjarController controller.TahunAjarController,
	kelasController controller.KelasController,
	guruController controller.GuruController,
	siswaController controller.SiswaController,
	dataPoinController controller.DataPoinController,
	poinSiswaController controller.PoinSiswaController,
	settingController controller.SettingController,
	jwtService service.JWTService,
) *Router {
	return &Router{
		server,
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
		jwtService,
	}
}

func (r *Router) Init() {
	basePath := r.server.Group("/api")

	basePath.GET("", r.AuthController.Tes)

	guestPath := basePath.Group("", r.JWTService.IsGuest)
	{
		guestPath.POST("login", r.AuthController.Login)
	}

	loggedPath := basePath.Group("", r.JWTService.IsLoggedIn)
	{
		loggedPath.GET("me", r.ProfileController.GetMyProfile)
		loggedPath.PATCH("profile", r.ProfileController.UpdateProfile)

		jurusan := loggedPath.Group("jurusan", r.JWTService.IsAdmin)
		{
			jurusan.GET("", r.JurusanController.GetJurusan)
			jurusan.POST("", r.JurusanController.CreateJurusan)
			jurusan.PATCH(":jurusan_id", r.JurusanController.UpdateJurusan)
			jurusan.DELETE(":jurusan_id", r.JurusanController.DeleteJurusan)
		}

		tahunAjar := loggedPath.Group("tahun-ajar", r.JWTService.IsAdmin)
		{
			tahunAjar.GET("", r.TahunAjarController.GetTahunAjar)
			tahunAjar.POST("", r.TahunAjarController.CreateTahunAjar)
			tahunAjar.PATCH(":tahun_ajar_id", r.TahunAjarController.UpdateTahunAjar)
			tahunAjar.DELETE(":tahun_ajar_id", r.TahunAjarController.DeleteTahunAjar)
			tahunAjar.PATCH(":tahun_ajar_id/set-active", r.TahunAjarController.SetActiveTahunAjar)
		}

		kelas := loggedPath.Group("kelas")
		{
			kelas.GET("", r.KelasController.GetKelas)
			kelas.GET(":kelas_id", r.KelasController.GetKelasByID)

			kelasData := kelas.Group("", r.JWTService.IsAdmin)
			{
				kelasData.POST("", r.KelasController.CreateKelas)
				kelasData.PATCH(":kelas_id", r.KelasController.UpdateKelas)
				kelasData.DELETE(":kelas_id", r.KelasController.DeleteKelas)
				kelasData.POST(":kelas_id/add-siswa", r.KelasController.AddSiswaToKelas)
				kelasData.DELETE(":kelas_id/remove-siswa", r.KelasController.RemoveSiswaFromKelas)
			}
		}

		guru := loggedPath.Group("guru", r.JWTService.IsAdmin)
		{
			guru.GET("", r.GuruController.GetGuru)
			guru.POST("", r.GuruController.CreateGuru)
			guru.PATCH(":guru_id", r.GuruController.UpdateGuru)
			guru.DELETE(":guru_id", r.GuruController.DeleteGuru)
		}

		siswa := loggedPath.Group("siswa")
		{
			siswa.GET("", r.SiswaController.GetSiswa)

			siswaData := siswa.Group("", r.JWTService.IsAdmin)
			{
				siswaData.POST("", r.JWTService.IsAdmin, r.SiswaController.CreateSiswa)
				siswaData.PATCH(":siswa_id", r.JWTService.IsAdmin, r.SiswaController.UpdateSiswa)
				siswaData.DELETE(":siswa_id", r.JWTService.IsAdmin, r.SiswaController.DeleteSiswa)
			}
		}

		dataPoin := loggedPath.Group("data-poin", r.JWTService.IsAdmin)
		{
			dataPoin.GET("", r.DataPoinController.GetDataPoin)
			dataPoin.POST("", r.DataPoinController.CreateDataPoin)
			dataPoin.PATCH(":data_poin_id", r.DataPoinController.UpdateDataPoin)
			dataPoin.DELETE(":data_poin_id", r.DataPoinController.DeleteDataPoin)
		}

		poinSiswa := loggedPath.Group("poin", r.JWTService.IsGuru)
		{
			//dataPoin.GET("", r.DataPoinController.GetDataPoin)
			poinSiswa.POST("", r.PoinSiswaController.AddPoinSiswa)
			poinSiswa.PATCH("log/:poin_log_id", r.PoinSiswaController.UpdatePoinSiswa)
			poinSiswa.DELETE("log/:poin_log_id", r.PoinSiswaController.DeletePoinSiswa)
		}

		setting := loggedPath.Group("setting", r.JWTService.IsAdmin)
		{
			setting.GET("", r.SettingController.GetSetting)
			setting.POST("", r.SettingController.CreateSetting)
			setting.PATCH(":setting_id", r.SettingController.UpdateSetting)
			setting.DELETE(":setting_id", r.SettingController.DeleteSetting)
		}
	}
}
