package models

import "time"

// Department — модель подразделения
type Department struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	ParentID  *uint  `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Связи
	Children  []Department `gorm:"foreignKey:ParentID"`
	Employees []Employee   `gorm:"foreignKey:DepartmentID"`
}
