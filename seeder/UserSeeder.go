package seeder

import (
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
	"time"
)

func CreateUser(
	db *gorm.DB,
	username string,
	password string,
	role string,
	nama string,
	tempatLahir string,
	tanggalLahir string,
	jenisKelamin string,
) {
	hashedPassword, _ := helper.HashPassword(password)

	user := entity.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}
	db.FirstOrCreate(&user, entity.User{
		Username: username,
	})

	dateString := tanggalLahir
	date, _ := time.Parse("2006-01-02", dateString)
	db.FirstOrCreate(&entity.Profile{
		Nama:         nama,
		TanggalLahir: date,
		TempatLahir:  tempatLahir,
		JenisKelamin: jenisKelamin,
	}, entity.Profile{
		UserID: user.ID,
	})
}
