package config

import (
	"fmt"
	"github/billing-engine-main/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Could not load .env file, proceeding with other configurations")
	}
	return nil
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s ",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DBNAME"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_TIMEZONE"))
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the schemas
	err = db.AutoMigrate(
		&models.Loan{},
		&models.Schedule{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
