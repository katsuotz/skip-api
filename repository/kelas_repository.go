package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type KelasRepository interface {
	GetKelas(ctx context.Context) []entity.Kelas
	CreateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error)
	UpdateKelas(ctx context.Context, kelas entity.Kelas) (entity.Kelas, error)
	DeleteKelas(ctx context.Context, kelasID int) error
}

type kelasRepository struct {
	db *gorm.DB
}

func NewKelasRepository(db *gorm.DB) KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) GetKelas(ctx context.Context) []entity.Kelas {
	var kelas []entity.Kelas
	r.db.Find(&kelas)
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
