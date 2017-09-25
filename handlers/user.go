package handlers

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
)

func hashPassword(password string) string {
	password_hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", password_hash)
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

func CreateUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	var user models.User
	e := c.BindJSON(&user)

	if e == nil {
		if !isValidEmail(user.Email) {
			c.JSON(400, gin.H{"error": "Email address is not valid"})
			return
		}
		user.Password = hashPassword(user.Password)

		var exisingUser models.User

		db.Where("email = ?", user.Email).First(&exisingUser)

		if exisingUser.ID != 0 {
			c.JSON(409, gin.H{"error": "User with this email already exists"})
		} else {
			db.Create(&user)
			c.JSON(201, user)
		}

	} else {
		c.JSON(422, gin.H{"error": e})
	}
}

func GetUsers(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	var users []models.User

	currentUser := c.MustGet("CurrentUser").(models.User)

	db.Order("id asc").Where("id = ?", currentUser.ID).Find(&users)

	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	userId, _ := strconv.ParseUint(c.Params.ByName("userId"), 0, 64)

	currentUser := c.MustGet("CurrentUser").(models.User)

	if uint64(userId) != currentUser.ID {
		c.JSON(403, gin.H{"error": "Access denied"})
	}

	var user models.User

	db.First(&user, userId)

	if user.ID != 0 {
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func UpdateUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	userId, _ := strconv.ParseUint(c.Params.ByName("userId"), 0, 64)

	currentUser := c.MustGet("CurrentUser").(models.User)

	if uint64(userId) != currentUser.ID {
		c.JSON(403, gin.H{"error": "Access denied"})
	}

	var newUser models.User
	e := c.BindJSON(&newUser)

	if e != nil {
		c.JSON(422, e)
	}

	if !isValidEmail(newUser.Email) {
		c.JSON(400, gin.H{"error": "Email address is not valid"})
		return
	}

	var user models.User

	db.First(&user, userId)

	if user.ID != 0 {
		user = models.User{
			ID:        user.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			Password:  hashPassword(newUser.Password),
			CreatedAt: user.CreatedAt,
			UpdatedAt: time.Now(),
		}

		db.Save(&user)

		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func DeleteUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	userId, _ := strconv.ParseUint(c.Params.ByName("userId"), 0, 64)

	currentUser := c.MustGet("CurrentUser").(models.User)

	if uint64(userId) != currentUser.ID {
		c.JSON(403, gin.H{"error": "Access denied"})
	}

	var user models.User

	db.First(&user, userId)

	if user.ID != 0 {
		db.Delete(&user)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}
