package apperrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type AppError struct {
	Message    string `json:"message"`
	HTTPStatus int    `json:"status"`
	Code       string `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(httpStatus int, code, message string) *AppError {
	return &AppError{
		Message:    message,
		HTTPStatus: httpStatus,
		Code:       code,
	}
}

func Newf(httpStatus int, code, format string, args ...interface{}) *AppError {
	return &AppError{
		Message:    fmt.Sprintf(format, args...),
		HTTPStatus: httpStatus,
		Code:       code,
	}
}

func NewNotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func NewBadRequest(message string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func NewUnauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func NewInternalServerError(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
		return NewBadRequest(err.Error())
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return NewBadRequest("invalid or empty request body")
	}

	return NewInternalServerError(err.Error())
}
