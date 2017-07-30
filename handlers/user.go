package handlers

import (
	"time"
	"crypto/sha256"
	"fmt"
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
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func CreateUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	var user models.User
	e := c.BindJSON(&user)

	if e == nil {
		user.Password = hashPassword(user.Password)

		db.Create(&user)
		c.JSON(201, user)
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

	db.Order("id asc").Find(&users)

	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	userId := c.Params.ByName("userId")

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

	userId := c.Params.ByName("userId")

	var newUser models.User
	e := c.BindJSON(&newUser)

	if e != nil {
		c.JSON(422, e)
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

	userId := c.Params.ByName("userId")

	var user models.User

	db.First(&user, userId)

	if user.ID != 0 {
		db.Delete(&user)

		c.Writer.WriteHeader(204)
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}
