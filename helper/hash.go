package helper

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func BirthDateToPassword(birthDate time.Time) string {
	year := birthDate.Format("2006")
	month := birthDate.Format("01")
	day := birthDate.Format("02")

	password, _ := HashPassword(day + month + year)

	return password
}
