package models

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
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
