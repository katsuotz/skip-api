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
	PoinLogController   controller.PoinLogController
	SettingController   controller.SettingController
	InfoController      controller.InfoController
	FileController      controller.FileController
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
	poinLogController controller.PoinLogController,
	settingController controller.SettingController,
	infoController controller.InfoController,
	fileController controller.FileController,
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
		poinLogController,
		settingController,
		infoController,
		fileController,
		jwtService,
	}
}

func (r *Router) Init() {

	/* Static files from storage dir */

	r.server.Static("/storage", "./storage")

	base := r.server.Group("/api")
	{
		base.GET("", r.AuthController.Tes)

		/* Guest User Only */

		guest := base.Group("", r.JWTService.IsGuest)
		{
			guest.POST("login", r.AuthController.Login)
		}

		/* Authorized User Only */

		authorized := base.Group("", r.JWTService.IsLoggedIn)
		{
			/* Auth & Profile API */

			authorized.GET("me", r.ProfileController.GetMyProfile)
			authorized.PATCH("profile", r.ProfileController.UpdateProfile)
			authorized.PATCH("password", r.AuthController.UpdatePassword)

			/* Jurusan CRUD */

			jurusan := authorized.Group("jurusan", r.JWTService.IsAdmin)
			{
				jurusan.GET("", r.JurusanController.GetJurusan)
				jurusan.POST("", r.JurusanController.CreateJurusan)
				jurusan.PATCH(":jurusan_id", r.JurusanController.UpdateJurusan)
				jurusan.DELETE(":jurusan_id", r.JurusanController.DeleteJurusan)
			}

			/* Tahun Ajar CRUD */

			authorized.GET("tahun-ajar", r.TahunAjarController.GetTahunAjar)

			tahunAjar := authorized.Group("tahun-ajar", r.JWTService.IsAdmin)
			{
				tahunAjar.POST("", r.TahunAjarController.CreateTahunAjar)
				tahunAjar.PATCH(":tahun_ajar_id", r.TahunAjarController.UpdateTahunAjar)
				tahunAjar.DELETE(":tahun_ajar_id", r.TahunAjarController.DeleteTahunAjar)
				tahunAjar.PATCH(":tahun_ajar_id/set-active", r.TahunAjarController.SetActiveTahunAjar)
			}

			/* Kelas CRUD */

			kelas := authorized.Group("kelas")
			{
				kelas.GET("", r.KelasController.GetKelas)
				kelas.GET(":kelas_id", r.KelasController.GetKelasByID)

				kelasAdmin := kelas.Group("", r.JWTService.IsAdmin)
				{
					kelasAdmin.POST("", r.KelasController.CreateKelas)
					kelasAdmin.PATCH(":kelas_id", r.KelasController.UpdateKelas)
					kelasAdmin.DELETE(":kelas_id", r.KelasController.DeleteKelas)
					kelasAdmin.POST(":kelas_id/add-siswa", r.KelasController.AddSiswaToKelas)
					kelasAdmin.POST(":kelas_id/naik-kelas", r.KelasController.SiswaNaikKelas)
					kelasAdmin.DELETE(":kelas_id/remove-siswa", r.KelasController.RemoveSiswaFromKelas)
				}
			}

			/* Guru CRUD */

			guru := authorized.Group("guru", r.JWTService.IsAdmin)
			{
				guru.GET("", r.GuruController.GetGuru)
				guru.POST("", r.GuruController.CreateGuru)
				guru.PATCH(":guru_id", r.GuruController.UpdateGuru)
				guru.DELETE(":guru_id", r.GuruController.DeleteGuru)
			}

			/* Siswa CRUD */

			siswa := authorized.Group("siswa")
			{
				siswa.GET("", r.JWTService.IsNotSiswa, r.SiswaController.GetSiswa)
				siswa.GET(":nis/log", r.SiswaController.GetSiswaDetailByNIS)

				siswaAdmin := siswa.Group("", r.JWTService.IsAdmin)
				{
					siswaAdmin.POST("", r.JWTService.IsAdmin, r.SiswaController.CreateSiswa)
					siswaAdmin.PATCH(":siswa_id", r.JWTService.IsAdmin, r.SiswaController.UpdateSiswa)
					siswaAdmin.DELETE(":siswa_id", r.JWTService.IsAdmin, r.SiswaController.DeleteSiswa)
				}
			}

			/* Data Poin CRUD */

			dataPoin := authorized.Group("data-poin")
			{
				dataPoin.GET("", r.DataPoinController.GetDataPoin)

				dataPoinAdmin := dataPoin.Group("", r.JWTService.IsAdmin)
				{
					dataPoinAdmin.POST("", r.DataPoinController.CreateDataPoin)
					dataPoinAdmin.PATCH(":data_poin_id", r.DataPoinController.UpdateDataPoin)
					dataPoinAdmin.DELETE(":data_poin_id", r.DataPoinController.DeleteDataPoin)
				}
			}

			/* Poin Siswa & Log API */

			poinSiswa := authorized.Group("poin")
			{
				/* Get Poin Result of Siswa, Kelas, Jurusan - Only for Admin */

				poinSiswaAdmin := poinSiswa.Group("", r.JWTService.IsAdmin)
				{
					poinSiswaAdmin.GET("siswa/:siswa_kelas_id", r.PoinSiswaController.GetPoinSiswa)
					poinSiswaAdmin.GET("kelas/:kelas_id", r.PoinSiswaController.GetPoinKelas)
					poinSiswaAdmin.GET("jurusan/:jurusan_id/:tahun_ajar_id", r.PoinSiswaController.GetPoinJurusan)
				}

				/* Poin Siswa Transaction - Only for Guru */

				poinSiswaGuru := authorized.Group("", r.JWTService.IsGuru)
				{
					poinSiswaGuru.POST("", r.PoinSiswaController.AddPoinSiswa)
					poinSiswaGuru.PATCH("log/:poin_log_id", r.PoinSiswaController.UpdatePoinSiswa)
					poinSiswaGuru.DELETE("log/:poin_log_id", r.PoinSiswaController.DeletePoinSiswa)
				}

				/* Poin Log API - Only for Admin */

				poinSiswaLog := poinSiswa.Group("log", r.JWTService.IsAdmin)
				{
					poinSiswaLog.GET("", r.PoinLogController.GetPoinLog)
					poinSiswaLog.GET(":siswa_kelas_id", r.PoinLogController.GetPoinSiswaLog)
				}
			}

			/* Info API - For statistics & metrics */

			info := authorized.Group("info", r.JWTService.IsAdmin)
			{

				/* Info Poin Siswa & Log */

				infoMetrics := info.Group("", r.JWTService.IsNotSiswa)
				{
					infoMetrics.GET("poin/count", r.InfoController.CountPoin)
					infoMetrics.GET("poin/max", r.InfoController.MaxPoin)
					infoMetrics.GET("poin/min", r.InfoController.MinPoin)
					infoMetrics.GET("poin/avg", r.InfoController.AvgPoin)
					infoMetrics.GET("poin/list", r.InfoController.ListPoinSiswa)
					infoMetrics.GET("poin/list/count", r.InfoController.ListCountPoinLog)
					infoMetrics.GET("poin/graph/count", r.InfoController.GraphCountPoinLog)
				}

				/* Info Auth */

				infoAdmin := info.Group("", r.JWTService.IsAdmin)
				{
					infoAdmin.GET("login", r.AuthController.GetLog)
				}
			}

			/* Setting CRUD - unused, maybe later */

			setting := authorized.Group("setting", r.JWTService.IsAdmin)
			{
				setting.GET("", r.SettingController.GetSetting)
				setting.POST("", r.SettingController.CreateSetting)
				setting.PATCH(":setting_id", r.SettingController.UpdateSetting)
				setting.DELETE(":setting_id", r.SettingController.DeleteSetting)
			}

			/* Upload File API */

			authorized.POST("upload", r.FileController.Upload)
		}
	}
}
