package handlers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/jinzhu/gorm"
)


func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Auth-Token")) == 0 {
			c.JSON(403, gin.H{"error": "Auth-Token is required"})
			c.AbortWithStatus(403)
			return
		}

		user, ok := validateToken(db, c.GetHeader("Auth-Token"))

		if !ok {
			c.JSON(403, gin.H{"error": "Auth-Token is invalid"})
			c.AbortWithStatus(403)
			return
		}

		c.Set("CurrentUser", user)
		c.Next()
	}
}

func DefineRoutes(db *gorm.DB, r *gin.Engine) {
	v1 := r.Group("api/v1")
	{
		v1.OPTIONS("/user/", OptionsUser)
		v1.POST("/user/", CreateUser)

		v1.POST("/auth/login/", Login)
		v1.OPTIONS("/auth/login/", OptionsLogin)

		v1.Use(AuthMiddleware(db))
		{
			v1.OPTIONS("/user/:userId/", OptionsUser)
			v1.GET("/user/", GetUsers)
			v1.GET("/user/:userId/", GetUser)
			v1.PUT("/user/:userId/", UpdateUser)
			v1.DELETE("/user/:userId/", DeleteUser)

			v1.OPTIONS("/list/", OptionsTodoList)
			v1.OPTIONS("/list/:listId/", OptionsTodoList)
			v1.POST("/list/", CreateTodoList)
			v1.GET("/list/", GetTodoLists)
			v1.GET("/list/:listId/", GetTodoList)
			v1.PUT("/list/:listId/", UpdateTodoList)
			v1.DELETE("/list/:listId/", DeleteTodoList)

			v1.OPTIONS("/list/:listId/todo/", OptionsTodo)
			v1.OPTIONS("/list/:listId/todo/:todoId/", OptionsTodo)
			v1.POST("/list/:listId/todo/", CreateTodo)
			v1.GET("/list/:listId/todo/", GetTodos)
			v1.GET("/list/:listId/todo/:todoId/", GetTodo)
			v1.PUT("/list/:listId/todo/:todoId/", UpdateTodo)
			v1.DELETE("/list/:listId/todo/:todoId/", DeleteTodo)
		}
	}
}
