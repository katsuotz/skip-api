package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
	"sort"
)

type PoinSiswaRepository interface {
	GetPoinSiswa(ctx context.Context, siswaKelasID int) dto.PoinSiswaResponse
	GetPoinSiswaLog(ctx context.Context, page int, perPage int, siswaKelasID int) dto.PoinLogPagination
	GetPoinLogSiswaByKelas(ctx context.Context, nis string) []dto.PoinLogSiswaByKelas
	AddPoinSiswa(ctx context.Context, req dto.PoinSiswaRequest) error
	UpdatePoinSiswa(ctx context.Context, poinLog entity.PoinLog) error
	DeletePoinSiswa(ctx context.Context, poinLogID int) error
}

type poinSiswaRepository struct {
	db *gorm.DB
}

func NewPoinSiswaRepository(db *gorm.DB) PoinSiswaRepository {
	return &poinSiswaRepository{db: db}
}

func (r *poinSiswaRepository) GetPoinSiswa(ctx context.Context, siswaKelasID int) dto.PoinSiswaResponse {
	result := dto.PoinSiswaResponse{}
	poinLog := entity.PoinSiswa{}
	temp := r.db.Model(&poinLog)

	temp.Select("nis, nama, nama_kelas, poin, poin_siswa.created_at, poin_siswa.updated_at")
	temp.Where("siswa_kelas.id = ?", siswaKelasID)
	temp.Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id")
	temp.Joins("join siswa on siswa.id = siswa_kelas.siswa_id")
	temp.Joins("join users on users.id = siswa.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")
	temp.Joins("join kelas on kelas.id = siswa_kelas.kelas_id")
	temp.First(&result)

	return result
}

func (r *poinSiswaRepository) GetPoinSiswaLog(ctx context.Context, page int, perPage int, siswaKelasID int) dto.PoinLogPagination {
	result := dto.PoinLogPagination{}
	poinLog := entity.PoinLog{}
	temp := r.db.Model(&poinLog)

	temp.Select("poin_log.id as id, title, description, poin_log.poin, type, guru_id, nip, profiles.nama as nama_guru, poin_log.created_at, poin_log.updated_at")
	temp.Where("siswa_kelas.id = ?", siswaKelasID)
	temp.Joins("join guru on guru.id = poin_log.guru_id")
	temp.Joins("join users on users.id = guru.user_id")
	temp.Joins("join profiles on profiles.user_id = users.id")
	temp.Joins("join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id")
	temp.Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id")
	//temp.Joins("join siswa on siswa.id = siswa_kelas.siswa_id")
	temp.Order("poin_log.created_at desc")
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

func (r *poinSiswaRepository) GetPoinLogSiswaByKelas(ctx context.Context, nis string) []dto.PoinLogSiswaByKelas {
	var result []dto.PoinLogSiswaByKelas

	var siswaKelas []entity.SiswaKelas
	r.db.Model(&siswaKelas).
		Where("siswa.nis = ?", nis).
		Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
		Find(&siswaKelas)

	for _, siswa := range siswaKelas {
		data := dto.PoinLogSiswaByKelas{}

		kelas := entity.Kelas{}

		r.db.Model(&kelas).
			Select("kelas.*, tahun_ajar.tahun_ajar").
			Where("kelas.id = ?", siswa.KelasID).
			Joins("join tahun_ajar on tahun_ajar.id = kelas.tahun_ajar_id").
			First(&data.Kelas)

		poinLog := entity.PoinLog{}

		r.db.Model(&poinLog).
			Select("poin_log.id as id, title, description, poin_log.poin, type, guru_id, nip, profiles.nama as nama_guru, poin_log.created_at, poin_log.updated_at").
			Where("siswa_kelas.id = ?", siswa.ID).
			Joins("join guru on guru.id = poin_log.guru_id").
			Joins("join users on users.id = guru.user_id").
			Joins("join profiles on profiles.user_id = users.id").
			Joins("join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id").
			Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
			Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
			Order("poin_log.created_at desc").
			Find(&data.Data)

		result = append(result, data)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Kelas.TahunAjar > result[j].Kelas.TahunAjar
	})

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
