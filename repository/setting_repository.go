package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type SettingRepository interface {
	GetSetting(ctx context.Context) []entity.Setting
	CreateSetting(ctx context.Context, setting entity.Setting) (entity.Setting, error)
	UpdateSetting(ctx context.Context, setting entity.Setting) (entity.Setting, error)
	DeleteSetting(ctx context.Context, settingID int) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) GetSetting(ctx context.Context) []entity.Setting {
	var setting []entity.Setting
	r.db.Find(&setting)
	return setting
}

func (r *settingRepository) CreateSetting(ctx context.Context, setting entity.Setting) (entity.Setting, error) {
	err := r.db.Create(&setting).Error
	return setting, err
}

func (r *settingRepository) UpdateSetting(ctx context.Context, setting entity.Setting) (entity.Setting, error) {
	err := r.db.Updates(&setting).Error
	return setting, err
}

func (r *settingRepository) DeleteSetting(ctx context.Context, settingID int) error {
	err := r.db.
		Where("id = ?", settingID).
		Where("can_delete = ?", true).
		Delete(&entity.Setting{}).Error
	return err
}
