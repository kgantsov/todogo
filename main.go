package main

import (
	"github.com/kgantsov/todogo/handlers"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/jinzhu/gorm"
)

func DBMiddleware(db gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}


func main() {
	db := models.InitDb()
	models.CreateTables(db)

	r := gin.Default()

	r.Use(DBMiddleware(*db))

	v1 := r.Group("api/v1")
	{
		v1.POST("/list/", handlers.CreateTodoList)
		v1.GET("/list/", handlers.GetTodoLists)
		v1.GET("/list/:listId/", handlers.GetTodoList)
		v1.PUT("/list/:listId/", handlers.UpdateTodoList)
		v1.DELETE("/list/:listId/", handlers.DeleteTodoList)

		v1.POST("/list/:listId/todo/", handlers.CreateTodo)
		v1.GET("/list/:listId/todo/", handlers.GetTodos)
		v1.GET("/list/:listId/todo/:todoId/", handlers.GetTodo)
		v1.PUT("/list/:listId/todo/:todoId/", handlers.UpdateTodo)
		v1.DELETE("/list/:listId/todo/:todoId/", handlers.DeleteTodo)
	}
	r.Run(":8080")
}
