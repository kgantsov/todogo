package models

import (
	"time"
)

type Todo struct {
	ID         uint      `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	CreatedAt  time.Time `gorm:"not null" form:"created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null" form:"updated_at" json:"updated_at"`
	Title      string    `gorm:"not null" form:"title" json:"title" binding:"required"`
	Completed  bool      `gorm:"not null" form:"completed" json:"completed"`
	Note       string    `gorm:"not null" form:"note" json:"note"`
	TodoListID uint      `gorm:"index" form:"todo_list_id" json:"todo_list_id"`
	UserID     uint      `gorm:"index" form:"user_id" json:"user_id"`
	DeadLineAt time.Time `gorm:"index" form:"dead_line_at" json:"dead_line_at"`
}
