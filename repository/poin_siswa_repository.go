package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type PoinSiswaRepository interface {
	GetPoinSiswa(ctx context.Context, siswaKelasID int) dto.PoinSiswaResponse
	GetPoinKelas(ctx context.Context, kelasID int) dto.PoinKelasResponse
	GetPoinJurusan(ctx context.Context, jurusanID int, tahunAjarID int) dto.PoinJurusanResponse
	AddPoinSiswa(ctx context.Context, req dto.PoinSiswaRequest) error
	UpdatePoinSiswa(ctx context.Context, poinLog entity.PoinLog) error
	DeletePoinSiswa(ctx context.Context, poinLogID int) error
	GetPoinSiswaPagination(ctx context.Context, page int, perPage int, order string, orderBy string, search string, tahunAjarID string) dto.PoinSiswaPagination
}

type poinSiswaRepository struct {
	db *gorm.DB
}

func NewPoinSiswaRepository(db *gorm.DB) PoinSiswaRepository {
	return &poinSiswaRepository{db: db}
}

func (r *poinSiswaRepository) GetPoinSiswa(ctx context.Context, siswaKelasID int) dto.PoinSiswaResponse {
	result := dto.PoinSiswaResponse{}
	poinSiswa := entity.PoinSiswa{}
	temp := r.db.Model(&poinSiswa)

	temp.Select("nis, nama, nama_kelas, poin, poin_siswa.created_at, poin_siswa.updated_at").
		Where("siswa_kelas.id = ?", siswaKelasID).
		Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
		Joins("join users on users.id = siswa.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Joins("join kelas on kelas.id = siswa_kelas.kelas_id")
	temp.First(&result)

	return result
}

func (r *poinSiswaRepository) GetPoinKelas(ctx context.Context, kelasID int) dto.PoinKelasResponse {
	result := dto.PoinKelasResponse{}
	kelas := entity.Kelas{}
	temp := r.db.Model(&kelas)

	temp.Select("kelas.id, tahun_ajar.id as tahun_ajar_id, nama_kelas, tahun_ajar.tahun_ajar, avg(poin) as poin").
		Where("kelas.id = ?", kelasID).
		Joins("join siswa_kelas on siswa_kelas.kelas_id = kelas.id").
		Joins("join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id").
		Joins("join tahun_ajar on tahun_ajar.id = kelas.tahun_ajar_id").
		Group("kelas.id, tahun_ajar.id").
		First(&result)

	return result
}

func (r *poinSiswaRepository) GetPoinJurusan(ctx context.Context, jurusanID int, tahunAjarID int) dto.PoinJurusanResponse {
	result := dto.PoinJurusanResponse{}
	jurusan := entity.Jurusan{}
	temp := r.db.Model(&jurusan)

	temp.Select("jurusan.id, tahun_ajar.id as tahun_ajar_id, nama_jurusan, tahun_ajar.tahun_ajar, avg(poin) as poin").
		Where("kelas.jurusan_id = ?", jurusanID).
		Where("kelas.tahun_ajar_id = ?", tahunAjarID).
		Joins("join kelas on kelas.jurusan_id = jurusan.id").
		Joins("join siswa_kelas on siswa_kelas.kelas_id = kelas.id").
		Joins("join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id").
		Joins("join tahun_ajar on tahun_ajar.id = kelas.tahun_ajar_id").
		Group("jurusan.id, nama_jurusan, tahun_ajar.id").
		First(&result)

	return result
}

func (r *poinSiswaRepository) AddPoinSiswa(ctx context.Context, req dto.PoinSiswaRequest) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	poinSiswa := entity.PoinSiswa{}

	err := tx.Where("siswa_kelas_id = ?", req.SiswaKelasID).First(&poinSiswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	poinBefore := poinSiswa.Poin
	poinAfter := poinBefore + req.Poin

	err = tx.Model(&entity.PoinSiswa{}).
		Where("id = ?", poinSiswa.ID).
		Update("poin", poinAfter).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	poinLog := entity.PoinLog{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Poin:        req.Poin,
		PoinBefore:  poinBefore,
		PoinAfter:   poinAfter,
		GuruID:      req.GuruID,
		PoinSiswaID: poinSiswa.ID,
	}

	err = tx.Create(&poinLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *poinSiswaRepository) UpdatePoinSiswa(ctx context.Context, poinLog entity.PoinLog) error {
	err := r.db.Model(&entity.PoinLog{}).
		Where("id = ?", poinLog.ID).
		Update("description", poinLog.Description).
		Error
	return err
}

func (r *poinSiswaRepository) DeletePoinSiswa(ctx context.Context, poinLogID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	poinLog := entity.PoinLog{
		ID: poinLogID,
	}

	err := tx.First(&poinLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	poinSiswa := entity.PoinSiswa{
		SiswaKelasID: poinLog.PoinSiswaID,
	}

	err = tx.First(&poinSiswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&entity.PoinSiswa{}).
		Where("id = ?", poinSiswa.ID).
		Update("poin", poinSiswa.Poin-poinLog.Poin).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&poinLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *poinSiswaRepository) GetPoinSiswaPagination(ctx context.Context, page int, perPage int, order string, orderBy string, search string, tahunAjarID string) dto.PoinSiswaPagination {
	result := dto.PoinSiswaPagination{}
	poinSiswa := entity.PoinSiswa{}
	temp := r.db.Model(&poinSiswa)
	if search != "" {
		search = "%" + search + "%"
		temp.Where("nama ilike ? or nis ilike ?", search, search)
	}

	if tahunAjarID != "" {
		temp.Where("kelas.tahun_ajar_id = ?", tahunAjarID)
	}

	temp.Select("nis, nama, nama_kelas, poin, poin_siswa.created_at, poin_siswa.updated_at").
		Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
		Joins("join users on users.id = siswa.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Joins("join kelas on kelas.id = siswa_kelas.kelas_id").
		Order(orderBy + " " + order)

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
