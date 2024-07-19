package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	OwnerID     uint
	Owner       User   `gorm:"foreignKey:OwnerID"`
	Users       []User `gorm:"many2many:project_users;"`
}
