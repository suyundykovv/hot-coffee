package models

import "fmt"

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

func NewErrorResponse(statusCode int, message string, err interface{}) error {
	return fmt.Errorf(`{"status_code": %d, "message": "%s", "error": %v}`, statusCode, message, err)
}

func NewSuccessResponse(statusCode int, message string, data interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
