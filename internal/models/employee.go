package models

import "time"

// Employee — модель сотрудника
type Employee struct {
	ID           uint      `gorm:"primaryKey"`
	FullName     string    `gorm:"type:varchar(200);not null"`
	Position     string    `gorm:"type:varchar(200);not null"`
	Email        string    `gorm:"type:varchar(100)"`
	DepartmentID uint      `gorm:"not null"`
	HiredAt      time.Time `gorm:"not null"` // новое поле
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Связь с департаментом
	Department Department `gorm:"foreignKey:DepartmentID"`
}
