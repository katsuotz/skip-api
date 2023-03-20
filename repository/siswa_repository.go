package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	GetSiswa(ctx context.Context, kelasID int) []entity.Siswa
	CreateSiswa(ctx context.Context, siswa dto.CreateSiswaRequest) error
	UpdateSiswa(ctx context.Context, siswa entity.Siswa) (entity.Siswa, error)
	DeleteSiswa(ctx context.Context, siswaID int) error
}

type siswaRepository struct {
	db *gorm.DB
}

func NewSiswaRepository(db *gorm.DB) SiswaRepository {
	return &siswaRepository{db: db}
}

func (r *siswaRepository) GetSiswa(ctx context.Context, kelasID int) []entity.Siswa {
	var siswa []entity.Siswa
	r.db.
		Where("detail_kelas.kelas_id = ?", kelasID).
		Joins("detail_kelas on detail_kelas.siswa_id = siswa.id").
		Find(&siswa)
	return siswa
}

func (r *siswaRepository) CreateSiswa(ctx context.Context, req dto.CreateSiswaRequest) error {
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

func (r *siswaRepository) UpdateSiswa(ctx context.Context, siswa entity.Siswa) (entity.Siswa, error) {
	err := r.db.Updates(&siswa).Error
	return siswa, err
}

func (r *siswaRepository) DeleteSiswa(ctx context.Context, siswaID int) error {
	err := r.db.Where("id = ?", siswaID).Delete(&entity.Siswa{}).Error
	return err
}
