package handlers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gorm.io/gorm"
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
		v1.OPTIONS("/user/:userID/", OptionsUser)
		v1.OPTIONS("/auth/login/", OptionsLogin)
		v1.OPTIONS("/list/", OptionsTodoList)
		v1.OPTIONS("/list/:listID/", OptionsTodoList)
		v1.OPTIONS("/list/:listID/todo/", OptionsTodo)
		v1.OPTIONS("/list/:listID/todo/:todoID/", OptionsTodo)

		v1.POST("/user/", CreateUser)
		v1.POST("/auth/login/", Login)

		v1.Use(AuthMiddleware(db))
		{
			v1.GET("/user/", GetUsers)
			v1.GET("/user/:userID/", GetUser)
			v1.PUT("/user/:userID/", UpdateUser)
			v1.DELETE("/user/:userID/", DeleteUser)

			v1.POST("/list/", CreateTodoList)
			v1.GET("/list/", GetTodoLists)
			v1.GET("/list/:listID/", GetTodoList)
			v1.PUT("/list/:listID/", UpdateTodoList)
			v1.DELETE("/list/:listID/", DeleteTodoList)

			v1.POST("/list/:listID/todo/", CreateTodo)
			v1.GET("/list/:listID/todo/", GetTodos)
			v1.GET("/list/:listID/todo/:todoID/", GetTodo)
			v1.PUT("/list/:listID/todo/:todoID/", UpdateTodo)
			v1.DELETE("/list/:listID/todo/:todoID/", DeleteTodo)
		}
	}

	graphqlV1 := r.Group("graphql")
	{
		graphqlV1.Use(AuthMiddleware(db))
		{
			graphqlV1.GET("/user/", UserGraphql)
			graphqlV1.GET("/list/", TodoListGraphql)
			graphqlV1.GET("/todo/", TodoGraphql)
		}
	}
}
