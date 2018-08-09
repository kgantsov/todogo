package models

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

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
	UserType.AddFieldConfig("todo_lists", &graphql.Field{
		Type:        graphql.NewList(TodoListType),
		Description: "The friends of the character, or an empty list if they have none.",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var todoList []TodoList

			db := params.Context.Value("db").(*gorm.DB)
			// if db := params.Context.Value(k); v != nil {
			// 	fmt.Println("found value:", v)
			// 	return
			// }

			// db = db.(gorm.DB)

			db.Order("id asc").Where("user_id = ?", params.Source.(User).ID).Find(&todoList)
			return todoList, nil
		},
	})
}
