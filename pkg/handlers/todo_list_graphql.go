package handlers

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

func TodoListGraphql(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"todo_list": &graphql.Field{
					Type:        models.TodoListType,
					Description: "Get todo_list by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, _ := p.Args["id"].(string)

						currentUser := c.MustGet("CurrentUser").(models.User)

						var todoList models.TodoList

						db.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&todoList)

						if todoList.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
							return todoList, nil
						}
						return nil, errors.New("TodoList not found")
					},
				},
				"list": &graphql.Field{
					Type:        graphql.NewList(models.TodoListType),
					Description: "Get list of todo_list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {

						var todoLists []models.TodoList

						currentUser := c.MustGet("CurrentUser").(models.User)

						db.Order("created_at asc").Where("user_id = ?", currentUser.ID).Find(&todoLists)
						return todoLists, nil
					},
				},
			},
		},
	)

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)

	result := executeQuery(&db, c.Query("query"), schema)
	c.JSON(200, result)
}
