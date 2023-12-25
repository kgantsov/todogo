package handlers

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func hashPassword(password string) string {
	passwordHash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", passwordHash)
}

func OptionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Auth-Token")
	c.Next()
}

func isValidEmail(email string) bool {
	regex, _ := regexp.Compile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`)

	if regex.MatchString(email) {
		return true
	}
	return false
}

// Create godoc
// @Summary Create a user
// @Schemes
// @Description Returns an newly created user
// @Tags registration
// @Accept json
// @Produce json
// @Param        body  body     models.User  true  "User object"
// @Success      200  {object}  models.User
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Router       /user/ [post]
func CreateUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	var user models.User
	e := c.BindJSON(&user)

	if e == nil {
		if !isValidEmail(user.Email) {
			c.JSON(400, ErrorSchema{Error: "Email address is not valid"})
			return
		}
		user.Password = hashPassword(user.Password)

		var exisingUser models.User

		db.Where("email = ?", user.Email).First(&exisingUser)

		if exisingUser.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
			c.JSON(409, ErrorSchema{Error: "User with this email already exists"})
		} else {
			user.ID = uuid.NewV4()
			db.Create(&user)
			c.JSON(201, user)
		}

	} else {
		c.JSON(422, ErrorSchema{Error: e.Error()})
	}
}

// Create godoc
// @Summary Get list of users
// @Schemes
// @Description Returns a list of users
// @Tags users
// @Accept json
// @Produce json
// @Success      200  {object}  []models.User
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /user/ [get]
func GetUsers(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	var users []models.User

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Order("id asc").Where("id = ?", currentUser.ID).Find(&users)

	c.JSON(200, users)
}

// Create godoc
// @Summary Get a user
// @Schemes
// @Description Returns a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param        userID    path     string  true  "ID of a User"
// @Success      200  {object}  models.User
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /user/{userID}/ [get]
func GetUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	userID := uuid.FromStringOrNil(c.Params.ByName("userID"))

	currentUser := c.MustGet("CurrentUser").(models.User)

	if userID != currentUser.ID {
		c.JSON(403, ErrorSchema{Error: "Access denied"})
		return
	}

	var user models.User

	db.Where("id = ?", userID).First(&user)

	if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		c.JSON(200, user)
	} else {
		c.JSON(404, ErrorSchema{Error: "User not found"})
	}
}

// Create godoc
// @Summary Update a user
// @Schemes
// @Description Updates a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param        userID    path     string  true  "ID of a User"
// @Param        body  body     models.User  true  "User object"
// @Success      200  {object}  models.User
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /user/{userID}/ [put]
func UpdateUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	userID := uuid.FromStringOrNil(c.Params.ByName("userID"))

	currentUser := c.MustGet("CurrentUser").(models.User)

	if userID != currentUser.ID {
		c.JSON(403, ErrorSchema{Error: "Access denied"})
		return
	}

	var newUser models.User
	e := c.BindJSON(&newUser)

	if e != nil {
		c.JSON(422, e)
		return
	}

	if !isValidEmail(newUser.Email) {
		c.JSON(400, ErrorSchema{Error: "Email address is not valid"})
		return
	}

	var user models.User

	db.Where("id = ?", userID).First(&user)

	if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		updatedAt := time.Now()
		user = models.User{
			ID:        user.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			Password:  hashPassword(newUser.Password),
			CreatedAt: user.CreatedAt,
			UpdatedAt: &updatedAt,
		}

		db.Save(&user)

		c.JSON(200, user)
	} else {
		c.JSON(404, ErrorSchema{Error: "User not found"})
	}
}

// Create godoc
// @Summary Delete a user
// @Schemes
// @Description Deletes a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param        userID    path     string  true  "ID of a User"
// @Success      204  {object}  models.User
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Security     HttpBearer
// @Router       /user/{userID}/ [delete]
func DeleteUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	userID := uuid.FromStringOrNil(c.Params.ByName("userID"))

	currentUser := c.MustGet("CurrentUser").(models.User)

	if userID != currentUser.ID {
		c.JSON(403, ErrorSchema{Error: "Access denied"})
		return
	}

	var user models.User

	db.Where("id = ?", userID).First(&user)

	if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		db.Delete(&user)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, ErrorSchema{Error: "User not found"})
	}
}
