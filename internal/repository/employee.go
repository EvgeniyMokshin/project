package repository

import (
	"errors"
	"fmt"
	"time"

	"project/internal/models"

	"gorm.io/gorm"
)

// EmployeeRepository — интерфейс репозитория сотрудников
type EmployeeRepository interface {
	Create(employee *models.Employee) error
	GetByID(id uint) (*models.Employee, error)
	GetAll() ([]models.Employee, error)
	Update(employee *models.Employee) error
	Delete(id uint) error
}

// employeeRepository — реализация репозитория
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository создаёт новый экземпляр репозитория сотрудников
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

// Create создаёт нового сотрудника в БД
func (r *employeeRepository) Create(employee *models.Employee) error {
	if employee.FullName == "" || employee.Position == "" {
		return errors.New("full name and position cannot be empty")
	}

	// Инициализируем HiredAt, если не задано
	if employee.HiredAt.IsZero() {
		employee.HiredAt = time.Now()
	}

	result := r.db.Create(employee)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID получает сотрудника по ID
func (r *employeeRepository) GetByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.Preload("Department").First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("employee with id %d not found", id)
		}
		return nil, err
	}
	return &employee, nil
}

// GetAll получает всех сотрудников
func (r *employeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	if err := r.db.Preload("Department").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

// Update обновляет существующего сотрудника
func (r *employeeRepository) Update(employee *models.Employee) error {
	existing, err := r.GetByID(employee.ID)
	if err != nil {
		return err
	}

	if employee.FullName != "" {
		existing.FullName = employee.FullName
	}
	if employee.Position != "" {
		existing.Position = employee.Position
	}
	if employee.Email != "" {
		existing.Email = employee.Email
	}
	if employee.DepartmentID != 0 {
		existing.DepartmentID = employee.DepartmentID
	}
	// Обновляем HiredAt, если передано новое значение
	if !employee.HiredAt.IsZero() {
		existing.HiredAt = employee.HiredAt
	}

	result := r.db.Save(existing)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete удаляет сотрудника по ID
func (r *employeeRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Employee{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee with id %d not found", id)
	}
	return nil
}
