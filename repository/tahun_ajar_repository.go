package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type TahunAjarRepository interface {
	GetTahunAjar(ctx context.Context, page int, perPage int, search string) dto.TahunAjarPagination
	CreateTahunAjar(ctx context.Context, tahunAjar entity.TahunAjar) (entity.TahunAjar, error)
	UpdateTahunAjar(ctx context.Context, tahunAjar entity.TahunAjar) (entity.TahunAjar, error)
	DeleteTahunAjar(ctx context.Context, tahunAjarID int) error
	SetActiveTahunAjar(ctx context.Context, tahunAjarID int) error
}

type tahunAjarRepository struct {
	db *gorm.DB
}

func NewTahunAjarRepository(db *gorm.DB) TahunAjarRepository {
	return &tahunAjarRepository{db: db}
}

func (r *tahunAjarRepository) GetTahunAjar(ctx context.Context, page int, perPage int, search string) dto.TahunAjarPagination {
	result := dto.TahunAjarPagination{}
	tahunAjar := entity.TahunAjar{}
	temp := r.db.Model(&tahunAjar)
	if search != "" {
		search = "%" + search + "%"
		temp.Where("tahun_ajar ilike ?", search)
		//temp.Or("semester ilike ?", search)
	}

	temp.Order("tahun_ajar desc")
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

func (r *tahunAjarRepository) CreateTahunAjar(ctx context.Context, tahunAjar entity.TahunAjar) (entity.TahunAjar, error) {
	err := r.db.Create(&tahunAjar).Error
	return tahunAjar, err
}

func (r *tahunAjarRepository) UpdateTahunAjar(ctx context.Context, tahunAjar entity.TahunAjar) (entity.TahunAjar, error) {
	err := r.db.Updates(&tahunAjar).Error
	return tahunAjar, err
}

func (r *tahunAjarRepository) DeleteTahunAjar(ctx context.Context, tahunAjarID int) error {
	err := r.db.Where("id = ?", tahunAjarID).Delete(&entity.TahunAjar{}).Error
	return err
}

func (r *tahunAjarRepository) SetActiveTahunAjar(ctx context.Context, tahunAjarID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(&entity.TahunAjar{}).Where("is_active = ?", true).Update("is_active", false).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = r.db.Model(&entity.TahunAjar{}).Where("id = ?", tahunAjarID).Update("is_active", true).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
