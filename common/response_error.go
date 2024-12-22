package common

import (
	"encoding/json"
	"net/http"
)

const (
	// Error messages
	BadRequestMsg             = "Bad request"
	InternalServerErrorMsg    = "Internal server error"
	FailedToCreateResponseMsg = "Failed to create a response"
)

type ResponseError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func NewResponseError(statusCode int, msg string) *ResponseError {
	return &ResponseError{
		StatusCode: statusCode,
		Message:    msg,
	}
}

func (re *ResponseError) WriteError(w http.ResponseWriter) {
	reJson, _ := json.Marshal(re)
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(re.StatusCode)
	w.Write(reJson)
}

func (re *ResponseError) Error() string {
	return re.Message
}
