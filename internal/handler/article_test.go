package handler_test

import (
	"article/internal/handler"
	"article/internal/models"
	"article/internal/response"
	"article/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateArticle(t *testing.T) {
	type args struct {
		req handler.ArticleRequest
	}

	tests := []struct {
		name         string
		args         args
		mockDB       func() *handler.Application
		wantResp     handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name: "success",
			args: args{req: handler.ArticleRequest{Title: "Test title", Content: "Test content", Author: "Test Author"}},
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().Store(mock.Anything).Return(1, nil)

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantResp:     handler.ArticleResponse{ID: 1},
			wantRespBody: response.Body{Status: http.StatusCreated, Message: response.StatusSuccess},
		},
		{
			name: "validation error",
			args: args{req: handler.ArticleRequest{Title: " ", Content: "Test content", Author: "Test Author"}},
			mockDB: func() *handler.Application {
				return handler.New(&models.Models{})
			},
			wantRespBody: response.Body{Status: http.StatusBadRequest, Message: "Field validation for 'Title' failed on the 'required' tag"},
		},
		{
			name: "error",
			args: args{req: handler.ArticleRequest{Title: "Test title", Content: "Test content", Author: "Test Author"}},
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().Store(mock.Anything).Return(0, errors.New("error storing article"))

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantRespBody: response.Body{Status: http.StatusInternalServerError, Message: "error storing article"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock database calls
			app := tt.mockDB()

			handlerFunc := app.CreateArticle()
			resp, err := callEndpoint(t, &tt.args.req, handlerFunc, nil)
			if err != nil {
				t.Errorf("error in call endpoint : %v", err)
			}

			// convert response data into struct
			var gotResp handler.ArticleResponse
			aa, err := json.Marshal(resp.Data)
			if err != nil {
				t.Error("error marshalling response data to bytes", err)
			}

			err = json.Unmarshal(aa, &gotResp)
			if err != nil {
				t.Error("error unmarshalling response data", err)
			}

			assert.Equal(t, gotResp.ID, tt.wantResp.ID)
			assert.Equal(t, resp.Status, tt.wantRespBody.Status)
			assert.Equal(t, resp.Message, tt.wantRespBody.Message)
		})
	}
}

func Test_GetArticle(t *testing.T) {
	tests := []struct {
		name         string
		urlParams    map[string]string
		mockDB       func() *handler.Application
		wantResp     handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name:      "success",
			urlParams: map[string]string{"article_id": "1"},
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().GetByID(mock.Anything).Return(&models.Article{
					ID:      1,
					Title:   "Test title",
					Content: "Test content",
					Author:  "Test author",
				}, nil)

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantResp:     handler.ArticleResponse{ID: 1, Title: "Test title", Content: "Test content", Author: "Test author"},
			wantRespBody: response.Body{Status: http.StatusOK, Message: response.StatusSuccess},
		},
		{
			name:      "error : empty articleID",
			urlParams: map[string]string{"article_id": ""},
			mockDB: func() *handler.Application {
				return handler.New(nil)
			},
			wantRespBody: response.Body{Status: http.StatusBadRequest, Message: "please provide article id"},
		},
		{
			name:      "error : invalid article id",
			urlParams: map[string]string{"article_id": "aa"},
			mockDB: func() *handler.Application {
				return handler.New(&models.Models{})
			},
			wantRespBody: response.Body{Status: http.StatusInternalServerError, Message: "error converting article id"},
		},
		{
			name:      "error : database error",
			urlParams: map[string]string{"article_id": "1"},
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().GetByID(mock.Anything).Return(nil, errors.New("error fetching article by articleID"))

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantRespBody: response.Body{Status: http.StatusInternalServerError, Message: "error fetching article by articleID"},
		},
		{
			name:      "error : article id not found",
			urlParams: map[string]string{"article_id": "1"},
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().GetByID(mock.Anything).Return(&models.Article{}, nil)

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantRespBody: response.Body{Status: http.StatusBadRequest, Message: "invalid article id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock database calls
			app := tt.mockDB()

			handlerFunc := app.GetArticle()
			resp, err := callEndpoint(t, nil, handlerFunc, tt.urlParams)
			if err != nil {
				t.Errorf("error in call endpoint : %v", err)
			}

			// convert response data into struct
			var gotResp []*handler.ArticleResponse
			aa, err := json.Marshal(resp.Data)
			if err != nil {
				t.Error("error marshalling response data to bytes", err)
			}

			err = json.Unmarshal(aa, &gotResp)
			if err != nil {
				t.Error("error unmarshalling response data", err)
			}

			assert.Equal(t, resp.Status, tt.wantRespBody.Status)
			assert.Equal(t, resp.Message, tt.wantRespBody.Message)
			if resp.Message == "Success" {
				assert.Equal(t, gotResp[0].ID, tt.wantResp.ID)
				assert.Equal(t, gotResp[0].Title, tt.wantResp.Title)
			}
		})
	}
}

func Test_GetArticles(t *testing.T) {
	tests := []struct {
		name         string
		urlParams    map[string]string
		mockDB       func() *handler.Application
		wantResp     handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name: "success",
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().GetAll().Return([]*models.Article{
					{
						ID:      1,
						Title:   "Test title",
						Content: "Test content",
						Author:  "Test author",
					},
				}, nil)

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantResp:     handler.ArticleResponse{ID: 1, Title: "Test title", Content: "Test content", Author: "Test author"},
			wantRespBody: response.Body{Status: http.StatusOK, Message: response.StatusSuccess},
		},
		{
			name: "error : database error",
			mockDB: func() *handler.Application {
				articleMock := mocks.NewArticleStore(t)
				articleMock.EXPECT().GetAll().Return(nil, errors.New("error fetching all articles"))

				m := models.Models{
					Article: articleMock,
				}

				return handler.New(&m)
			},
			wantRespBody: response.Body{Status: http.StatusInternalServerError, Message: "error fetching all articles"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock database calls
			app := tt.mockDB()

			handlerFunc := app.GetArticles()

			resp, err := callEndpoint(t, nil, handlerFunc, nil)
			if err != nil {
				t.Errorf("error in call endpoint : %v", err)
			}

			// convert response data into struct
			var gotResp []*handler.ArticleResponse
			aa, err := json.Marshal(resp.Data)
			if err != nil {
				t.Error("error marshalling response data to bytes", err)
			}

			err = json.Unmarshal(aa, &gotResp)
			if err != nil {
				t.Error("error unmarshalling response data", err)
			}

			assert.Equal(t, resp.Status, tt.wantRespBody.Status)
			assert.Equal(t, resp.Message, tt.wantRespBody.Message)
			if resp.Message == "Success" {
				assert.Equal(t, gotResp[0].ID, tt.wantResp.ID)
				assert.Equal(t, gotResp[0].Title, tt.wantResp.Title)
			}
		})
	}
}

// callEndpoint creates a request and make a http call
func callEndpoint(t *testing.T, req *handler.ArticleRequest, handlerFunc http.HandlerFunc, urlParams map[string]string) (*response.Body, error) {
	w := httptest.NewRecorder()

	rawReq, _ := json.Marshal(req)

	// create a request
	r, err := http.NewRequest(mock.Anything, "/articles", bytes.NewBuffer(rawReq))
	if err != nil {
		t.Fatal(err)
	}

	// appends a urlParams at the end of route
	r = setURLParams(r, urlParams)

	// server http call
	handlerFunc.ServeHTTP(w, r)

	resp := response.Body{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("error unmarshalling response : %v", err)
	}

	return &resp, nil
}

// setURLParams appends a urlParams at the end of route
func setURLParams(req *http.Request, urlParams map[string]string) *http.Request {
	if len(urlParams) > 0 {
		routeContext := chi.NewRouteContext()
		for k, v := range urlParams {
			routeContext.URLParams.Add(k, v)
		}

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
	}

	return req
}
