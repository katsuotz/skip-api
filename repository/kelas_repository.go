package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type KelasRepository interface {
	GetKelas(ctx context.Context, jurusanID string, tahunAjarID string) []dto.KelasResponse
	CreateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error)
	UpdateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error)
	DeleteKelas(ctx context.Context, kelasID int) error
	AddSiswaToKelas(ctx context.Context, kelasID int, siswaIDs []int) error
	RemoveSiswaFromKelas(ctx context.Context, kelasID int, siswaIDs []int) error
}

type kelasRepository struct {
	db *gorm.DB
}

func NewKelasRepository(db *gorm.DB) KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) GetKelas(ctx context.Context, jurusanID string, tahunAjarID string) []dto.KelasResponse {
	var kelas []dto.KelasResponse
	temp := r.db.Model(&entity.Kelas{})

	if jurusanID != "" {
		temp.Where("jurusan_id = ?", jurusanID)
	}

	if tahunAjarID != "" {
		temp.Where("tahun_ajar_id = ?", tahunAjarID)
	}

	temp.Select("kelas.id as id, nama_kelas, jurusan_id, tahun_ajar_id, guru_id, nip, tipe_guru, nama")
	temp.Joins("join guru on guru.id = kelas.guru_id")
	temp.Joins("join users on users.id = guru.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")
	temp.Find(&kelas)

	return kelas
}

func (r *kelasRepository) CreateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error) {
	err := r.db.Create(&kelas).Error
	return kelas, err
}

func (r *kelasRepository) UpdateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error) {
	err := r.db.Updates(&kelas).Error
	return kelas, err
}

func (r *kelasRepository) DeleteKelas(ctx context.Context, kelasID int) error {
	err := r.db.Where("id = ?", kelasID).Delete(&entity.Kelas{}).Error
	return err
}

func (r *kelasRepository) AddSiswaToKelas(ctx context.Context, kelasID int, siswaIDs []int) error {
	var siswaKelasInsert []entity.SiswaKelas

	var addedSiswaIDs []int

	r.db.
		Model(&entity.SiswaKelas{}).
		Where("kelas_id = ?", kelasID).
		Where("siswa_id in ?", siswaIDs).
		Pluck("siswa_id", &addedSiswaIDs)

	for _, siswaID := range siswaIDs {
		if !helper.IsInArray(addedSiswaIDs, siswaID) {
			siswaKelasInsert = append(siswaKelasInsert, entity.SiswaKelas{
				KelasID: kelasID,
				SiswaID: siswaID,
			})
		}
	}

	if len(siswaKelasInsert) == 0 {
		return nil
	}

	err := r.db.Create(&siswaKelasInsert).Error
	return err
}

func (r *kelasRepository) RemoveSiswaFromKelas(ctx context.Context, kelasID int, siswaIDs []int) error {
	err := r.db.
		Where("kelas_id = ?", kelasID).
		Where("siswa_id in ?", siswaIDs).
		Delete(&entity.SiswaKelas{}).Error

	return err
}
