package repository

import (
	"errors"
	"fmt"

	"project/internal/models"

	"gorm.io/gorm"
)

// DepartmentRepository — интерфейс репозитория подразделений
type DepartmentRepository interface {
	Create(department *models.Department) error
	GetByID(id uint) (*models.Department, error)
	GetAll() ([]models.Department, error)
	Update(department *models.Department) error
	Delete(id uint) error
	GetWithSubDepartments(id uint) (*models.Department, error)
}

// departmentRepository — реализация репозитория
type departmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository создаёт новый экземпляр репозитория подразделений
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

// Create создаёт новое подразделение в БД
func (r *departmentRepository) Create(department *models.Department) error {
	if department.Name == "" {
		return errors.New("department name cannot be empty")
	}

	result := r.db.Create(department)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID получает подразделение по ID
func (r *departmentRepository) GetByID(id uint) (*models.Department, error) {
	var department models.Department
	if err := r.db.First(&department, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department with id %d not found", id)
		}
		return nil, err
	}
	return &department, nil
}

// GetAll получает все подразделения
func (r *departmentRepository) GetAll() ([]models.Department, error) {
	var departments []models.Department
	if err := r.db.Preload("Parent").Preload("SubDepartments").Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

// Update обновляет существующее подразделение
func (r *departmentRepository) Update(department *models.Department) error {
	// Проверяем существование подразделения
	existing, err := r.GetByID(department.ID)
	if err != nil {
		return err
	}

	if department.Name != "" {
		existing.Name = department.Name
	}
	if department.ParentID != nil {
		existing.ParentID = department.ParentID
	}

	result := r.db.Save(existing)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete удаляет подразделение по ID
func (r *departmentRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Department{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("department with id %d not found", id)
	}
	return nil
}

// GetWithSubDepartments получает подразделение с деревом подчинённых подразделений
func (r *departmentRepository) GetWithSubDepartments(id uint) (*models.Department, error) {
	var department models.Department

	// Загружаем подразделение с родительским подразделением и дочерними подразделениями
	err := r.db.
		Preload("Parent").
		Preload("SubDepartments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("SubDepartments") // Рекурсивная загрузка дочерних подразделений
		}).
		First(&department, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department with id %d not found", id)
		}
		return nil, err
	}

	return &department, nil
}
