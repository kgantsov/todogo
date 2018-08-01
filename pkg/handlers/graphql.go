package handlers

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

func executeQuery(db *gorm.DB, query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.Background(), "db", db),
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
