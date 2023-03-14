package models

import "database/sql"

// Application holds database object
type Application struct {
	db *sql.DB
}

// Models holds article interface
type Models struct {
	Article ArticleStore
}

// NewModels store db object and return models
func NewModels(db *sql.DB) *Models {
	app := Application{db: db}

	return &Models{
		Article: &article{app: &app},
	}
}
