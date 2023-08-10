package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, email string) dto.UserResponse
	FindByID(ctx context.Context, userID int) dto.UserResponse
	UpdatePassword(ctx context.Context, userID int, password string) error
	LoginLog(ctx *gin.Context, userID int, message string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) dto.UserResponse {
	user := dto.UserResponse{}
	r.db.
		Model(&entity.User{}).
		Select("users.id as id, nis, nip, nama, jenis_kelamin, tempat_lahir, tanggal_lahir, username, role, password, foto, pegawai.id as pegawai_id, siswa.id as siswa_id").
		Where("username = ?", username).
		Joins("left join profiles on profiles.user_id = users.id").
		Joins("left join pegawai on pegawai.user_id = users.id").
		Joins("left join siswa on siswa.user_id = users.id").
		First(&user)
	return user
}

func (r *userRepository) FindByID(ctx context.Context, userID int) dto.UserResponse {
	user := dto.UserResponse{}
	r.db.
		Model(&entity.User{}).
		Select("users.id as id, nis, nip, nama, jenis_kelamin, tempat_lahir, tanggal_lahir, username, role, password, foto, siswa.id as siswa_id, pegawai.id as pegawai_id").
		Where("users.id = ?", userID).
		Joins("left join profiles on profiles.user_id = users.id").
		Joins("left join pegawai on pegawai.user_id = users.id").
		Joins("left join siswa on siswa.user_id = users.id").
		First(&user)
	return user
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID int, password string) error {
	err := r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("password", password).
		Error
	return err
}

func (r *userRepository) LoginLog(ctx *gin.Context, userID int, message string) error {
	//location := helper.GetUserLocation(ctx)

	uaString := ctx.Request.UserAgent()
	ua := useragent.New(uaString)
	os := ua.OS()

	browser, version := ua.Browser()

	log := entity.LoginLog{
		UserID:    userID,
		Action:    message,
		UserAgent: uaString,
		OS:        os,
		Browser:   browser + " " + version,
		Location:  "",
	}

	err := r.db.Create(&log).Error
	return err
}
