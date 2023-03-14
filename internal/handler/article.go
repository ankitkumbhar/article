package handler

import (
	"article/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// ArticleRequest used in request
type ArticleRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
}

// ArticleResponse used in response
type ArticleResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Author  string `json:"author,omitempty"`
}

// CreateArticle stores an article with given details
func (app *Application) CreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ArticleRequest

		// validate request body
		err := app.validateRequest(w, r, &req)
		if err != nil {
			return
		}

		// prepare article model
		article := models.Article{
			Title:   req.Title,
			Content: req.Content,
			Author:  req.Author,
		}

		// store article
		insertedID, err := app.models.Article.Store(&article)
		if err != nil {
			app.logger.Println("error storing article : ", err)
			app.response.InternalServerError(w, "error storing article")

			return
		}

		// prepare response
		resp := ArticleResponse{
			ID: insertedID,
		}

		app.response.Created(w, resp)
	}
}

// GetArticle fetch an article using articleID
func (app *Application) GetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch articleID from url params
		articleID := chi.URLParam(r, "article_id")
		if articleID == "" {
			app.logger.Println("article id not passed")
			app.response.BadRequest(w, "please provide article id")

			return
		}

		// convert articleID from string to integer
		id, err := strconv.Atoi(articleID)
		if err != nil {
			app.logger.Println("error converting articleID from string to integer")
			app.response.InternalServerError(w, "error converting article id")

			return
		}

		// get article by id
		article, err := app.models.Article.GetByID(id)
		if err != nil {
			app.logger.Println("error fetching article by articleID : ", err)
			app.response.InternalServerError(w, "error fetching article by articleID")

			return
		}

		if article.ID == 0 {
			app.logger.Println("invalid article id")
			app.response.BadRequest(w, "invalid article id")

			return
		}

		// prepare response
		resp := []ArticleResponse{
			{
				ID:      int64(article.ID),
				Title:   article.Title,
				Content: article.Content,
				Author:  article.Author,
			},
		}

		app.response.Success(w, resp)
	}
}

// GetArticles fetchs all article
func (app *Application) GetArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get article by id
		articles, err := app.models.Article.GetAll()
		if err != nil {
			app.logger.Println("error fetching all article : ", err)
			app.response.InternalServerError(w, "error fetching all articles")

			return
		}

		// prepare response
		resp := []ArticleResponse{}

		for _, val := range articles {
			a := ArticleResponse{
				ID:      int64(val.ID),
				Title:   val.Title,
				Content: val.Content,
				Author:  val.Author,
			}

			resp = append(resp, a)
		}

		app.response.Success(w, resp)
	}
}

// validateRequest validates request body
func (app *Application) validateRequest(w http.ResponseWriter, r *http.Request, req *ArticleRequest) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.logger.Println("error decoding request body : ", err)
		app.response.BadRequest(w, "invalid request")

		return err
	}

	// remove white space
	// trimspace removes all leading and trailing space
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.Author = strings.TrimSpace(req.Author)

	// validate request body
	err = app.validate.Struct(req)
	if err != nil {
		app.logger.Println("error validating request : ", err)

		var errorBag []string
		for _, v := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, strings.Split(v.Error(), "Error:")[1])
		}

		app.response.BadRequest(w, fmt.Sprint(strings.Join(errorBag[:], ", ")))

		return err
	}

	return nil
}
