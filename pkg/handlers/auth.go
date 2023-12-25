package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"gorm.io/gorm"
)

var mySigningKey = []byte("secret")

func OptionsLogin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Auth-Token")
	c.Next()
}

// @BasePath /api/v1

// Create godoc
// @Summary Create auth token
// @Schemes
// @Description Returns an newly created authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Param        body  body     LoginSchema  true  "User credentials"
// @Success      200  {object}  TokenSchema
// @Failure      401  {object}  ErrorSchema
// @Failure      500  {object}  ErrorSchema
// @Router       /auth/login/ [post]
func Login(c *gin.Context) {
	db, ok := c.MustGet("db").(gorm.DB)
	if !ok {
		c.JSON(500, ErrorSchema{Error: "No connection to DB"})
		return
	}

	var loginSchema LoginSchema
	e := c.BindJSON(&loginSchema)

	if e != nil {
		c.JSON(401, ErrorSchema{Error: e.Error()})
		return
	}

	var user models.User

	db.Where("email = ?", loginSchema.Email).Find(&user)

	if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		if hashPassword(loginSchema.Password) == user.Password {
			token, err := createToken(user.ID)
			if err == nil {
				c.JSON(201, TokenSchema{Token: token})
				return
			}
		}
	}
	c.JSON(401, ErrorSchema{Error: "Login or password is incorrect"})
}

func createToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 31).Unix(),
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
			userID, _ := claims["user_id"].(string)

			db.Where("id = ?", userID).First(&user)

			if user.ID != uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
				return user, true
			}
		}
	}
	return user, false
}
