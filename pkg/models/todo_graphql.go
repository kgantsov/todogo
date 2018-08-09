package models

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

var TodoType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TodoType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"completed": &graphql.Field{
				Type: graphql.Boolean,
			},
			"note": &graphql.Field{
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
	TodoType.AddFieldConfig("todo_list", &graphql.Field{
		Type:        TodoListType,
		Description: "Get Todo todo_list",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var todoList TodoList

			db := params.Context.Value("db").(*gorm.DB)

			db.Order("id asc").Where("id = ?", params.Source.(Todo).TodoListID).First(&todoList)
			return todoList, nil
		},
	})
	TodoType.AddFieldConfig("user", &graphql.Field{
		Type:        UserType,
		Description: "Get Todo user",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var user User

			db := params.Context.Value("db").(*gorm.DB)

			db.Order("id asc").Where("id = ?", params.Source.(Todo).UserID).First(&user)
			return user, nil
		},
	})
}
