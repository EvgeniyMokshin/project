package database

import (
	"fmt"
	"log"
	"project/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect устанавливает соединение с базой данных PostgreSQL и выполняет автомиграцию моделей
func Connect() (*gorm.DB, error) {
	// Строка подключения к PostgreSQL
	dsn := "host=db user=postgres password=postgres dbname=org_structure port=5432 sslmode=disable"

	// Открываем соединение с БД через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database")

	// Выполняем автомиграцию — создаём/обновляем таблицы для моделей Department и Employee
	err = db.AutoMigrate(&models.Department{}, &models.Employee{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return db, nil
}
