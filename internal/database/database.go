package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	username = "root"
	password = "root"
	database = "article"
	host     = "localhost"
)

// InitDB used to connect mysql database
func InitDB() (*sql.DB, error) {
	logger := log.New(log.Default().Writer(), "", 1)

	// get username from env
	if envUsername := os.Getenv("DB_USERNAME"); envUsername != "" {
		username = envUsername
	}

	// get password from env
	if envPassword := os.Getenv("DB_PASSWORD"); envPassword != "" {
		password = envPassword
	}

	// get host from env
	if envHost := os.Getenv("DB_HOST"); envHost != "" {
		host = envHost
	}

	// open database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, host, database))
	if err != nil {
		logger.Println("error connecting database : ", err)

		return nil, err
	}

	return db, nil
}
