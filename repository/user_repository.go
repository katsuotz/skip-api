package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, email string) entity.User
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) entity.User {
	user := entity.User{}
	r.db.
		Select("id, username, role, password").
		Where("username = ?", username).
		Preload("Profile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nama, tempat_lahir, tanggal_lahir, jenis_kelamin, user_id")
		}).
		First(&user)
	return user
}
