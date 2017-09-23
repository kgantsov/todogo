package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/handlers"
	"github.com/kgantsov/todogo/models"
	"github.com/newrelic/go-agent"
	"gopkg.in/gin-gonic/gin.v1"
)

func DBMiddleware(db gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func NewRelicMiddleware(appName, newRelicKey string) gin.HandlerFunc {
	config := newrelic.NewConfig(appName, newRelicKey)
	app, _ := newrelic.NewApplication(config)

	return func(c *gin.Context) {
		name := c.HandlerName()

		txn := app.StartTransaction(name, c.Writer, c.Request)
		defer txn.End()

		c.Next()
	}
}

func main() {
	debug := flag.Bool("debug", false, "Debug flag")
	port := flag.Int("port", 8080, "Port")

	flag.Parse()

	db := models.InitDbURI(
		os.Getenv("DB_URI"),
		*debug,
	)
	models.CreateTables(db)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(NewRelicMiddleware(os.Getenv("NEWRELIC_APP_NAME"), os.Getenv("NEWRELIC_APP_KEY")))

	if *debug {
		r.Use(gin.Logger())
		r.Use(CorsMiddleware())
	}

	r.Use(DBMiddleware(*db))

	handlers.DefineRoutes(db, r)
	r.Run(fmt.Sprintf(":%d", *port))
}
