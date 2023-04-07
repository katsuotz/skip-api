package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	GetSiswa(ctx context.Context, page int, perPage int, search string, kelasID string, tahunAjarActive string) dto.SiswaPagination
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

func (r *siswaRepository) GetSiswa(ctx context.Context, page int, perPage int, search string, kelasID string, tahunAjarActive string) dto.SiswaPagination {
	result := dto.SiswaPagination{}
	siswa := entity.Siswa{}
	temp := r.db.Model(&siswa)
	if search != "" {
		search = "%" + search + "%"
		temp.Where(temp.Where("nama ilike ?", search).Or("nis ilike ?", search))
	}

	selectQuery := "siswa.id as id, siswa.user_id as user_id, nis, profiles.id as profile_id, nama, jenis_kelamin, tanggal_lahir, tempat_lahir"

	temp.Joins("join users on users.id = siswa.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")

	if tahunAjarActive == "true" {
		activeTahunAjar := entity.TahunAjar{}

		r.db.Where("is_active = ?", true).First(&activeTahunAjar)

		if activeTahunAjar.ID == 0 {
			return result
		}

		temp.Joins("right join siswa_kelas on siswa_kelas.siswa_id = siswa.id")
		temp.Joins("join kelas on kelas.id = siswa_kelas.kelas_id")
		temp.Where("kelas.tahun_ajar_id = ?", activeTahunAjar.ID)

		temp.Joins("left join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id")
		selectQuery += ", poin_siswa.poin as poin, siswa_kelas.id as siswa_kelas_id"
	} else {
		if kelasID != "" {
			temp.Joins("left join siswa_kelas on siswa_kelas.siswa_id = siswa.id")

			if kelasID != "0" {
				temp.Where("siswa_kelas.deleted_at is NULL")
				temp.Where("siswa_kelas.kelas_id = ?", kelasID)
				temp.Joins("left join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id")
				selectQuery += ", poin_siswa.poin as poin, siswa_kelas.id as siswa_kelas_id"
			} else {
				temp.Preload("SiswaKelas")
				//temp.Where(temp.Where("siswa_kelas.deleted_at is NOT NULL").Or("siswa_kelas.id is NULL"))
				temp.Group("siswa_kelas.siswa_id, siswa.id, profiles.id")
			}
		}
	}

	temp.Select(selectQuery)
	temp.Order("nama asc")
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

func (r *siswaRepository) GetSiswaByNIS(ctx context.Context, nis string) dto.SiswaResponse {
	result := dto.SiswaResponse{}
	siswa := entity.Siswa{}

	temp := r.db.Model(&siswa)
	temp.Select("siswa.id as id, siswa.user_id as user_id, nis, profiles.id as profile_id, nama, jenis_kelamin, tanggal_lahir, tempat_lahir")
	temp.Where("nis = ?", nis)
	temp.Joins("join users on users.id = siswa.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")
	temp.First(&result)

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

	user := entity.User{
		Username: req.Nis,
		Password: helper.BirthDateToPassword(tanggalLahir),
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
	tx.Find(&findSiswa)

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
	}

	err = tx.Where("user_id = ?", findSiswa.UserID).Updates(&profile).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *siswaRepository) DeleteSiswa(ctx context.Context, siswaID int) error {
	err := r.db.Where("id = ?", siswaID).Delete(&entity.Siswa{}).Error

	if err == nil {
		err = r.db.Where("siswa_id = ?", siswaID).Delete(&entity.SiswaKelas{}).Error
	}

	return err
}
