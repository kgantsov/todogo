package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func OptionsTodoList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Auth-Token")
	c.Next()
}

// Create godoc
// @Summary Create a TODO list
// @Schemes
// @Description Returns an newly created TODO list
// @Tags TODO list
// @Accept json
// @Produce json
// @Param        body  body     models.TodoList  true  "TODO list"
// @Success      200  {object}  models.TodoList
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/ [post]
func CreateTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	var todoList models.TodoList
	e := c.BindJSON(&todoList)

	if e == nil {
		currentUser := c.MustGet("CurrentUser").(models.User)
		todoList.ID = uuid.NewV4()
		todoList.UserID = currentUser.ID

		db.Create(&todoList)
		c.JSON(201, todoList)
	} else {
		c.JSON(422, gin.H{"error": e})
	}
}

// Create godoc
// @Summary Get a TODO lists
// @Schemes
// @Description Returns TODO lists
// @Tags TODO list
// @Accept json
// @Produce json
// @Success      200  {object}  []models.TodoList
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/ [get]
func GetTodoLists(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	var todoLists []models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Order("created_at asc").Where("user_id = ?", currentUser.ID).Find(&todoLists)

	c.JSON(200, todoLists)
}

// Create godoc
// @Summary Get a TODO list
// @Schemes
// @Description Returns a particular TODO list by its ID
// @Tags TODO list
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Success      200  {object}  models.TodoList
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/ [get]
func GetTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	listID := c.Params.ByName("listID")

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("id = ? AND user_id = ?", listID, currentUser.ID).First(&todoList)

	if todoList.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(200, todoList)
	} else {
		c.JSON(404, gin.H{"error": "TODO list not found"})
	}
}

// Create godoc
// @Summary Update a TODO list
// @Schemes
// @Description Updates a particular TODO list by its ID
// @Tags TODO list
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Param        body  body     models.TodoList  true  "TODO list"
// @Success      200  {object}  models.TodoList
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/ [put]
func UpdateTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	listID := c.Params.ByName("listID")

	var newTodoList models.TodoList
	e := c.BindJSON(&newTodoList)

	if e != nil {
		c.JSON(422, e)
		return
	}

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("id = ? AND user_id = ?", listID, currentUser.ID).First(&todoList)

	if todoList.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		updatedAt := time.Now()
		todoList = models.TodoList{
			ID:        todoList.ID,
			Title:     newTodoList.Title,
			CreatedAt: todoList.CreatedAt,
			UpdatedAt: &updatedAt,
		}

		db.Save(&todoList)

		c.JSON(200, todoList)
	} else {
		c.JSON(404, gin.H{"error": "Todo List not found"})
	}
}

// Create godoc
// @Summary Delte a TODO list
// @Schemes
// @Description Deletes a particular TODO list by its ID
// @Tags TODO list
// @Accept json
// @Produce json
// @Param        listID    path     string  true  "ID of a TODO list"
// @Success      200  {object}  models.TodoList
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /list/{listID}/ [delete]
func DeleteTodoList(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	listID := c.Params.ByName("listID")

	var todoList models.TodoList

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Where("id = ? AND user_id = ?", listID, currentUser.ID).First(&todoList)

	if todoList.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		db.Delete(&todoList)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, gin.H{"error": "Todo List not found"})
	}
}
