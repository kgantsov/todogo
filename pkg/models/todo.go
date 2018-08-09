package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type PRIORITY int

const (
	PRIORITY_NIL PRIORITY = iota
	PRIORITY_IRRELEVANT
	PRIORITY_EXTRA_LOW
	PRIORITY_LOW
	PRIORITY_NORMAL
	PRIORITY_HIGH
	PRIORITY_URGENT
	PRIORITY_SUPER_URGENT
	PRIORITY_IMMEDIATE
)

type Todo struct {
	ID         uuid.UUID  `gorm:"primary_key" form:"id" json:"id"`
	CreatedAt  *time.Time `gorm:"not null" form:"created_at" json:"created_at,omitempty"`
	UpdatedAt  *time.Time `gorm:"not null" form:"updated_at" json:"updated_at,omitempty"`
	Title      string     `gorm:"not null" form:"title" json:"title" binding:"required"`
	Completed  bool       `gorm:"not null" form:"completed" json:"completed"`
	Note       string     `gorm:"not null" form:"note" json:"note"`
	TodoListID uuid.UUID  `gorm:"index" form:"todo_list_id" json:"todo_list_id"`
	UserID     uuid.UUID  `gorm:"index" form:"user_id" json:"user_id"`
	DeadLineAt *time.Time `gorm:"index" form:"dead_line_at" json:"dead_line_at,omitempty"`
	Priority   PRIORITY   `gorm:"index" form:"priority" json:"priority"`
}
