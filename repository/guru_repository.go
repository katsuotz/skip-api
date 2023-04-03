package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type GuruRepository interface {
	GetGuru(ctx context.Context, page int, perPage int, search string) dto.GuruPagination
	CreateGuru(ctx context.Context, guru dto.GuruRequest) error
	UpdateGuru(ctx context.Context, guru dto.GuruRequest, guruID int) error
	DeleteGuru(ctx context.Context, guruID int) error
}

type guruRepository struct {
	db *gorm.DB
}

func NewGuruRepository(db *gorm.DB) GuruRepository {
	return &guruRepository{db: db}
}

func (r *guruRepository) GetGuru(ctx context.Context, page int, perPage int, search string) dto.GuruPagination {
	result := dto.GuruPagination{}
	guru := entity.Guru{}
	temp := r.db.Model(&guru)
	if search != "" {
		search = "%" + search + "%"
		temp.Where("nama ilike ?", search)
	}

	temp.Select("guru.id as id, guru.user_id as user_id, nip, tipe_guru, nama, jenis_kelamin, tanggal_lahir, tempat_lahir")
	temp.Joins("join users on users.id = guru.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")
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

func (r *guruRepository) CreateGuru(ctx context.Context, req dto.GuruRequest) error {
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
		Username: req.Nip,
		Password: helper.BirthDateToPassword(tanggalLahir),
		Role:     "guru",
	}

	err := tx.Create(&user).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	guru := entity.Guru{
		Nip:      req.Nip,
		TipeGuru: req.TipeGuru,
		UserID:   user.ID,
	}

	err = tx.Create(&guru).Error

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

func (r *guruRepository) UpdateGuru(ctx context.Context, req dto.GuruRequest, guruID int) error {
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

	findGuru := entity.Guru{
		ID: guruID,
	}
	tx.Find(&findGuru)

	user := entity.User{
		Username: req.Nip,
		Password: helper.BirthDateToPassword(tanggalLahir),
	}

	err := tx.Where("id = ?", findGuru.UserID).Updates(&user).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	guru := entity.Guru{
		Nip:      req.Nip,
		TipeGuru: req.TipeGuru,
	}

	err = tx.Where("id = ?", findGuru.ID).Updates(&guru).Error

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

	err = tx.Where("user_id = ?", findGuru.UserID).Updates(&profile).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *guruRepository) DeleteGuru(ctx context.Context, guruID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	findGuru := entity.Guru{
		ID: guruID,
	}
	tx.Find(&findGuru)

	err := tx.Where("id = ?", findGuru.ID).Delete(&entity.Guru{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("user_id = ?", findGuru.UserID).Delete(&entity.Profile{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("id = ?", findGuru.UserID).Delete(&entity.User{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
