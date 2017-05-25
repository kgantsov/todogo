package models

import (
	"time"
)

type TodoList struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	CreatedAt time.Time `gorm:"not null" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" form:"updated_at" json:"updated_at"`
	Title     string    `gorm:"not null" form:"title" json:"title" binding:"required"`
}
