package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func UserGraphql(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type:        models.UserType,
					Description: "Get user by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, _ := p.Args["id"].(string)

						currentUser := c.MustGet("CurrentUser").(models.User)
						var user models.User

						if uuid.FromStringOrNil(id) != currentUser.ID {
							c.JSON(403, gin.H{"error": "Access denied"})
							return nil, errors.New("User not found")
						}

						db.Where("id = ?", id).First(&user)

						if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
							return user, nil
						}
						return nil, errors.New("User not found")
					},
				},
				"list": &graphql.Field{
					Type:        graphql.NewList(models.UserType),
					Description: "Get user list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {

						var users []models.User

						currentUser := c.MustGet("CurrentUser").(models.User)

						db.Order("id asc").Where("id = ?", currentUser.ID).Find(&users)
						return users, nil
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
