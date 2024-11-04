package utils

import (
	"net/http"

	"example.com/RMS/errors"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}
type ErrorContent struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Args    interface{} `json:"details"`
}

type ErrorResponse struct {
	Error ErrorContent `json:"error"`
}

func RenderSuccess(data interface{}) SuccessResponse {
	response := SuccessResponse{
		Data: data,
	}
	return response
}

func renderErrorMsg(code string, message string, args interface{}) ErrorResponse {
	errorContent := ErrorContent{
		Code:    code,
		Message: message,
		Args:    args,
	}
	response := ErrorResponse{
		Error: errorContent,
	}
	return response
}

func RenderError(err error, args interface{}, customMessage ...string) (int, ErrorResponse) {
	var httpStatus int
	var code string
	var message string

	switch err {
	case errors.ErrInvalidRequest:
		code = err.Error()
		httpStatus = http.StatusBadRequest
		message = "invalid request"
	case errors.ErrUnauthorized:
		code = err.Error()
		httpStatus = http.StatusUnauthorized
		message = "unauthorized access"
	case errors.ErrNotFound:
		code = err.Error()
		httpStatus = http.StatusNotFound
		message = "record not found"
	case errors.ErrDatabase:
		code = err.Error()
		httpStatus = http.StatusInternalServerError
		message = "database error"
	case errors.ErrConflict:
		code = err.Error()
		httpStatus = http.StatusConflict
		message = "Record Exists"
	case errors.ErrForbidden:
		code = err.Error()
		httpStatus = http.StatusForbidden
		message = "Access denied"
	default:
		code = errors.SystemError
		httpStatus = http.StatusInternalServerError
		message = "internal server error"
	}

	if len(customMessage) > 0 && customMessage[0] != "" {
		message = customMessage[0]
	}

	if args == nil || args == "" {
		args = err.Error()
	}

	return httpStatus, renderErrorMsg(code, message, args)
}
