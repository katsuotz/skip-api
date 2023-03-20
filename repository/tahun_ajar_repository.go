package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type TahunAjarRepository interface {
	GetTahunAjar(ctx context.Context) []entity.TahunAjar
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

func (r *tahunAjarRepository) GetTahunAjar(ctx context.Context) []entity.TahunAjar {
	var tahunAjar []entity.TahunAjar
	r.db.Find(&tahunAjar)
	return tahunAjar
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
	err := r.db.Where("tahun_ajar_id = ?", tahunAjarID).Delete(&entity.TahunAjar{}).Error
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
	err = r.db.Model(&entity.TahunAjar{}).Where("tahun_ajar_id = ?", tahunAjarID).Update("is_active", true).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
