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
	UpdateScoreSiswa(ctx context.Context, scoreSiswa entity.ScoreSiswa) (entity.ScoreSiswa, error)
	DeleteScoreSiswa(ctx context.Context, scoreSiswaID int) error
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
		Score:        50,
		SiswaKelasID: req.SiswaKelasID,
	}

	err := tx.FirstOrCreate(&scoreSiswa, entity.ScoreSiswa{
		SiswaKelasID: req.SiswaKelasID,
	}).Error

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

	return nil
}

func (r *scoreSiswaRepository) UpdateScoreSiswa(ctx context.Context, scoreSiswa entity.ScoreSiswa) (entity.ScoreSiswa, error) {
	err := r.db.Updates(&scoreSiswa).Error
	return scoreSiswa, err
}

func (r *scoreSiswaRepository) DeleteScoreSiswa(ctx context.Context, scoreSiswaID int) error {
	err := r.db.Where("id = ?", scoreSiswaID).Delete(&entity.ScoreSiswa{}).Error
	return err
}
