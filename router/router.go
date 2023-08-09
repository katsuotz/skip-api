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
	PegawaiController   controller.PegawaiController
	SiswaController     controller.SiswaController
	DataPoinController  controller.DataPoinController
	PoinSiswaController controller.PoinSiswaController
	PoinLogController   controller.PoinLogController
	SettingController   controller.SettingController
	InfoController      controller.InfoController
	FileController      controller.FileController
	SitiController      controller.SitiController
	JWTService          service.JWTService
}

func NewRouter(server *gin.Engine,
	authController controller.AuthController,
	profileController controller.ProfileController,
	jurusanController controller.JurusanController,
	tahunAjarController controller.TahunAjarController,
	kelasController controller.KelasController,
	pegawaiController controller.PegawaiController,
	siswaController controller.SiswaController,
	dataPoinController controller.DataPoinController,
	poinSiswaController controller.PoinSiswaController,
	poinLogController controller.PoinLogController,
	settingController controller.SettingController,
	infoController controller.InfoController,
	fileController controller.FileController,
	sitiController controller.SitiController,
	jwtService service.JWTService,
) *Router {
	return &Router{
		server,
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
		sitiController,
		jwtService,
	}
}

func (r *Router) Init() {

	/* Static files from storage dir */

	r.server.Static("/storage", "./storage")

	r.server.GET("", r.AuthController.Tes)

	/* Guest User Only */

	guest := r.server.Group("", r.JWTService.IsGuest)
	{
		guest.POST("login", r.AuthController.Login)
	}

	r.server.GET("siswa/:nis/log", r.SiswaController.GetSiswaDetailByNIS)

	/* Authorized User Only */

	authorized := r.server.Group("", r.JWTService.IsLoggedIn)
	{
		/* Auth & Profile API */

		authorized.GET("me", r.ProfileController.GetMyProfile)
		authorized.PATCH("profile", r.ProfileController.UpdateProfile)
		authorized.PATCH("password", r.AuthController.UpdatePassword)

		/* Jurusan CRUD */

		jurusan := authorized.Group("jurusan")
		{
			jurusan.GET("", r.JurusanController.GetJurusan)

			jurusanAdmin := jurusan.Group("", r.JWTService.IsAdmin)
			{
				jurusanAdmin.POST("", r.JurusanController.CreateJurusan)
				jurusanAdmin.PATCH(":jurusan_id", r.JurusanController.UpdateJurusan)
				jurusanAdmin.DELETE(":jurusan_id", r.JurusanController.DeleteJurusan)
			}
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

		/* Pegawai CRUD */

		pegawai := authorized.Group("pegawai", r.JWTService.IsAdmin)
		{
			pegawai.GET("", r.PegawaiController.GetPegawai)
			pegawai.POST("", r.PegawaiController.CreatePegawai)
			pegawai.PATCH(":pegawai_id", r.PegawaiController.UpdatePegawai)
			pegawai.DELETE(":pegawai_id", r.PegawaiController.DeletePegawai)
		}

		/* Siswa CRUD */

		siswa := authorized.Group("siswa")
		{
			siswa.GET("", r.JWTService.IsNotSiswa, r.SiswaController.GetSiswa)
			//siswa.GET(":nis/log", r.SiswaController.GetSiswaDetailByNIS)

			siswaAdmin := siswa.Group("", r.JWTService.IsAdmin)
			{
				siswaAdmin.POST("", r.SiswaController.CreateSiswa)
				siswaAdmin.PATCH(":siswa_id", r.SiswaController.UpdateSiswa)
				siswaAdmin.DELETE(":siswa_id", r.SiswaController.DeleteSiswa)
			}
		}

		/* Data Poin CRUD */

		dataPoin := authorized.Group("data-poin")
		{
			dataPoin.GET("", r.DataPoinController.GetDataPoin)
			dataPoin.GET(":data_poin_id", r.DataPoinController.GetDataPoinByID)

			dataPoinAdmin := dataPoin.Group("", r.JWTService.IsStaff)
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

			poinSiswaAdmin := poinSiswa.Group("", r.JWTService.IsStaff)
			{
				poinSiswaAdmin.GET("kelas/:kelas_id", r.PoinSiswaController.GetPoinKelas)
				poinSiswaAdmin.GET("jurusan/:jurusan_id", r.PoinSiswaController.GetPoinJurusan)

				// update tindak lanjut
				poinSiswaAdmin.PATCH("log/:poin_log_id", r.PoinSiswaController.UpdatePoinSiswa)
			}

			/* Poin Siswa Transaction - Only for Pegawai */

			poinSiswaPegawai := poinSiswa.Group("", r.JWTService.IsPegawai)
			{
				poinSiswaPegawai.GET("siswa/:siswa_kelas_id", r.PoinSiswaController.GetPoinSiswa)
				poinSiswaPegawai.POST("", r.PoinSiswaController.AddPoinSiswa)
				poinSiswaPegawai.DELETE("log/:poin_log_id", r.PoinSiswaController.DeletePoinSiswa)
			}

			/* Poin Log API - Only for Admin */

			poinSiswaLog := poinSiswa.Group("log", r.JWTService.IsPegawai)
			{
				poinSiswaLog.GET("", r.PoinLogController.GetPoinLog)
				poinSiswaLog.GET(":siswa_kelas_id", r.PoinLogController.GetPoinSiswaLog)
			}
		}

		/* Info API - For statistics & metrics */

		info := authorized.Group("info", r.JWTService.IsPegawai)
		{

			/* Info Poin Siswa & Log */

			infoMetrics := info.Group("")
			{
				// count poin log by category
				infoMetrics.GET("poin/count", r.InfoController.CountPoinLog)
				// get max poin siswa
				infoMetrics.GET("poin/max", r.InfoController.MaxPoinSiswa)
				// get min poin siswa
				infoMetrics.GET("poin/min", r.InfoController.MinPoinSiswa)
				// get avg poin siswa
				infoMetrics.GET("poin/avg", r.InfoController.AvgPoinSiswa)
				// get poin list
				infoMetrics.GET("poin/list", r.InfoController.ListPoinSiswa)
				// get poin siswa list grouped by category
				infoMetrics.GET("poin/list/count", r.InfoController.ListCountPoinLog)
				// get poin log list grouped by month
				infoMetrics.GET("poin/graph/count", r.InfoController.GraphCountPoinLog)
				// get poin siswa total with filter
				infoMetrics.GET("poin/total", r.InfoController.CountPoinSiswaTotal)
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

	r.server.POST("/sync", r.SitiController.Sync)
	r.server.POST("/sync/password", r.SitiController.SyncPassword)
}
