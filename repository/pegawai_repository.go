package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type PegawaiRepository interface {
	GetPegawai(ctx context.Context, page int, perPage int, search string) dto.PegawaiPagination
	CreatePegawai(ctx context.Context, pegawai dto.PegawaiRequest) error
	UpdatePegawai(ctx context.Context, pegawai dto.PegawaiRequest, pegawaiID int) error
	DeletePegawai(ctx context.Context, pegawaiID int) error
}

type pegawaiRepository struct {
	db *gorm.DB
}

func NewPegawaiRepository(db *gorm.DB) PegawaiRepository {
	return &pegawaiRepository{db: db}
}

func (r *pegawaiRepository) GetPegawai(ctx context.Context, page int, perPage int, search string) dto.PegawaiPagination {
	result := dto.PegawaiPagination{}
	temp := r.db.Model(&entity.Pegawai{})
	if search != "" {
		search = "%" + search + "%"
		temp.Where("nama ilike ?", search)
	}

	temp.Select("pegawai.id as id, pegawai.user_id as user_id, nip, tipe_pegawai, nama, jenis_kelamin, tanggal_lahir, tempat_lahir, username, foto").
		Joins("join users on users.id = pegawai.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Order("nama asc")
	temp.Offset(perPage * (page - 1)).Limit(perPage).Find(&result.Data)

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

func (r *pegawaiRepository) CreatePegawai(ctx context.Context, req dto.PegawaiRequest) error {
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
	password, _ := helper.HashPassword(req.Password)

	user := entity.User{
		Username: req.Username,
		Password: password,
		Role:     getPegawaiRole(req.TipePegawai),
	}

	err := tx.Create(&user).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	pegawai := entity.Pegawai{
		Nip:         req.Nip,
		TipePegawai: req.TipePegawai,
		UserID:      user.ID,
	}

	err = tx.Create(&pegawai).Error

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

func (r *pegawaiRepository) UpdatePegawai(ctx context.Context, req dto.PegawaiRequest, pegawaiID int) error {
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
	password, _ := helper.HashPassword(req.Password)

	findPegawai := entity.Pegawai{
		ID: pegawaiID,
	}
	tx.First(&findPegawai)

	user := entity.User{
		Username: req.Username,
		Password: password,
		Role:     getPegawaiRole(req.TipePegawai),
	}

	err := tx.Where("id = ?", findPegawai.UserID).Updates(&user).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	pegawai := entity.Pegawai{
		Nip:         req.Nip,
		TipePegawai: req.TipePegawai,
	}

	err = tx.Where("id = ?", findPegawai.ID).Updates(&pegawai).Error

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

	err = tx.Where("user_id = ?", findPegawai.UserID).Updates(&profile).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *pegawaiRepository) DeletePegawai(ctx context.Context, pegawaiID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	findPegawai := entity.Pegawai{
		ID: pegawaiID,
	}
	tx.First(&findPegawai)

	err := tx.Where("id = ?", findPegawai.ID).Delete(&entity.Pegawai{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("user_id = ?", findPegawai.UserID).Delete(&entity.Profile{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("id = ?", findPegawai.UserID).Delete(&entity.User{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func getPegawaiRole(tipePegawai string) string {
	switch tipePegawai {
	case "Staff ICT":
		return "staff-ict"
	case "Guru BK":
		return "guru-bk"
	case "Tata Usaha":
		return "tata-usaha"
	}

	return "guru"
}
