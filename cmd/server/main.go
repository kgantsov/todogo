package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/todogo/pkg/handlers"
	"github.com/kgantsov/todogo/pkg/models"
	newrelic "github.com/newrelic/go-agent"
	"gorm.io/gorm"

	docs "github.com/kgantsov/todogo/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title TODOGO API
// @version 1.0
// @description This is a sample TODO application.

// @BasePath /
// @securityDefinitions.apikey  HttpBearer
// @in                          header
// @name                        Auth-Token
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

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	NewrelicAppName := os.Getenv("NEWRELIC_APP_NAME")
	NewrelicAppKey := os.Getenv("NEWRELIC_APP_KEY")

	if len(NewrelicAppName) > 0 && len(NewrelicAppKey) > 0 {
		r.Use(NewRelicMiddleware(os.Getenv("NEWRELIC_APP_NAME"), os.Getenv("NEWRELIC_APP_KEY")))
	}

	if *debug {
		r.Use(gin.Logger())
		r.Use(CorsMiddleware())
	}

	r.Use(DBMiddleware(*db))

	handlers.DefineRoutes(db, r)
	r.Run(fmt.Sprintf(":%d", *port))
}
