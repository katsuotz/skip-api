package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type PoinSiswaRepository interface {
	GetPoinSiswa(ctx context.Context) []entity.PoinSiswa
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

func (r *poinSiswaRepository) GetPoinSiswa(ctx context.Context) []entity.PoinSiswa {
	var poinSiswa []entity.PoinSiswa
	r.db.Find(&poinSiswa)
	return poinSiswa
}

func (r *poinSiswaRepository) AddPoinSiswa(ctx context.Context, req dto.PoinSiswaRequest) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	poinSiswa := entity.PoinSiswa{
		SiswaKelasID: req.SiswaKelasID,
	}

	err := tx.First(&poinSiswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&entity.PoinSiswa{}).
		Where("id = ?", poinSiswa.ID).
		Update("poin", poinSiswa.Poin+req.Poin).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	poinLog := entity.PoinLog{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Poin:        req.Poin,
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
