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
		&entity.Siswa{},
		&entity.TahunAjar{},
		&entity.Jurusan{},
		&entity.Kelas{},
		&entity.Guru{},
		&entity.SiswaKelas{},
		&entity.DataScore{},
		&entity.ScoreSiswa{},
		&entity.ScoreLog{},
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

	err = CreateTriggerAndFunction(db)
	if err != nil {
		return nil
	}

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println("Failed disconnect database")
	}
	dbSQL.Close()
}

func CreateTriggerAndFunction(db *gorm.DB) error {
	// Check if the trigger already exists
	//var count int64
	//db.Raw("SELECT count(*) FROM pg_trigger WHERE tgname = 'siswa_kelas_insert_trigger'").Count(&count)
	//if count > 0 {
	//	Trigger already exists, do nothing
	//return nil
	//}

	// Create the function
	db.Exec("CREATE OR REPLACE FUNCTION siswa_kelas_insert() RETURNS TRIGGER AS $siswa_kelas_insert$ BEGIN " +
		"INSERT INTO score_siswa (siswa_kelas_id, score, created_at, updated_at) VALUES (NEW.id, 50, NOW(), NOW())" +
		"; RETURN NEW; END; $siswa_kelas_insert$ LANGUAGE plpgsql;")

	// Create the trigger
	db.Exec("CREATE OR REPLACE TRIGGER siswa_kelas_insert_trigger AFTER INSERT ON siswa_kelas FOR EACH ROW EXECUTE FUNCTION siswa_kelas_insert();")

	return nil
}
