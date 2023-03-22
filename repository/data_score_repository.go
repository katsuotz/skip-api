package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type DataScoreRepository interface {
	GetDataScore(ctx context.Context) []entity.DataScore
	CreateDataScore(ctx context.Context, dataScore entity.DataScore) (entity.DataScore, error)
	UpdateDataScore(ctx context.Context, dataScore entity.DataScore) (entity.DataScore, error)
	DeleteDataScore(ctx context.Context, dataScoreID int) error
}

type dataScoreRepository struct {
	db *gorm.DB
}

func NewDataScoreRepository(db *gorm.DB) DataScoreRepository {
	return &dataScoreRepository{db: db}
}

func (r *dataScoreRepository) GetDataScore(ctx context.Context) []entity.DataScore {
	var dataScore []entity.DataScore
	r.db.Find(&dataScore)
	return dataScore
}

func (r *dataScoreRepository) CreateDataScore(ctx context.Context, dataScore entity.DataScore) (entity.DataScore, error) {
	err := r.db.Create(&dataScore).Error
	return dataScore, err
}

func (r *dataScoreRepository) UpdateDataScore(ctx context.Context, dataScore entity.DataScore) (entity.DataScore, error) {
	err := r.db.Updates(&dataScore).Error
	return dataScore, err
}

func (r *dataScoreRepository) DeleteDataScore(ctx context.Context, dataScoreID int) error {
	err := r.db.Where("id = ?", dataScoreID).Delete(&entity.DataScore{}).Error
	return err
}
