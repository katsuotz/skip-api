package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type JurusanRepository interface {
	GetJurusan(ctx context.Context) []entity.Jurusan
	CreateJurusan(ctx context.Context, jurusan entity.Jurusan) (entity.Jurusan, error)
	UpdateJurusan(ctx context.Context, jurusan entity.Jurusan) (entity.Jurusan, error)
	DeleteJurusan(ctx context.Context, jurusanID int) error
}

type jurusanRepository struct {
	db *gorm.DB
}

func NewJurusanRepository(db *gorm.DB) JurusanRepository {
	return &jurusanRepository{db: db}
}

func (r *jurusanRepository) GetJurusan(ctx context.Context) []entity.Jurusan {
	var jurusan []entity.Jurusan
	r.db.Order("nama_jurusan asc").Find(&jurusan)
	return jurusan
}

func (r *jurusanRepository) CreateJurusan(ctx context.Context, jurusan entity.Jurusan) (entity.Jurusan, error) {
	err := r.db.Create(&jurusan).Error
	return jurusan, err
}

func (r *jurusanRepository) UpdateJurusan(ctx context.Context, jurusan entity.Jurusan) (entity.Jurusan, error) {
	err := r.db.Updates(&jurusan).Error
	return jurusan, err
}

func (r *jurusanRepository) DeleteJurusan(ctx context.Context, jurusanID int) error {
	err := r.db.Where("id = ?", jurusanID).Delete(&entity.Jurusan{}).Error
	return err
}
