package handlers

import (
	"fmt"
	"time"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
)

var mySigningKey = []byte("secret")

type LoginForm struct {
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"Password" json:"Password" binding:"required"`
}

func OptionsLogin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func Login(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "No connection to DB"})
	}

	var loginForm LoginForm
	e := c.BindJSON(&loginForm)

	if e != nil {
		c.JSON(403, gin.H{"error": e})
	}

	var user models.User

	db.Where("email = ?", loginForm.Email).Find(&user)

	if user.ID != 0 {
		if hashPassword(loginForm.Password) == user.Password {
			token, err := createToken(user.ID)
			if err == nil {
				c.JSON(201, gin.H{"token": token})
				return
			}
		}
	}
	c.JSON(403, gin.H{"error": "Login or password is incorrect"})
}

func createToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
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
			db.First(&user, claims["user_id"])

			if user.ID != 0 {
				return user, true
			}
		}
	}
	return user, false
}
