package http

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Status:  "error",
		Message: message,
	}
}
