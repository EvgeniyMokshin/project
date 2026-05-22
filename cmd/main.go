package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"    // для MySQL
	"gorm.io/driver/postgres" // для PostgreSQL

	"project/internal/handlers"
	"project/internal/models"
	"project/pkg/logger"

	"github.com/joho/godotenv"
)

// Config — структура конфигурации приложения
type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Dialect  string
		Host     string
		Port     string
		Username string
		Password string
		Name     string
		SSLMode  string
	}
}

func main() {
	// Инициализируем логгер
	logg, err := logger.InitFileLogger(logger.INFO, "logs/app.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logg.Close()

	// Загружаем конфигурацию (в реальном проекте — из config.yaml)
	config := loadConfig()

	// Подключаемся к БД
	db, err := connectToDatabase(config)
	if err != nil {
		logg.Error("Failed to connect to database: %v", err)
		return
	}

	// Выполняем миграции
	logg.Info("Starting database migrations...")
	err = runMigrations(db)
	if err != nil {
		logg.Error("Migration failed: %v", err)
		return
	}
	logg.Info("Database migrations completed successfully!")

	// Инициализируем обработчики
	departmentHandler := handlers.NewDepartmentHandler(db)
	employeeHandler := handlers.NewEmployeeHandler(db)

	// Настраиваем Gin
	router := gin.Default()

	// API endpoints
	api := router.Group("/api")
	{
		departments := api.Group("/departments")
		{
			departments.GET("", departmentHandler.GetAll)
			departments.GET("/:id", departmentHandler.GetByID)
			departments.POST("", departmentHandler.Create)
			departments.PUT("/:id", departmentHandler.Update)
			departments.DELETE("/:id", departmentHandler.Delete)
		}

		employees := api.Group("/employees")
		{
			employees.GET("", employeeHandler.GetAll)
			employees.GET("/:id", employeeHandler.GetByID)
			employees.GET("/search", employeeHandler.Search)
			employees.POST("", employeeHandler.Create)
			employees.PUT("/:id", employeeHandler.Update)
			employees.DELETE("/:id", employeeHandler.Delete)
		}
	}

	// Запускаем сервер
	server := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	logg.Info("Server starting on port %s", config.Server.Port)
	log.Printf("Server starting on port %s\n", config.Server.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logg.Error("Server failed to start: %v", err)
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loadConfig загружает конфигурацию приложения
// В реальном проекте — чтение из config.yaml
// loadConfig загружает конфигурацию приложения из .env файла
func loadConfig() Config {
	// Загружаем переменные из .env (если файл есть)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	return Config{
		Server: struct {
			Port string
		}{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: struct {
			Dialect  string
			Host     string
			Port     string
			Username string
			Password string
			Name     string
			SSLMode  string
		}{
			Dialect:  getEnv("DB_DIALECT", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "your_username"),
			Password: getEnv("DB_PASSWORD", "your_password"),
			Name:     getEnv("DB_NAME", "employee_management"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// connectToDatabase подключает к базе данных
func connectToDatabase(config Config) (*gorm.DB, error) {
	var dsn string

	switch config.Database.Dialect {
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Database.Host,
			config.Database.Port,
			config.Database.Username,
			config.Database.Password,
			config.Database.Name,
			config.Database.SSLMode,
		)
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})

	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})

	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", config.Database.Dialect)
	}
}

// runMigrations выполняет миграции базы данных
func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Автомиграция для моделей
	err := db.AutoMigrate(
		&models.Department{},
		&models.Employee{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations completed successfully!")
	return nil
}
