package repository

import (
	"context"
	"errors"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
	"math"
	"time"
	_ "time/tzdata"
)

type PoinSiswaRepository interface {
	GetPoinSiswa(ctx context.Context, siswaKelasID int) dto.PoinSiswaResponse
	GetPoinKelas(ctx context.Context, kelasID int) dto.PoinKelasResponse
	GetPoinJurusan(ctx context.Context, jurusanID string, tahunAjarID string) dto.PoinJurusanWithKelasResponse
	AddPoinSiswa(ctx context.Context, req dto.PoinSiswaRequest) error
	UpdatePoinSiswa(ctx context.Context, poinLog entity.PoinLog) error
	DeletePoinSiswa(ctx context.Context, poinLogID int) error
	GetPoinSiswaPagination(ctx context.Context, page int, perPage int, order string, orderBy string, search string, tahunAjarID string, pegawaiID int, maxPoin string) dto.PoinSiswaPagination
	CountPoinSiswa(ctx context.Context, countType string, kelasID string, jurusanID string, tahunAjarID string, pegawaiID int) dto.CountResponse
	CountPoinSiswaTotal(ctx context.Context, kelasID string, jurusanID string, tahunAjarID string, pegawaiID int, maxPoin string) dto.CountResponse
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

	temp.Select("nis, nama, nama_kelas, poin, poin_siswa.created_at, poin_siswa.updated_at, foto").
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

func (r *poinSiswaRepository) GetPoinJurusan(ctx context.Context, jurusanID string, tahunAjarID string) dto.PoinJurusanWithKelasResponse {
	result := dto.PoinJurusanWithKelasResponse{}

	r.db.Model(&entity.Jurusan{}).
		Select("nama_jurusan").
		Where("id = ?", jurusanID).
		First(&result.Jurusan)

	tahunAjar := entity.TahunAjar{}

	r.db.Model(&tahunAjar).
		Where("id = ?", tahunAjarID).
		First(&tahunAjar)

	result.Jurusan.TahunAjar = tahunAjar.TahunAjar

	r.db.Model(&entity.Kelas{}).Select("kelas.id, nama_kelas, avg(poin) as poin").
		Where("kelas.jurusan_id = ?", jurusanID).
		Where("kelas.tahun_ajar_id = ?", tahunAjarID).
		Joins("join siswa_kelas on siswa_kelas.kelas_id = kelas.id").
		Joins("join poin_siswa on poin_siswa.siswa_kelas_id = siswa_kelas.id").
		Order("nama_kelas").
		Group("kelas.id").
		Find(&result.Data)

	var poin float64 = 0

	for i, _ := range result.Data {
		result.Data[i].Poin = math.Round(result.Data[i].Poin*100) / 100
		poin += result.Data[i].Poin
	}

	result.Jurusan.Poin = poin

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

	jakarta, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now().In(jakarta)
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location())

	poinLogBefore := entity.PoinLog{}
	err = tx.
		Where("poin_siswa_id = ?", poinSiswa.ID).
		Where("created_at >= ?", startOfDay).
		Where("created_at < ?", endOfDay).
		Where("title = ?", req.Title).
		First(&poinLogBefore).Error

	if poinLogBefore.ID != 0 {
		tx.Rollback()
		return errors.New(req.Type + " yang sama sudah dilaporkan untuk hari ini")
	}

	poinBefore := poinSiswa.Poin
	poinAfter := poinBefore

	if req.Type == "Pelanggaran" {
		poinAfter -= req.Poin
	} else if req.Type == "Penghargaan" {
		poinAfter += req.Poin
	}

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
		Penanganan:  req.Penanganan,
		Type:        req.Type,
		Poin:        req.Poin,
		PoinBefore:  poinBefore,
		PoinAfter:   poinAfter,
		File:        req.File,
		PegawaiID:   req.PegawaiID,
		DataPoinID:  req.DataPoinID,
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
		Update("tindak_lanjut", poinLog.TindakLanjut).
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
		ID: poinLog.PoinSiswaID,
	}

	err = tx.First(&poinSiswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	currentPoin := poinSiswa.Poin

	if poinLog.Type == "Pelanggaran" {
		currentPoin += poinLog.Poin
	} else if poinLog.Type == "Penghargaan" {
		currentPoin -= poinLog.Poin
	}

	err = tx.Model(&entity.PoinSiswa{}).
		Where("id = ?", poinSiswa.ID).
		Update("poin", currentPoin).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	var oldPoinLog []entity.PoinLog
	err = r.db.
		Where("poin_siswa_id = ?", poinLog.PoinSiswaID).
		Where("created_at > ?", poinLog.CreatedAt).
		Find(&oldPoinLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, oldLog := range oldPoinLog {
		if poinLog.Type == "Pelanggaran" {
			oldLog.PoinBefore += poinLog.Poin
			oldLog.PoinAfter += poinLog.Poin
		} else if poinLog.Type == "Penghargaan" {
			oldLog.PoinBefore -= poinLog.Poin
			oldLog.PoinAfter -= poinLog.Poin
		}
		err = tx.Save(oldLog).Error

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Delete(&poinLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *poinSiswaRepository) GetPoinSiswaPagination(ctx context.Context, page int, perPage int, order string, orderBy string, search string, tahunAjarID string, pegawaiID int, maxPoin string) dto.PoinSiswaPagination {
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

	if pegawaiID != 0 {
		temp.Where("kelas.pegawai_id = ?", pegawaiID)
	}

	if maxPoin != "" {
		temp.Where("poin <= ?", maxPoin)
	}

	temp.Select("nis, nama, foto, nama_kelas, poin, poin_siswa.created_at, poin_siswa.updated_at").
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

func (r *poinSiswaRepository) CountPoinSiswa(ctx context.Context, countType string, kelasID string, jurusanID string, tahunAjarID string, pegawaiID int) dto.CountResponse {
	result := dto.CountResponse{}

	temp := r.db.Model(&entity.PoinSiswa{})

	if countType == "max" {
		temp.Select("max(poin_siswa.poin)")
	} else if countType == "min" {
		temp.Select("min(poin_siswa.poin)")
	} else if countType == "avg" {
		temp.Select("avg(poin_siswa.poin)")
	}

	if kelasID != "" {
		temp.Where("kelas.id = ?", kelasID)
	}

	if tahunAjarID != "" {
		temp.Where("kelas.tahun_ajar_id = ?", tahunAjarID)
	}

	if jurusanID != "" {
		temp.Where("kelas.jurusan_id = ?", jurusanID)
	}

	if pegawaiID != 0 {
		temp.Where("kelas.pegawai_id = ?", pegawaiID)
	}

	temp.
		Where("poin_siswa.deleted_at is NULL").
		Where("siswa_kelas.deleted_at is NULL").
		Where("kelas.deleted_at is NULL").
		Joins("left join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("left join kelas on kelas.id = siswa_kelas.kelas_id")

	temp.Scan(&result.Total)

	return result
}

func (r *poinSiswaRepository) CountPoinSiswaTotal(ctx context.Context, kelasID string, jurusanID string, tahunAjarID string, pegawaiID int, maxPoin string) dto.CountResponse {
	result := dto.CountResponse{}

	temp := r.db.Model(&entity.PoinSiswa{})

	temp.Select("count(*)")

	if kelasID != "" {
		temp.Where("kelas.id = ?", kelasID)
	}

	if tahunAjarID != "" {
		temp.Where("kelas.tahun_ajar_id = ?", tahunAjarID)
	}

	if jurusanID != "" {
		temp.Where("kelas.jurusan_id = ?", jurusanID)
	}

	if pegawaiID != 0 {
		temp.Where("kelas.pegawai_id = ?", pegawaiID)
	}

	if maxPoin != "" {
		temp.Where("poin_siswa.poin <= ?", maxPoin)
	}

	temp.
		Where("poin_siswa.deleted_at is NULL").
		Where("siswa_kelas.deleted_at is NULL").
		Where("kelas.deleted_at is NULL").
		Joins("left join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("left join kelas on kelas.id = siswa_kelas.kelas_id")

	temp.Scan(&result.Total)

	return result
}
