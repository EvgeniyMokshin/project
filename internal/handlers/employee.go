package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"project/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmployeeHandler — структура обработчика для сотрудников
type EmployeeHandler struct {
	db *gorm.DB
}

// NewEmployeeHandler создаёт новый экземпляр обработчика сотрудников
func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{db: db}
}

// GetAll получает всех сотрудников
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	var employees []models.Employee
	if err := h.db.Preload("Department").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// GetByID получает сотрудника по ID
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var employee models.Employee
	if err := h.db.Preload("Department").First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, employee)
}

// Search ищет сотрудников по имени или должности
func (h *EmployeeHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	searchQuery := "%" + query + "%"
	var employees []models.Employee
	if err := h.db.Preload("Department").
		Where("full_name LIKE ? OR position LIKE ?", searchQuery, searchQuery).
		Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// Create создаёт нового сотрудника
func (h *EmployeeHandler) Create(c *gin.Context) {
	var input struct {
		FullName     string  `json:"FullName" binding:"required"`
		Position     string  `json:"Position" binding:"required"`
		Email        *string `json:"Email"`
		DepartmentID uint    `json:"DepartmentID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование департамента
	var department models.Department
	if err := h.db.First(&department, input.DepartmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Обрабатываем поле Email с проверкой на nil
	var email string
	if input.Email != nil {
		email = *input.Email
	} else {
		email = "" // значение по умолчанию, если Email не передан
	}

	employee := models.Employee{
		FullName:     input.FullName,
		Position:     input.Position,
		Email:        email,
		DepartmentID: input.DepartmentID,
	}

	if err := h.db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// Update обновляет сотрудника
func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input struct {
		FullName     *string `json:"full_name"`
		Position     *string `json:"position"`
		Email        *string `json:"email"`
		DepartmentID *uint   `json:"department_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var employee models.Employee
	if err := h.db.First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if input.FullName != nil {
		employee.FullName = *input.FullName
	}
	if input.Position != nil {
		employee.Position = *input.Position
	}
	if input.Email != nil {
		employee.Email = *input.Email
	}
	if input.DepartmentID != nil {
		// Проверяем существование нового департамента
		var department models.Department
		if err := h.db.First(&department, *input.DepartmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "New department not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		employee.DepartmentID = *input.DepartmentID
	}

	if err := h.db.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// Delete удаляет сотрудника
func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := h.db.Delete(&models.Employee{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
