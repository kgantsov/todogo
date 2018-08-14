package models

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

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

			db.Order("created_at asc").Where(
				"todo_list_id = ?",
				params.Source.(TodoList).ID,
			).Find(&todos)
			return todos, nil
		},
	})
}
