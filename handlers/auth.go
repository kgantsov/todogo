package handlers

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"gopkg.in/gin-gonic/gin.v1"
)

var mySigningKey = []byte("secret")

type LoginForm struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"Password" json:"Password" binding:"required"`
}

func OptionsLogin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Auth-Token")
	c.Next()
}

func Login(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
		return
	}

	var loginForm LoginForm
	e := c.BindJSON(&loginForm)

	if e != nil {
		c.JSON(401, gin.H{"error": e})
		return
	}

	var user models.User

	db.Where("email = ?", loginForm.Email).Find(&user)

	if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		if hashPassword(loginForm.Password) == user.Password {
			token, err := createToken(user.ID)
			if err == nil {
				c.JSON(201, gin.H{"token": token})
				return
			}
		}
	}
	c.JSON(401, gin.H{"error": "Login or password is incorrect"})
}

func createToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)

	return tokenString, err
}

func validateToken(db *gorm.DB, tokenString string) (models.User, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	var user models.User

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId, _ := claims["user_id"].(string)

			db.Where("id = ?", userId).First(&user)

			if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
				return user, true
			}
		}
	}
	return user, false
}
