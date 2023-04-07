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
	InfoController      controller.InfoController
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
	infoController controller.InfoController,
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
		infoController,
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
				kelasData.POST(":kelas_id/naik-kelas", r.KelasController.SiswaNaikKelas)
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
			siswa.GET(":nis/log", r.SiswaController.GetSiswaDetailByNIS)

			siswaData := siswa.Group("", r.JWTService.IsAdmin)
			{
				siswaData.POST("", r.JWTService.IsAdmin, r.SiswaController.CreateSiswa)
				siswaData.PATCH(":siswa_id", r.JWTService.IsAdmin, r.SiswaController.UpdateSiswa)
				siswaData.DELETE(":siswa_id", r.JWTService.IsAdmin, r.SiswaController.DeleteSiswa)
			}
		}

		loggedPath.GET("data-poin", r.DataPoinController.GetDataPoin)

		dataPoin := loggedPath.Group("data-poin", r.JWTService.IsAdmin)
		{
			dataPoin.POST("", r.DataPoinController.CreateDataPoin)
			dataPoin.PATCH(":data_poin_id", r.DataPoinController.UpdateDataPoin)
			dataPoin.DELETE(":data_poin_id", r.DataPoinController.DeleteDataPoin)
		}

		poinSiswaAdmin := loggedPath.Group("poin",
			func(context *gin.Context) {
				r.JWTService.IsAdmin(context)
				return
			})
		{
			poinSiswaAdmin.GET("siswa/:siswa_kelas_id", r.PoinSiswaController.GetPoinSiswa)
			poinSiswaAdmin.GET("kelas/:kelas_id", r.PoinSiswaController.GetPoinKelas)
			poinSiswaAdmin.GET("jurusan/:jurusan_id/:tahun_ajar_id", r.PoinSiswaController.GetPoinJurusan)
			poinSiswaAdmin.GET("log/:siswa_kelas_id", r.PoinSiswaController.GetPoinSiswaLog)
		}

		poinSiswaGuru := loggedPath.Group("poin", r.JWTService.IsGuru)
		{
			poinSiswaGuru.POST("", r.PoinSiswaController.AddPoinSiswa)
			poinSiswaGuru.PATCH("log/:poin_log_id", r.PoinSiswaController.UpdatePoinSiswa)
			poinSiswaGuru.DELETE("log/:poin_log_id", r.PoinSiswaController.DeletePoinSiswa)
		}

		setting := loggedPath.Group("setting", r.JWTService.IsAdmin)
		{
			setting.GET("", r.SettingController.GetSetting)
			setting.POST("", r.SettingController.CreateSetting)
			setting.PATCH(":setting_id", r.SettingController.UpdateSetting)
			setting.DELETE(":setting_id", r.SettingController.DeleteSetting)
		}

		info := loggedPath.Group("info", r.JWTService.IsAdmin)
		{
			info.GET("poin/count", r.InfoController.CountPoin)
			info.GET("poin/list", r.InfoController.ListPoin)
		}
	}
}
