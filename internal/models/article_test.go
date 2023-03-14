package models_test

import (
	"article/internal/handler"
	"article/internal/models"
	"article/internal/response"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Store(t *testing.T) {

	tests := []struct {
		name         string
		mockDB       func() *sql.DB
		wantResp     handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name: "success",
			mockDB: func() *sql.DB {
				// create sql mock database connection
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error opening a stub database connection %v", err)
				}

				// mock expected query
				mock.ExpectExec("INSERT INTO article").WillReturnResult(sqlmock.NewResult(1, 1))

				return db
			},
		},
		{
			name: "error",
			mockDB: func() *sql.DB {
				// create sql mock database connection
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error opening a stub database connection %v", err)
				}

				// mock expected query
				mock.ExpectExec("INSERT INTO article").WillReturnError(errors.New("db error"))

				return db
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mockDB()

			// store mocked db object in models
			a := models.NewModels(db)

			// call model function
			gotID, err := a.Article.Store(&models.Article{})
			if err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "db error")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, gotID, int64(1))
			}
		})
	}

}

func Test_GetByID(t *testing.T) {

	tests := []struct {
		name         string
		mockDB       func() *sql.DB
		wantResp     handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name: "success",
			mockDB: func() *sql.DB {
				// create sql mock database connection
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error opening a stub database connection %v", err)
				}

				// mock return valid rows
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author"}).AddRow(int64(1), "Test title", "Test content", "Test author")
				mock.ExpectQuery("SELECT id, title, content, author FROM article").WillReturnRows(rows)

				return db
			},
			wantResp: handler.ArticleResponse{ID: 1, Title: "Test title", Content: "Test content", Author: "Test author"},
		},
		{
			name: "error : select query error",
			mockDB: func() *sql.DB {
				// create sql mock database connection
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error opening a stub database connection %v", err)
				}

				// mock return error
				mock.ExpectQuery("SELECT id, title, content, author FROM article").WillReturnError(errors.New("db error"))

				return db
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mockDB()

			// store mocked db object in models
			a := models.NewModels(db)

			// call model function
			gotResp, err := a.Article.GetByID(int(tt.wantResp.ID))
			if err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "db error")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, gotResp.ID, int(tt.wantResp.ID))
				assert.Equal(t, gotResp.Title, tt.wantResp.Title)
				assert.Equal(t, gotResp.Content, tt.wantResp.Content)
				assert.Equal(t, gotResp.Author, tt.wantResp.Author)
			}
		})
	}

}

func Test_GetAll(t *testing.T) {

	tests := []struct {
		name         string
		mockDB       func() *sql.DB
		wantResp     []handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name: "success",
			mockDB: func() *sql.DB {
				// create sql mock database connection
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error opening a stub database connection %v", err)
				}

				// mock return valid rows
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author"}).AddRow(int64(1), "Test title", "Test content", "Test author")
				mock.ExpectQuery("SELECT id, title, content, author FROM article").WillReturnRows(rows)

				return db
			},
			wantResp: []handler.ArticleResponse{
				{
					ID:      1,
					Title:   "Test title",
					Content: "Test content",
					Author:  "Test author"},
			},
		},
		{
			name: "error : select query error",
			mockDB: func() *sql.DB {
				// create sql mock database connection
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error opening a stub database connection %v", err)
				}

				// mock return error
				mock.ExpectQuery("SELECT id, title, content, author FROM article").WillReturnError(errors.New("db error"))

				return db
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mockDB()

			// store mocked db object in models
			a := models.NewModels(db)

			// call model function
			gotResp, err := a.Article.GetAll()
			if err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "db error")
			} else {
				assert.Nil(t, err)

				for key, val := range gotResp {
					assert.Equal(t, val.ID, int(tt.wantResp[key].ID))
					assert.Equal(t, val.Title, tt.wantResp[key].Title)
					assert.Equal(t, val.Content, tt.wantResp[key].Content)
					assert.Equal(t, val.Author, tt.wantResp[key].Author)
				}
			}
		})
	}
}
