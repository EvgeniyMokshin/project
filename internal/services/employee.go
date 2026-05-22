package services

import (
	"errors"
	"time"

	"project/internal/models"

	"gorm.io/gorm"
)

// EmployeeService — структура сервиса для работы с сотрудниками
type EmployeeService struct {
	db *gorm.DB
}

// NewEmployeeService — конструктор для создания сервиса сотрудников
func NewEmployeeService(db *gorm.DB) *EmployeeService {
	return &EmployeeService{db: db}
}

// CreateEmployee создаёт нового сотрудника в указанном подразделении
func (s *EmployeeService) CreateEmployee(departmentID uint, fullName, position string, hiredAt time.Time) (*models.Employee, error) {
	// Валидация обязательных полей
	if fullName == "" {
		return nil, errors.New("full name is required")
	}
	if len(fullName) > 200 {
		return nil, errors.New("full name cannot exceed 200 characters")
	}
	if position == "" {
		return nil, errors.New("position is required")
	}
	if len(position) > 200 {
		return nil, errors.New("position cannot exceed 200 characters")
	}

	// Проверяем существование подразделения
	var department models.Department
	if err := s.db.First(&department, departmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}

	// Создаём запись сотрудника
	employee := &models.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.db.Create(employee).Error; err != nil {
		return nil, err
	}

	return employee, nil
}

// GetEmployeeByID получает сотрудника по ID с загрузкой связанного подразделения
func (s *EmployeeService) GetEmployeeByID(id int) (*models.Employee, error) {
	var employee models.Employee
	if err := s.db.Preload("Department").First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return &employee, nil
}

// UpdateEmployee обновляет данные сотрудника (только не-nil поля)
func (s *EmployeeService) UpdateEmployee(id int, fullName, position *string, hiredAt *time.Time) (*models.Employee, error) {
	var employee models.Employee
	if err := s.db.First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}

	// Обновляем только не-nil поля с валидацией
	if fullName != nil {
		if *fullName == "" {
			return nil, errors.New("full name cannot be empty")
		}
		if len(*fullName) > 200 {
			return nil, errors.New("full name cannot exceed 200 characters")
		}
		employee.FullName = *fullName
	}
	if position != nil {
		if *position == "" {
			return nil, errors.New("position cannot be empty")
		}
		if len(*position) > 200 {
			return nil, errors.New("position cannot exceed 200 characters")
		}
		employee.Position = *position
	}
	if hiredAt != nil {
		employee.HiredAt = *hiredAt // разыменование указателя
	}

	if err := s.db.Save(&employee).Error; err != nil {
		return nil, err
	}

	return &employee, nil
}

// DeleteEmployee удаляет сотрудника по ID
func (s *EmployeeService) DeleteEmployee(id int) error {
	result := s.db.Delete(&models.Employee{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("employee not found")
	}
	return nil
}

// GetEmployeesByDepartment получает всех сотрудников подразделения с сортировкой
func (s *EmployeeService) GetEmployeesByDepartment(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee

	err := s.db.
		Preload("Department").
		Where("department_id = ?", departmentID).
		Order("created_at ASC, full_name ASC").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}
	return employees, nil
}

// SearchEmployees ищет сотрудников по имени или должности
func (s *EmployeeService) SearchEmployees(query string) ([]models.Employee, error) {
	var employees []models.Employee

	searchQuery := "%" + query + "%"
	err := s.db.
		Preload("Department").
		Where("full_name LIKE ? OR position LIKE ?", searchQuery, searchQuery).
		Order("full_name ASC").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}
	return employees, nil
}
