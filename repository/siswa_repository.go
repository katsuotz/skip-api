package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	GetSiswa(ctx context.Context, page int, perPage int, search string, kelasID string, tahunAjarID string, jurusanID string, tahunAjarActive string, summary bool) dto.SiswaPagination
	GetSiswaByNIS(ctx context.Context, nis string) dto.SiswaResponse
	CreateSiswa(ctx context.Context, siswa dto.SiswaRequest) error
	UpdateSiswa(ctx context.Context, siswa dto.SiswaRequest, siswaID int) error
	DeleteSiswa(ctx context.Context, siswaID int) error
}

type siswaRepository struct {
	db *gorm.DB
}

func NewSiswaRepository(db *gorm.DB) SiswaRepository {
	return &siswaRepository{db: db}
}

func (r *siswaRepository) GetSiswa(ctx context.Context, page int, perPage int, search string, kelasID string, tahunAjarID string, jurusanID string, tahunAjarActive string, summary bool) dto.SiswaPagination {
	result := dto.SiswaPagination{}
	siswa := entity.Siswa{}
	temp := r.db.Model(&siswa)
	if search != "" {
		search = "%" + search + "%"
		temp.Where("(nama ilike ? or nis ilike ?)", search, search)
	}

	selectQuery := "siswa.id as id, siswa.user_id as user_id, nis, profiles.id as profile_id, foto, nama, jenis_kelamin, tanggal_lahir, tempat_lahir, foto"

	if summary {
		selectQuery += ",(SELECT sum(poin_log.poin) FROM \"poin_log\" WHERE poin_log.poin_siswa_id = poin_siswa.id AND type = 'Penghargaan' AND \"poin_log\".\"deleted_at\" IS NULL) as total_penghargaan"
		selectQuery += ",(SELECT sum(poin_log.poin) FROM \"poin_log\" WHERE poin_log.poin_siswa_id = poin_siswa.id AND type = 'Pelanggaran' AND \"poin_log\".\"deleted_at\" IS NULL) as total_pelanggaran"
	}

	temp.Joins("join users on users.id = siswa.user_id").
		Joins("join profiles on profiles.user_id = users.id")

	if tahunAjarActive == "true" {
		activeTahunAjar := entity.TahunAjar{}

		r.db.Where("is_active = ?", true).First(&activeTahunAjar)

		if activeTahunAjar.ID == 0 {
			return result
		}

		temp.Where("kelas.tahun_ajar_id = ?", activeTahunAjar.ID).
			Joins("join siswa_kelas on siswa_kelas.siswa_id = siswa.id").
			Joins("join kelas on kelas.id = siswa_kelas.kelas_id").
			Joins("join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id")
		selectQuery += ", poin_siswa.poin as poin, siswa_kelas.id as siswa_kelas_id, nama_kelas"
	} else {
		if kelasID != "" {
			if kelasID != "0" {
				r.WhereSiswaKelas(temp, kelasID, tahunAjarID, jurusanID)
				selectQuery += ", poin_siswa.poin as poin, siswa_kelas.id as siswa_kelas_id"
			} else {
				temp.Joins("left join siswa_kelas on siswa_kelas.siswa_id = siswa.id").
					Preload("SiswaKelas").
					Group("siswa_kelas.siswa_id, siswa.id, profiles.id")
				//temp.Where(temp.Where("siswa_kelas.deleted_at is NOT NULL").Or("siswa_kelas.id is NULL"))
			}
		} else if tahunAjarID != "" || jurusanID != "" {
			r.WhereSiswaKelas(temp, kelasID, tahunAjarID, jurusanID)
			selectQuery += ", poin_siswa.poin as poin, siswa_kelas.id as siswa_kelas_id"
		}
	}

	temp.Select(selectQuery)
	temp.Order("nama asc")
	temp.Offset(perPage * (page - 1)).Limit(perPage).Find(&result.Data)

	var totalItem int64
	var totalPage int64
	if perPage != -1 {
		temp.Offset(-1).Limit(-1).Count(&totalItem)
		result.Pagination.Page = page
		totalPage = totalItem / int64(perPage)
		if totalItem%int64(perPage) > 0 {
			totalPage++
		}
	} else {
		totalItem = int64(len(result.Data))
		totalPage = int64(1)
		perPage = int(totalItem)
	}
	result.Pagination.TotalItem = totalItem
	result.Pagination.TotalPage = totalPage
	result.Pagination.PerPage = perPage

	return result
}

func (r *siswaRepository) WhereSiswaKelas(db *gorm.DB, kelasID string, tahunAjarID string, jurusanID string) {
	db.Joins("left join siswa_kelas on siswa_kelas.siswa_id = siswa.id").
		Joins("left join kelas on kelas.id = siswa_kelas.kelas_id")

	db.Where("siswa_kelas.deleted_at is NULL").
		Joins("left join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id")

	if kelasID != "" {
		db.Where("siswa_kelas.kelas_id = ?", kelasID)
	}
	if tahunAjarID != "" {
		db.Where("kelas.tahun_ajar_id = ?", tahunAjarID)
	}
	if jurusanID != "" {
		db.Where("kelas.jurusan_id = ?", jurusanID)
	}
}

func (r *siswaRepository) GetSiswaByNIS(ctx context.Context, nis string) dto.SiswaResponse {
	result := dto.SiswaResponse{}
	siswa := entity.Siswa{}

	temp := r.db.Model(&siswa)
	temp.Select("siswa.id as id, siswa.user_id as user_id, nis, profiles.id as profile_id, nama, jenis_kelamin, tanggal_lahir, tempat_lahir, foto").
		Where("nis = ?", nis).
		Joins("join users on users.id = siswa.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		First(&result)

	return result
}

func (r *siswaRepository) CreateSiswa(ctx context.Context, req dto.SiswaRequest) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	tanggalLahir, _ := helper.StringToDate(req.TanggalLahir)
	tanggalLahirPassword, _ := helper.DatePassword(req.TanggalLahir)
	password, _ := helper.HashPassword(tanggalLahirPassword)

	user := entity.User{
		Username: req.Nis,
		Password: password,
		Role:     "siswa",
	}

	err := tx.Create(&user).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	siswa := entity.Siswa{
		Nis:    req.Nis,
		UserID: user.ID,
	}

	err = tx.Create(&siswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	profile := entity.Profile{
		Nama:         req.Nama,
		JenisKelamin: req.JenisKelamin,
		TanggalLahir: tanggalLahir,
		TempatLahir:  req.TempatLahir,
		Foto:         req.Foto,
		UserID:       user.ID,
	}

	err = tx.Create(&profile).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *siswaRepository) UpdateSiswa(ctx context.Context, req dto.SiswaRequest, siswaID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	tanggalLahir, _ := helper.StringToDate(req.TanggalLahir)

	findSiswa := entity.Siswa{
		ID: siswaID,
	}
	tx.First(&findSiswa)

	user := entity.User{
		Username: req.Nis,
		Password: helper.BirthDateToPassword(tanggalLahir),
	}

	err := tx.Where("id = ?", findSiswa.UserID).Updates(&user).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	siswa := entity.Siswa{
		Nis: req.Nis,
	}

	err = tx.Where("id = ?", findSiswa.ID).Updates(&siswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	profile := entity.Profile{
		Nama:         req.Nama,
		JenisKelamin: req.JenisKelamin,
		TanggalLahir: tanggalLahir,
		TempatLahir:  req.TempatLahir,
		Foto:         req.Foto,
	}

	err = tx.Where("user_id = ?", findSiswa.UserID).Updates(&profile).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *siswaRepository) DeleteSiswa(ctx context.Context, siswaID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where("id = ?", siswaID).Delete(&entity.Siswa{}).Error

	if err == nil {
		siswaKelas := entity.SiswaKelas{}
		tx.Where("siswa_id = ?", siswaID).First(&siswaKelas)

		err = tx.Where("siswa_kelas_id = ?", siswaKelas.ID).Delete(&entity.PoinSiswa{}).Error

		if err == nil {
			err = tx.Delete(&siswaKelas).Error
		}
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
