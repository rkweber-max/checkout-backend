package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email;unique"`
	Password  string `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
