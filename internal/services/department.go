package services

import (
	"errors"
	"fmt"

	"project/internal/models"
	"project/internal/repository"
)

// DepartmentService — сервис для работы с подразделениями
type DepartmentService struct {
	repo repository.DepartmentRepository
}

// NewDepartmentService создаёт новый экземпляр сервиса подразделений
func NewDepartmentService(repo repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// CreateDepartment создаёт новое подразделение
func (s *DepartmentService) CreateDepartment(name string, parentID *uint) (*models.Department, error) {
	if name == "" {
		return nil, errors.New("department name cannot be empty")
	}

	department := &models.Department{
		Name:     name,
		ParentID: parentID, // теперь тип совпадает: *uint
	}

	if err := s.repo.Create(department); err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return department, nil
}

// GetDepartmentByID получает подразделение по ID
func (s *DepartmentService) GetDepartmentByID(id uint) (*models.Department, error) {
	department, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	return department, nil
}

// UpdateDepartment обновляет подразделение
func (s *DepartmentService) UpdateDepartment(id uint, name string, parentID *uint) error {
	// Получаем существующее подразделение
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Обновляем поля, если они переданы
	if name != "" {
		existing.Name = name
	}
	if parentID != nil {
		existing.ParentID = parentID // тип *uint соответствует модели
	}

	if err := s.repo.Update(existing); err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	return nil
}

// DeleteDepartment удаляет подразделение
func (s *DepartmentService) DeleteDepartment(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}
	return nil
}

// GetDepartmentWithSubDepartments получает подразделение с деревом подчинённых подразделений
func (s *DepartmentService) GetDepartmentWithSubDepartments(id uint) (*models.Department, error) {
	department, err := s.repo.GetWithSubDepartments(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department with sub-departments: %w", err)
	}
	return department, nil
}
