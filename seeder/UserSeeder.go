package seeder

import (
	"fmt"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username string, password string, role string) {
	hashedPassword, _ := helper.HashPassword(password)

	db.FirstOrCreate(&entity.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}, entity.User{
		Username: username,
	})

	fmt.Println(hashedPassword)
}
