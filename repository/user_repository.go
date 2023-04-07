package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, email string) entity.User
	LoginLog(ctx context.Context, log entity.LoginLog) error
	SuccessfulLoginLog(ctx *gin.Context, userID int) error
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

func (r *userRepository) LoginLog(ctx context.Context, log entity.LoginLog) error {
	err := r.db.Create(&log).Error
	return err
}

func (r *userRepository) SuccessfulLoginLog(ctx *gin.Context, userID int) error {
	location := helper.GetUserLocation(ctx)

	log := entity.LoginLog{
		UserID:    userID,
		Action:    "Successful Login",
		UserAgent: ctx.Request.UserAgent(),
		Location:  location,
	}

	return r.LoginLog(ctx, log)
}
