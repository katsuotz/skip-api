package repository

import (
	"context"
	"fmt"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/entity/siti_entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type SyncRepository interface {
	GetSync(ctx context.Context, page int, perPage int) dto.SyncPagination
	IsOnProgress(ctx context.Context, syncType string) bool
	Sync(ctx context.Context)
	SyncPassword(ctx context.Context)
}

type syncRepository struct {
	db     *gorm.DB
	sitiDb *gorm.DB
}

func NewSyncRepository(db *gorm.DB, sitiDb *gorm.DB) SyncRepository {
	return &syncRepository{
		db,
		sitiDb,
	}
}

func (r *syncRepository) GetSync(ctx context.Context, page int, perPage int) dto.SyncPagination {
	result := dto.SyncPagination{}
	sync := entity.Sync{}
	temp := r.db.Model(&sync)

	temp.Order("created_at desc")
	temp.Offset(perPage * (page - 1)).Limit(perPage)
	temp.Find(&result.Data)

	var totalItem int64
	temp.Offset(-1).Limit(-1).Count(&totalItem)
	result.Pagination.TotalItem = totalItem
	result.Pagination.Page = page
	totalPage := totalItem / int64(perPage)
	if totalItem%int64(perPage) > 0 {
		totalPage++
	}
	result.Pagination.TotalPage = totalPage
	result.Pagination.PerPage = perPage

	return result
}

func (r *syncRepository) IsOnProgress(ctx context.Context, syncType string) bool {
	sync := entity.Sync{}
	r.db.
		Where("status = ?", "on progress").
		Where("type = ?", syncType).
		First(&sync)
	return sync.ID != 0
}

func (r *syncRepository) Sync(ctx context.Context) {
	sync := entity.Sync{
		Type:        "siti",
		Status:      "on progress",
		Description: "Start synchronizing - synchronizing tahun ajar",
	}
	r.db.Create(&sync)

	var tahunPelajaran siti_entity.TahunPelajaran
	r.sitiDb.
		Where("aktif = ?", true).
		First(&tahunPelajaran)

	tahunAjar := entity.TahunAjar{
		ID:        tahunPelajaran.IDTahunPelajaran,
		TahunAjar: tahunPelajaran.TahunPelajaran,
	}

	err := r.db.First(&tahunAjar).Error

	if err != nil {
		r.db.Create(&tahunAjar)
	}

	sync.Description = "Synchronizing jurusan"
	r.db.Save(&sync)

	sync.Description = "Synchronizing guru"
	r.db.Save(&sync)

	var sitiGuru []siti_entity.GuruBio
	r.sitiDb.
		Model(&siti_entity.Guru{}).
		Select("*").
		Where("aktif = ?", 1).
		Joins("join t_login on t_login.id_rujuk = t_guru.grnrs").
		Joins("join t_bio on t_bio.id_bio = t_guru.id_bio").
		Find(&sitiGuru)

	for _, item := range sitiGuru {
		tipePegawai := "Guru"
		role := "guru"

		switch item.KategoriKaryawan {
		case "Tata Usaha":
		case "Caraka":
			tipePegawai = "Tata Usaha"
			role = "tata-usaha"
			break
		}

		user := entity.User{}

		r.db.Where("username = ?", item.Username).First(&user)

		if user.ID == 0 {
			password, _ := helper.HashPassword(item.Password)

			user = entity.User{
				ID:       item.IDLogin,
				Username: item.Username,
				Password: password,
				Role:     role,
			}

			r.db.Create(&user)

			pegawai := entity.Pegawai{
				ID:          item.IDGuru,
				Nip:         item.GRNRS,
				TipePegawai: tipePegawai,
				UserID:      user.ID,
			}

			r.db.Create(&pegawai)

			profile := entity.Profile{
				Nama:         item.NamaLengkap,
				JenisKelamin: item.JenisKelamin,
				TanggalLahir: item.TanggalLahir,
				TempatLahir:  item.TempatLahir,
				UserID:       user.ID,
			}

			r.db.Create(&profile)
		}
	}

	sync.Description = "Synchronizing kelas"
	r.db.Save(&sync)

	var sitiKelas []siti_entity.Kelas
	r.sitiDb.Where("id_tahun_pelajaran = ?", tahunPelajaran.IDTahunPelajaran).Find(&sitiKelas)

	tahunPelajaranString := strings.Split(tahunAjar.TahunAjar, "/")
	fmt.Println(tahunPelajaranString)
	tahunPelajaran1, _ := strconv.Atoi(tahunPelajaranString[0])
	tahunPelajaran2, _ := strconv.Atoi(tahunPelajaranString[1])

	tahunPelajaran1 -= 1
	tahunPelajaran2 -= 1

	prevTahunPelajaran := strconv.Itoa(tahunPelajaran1) + "/" + strconv.Itoa(tahunPelajaran2)

	prevTahunAjar := entity.TahunAjar{}

	r.db.Where("tahun_ajar = ?", prevTahunPelajaran).First(&prevTahunAjar)

	for _, item := range sitiKelas {
		sync.Description = "Synchronizing kelas - " + item.NamaKelas
		r.db.Save(&sync)

		var keahlian siti_entity.Keahlian
		r.sitiDb.Where("id_keahlian = ?", item.IDKeahlian).First(&keahlian)

		var jurusan entity.Jurusan
		r.db.Where("nama_jurusan = ?", keahlian.NamaKeahlian).First(&jurusan)

		if jurusan.ID == 0 {
			jurusan = entity.Jurusan{
				NamaJurusan: keahlian.NamaKeahlian,
			}
			r.db.Create(&jurusan)
		}

		kelas := entity.Kelas{
			ID:          item.IDKelas,
			NamaKelas:   item.NamaKelas,
			TahunAjarID: item.IDTahunPelajaran,
			JurusanID:   item.IDKeahlian,
			PegawaiID:   item.IDGuru,
			Tingkat:     item.Tingkat,
		}
		r.db.Create(&kelas)

		var siswa []siti_entity.DetKelasJoin
		r.sitiDb.
			Model(&siti_entity.DetKelas{}).
			Select("*").
			Where("id_kelas = ?", item.IDKelas).
			Joins("join t_siswa on t_siswa.nis = t_det_kelas.nis").
			Joins("join t_login on t_login.id_rujuk = t_siswa.nis").
			Joins("join t_bio on t_bio.id_bio = t_siswa.id_bio").
			Find(&siswa)

		for _, itemSiswa := range siswa {
			user := entity.User{}

			r.db.Where("username = ?", itemSiswa.Username).First(&user)

			if user.ID == 0 {
				password, _ := helper.HashPassword(itemSiswa.Password)

				user = entity.User{
					ID:       itemSiswa.IDLogin,
					Username: itemSiswa.Username,
					Password: password,
					Role:     "siswa",
				}

				r.db.Create(&user)

				siswaData := entity.Siswa{
					Nis:    itemSiswa.Nis,
					UserID: user.ID,
				}

				r.db.Create(&siswaData)

				profile := entity.Profile{
					Nama:         itemSiswa.NamaLengkap,
					JenisKelamin: itemSiswa.JenisKelamin,
					TanggalLahir: itemSiswa.TanggalLahir,
					TempatLahir:  itemSiswa.TempatLahir,
					UserID:       user.ID,
				}

				r.db.Create(&profile)
			}

			selectedSiswa := entity.Siswa{}

			r.db.Where("nis = ?", itemSiswa.Nis).First(&selectedSiswa)

			siswaKelas := entity.SiswaKelas{
				SiswaID: selectedSiswa.ID,
				KelasID: item.IDKelas,
			}

			r.db.
				Where("siswa_id = ?", selectedSiswa.ID).
				Where("kelas_id = ?", item.IDKelas).
				First(&siswaKelas)

			if siswaKelas.ID == 0 {
				siswaKelas = entity.SiswaKelas{
					SiswaID: selectedSiswa.ID,
					KelasID: item.IDKelas,
				}

				r.db.Create(&siswaKelas)
			}

			poinSiswa := entity.PoinSiswa{}

			r.db.Where("siswa_kelas_id = ?", siswaKelas.ID).First(&poinSiswa)

			if poinSiswa.ID == 0 {
				prevSiswa := entity.PoinSiswa{}
				poin := float64(200)

				if prevTahunAjar.ID != 0 {
					r.db.
						Where("kelas.tahun_ajar_id = ?", prevTahunAjar.ID).
						Where("siswa.nis = ?", itemSiswa.Nis).
						Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
						Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
						Joins("join kelas on kelas.id = siswa_kelas.kelas_id").
						First(&prevSiswa)

					if prevSiswa.ID != 0 {
						poin = prevSiswa.Poin
					}
				}

				poinSiswa = entity.PoinSiswa{
					SiswaKelasID: siswaKelas.ID,
					Poin:         float64(poin),
				}
				r.db.Create(&poinSiswa)
			}
		}
	}

	sync.Status = "done"
	sync.Description = "Synchronizing success"
	r.db.Save(&sync)
}

func (r *syncRepository) SyncPassword(ctx context.Context) {
	sync := entity.Sync{
		Type:        "password",
		Status:      "on progress",
		Description: "Start synchronizing - synchronizing password",
	}
	r.db.Create(&sync)

	var users []entity.User
	r.db.Find(&users)

	for _, user := range users {
		login := siti_entity.Login{}

		r.sitiDb.Where("username = ?", user.Username).First(&login)

		if login.IDLogin != 0 {
			password, _ := helper.HashPassword(login.Password)
			user.Password = password
			r.db.Save(&user)
		}
	}

	sync.Status = "done"
	sync.Description = "Synchronizing success"
	r.db.Save(&sync)
}
