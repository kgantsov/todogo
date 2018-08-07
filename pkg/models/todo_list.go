package models

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type TodoList struct {
	ID        uuid.UUID  `gorm:"primary_key" form:"id" json:"id"`
	CreatedAt *time.Time `gorm:"not null" form:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"not null" form:"updated_at" json:"updated_at,omitempty"`
	Title     string     `gorm:"not null" form:"title" json:"title" binding:"required"`
	UserID    uuid.UUID  `gorm:"index" form:"user_id" json:"user_id"`
}

var TodoListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TodoListType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func init() {
	TodoListType.AddFieldConfig("user", &graphql.Field{
		Type:        UserType,
		Description: "Get TodoList user",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var user User

			db := params.Context.Value("db").(*gorm.DB)
			// if db := params.Context.Value(k); v != nil {
			// 	fmt.Println("found value:", v)
			// 	return
			// }

			db.Order("id asc").Where("id = ?", params.Source.(TodoList).UserID).First(&user)
			return user, nil
		},
	})
	TodoListType.AddFieldConfig("todos", &graphql.Field{
		Type:        graphql.NewList(TodoType),
		Description: "Get users Todos",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var todos []Todo

			db := params.Context.Value("db").(*gorm.DB)
			// if db := params.Context.Value(k); v != nil {
			// 	fmt.Println("found value:", v)
			// 	return
			// }

			db.Order("created_at asc").Where("todo_list_id = ?", params.Source.(TodoList).ID).Find(&todos)
			return todos, nil
		},
	})
}
