package models

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
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
			// if db := params.Context.Value(k); v != nil {
			// 	fmt.Println("found value:", v)
			// 	return
			// }

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
			// if db := params.Context.Value(k); v != nil {
			// 	fmt.Println("found value:", v)
			// 	return
			// }

			db.Order("id asc").Where("id = ?", params.Source.(Todo).UserID).First(&user)
			return user, nil
		},
	})
}
