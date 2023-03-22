package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type ScoreSiswaRepository interface {
	GetScoreSiswa(ctx context.Context) []entity.ScoreSiswa
	AddScoreSiswa(ctx context.Context, req dto.ScoreSiswaRequest) error
	UpdateScoreSiswa(ctx context.Context, scoreLog entity.ScoreLog) error
	DeleteScoreSiswa(ctx context.Context, scoreLogID int) error
}

type scoreSiswaRepository struct {
	db *gorm.DB
}

func NewScoreSiswaRepository(db *gorm.DB) ScoreSiswaRepository {
	return &scoreSiswaRepository{db: db}
}

func (r *scoreSiswaRepository) GetScoreSiswa(ctx context.Context) []entity.ScoreSiswa {
	var scoreSiswa []entity.ScoreSiswa
	r.db.Find(&scoreSiswa)
	return scoreSiswa
}

func (r *scoreSiswaRepository) AddScoreSiswa(ctx context.Context, req dto.ScoreSiswaRequest) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	scoreSiswa := entity.ScoreSiswa{
		SiswaKelasID: req.SiswaKelasID,
	}

	err := tx.First(&scoreSiswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&entity.ScoreSiswa{}).
		Where("id = ?", scoreSiswa.ID).
		Update("score", scoreSiswa.Score+req.Score).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	scoreLog := entity.ScoreLog{
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Score:        req.Score,
		GuruID:       req.GuruID,
		ScoreSiswaID: scoreSiswa.ID,
	}

	err = tx.Create(&scoreLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *scoreSiswaRepository) UpdateScoreSiswa(ctx context.Context, scoreLog entity.ScoreLog) error {
	err := r.db.Model(&entity.ScoreLog{}).
		Where("id = ?", scoreLog.ID).
		Update("description", scoreLog.Description).
		Error
	return err
}

func (r *scoreSiswaRepository) DeleteScoreSiswa(ctx context.Context, scoreLogID int) error {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	scoreLog := entity.ScoreLog{
		ID: scoreLogID,
	}

	err := tx.First(&scoreLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	scoreSiswa := entity.ScoreSiswa{
		SiswaKelasID: scoreLog.ScoreSiswaID,
	}

	err = tx.First(&scoreSiswa).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&entity.ScoreSiswa{}).
		Where("id = ?", scoreSiswa.ID).
		Update("score", scoreSiswa.Score-scoreLog.Score).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&scoreLog).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
