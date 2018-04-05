package handlers

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
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
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, gin.H{"error": "Todo list not found"})
		return
	}

	var todo models.Todo
	e := c.BindJSON(&todo)

	if e == nil {
		if todo.Priority == 0 {
			todo.Priority = models.PRIORITY_NORMAL
		}

		currentUser := c.MustGet("CurrentUser").(models.User)
		todo.ID = uuid.NewV4()
		todo.UserID = currentUser.ID
		todo.TodoListID = listID

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
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, gin.H{"error": "Todo list not found"})
		return
	}

	var todos []models.Todo

	db.Order("completed asc, priority desc, created_at asc").Where(
		"user_id = ? AND todo_list_id = ?", currentUser.ID, listID,
	).Find(&todos)

	c.JSON(200, todos)
}

func GetTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))
	todoID := uuid.FromStringOrNil(c.Params.ByName("todoID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, gin.H{"error": "TODO list not found"})
		return
	}

	var todo models.Todo

	db.Where(
		"user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listID, todoID,
	).Find(&todo)

	if todo.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(200, todo)
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}

func UpdateTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))
	todoID := uuid.FromStringOrNil(c.Params.ByName("todoID"))

	var newTodo models.Todo
	e := c.BindJSON(&newTodo)

	if e != nil {
		c.JSON(422, e)
		return
	}

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, gin.H{"error": "Todo list not found"})
		return
	}

	var todo models.Todo

	db.Where(
		"user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listID, todoID,
	).Find(&todo)

	if todo.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		updatedAt := time.Now()
		todo = models.Todo{
			ID:         todo.ID,
			Title:      newTodo.Title,
			Completed:  newTodo.Completed,
			Note:       newTodo.Note,
			TodoListID: listID,
			UserID:     todo.UserID,
			CreatedAt:  todo.CreatedAt,
			UpdatedAt:  &updatedAt,
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
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))
	todoID := uuid.FromStringOrNil(c.Params.ByName("todoID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, gin.H{"error": "Todo list not found"})
		return
	}

	var todo models.Todo

	db.Where(
		"user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listID, todoID,
	).Find(&todo)

	if todo.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		db.Delete(&todo)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}
