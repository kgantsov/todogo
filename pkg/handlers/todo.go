package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func OptionsTodo(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Auth-Token")
	c.Next()
}

// Create godoc
// @Summary Create a TODO in a specific TODO list
// @Schemes
// @Description Returns an newly created TODO in a specific TODO list
// @Tags TODO
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Param        body  body     models.Todo  true  "TODO"
// @Success      200  {object}  models.Todo
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/todo/ [post]
func CreateTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, ErrorSchema{Error: "Todo list not found"})
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
		c.JSON(422, ErrorSchema{Error: e.Error()})
	}
}

// Create godoc
// @Summary Get TODOs in a specific TODO list
// @Schemes
// @Description Returns a list of TODOs in a specific TODO list
// @Tags TODO
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Success      200  {object}  []models.Todo
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/todo/ [get]
func GetTodos(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, ErrorSchema{Error: "Todo list not found"})
		return
	}

	var todos []models.Todo

	db.Order("completed asc, priority desc, created_at asc").Where(
		"user_id = ? AND todo_list_id = ?", currentUser.ID, listID,
	).Find(&todos)

	c.JSON(200, todos)
}

// Create godoc
// @Summary Get a TODO in a specific TODO list
// @Schemes
// @Description Returns a TODO in a specific TODO list
// @Tags TODO
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Param        todoID    path     string  true  "ID of a TODO"
// @Success      200  {object}  models.Todo
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/todo/{todoID}/ [get]
func GetTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))
	todoID := uuid.FromStringOrNil(c.Params.ByName("todoID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, ErrorSchema{Error: "TODO list not found"})
		return
	}

	var todo models.Todo

	db.Where(
		"user_id = ? AND todo_list_id = ? AND id = ?", currentUser.ID, listID, todoID,
	).Find(&todo)

	if todo.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(200, todo)
	} else {
		c.JSON(404, ErrorSchema{Error: "Todo not found"})
	}
}

// Create godoc
// @Summary Update a TODO in a specific TODO list
// @Schemes
// @Description Updates and returns a TODO in a specific TODO list
// @Tags TODO
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Param        todoID    path     string  true  "ID of a TODO"
// @Param        body  body     models.Todo  true  "TODO"
// @Success      200  {object}  models.Todo
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/todo/{todoID}/ [put]
func UpdateTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
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
		c.JSON(404, ErrorSchema{Error: "Todo list not found"})
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
		c.JSON(404, ErrorSchema{Error: "Todo not found"})
	}
}

// Create godoc
// @Summary Delete a TODO from a specific TODO list
// @Schemes
// @Description Deletes a TODO from a specific TODO list
// @Tags TODO
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Param        todoID    path     string  true  "ID of a TODO"
// @Success      200  {object}  models.Todo
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/todo/{todoID}/ [delete]
func DeleteTodo(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	listID := uuid.FromStringOrNil(c.Params.ByName("listID"))
	todoID := uuid.FromStringOrNil(c.Params.ByName("todoID"))

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("user_id = ? AND id = ?", currentUser.ID, listID).First(&todoList)

	if todoList.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(404, ErrorSchema{Error: "Todo list not found"})
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
		c.JSON(404, ErrorSchema{Error: "Todo not found"})
	}
}
