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
	r.db.Where("username = ?", username).First(&user)
	return user
}
