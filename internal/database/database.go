package database

import (
	"fmt"
	"log"
	"time"

	"book_management_system/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs setup
func InitDB() error {
	// Connect to database
	if err := connectDB(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed initial data
	if err := seedData(); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	log.Println("Database initialization completed successfully")
	return nil
}

// connectDB establishes the database connection
func connectDB() error {
	dsn := "postgresql://book_management_owner:KEY@ep-spring-tooth-a1h39tm9.ap-southeast-1.aws.neon.tech/book_management?sslmode=require"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	log.Println("Database connection established")
	return nil
}

// runMigrations performs database migrations
func runMigrations() error {
	if err := DB.AutoMigrate(&models.Book{}); err != nil {
		return err
	} //database tables creation and modification

	log.Println("Database migrations completed")
	return nil
}

// seedData populates the database with initial data if it's empty
func seedData() error {
	// Check if data already exists
	var count int64
	if err := DB.Model(&models.Book{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Database already contains data, skipping seed")
		return nil
	}

	// Initial books data
	books := []models.Book{
		{
			Title:       "The Go Programming Language",
			Author:      "Alan A. A. Donovan",
			ISBN:        "0134190440",
			PublishedAt: time.Date(2015, 11, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			ISBN:        "0132350882",
			PublishedAt: time.Date(2008, 8, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Design Patterns",
			Author:      "Erich Gamma",
			ISBN:        "0201633612",
			PublishedAt: time.Date(1994, 11, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	// Create books in a transaction
	return DB.Transaction(func(tx *gorm.DB) error {
		for _, book := range books {
			if err := tx.Create(&book).Error; err != nil {
				return err
			}
		}
		log.Println("Database seeded successfully")
		return nil
	})
}
