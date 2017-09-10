package handlers

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
)

func OptionsTodo(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Auth-Token")
	c.Next()
}

func CreateTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId, _ := strconv.ParseUint(c.Params.ByName("listId"), 0, 64)

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listId).First(&todoList)

	if todoList.ID == 0 {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}

	var todo models.Todo
	e := c.BindJSON(&todo)

	if e == nil {
		if todo.Priority == 0 {
			todo.Priority = models.PRIORITY_NORMAL
		}

		currentUser := c.MustGet("CurrentUser").(models.User)
		todo.UserID = currentUser.ID
		todo.TodoListID = uint(listId)

		db.Create(&todo)
		c.JSON(201, todo)
	} else {
		c.JSON(422, gin.H{"error": e})
	}
}

func GetTodos(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId := c.Params.ByName("listId")

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listId).First(&todoList)

	if todoList.ID == 0 {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}

	var todos []models.Todo

	db.Order("completed asc, priority desc, id asc").Where(
		"user_id = ? AND todo_list_id = ?", currentUser.ID, listId,
	).Find(&todos)

	c.JSON(200, todos)
}

func GetTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId := c.Params.ByName("listId")
	todoId := c.Params.ByName("todoId")

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listId).First(&todoList)

	if todoList.ID == 0 {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}

	var todo models.Todo

	db.Where("user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listId, todoId).Find(&todo)

	if todo.ID != 0 {
		c.JSON(200, todo)
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}

func UpdateTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId, _ := strconv.ParseUint(c.Params.ByName("listId"), 0, 64)
	todoId := c.Params.ByName("todoId")

	var newTodo models.Todo
	e := c.BindJSON(&newTodo)

	if e != nil {
		c.JSON(422, e)
	}

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listId).First(&todoList)

	if todoList.ID == 0 {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}

	var todo models.Todo

	db.Where("user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listId, todoId).Find(&todo)

	if todo.ID != 0 {
		todo = models.Todo{
			ID:         todo.ID,
			Title:      newTodo.Title,
			Completed:  newTodo.Completed,
			Note:       newTodo.Note,
			TodoListID: uint(listId),
			UserID:     todo.UserID,
			CreatedAt:  todo.CreatedAt,
			UpdatedAt:  time.Now(),
			DeadLineAt: newTodo.DeadLineAt,
			Priority:   newTodo.Priority,
		}

		db.Save(&todo)
		c.JSON(200, todo)
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}

func DeleteTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	listId := c.Params.ByName("listId")
	todoId := c.Params.ByName("todoId")

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listId).First(&todoList)

	if todoList.ID == 0 {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}

	var todo models.Todo

	db.Where("user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listId, todoId).Find(&todo)

	if todo.ID != 0 {
		db.Delete(&todo)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}
