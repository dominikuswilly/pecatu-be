package http

import "fmt"

type Response struct {
	Code    string      `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Code:    fmt.Sprintf("%d", code),
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    fmt.Sprintf("%d", code),
		Status:  "error",
		Message: message,
	}
}
