package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key" form:"id" json:"id"`
	CreatedAt time.Time `gorm:"not null" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" form:"updated_at" json:"updated_at"`
	Name      string    `gorm:"not null" form:"name" json:"name"`
	Email     string    `gorm:"type:varchar(100);unique_index" form:"email" json:"email" binding:"required"`
	Password  string    `gorm:"not null" form:"password" json:"password" binding:"required"`
}

func (User) TableName() string {
	return "profiles"
}
