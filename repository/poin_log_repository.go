package repository

import (
	"context"
	"fmt"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

type PoinLogRepository interface {
	GetPoinSiswaLog(ctx context.Context, page int, perPage int, siswaKelasID int) dto.PoinLogPagination
	GetPoinLogSiswaByKelas(ctx context.Context, nis string) []dto.PoinLogSiswaByKelas
	CountPoin(ctx context.Context, poinType string, kelasID string, jurusanID string, tahunAjarID string) dto.CountResponse
	GetPoinLogPagination(ctx context.Context, page int, perPage int, order string, orderBy string, tahunAjarID string) dto.PoinLogPagination
	GetCountPoinLogPagination(ctx context.Context, page int, perPage int, order string, orderBy string, groupBy string, tahunAjarID string, poinType string) dto.PoinLogCountPagination
	GetCountPoinLogPaginationByMonth(ctx context.Context, tahunAjarID string, poinType string) []dto.PoinLogCountGraphResponse
}

type poinLogRepository struct {
	db *gorm.DB
}

func NewPoinLogRepository(db *gorm.DB) PoinLogRepository {
	return &poinLogRepository{db: db}
}

func (r *poinLogRepository) GetPoinSiswaLog(ctx context.Context, page int, perPage int, siswaKelasID int) dto.PoinLogPagination {
	result := dto.PoinLogPagination{}
	poinLog := entity.PoinLog{}
	temp := r.db.Model(&poinLog)

	fmt.Println(1231232)

	temp.Select("poin_log.id as id, title, description, poin_log.poin, poin_before, poin_after, type, file, pegawai_id, nip, profiles.nama as nama_pegawai, poin_log.created_at, poin_log.updated_at").
		Where("siswa_kelas.id = ?", siswaKelasID).
		Joins("join pegawai on pegawai.id = poin_log.pegawai_id").
		Joins("join users on users.id = pegawai.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Joins("join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id").
		Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Order("poin_log.created_at desc")
	//temp.Joins("join siswa on siswa.id = siswa_kelas.siswa_id")
	temp.Offset(perPage * (page - 1)).Limit(perPage).
		Find(&result.Data)

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

func (r *poinLogRepository) GetPoinLogSiswaByKelas(ctx context.Context, nis string) []dto.PoinLogSiswaByKelas {
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
			Select("poin_log.id as id, title, description, poin_log.poin, poin_before, poin_after, type, pegawai_id, nip, profiles.nama as nama_pegawai, poin_log.created_at, poin_log.updated_at").
			Where("siswa_kelas.id = ?", siswa.ID).
			Joins("join pegawai on pegawai.id = poin_log.pegawai_id").
			Joins("join users on users.id = pegawai.user_id").
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

func (r *poinLogRepository) CountPoin(ctx context.Context, poinType string, kelasID string, jurusanID string, tahunAjarID string) dto.CountResponse {
	result := dto.CountResponse{}

	temp := r.db.Model(&entity.PoinLog{}).
		Select("count(*)")

	if poinType != "" {
		temp.Where("type = ?", poinType)
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

	temp.
		Where("poin_siswa.deleted_at is NULL").
		Where("siswa_kelas.deleted_at is NULL").
		Where("kelas.deleted_at is NULL").
		Joins("left join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id").
		Joins("left join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("left join kelas on kelas.id = siswa_kelas.kelas_id")

	temp.Scan(&result.Total)

	return result
}

func (r *poinLogRepository) GetPoinLogPagination(ctx context.Context, page int, perPage int, order string, orderBy string, tahunAjarID string) dto.PoinLogPagination {
	result := dto.PoinLogPagination{}
	poinLog := entity.PoinLog{}
	temp := r.db.Model(&poinLog)

	if tahunAjarID != "" {
		temp.Where("kelas.tahun_ajar_id = ?", tahunAjarID)
	}

	temp.Select("poin_log.id as id, title, description, poin_log.poin, poin_before, poin_after, type, poin_log.pegawai_id, nip, pg.nama as nama_pegawai, nis, ps.nama as nama, ps.foto as foto, file, poin_log.created_at, poin_log.updated_at").
		Joins("join pegawai on pegawai.id = poin_log.pegawai_id").
		Joins("join users ug on ug.id = pegawai.user_id").
		Joins("join profiles pg on pg.user_id = ug.id").
		Joins("join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id").
		Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("join kelas on kelas.id = siswa_kelas.kelas_id").
		Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
		Joins("join users us on us.id = siswa.user_id").
		Joins("join profiles ps on ps.user_id = us.id").
		Order("poin_log." + orderBy + " " + order)

	temp.Offset(perPage * (page - 1)).Limit(perPage).
		Find(&result.Data)

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

func (r *poinLogRepository) GetCountPoinLogPagination(ctx context.Context, page int, perPage int, order string, orderBy string, groupBy string, tahunAjarID string, poinType string) dto.PoinLogCountPagination {
	result := dto.PoinLogCountPagination{}
	poinLog := entity.PoinLog{}
	temp := r.db.Model(&poinLog)

	if tahunAjarID != "" {
		temp.Where("kelas.tahun_ajar_id = ?", tahunAjarID)
	}

	if poinType != "" {
		temp.Where("poin_log.type = ?", poinType)
	}

	temp.
		Joins("join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id").
		Joins("join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("join kelas on kelas.id = siswa_kelas.kelas_id")

	groupQuery := ""
	selectQuery := "count(*) as total, type"

	if groupBy == "siswa" {
		groupQuery += "poin_log.poin_siswa_id, profiles.nama, nis"
		selectQuery += ", nama, nis"
		temp.
			Joins("join siswa on siswa.id = siswa_kelas.siswa_id").
			Joins("join users on users.id = siswa.user_id").
			Joins("join profiles on profiles.user_id = users.id")
	} else if groupBy == "type" {
		groupQuery += "poin_log.title"
		selectQuery += ", poin_log.title"
	}

	groupQuery += ", poin_log.type"

	temp.
		Group(groupQuery)

	temp.Select(selectQuery).
		Order(orderBy + " " + order).Offset(perPage * (page - 1)).Limit(perPage).
		Find(&result.Data)

	var totalItem int64 = 1

	if len(result.Data) != 1 || page != 1 {
		temp.Select("count(*) as total").Offset(-1).Limit(-1).Count(&totalItem)
	}

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

func (r *poinLogRepository) GetCountPoinLogPaginationByMonth(ctx context.Context, tahunAjarID string, poinType string) []dto.PoinLogCountGraphResponse {
	var result []dto.PoinLogCountGraphResponse

	tahunAjar := entity.TahunAjar{}
	r.db.Model(&entity.TahunAjar{}).Where("id = ?", tahunAjarID).First(&tahunAjar)

	tahunAjarSplice := strings.Split(tahunAjar.TahunAjar, "/")

	startMonth := 7
	startYear, _ := strconv.Atoi(tahunAjarSplice[0])
	endMonth := 6
	endYear := startYear + 1

	poinLog := entity.PoinLog{}
	temp := r.db.Model(&poinLog)
	temp.Where("kelas.tahun_ajar_id = ?", tahunAjarID)

	if poinType != "" {
		temp.Where("poin_log.type = ?", poinType)
	}

	temp.
		Select("count(*) as total, EXTRACT(YEAR FROM poin_log.created_at) as year, EXTRACT(MONTH FROM poin_log.created_at) as month").
		Where("(EXTRACT(YEAR FROM poin_log.created_at) >= ? and EXTRACT(MONTH FROM poin_log.created_at) >= ?) or (EXTRACT(YEAR FROM poin_log.created_at) <= ? and EXTRACT(MONTH FROM poin_log.created_at) <= ?)", startYear, startMonth, endYear, endMonth).
		Joins("left join poin_siswa on poin_siswa.id = poin_log.poin_siswa_id").
		Joins("left join siswa_kelas on siswa_kelas.id = poin_siswa.siswa_kelas_id").
		Joins("left join kelas on kelas.id = siswa_kelas.kelas_id").
		Order("month, year").
		Group("year, month").
		Find(&result)

	return result
}
