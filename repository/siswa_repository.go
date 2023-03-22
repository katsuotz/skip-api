package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	GetSiswa(ctx context.Context, page int, perPage int, search string, kelasID int) dto.SiswaPagination
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

func (r *siswaRepository) GetSiswa(ctx context.Context, page int, perPage int, search string, kelasID int) dto.SiswaPagination {
	result := dto.SiswaPagination{}
	siswa := entity.Siswa{}
	temp := r.db.Model(siswa)
	if search != "" {
		search = "%" + search + "%"
		temp.Where("name ilike ?", search, search)
	}

	temp.Select("siswa.id as id, siswa.user_id as user_id, nis, nama, jenis_kelamin, tanggal_lahir, tempat_lahir")
	temp.Joins("join users on users.id = siswa.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")

	if kelasID != 0 {
		temp.Joins("join detail_kelas on detail_kelas.siswa_id = siswa.id")
		temp.Where("detail_kelas.kelas_id = ?", kelasID)
	}

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

	tx.Commit()

	return nil
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

	tx.Commit()

	return nil
}

func (r *siswaRepository) DeleteSiswa(ctx context.Context, siswaID int) error {
	err := r.db.Where("id = ?", siswaID).Delete(&entity.Siswa{}).Error
	return err
}
