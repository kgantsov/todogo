package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/handlers"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
	"os"
	"flag"
	"fmt"
)

func DBMiddleware(db gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}

func main() {
	debug := flag.Bool("debug", false, "Debug flag")
	port := flag.Int("port", 8080, "Port")

	flag.Parse()

	db := models.InitDb(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		*debug,
	)
	models.CreateTables(db)

	r := gin.New()

	r.Use(gin.Recovery())

	if *debug == true {
		r.Use(gin.Logger())
	}

	r.Use(DBMiddleware(*db))

	handlers.DefineRoutes(r)
	r.Run(fmt.Sprintf(":%d", *port))
}