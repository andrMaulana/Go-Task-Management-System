package database

import (
	"log"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
	)

	if err != nil {
		log.Fatalf("Could not run migrations: %v", err)
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
