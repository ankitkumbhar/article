package response_test

import (
	"article/internal/handler"
	"article/internal/response"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Created(t *testing.T) {
	tests := []struct {
		name         string
		mockResp     func() *httptest.ResponseRecorder
		wantResp     handler.ArticleResponse
		wantRespBody response.Body
	}{
		{
			name: "response created",
			mockResp: func() *httptest.ResponseRecorder {
				resp := response.New()
				w := httptest.NewRecorder()
				data := handler.ArticleRequest{Title: "Test title", Content: "Test content", Author: "Test Author"}
				resp.Created(w, data)

				return w
			},
			wantRespBody: response.Body{Status: http.StatusCreated, Message: response.StatusSuccess},
			wantResp:     handler.ArticleResponse{Title: "Test title", Content: "Test content", Author: "Test Author"},
		},
		{
			name: "response success",
			mockResp: func() *httptest.ResponseRecorder {
				resp := response.New()
				w := httptest.NewRecorder()
				data := handler.ArticleRequest{Title: "Test title", Content: "Test content", Author: "Test Author"}
				resp.Success(w, data)

				return w
			},
			wantRespBody: response.Body{Status: http.StatusOK, Message: response.StatusSuccess},
			wantResp:     handler.ArticleResponse{Title: "Test title", Content: "Test content", Author: "Test Author"},
		},
		{
			name: "error bad request",
			mockResp: func() *httptest.ResponseRecorder {
				resp := response.New()
				w := httptest.NewRecorder()
				resp.BadRequest(w, "error message - bad request")

				return w
			},
			wantRespBody: response.Body{Status: http.StatusBadRequest, Message: "error message - bad request"},
		},
		{
			name: "error internal server error",
			mockResp: func() *httptest.ResponseRecorder {
				resp := response.New()
				w := httptest.NewRecorder()
				resp.InternalServerError(w, "error message - internal server error")

				return w
			},
			wantRespBody: response.Body{Status: http.StatusInternalServerError, Message: "error message - internal server error"},
		},
		{
			name: "error not found",
			mockResp: func() *httptest.ResponseRecorder {
				resp := response.New()
				w := httptest.NewRecorder()
				resp.NotFound(w, "error message - not found")

				return w
			},
			wantRespBody: response.Body{Status: http.StatusNotFound, Message: "error message - not found"},
		},
		{
			name: "error method not allowed",
			mockResp: func() *httptest.ResponseRecorder {
				resp := response.New()
				w := httptest.NewRecorder()
				resp.NotAllowed(w, "error message - method not allowed")

				return w
			},
			wantRespBody: response.Body{Status: http.StatusMethodNotAllowed, Message: "error message - method not allowed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.mockResp()

			resp := response.Body{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			if err != nil {
				t.Errorf("error unmarshalling response : %v", err)
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

			assert.Equal(t, resp.Status, tt.wantRespBody.Status)
			assert.Equal(t, gotResp.ID, tt.wantResp.ID)
			assert.Equal(t, resp.Status, tt.wantRespBody.Status)
			assert.Equal(t, resp.Message, tt.wantRespBody.Message)
		})
	}

}
