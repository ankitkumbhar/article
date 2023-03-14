package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusSuccess = "Success"
	StatusError   = "error"
)

type Response struct{}

// Body response body structure for API
type Body struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// New retuns response obj
func New() *Response {
	return &Response{}
}

// SetStatus sets http status code in body
func (b *Body) SetStatus(statusCode int) {
	b.Status = statusCode
}

// SetMessage sets response message
func (b *Body) SetMessage(msg string) {
	b.Message = msg
}

// SetData sets response data
func (b *Body) SetData(data interface{}) {
	b.Data = data
}

// GetStatus gets response status
func (b *Body) GetStatus() int {
	return b.Status
}

// Created handles 201 success response
func (r *Response) Created(w http.ResponseWriter, data interface{}) {
	b := Body{}
	b.SetStatus(http.StatusCreated)
	b.SetMessage(StatusSuccess)

	SendResponse(w, &b, data)
}

// Success handles 200 success response
func (r *Response) Success(w http.ResponseWriter, data interface{}) {
	b := Body{}
	b.SetStatus(http.StatusOK)
	b.SetMessage(StatusSuccess)

	SendResponse(w, &b, data)
}

// BadRequest handles 400 error response
func (r *Response) BadRequest(w http.ResponseWriter, msg string, data ...interface{}) {
	b := Body{}
	b.SetStatus(http.StatusBadRequest)
	b.SetMessage(msg)

	SendResponse(w, &b, data)
}

// InternalServerError handles 500 error response
func (r *Response) InternalServerError(w http.ResponseWriter, msg string, data ...interface{}) {
	b := Body{}
	b.SetStatus(http.StatusInternalServerError)
	b.SetMessage(msg)

	SendResponse(w, &b, data)
}

// NotFound handles 404 error response
func (r *Response) NotFound(w http.ResponseWriter, msg string, data ...interface{}) {
	b := Body{}
	b.SetStatus(http.StatusNotFound)
	b.SetMessage(msg)

	SendResponse(w, &b, data)
}

// NotAllowed handles 405 error response
func (r *Response) NotAllowed(w http.ResponseWriter, msg string, data ...interface{}) {
	b := Body{}
	b.SetStatus(http.StatusMethodNotAllowed)
	b.SetMessage(msg)

	SendResponse(w, &b, data)
}

// SendResponse
func SendResponse(w http.ResponseWriter, b *Body, data interface{}) {
	b.SetData(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(b.GetStatus())

	// write json encoded to writer
	json.NewEncoder(w).Encode(b)
}
