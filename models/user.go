package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	CreatedAt time.Time `gorm:"not null" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" form:"updated_at" json:"updated_at"`
	Name      string    `gorm:"not null" form:"name" json:"name"`
	Email     string    `gorm:"type:varchar(100);unique_index" form:"email" json:"email" binding:"required"`
	Password  string    `gorm:"not null" form:"password" json:"password" binding:"required"`
}

func (User) TableName() string {
	return "profiles"
}
