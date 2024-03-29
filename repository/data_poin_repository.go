package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type DataPoinRepository interface {
	GetDataPoin(ctx context.Context, page int, perPage int, search string, poinType string, category string) dto.DataPoinPagination
	GetDataPoinByID(ctx context.Context, dataPoinID int) entity.DataPoin
	CreateDataPoin(ctx context.Context, dataPoin entity.DataPoin) (entity.DataPoin, error)
	UpdateDataPoin(ctx context.Context, dataPoin entity.DataPoin) (entity.DataPoin, error)
	DeleteDataPoin(ctx context.Context, dataPoinID int) error
}

type dataPoinRepository struct {
	db *gorm.DB
}

func NewDataPoinRepository(db *gorm.DB) DataPoinRepository {
	return &dataPoinRepository{db: db}
}

func (r *dataPoinRepository) GetDataPoin(ctx context.Context, page int, perPage int, search string, poinType string, category string) dto.DataPoinPagination {
	result := dto.DataPoinPagination{}
	dataPoin := entity.DataPoin{}
	temp := r.db.Model(&dataPoin)

	if search != "" {
		search = "%" + search + "%"
		temp.Where("title ilike ? or description ilike ?", search, search)
	}

	if poinType != "" {
		temp.Where("type = ?", poinType)
	}

	if category != "" {
		temp.Where("category = ?", category)
	}

	temp.Order("title asc")
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

func (r *dataPoinRepository) GetDataPoinByID(ctx context.Context, dataPoinID int) entity.DataPoin {
	dataPoin := entity.DataPoin{
		ID: dataPoinID,
	}
	r.db.First(&dataPoin)
	return dataPoin
}

func (r *dataPoinRepository) CreateDataPoin(ctx context.Context, dataPoin entity.DataPoin) (entity.DataPoin, error) {
	err := r.db.Create(&dataPoin).Error
	return dataPoin, err
}

func (r *dataPoinRepository) UpdateDataPoin(ctx context.Context, dataPoin entity.DataPoin) (entity.DataPoin, error) {
	err := r.db.Updates(&dataPoin).Error
	return dataPoin, err
}

func (r *dataPoinRepository) DeleteDataPoin(ctx context.Context, dataPoinID int) error {
	err := r.db.Where("id = ?", dataPoinID).Delete(&entity.DataPoin{}).Error
	return err
}
