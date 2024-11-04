package db

import (
	"example.com/RMS/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Import logger
	"log"
	"os"
	"time"
)

var (
	DB      *gorm.DB
	AuthKey string
	Apikey  string
)

func InitDB() error {
	err := godotenv.Load("prod.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	AuthKey = os.Getenv("AUTH_KEY")

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s",
		host, user, dbname, password, port, sslmode)

	// Set up GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Log format
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore not found errors
			Colorful:                  true,        // Enable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // Set logger
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return err
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error migrating User model: %v", err)
		return err

	}

	if err := DB.AutoMigrate(&models.Profile{}); err != nil {
		log.Fatalf("Error migrating Profile model: %v", err)
		return err

	}

	if err := DB.AutoMigrate(&models.Job{}); err != nil {
		log.Fatalf("Error migrating Job model: %v", err)
		return err

	}

	if err := DB.AutoMigrate(&models.Application{}); err != nil {
		log.Fatalf("Error migrating Application model: %v", err)
		return err

	}

	log.Println("DB migrated & Connected Successfully")
	return nil
}
