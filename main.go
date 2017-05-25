package main

import (
	"github.com/kgantsov/todogo/handlers"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	db := models.InitDb()
	defer db.Close()
	models.CreateTables(db)

	r := gin.New()

	v1 := r.Group("api/v1")
	{
		v1.POST("/list/", handlers.CreateTodoList)
		v1.GET("/list/", handlers.GetTodoLists)
		v1.GET("/list/:id/", handlers.GetTodoList)
		v1.PUT("/list/:id/", handlers.UpdateTodoList)
		v1.DELETE("/list/:id/", handlers.DeleteTodoList)
	}
	r.Run(":8080")
}
