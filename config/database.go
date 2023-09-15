package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/seeder"
	"gorm.io/driver/mysql"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
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
		&entity.Pegawai{},
		&entity.SiswaKelas{},
		&entity.DataPoin{},
		&entity.PoinSiswa{},
		&entity.PoinLog{},
		&entity.Setting{},
		&entity.LoginLog{},
		&entity.Sync{},
	)
	if err != nil {
		fmt.Println("Automigrate error")
		fmt.Println(err.Error())
	}

	CallSeeder(db)

	err = CreateTriggerAndFunction(db)
	if err != nil {
		return nil
	}

	return db
}

func SitiDatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("File env not found")
	}
	dbUser := os.Getenv("DB_SITI_USER")
	dbPass := os.Getenv("DB_SITI_PASS")
	dbHost := os.Getenv("DB_SITI_HOST")
	dbName := os.Getenv("DB_SITI_NAME")
	dbPort := os.Getenv("DB_SITI_PORT")

	if dbUser != "" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			fmt.Println("Failed connect siti database")
		}

		return db
	}

	return nil
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println("Failed disconnect database")
	}
	dbSQL.Close()
}

func CallSeeder(db *gorm.DB) {
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

	seeder.CreateSetting(
		db,
		"score_type",
		`["Perlombaan", "Penghargaan", "Keaktifan", "Pelanggaran"]`,
		false,
	)
}

func CreateTriggerAndFunction(db *gorm.DB) error {
	// Create the function
	//db.Exec("CREATE OR REPLACE FUNCTION siswa_kelas_insert() RETURNS TRIGGER AS $siswa_kelas_insert$ BEGIN " +
	//	"INSERT INTO poin_siswa (siswa_kelas_id, poin, created_at, updated_at) VALUES (NEW.id, 200, NOW(), NOW())" +
	//	"; RETURN NEW; END; $siswa_kelas_insert$ LANGUAGE plpgsql;")

	// Create the trigger
	//db.Exec("CREATE OR REPLACE TRIGGER siswa_kelas_insert_trigger AFTER INSERT ON siswa_kelas FOR EACH ROW EXECUTE FUNCTION siswa_kelas_insert();")

	return nil
}
