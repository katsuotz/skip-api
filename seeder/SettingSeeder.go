package seeder

import (
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

func CreateSetting(
	db *gorm.DB,
	key string,
	value string,
	canDelete bool,
) {
	setting := entity.Setting{
		Key:       key,
		Value:     value,
		CanDelete: canDelete,
	}
	if db.Model(&setting).Where("key = ?", key).Updates(&setting).RowsAffected == 0 {
		db.Create(&setting)
	}
}
