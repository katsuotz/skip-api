package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	FindByProfileID(ctx context.Context, userID int) entity.Profile
	UpdateProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error)
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) FindByProfileID(ctx context.Context, userID int) entity.Profile {
	profile := entity.Profile{}
	r.db.Where("user_id = ?", userID).First(&profile)
	return profile
}

func (r *profileRepository) UpdateProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {
	err := r.db.Updates(&profile).Error
	return profile, err
}
