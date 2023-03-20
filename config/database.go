package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/seeder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func SetupDatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("File env not found")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbZone := os.Getenv("DB_ZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", dbHost, dbUser, dbPass, dbName, dbPort, dbZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed connect database")
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Profile{},
		&entity.Kelas{},
		&entity.Siswa{},
		&entity.Guru{},
		&entity.TahunAjar{},
		&entity.Jurusan{},
	)
	if err != nil {
		fmt.Println("Automigrate error")
		fmt.Println(err.Error())
	}

	seeder.CreateUser(
		db,
		"admin",
		"admin",
		"admin",
		"Administrator",
		"Bandung",
		"2000-01-01",
		"L",
	)

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println("Failed disconnect database")
	}
	dbSQL.Close()
}
