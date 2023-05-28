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
	GetKelasByID(ctx context.Context, kelasID int) dto.KelasResponse
	CreateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error)
	UpdateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error)
	DeleteKelas(ctx context.Context, kelasID int) error
	AddSiswaToKelas(ctx context.Context, kelasID int, siswaIDs []int) error
	SiswaNaikKelas(ctx context.Context, kelasID int, siswaKelasIDs []int) error
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

	temp.Select("kelas.id as id, nama_kelas, jurusan_id, tahun_ajar_id, tahun_ajar, pegawai_id, nip, tipe_pegawai, nama").
		Joins("join pegawai on pegawai.id = kelas.pegawai_id").
		Joins("join users on users.id = pegawai.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Joins("join tahun_ajar on tahun_ajar.id = kelas.tahun_ajar_id").
		Order("nama_kelas asc").
		Find(&kelas)

	return kelas
}

func (r *kelasRepository) GetKelasByID(ctx context.Context, kelasID int) dto.KelasResponse {
	var kelas dto.KelasResponse
	r.db.Model(&entity.Kelas{}).
		Select("kelas.id as id, nama_kelas, jurusan_id, tahun_ajar_id, tahun_ajar, pegawai_id, nip, tipe_pegawai, nama").
		Where("kelas.id = ?", kelasID).
		Joins("join pegawai on pegawai.id = kelas.pegawai_id").
		Joins("join users on users.id = pegawai.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Joins("join tahun_ajar on tahun_ajar.id = kelas.tahun_ajar_id").
		Order("nama_kelas asc").
		First(&kelas)
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

	if err == nil {
		err = r.db.Where("kelas_id = ?", kelasID).Delete(&entity.SiswaKelas{}).Error
	}

	return err
}

func (r *kelasRepository) AddSiswaToKelas(ctx context.Context, kelasID int, siswaIDs []int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var siswaKelasInsert []entity.SiswaKelas

	var addedSiswaIDs []int
	var oldSiswa []int

	// to get added siswa to class
	tx.
		Model(&entity.SiswaKelas{}).
		Where("kelas_id = ?", kelasID).
		Where("siswa_id in ?", siswaIDs).
		Pluck("siswa_id", &addedSiswaIDs)

	// to get old siswa
	tx.
		Model(&entity.SiswaKelas{}).
		Where("siswa_id in ?", siswaIDs).
		Pluck("siswa_id", &oldSiswa)

	addedSiswaIDs = append(addedSiswaIDs, oldSiswa...)

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

	err := tx.Create(&siswaKelasInsert).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var poinSiswa []entity.PoinSiswa

	for _, siswaKelas := range siswaKelasInsert {
		poinSiswa = append(poinSiswa, entity.PoinSiswa{
			SiswaKelasID: siswaKelas.ID,
			Poin:         200,
		})
	}

	err = tx.Create(&poinSiswa).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *kelasRepository) SiswaNaikKelas(ctx context.Context, kelasID int, siswaKelasIDs []int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var siswaKelasInsert []entity.SiswaKelas

	var oldSiswa []dto.SiswaKelasPoin

	// get old poin of siswa
	tx.Model(&entity.SiswaKelas{}).
		Select("siswa_kelas_id, poin, siswa_id, kelas_id").
		Where("siswa_kelas.id in ?", siswaKelasIDs).
		Joins("join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id").
		Find(&oldSiswa)

	for _, siswa := range oldSiswa {
		siswaKelasInsert = append(siswaKelasInsert, entity.SiswaKelas{
			KelasID: kelasID,
			SiswaID: siswa.SiswaID,
		})
	}

	if len(siswaKelasInsert) == 0 {
		return nil
	}

	err := tx.Create(&siswaKelasInsert).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var poinSiswa []entity.PoinSiswa

	for i, siswaKelas := range siswaKelasInsert {
		poinSiswa = append(poinSiswa, entity.PoinSiswa{
			SiswaKelasID: siswaKelas.ID,
			Poin:         oldSiswa[i].Poin,
		})
	}

	err = tx.Create(&poinSiswa).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *kelasRepository) RemoveSiswaFromKelas(ctx context.Context, kelasID int, siswaIDs []int) error {
	err := r.db.
		Where("kelas_id = ?", kelasID).
		Where("siswa_id in ?", siswaIDs).
		Delete(&entity.SiswaKelas{}).Error

	return err
}
