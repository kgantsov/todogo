package handlers

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

func TodoGraphql(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"todo": &graphql.Field{
					Type:        models.TodoType,
					Description: "Get todo by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, _ := p.Args["id"].(string)

						currentUser := c.MustGet("CurrentUser").(models.User)

						var todo models.Todo

						db.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&todo)

						if todo.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
							return todo, nil
						}
						return nil, errors.New("Todo not found")
					},
				},
				"list": &graphql.Field{
					Type:        graphql.NewList(models.TodoType),
					Description: "Get list of todos",
					Args: graphql.FieldConfigArgument{
						"list_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						currentUser := c.MustGet("CurrentUser").(models.User)

						listID, _ := params.Args["list_id"].(string)

						var todoList models.TodoList

						db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

						if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
							c.JSON(404, gin.H{"error": "Todo list not found"})
							return nil, errors.New("TodoList not found")
						}

						var todos []models.Todo

						db.Order("completed asc, priority desc, created_at asc").Where(
							"user_id = ? AND todo_list_id = ?", currentUser.ID, listID,
						).Find(&todos)

						return todos, nil
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
