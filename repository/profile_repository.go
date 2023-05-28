package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	FindProfileByID(ctx context.Context, userID int) entity.Profile
	FindProfileWithJoinByID(ctx context.Context, userID int) dto.UserResponse
	UpdateProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error)
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) FindProfileByID(ctx context.Context, userID int) entity.Profile {
	profile := entity.Profile{}
	r.db.Where("user_id = ?", userID).First(&profile)
	return profile
}

func (r *profileRepository) FindProfileWithJoinByID(ctx context.Context, userID int) dto.UserResponse {
	profile := dto.UserResponse{}
	r.db.
		Model(&entity.Profile{}).
		Select("users.id as id, nis, nip, nama, jenis_kelamin, tempat_lahir, tanggal_lahir, username, role, foto").
		Where("profiles.user_id = ?", userID).
		Joins("join users on users.id = profiles.user_id").
		Joins("left join pegawai on pegawai.user_id = users.id").
		Joins("left join siswa on siswa.user_id = users.id").
		First(&profile)
	return profile
}

func (r *profileRepository) UpdateProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {
	err := r.db.Updates(&profile).Error
	return profile, err
}
