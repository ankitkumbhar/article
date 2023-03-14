package handler

import (
	"article/internal/models"
	"article/internal/response"
	"log"

	"github.com/go-playground/validator/v10"
)

// Application used to hold objects
type Application struct {
	models   *models.Models
	response response.Response
	validate *validator.Validate
	logger   *log.Logger
}

func New(models *models.Models) *Application {
	return &Application{
		models:   models,
		response: *response.New(),
		validate: validator.New(),
		logger:   log.New(log.Default().Writer(), "logger: ", 1),
	}
}
