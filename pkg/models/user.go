package models

import (
	"time"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"primary_key" form:"id" json:"id"`
	CreatedAt *time.Time `gorm:"not null" form:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"not null" form:"updated_at" json:"updated_at,omitempty"`
	Name      string     `gorm:"not null" form:"name" json:"name"`
	Email     string     `gorm:"type:varchar(100);unique_index" form:"email" json:"email" binding:"required"`
	Password  string     `gorm:"not null" form:"password" json:"password" binding:"required"`
}

func (User) TableName() string {
	return "profiles"
}

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
