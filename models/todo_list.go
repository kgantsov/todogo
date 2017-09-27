package models

import (
	"time"

	"github.com/satori/go.uuid"
)

type TodoList struct {
	ID        uuid.UUID `gorm:"primary_key" form:"id" json:"id"`
	CreatedAt time.Time `gorm:"not null" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" form:"updated_at" json:"updated_at"`
	Title     string    `gorm:"not null" form:"title" json:"title" binding:"required"`
	UserID    uuid.UUID `gorm:"index" form:"user_id" json:"user_id"`
}
