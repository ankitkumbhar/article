package main

import (
	"fmt"
	"log"
	"net/http"

	"article/internal/database"
	"article/internal/handler"
	"article/internal/models"
	"article/internal/routes"

	_ "github.com/go-sql-driver/mysql"
)

const port = 8080

var app *handler.Application

// init
func init() {
	// initialize database
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	models := models.NewModels(db)

	app = handler.New(models)
}

func main() {
	logger := log.New(log.Default().Writer(), "logger: ", 1)

	// register routes
	r := routes.InitRoutes(app)

	logger.Println("server listening on port :", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)
	if err != nil {
		panic(err)
	}
}
