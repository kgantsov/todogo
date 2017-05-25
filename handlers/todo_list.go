package handlers

import (
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
	"github.com/jinzhu/gorm"
)

func CreateTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	var todoList models.TodoList
	e := c.BindJSON(&todoList)

	if e == nil {
		db.Create(&todoList)
		c.JSON(201, todoList)
	} else {
		c.JSON(422, gin.H{"error": e})
	}
}

func GetTodoLists(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	var todoLists []models.TodoList

	db.Find(&todoLists)

	c.JSON(200, todoLists)
}

func GetTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId := c.Params.ByName("listId")

	var todoList models.TodoList

	db.First(&todoList, listId)

	if todoList.ID != 0 {
		c.JSON(200, todoList)
	} else {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}
}

func UpdateTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId := c.Params.ByName("listId")

	var newTodoList models.TodoList
	e := c.BindJSON(&newTodoList)

	if e != nil {
		c.JSON(422, e)
	}

	var todoList models.TodoList

	db.First(&todoList, listId)

	if todoList.ID != 0 {
		todoList = models.TodoList{
			ID:    todoList.ID,
			Title: newTodoList.Title,
			CreatedAt: todoList.CreatedAt,
			UpdatedAt: time.Now(),
		}

		db.Save(&todoList)

		c.JSON(200, todoList)
	} else {
		c.JSON(404, gin.H{"error": "Todo List not found"})
	}
}

func DeleteTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId := c.Params.ByName("listId")

	var todoList models.TodoList

	db.First(&todoList, listId)

	if todoList.ID != 0 {
		db.Delete(&todoList)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, gin.H{"error": "Todo List not found"})
	}
}
