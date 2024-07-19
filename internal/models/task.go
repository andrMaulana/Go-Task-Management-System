package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"not null;default:'pending'"`
	ProjectID   uint
	Project     Project `gorm:"foreignKey:ProjectID"`
	AssigneeID  uint
	Assignee    User `gorm:"foreignKey:AssigneeID"`
}
