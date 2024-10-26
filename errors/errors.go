package transport

import (
	"net/http"
)

type Error interface {
	Error() string
	GetCode() int
	GetMethod() string
}

type ErrMissingParams struct {
	Param  string
	Method string
}

type ErrUnknownParams struct {
	Param  string
	Method string
}

type ErrMethodNotAllowed struct {
	Method string
}

func (e *ErrMissingParams) Error() string {
	return "missing required parameter: " + e.Param
}

func (e *ErrUnknownParams) Error() string {
	return "provided unknown parameter: " + e.Param
}

func (e *ErrMethodNotAllowed) Error() string {
	return "method not allowed: " + e.Method
}

func (e *ErrMissingParams) GetCode() int {
	return http.StatusBadRequest
}

func (e *ErrUnknownParams) GetCode() int {
	return http.StatusBadRequest
}

func (e *ErrMethodNotAllowed) GetCode() int {
	return http.StatusMethodNotAllowed
}

func (e *ErrMissingParams) GetMethod() string {
	if e.Method != "" {
		return e.Method
	}
	return "GET"
}

func (e *ErrUnknownParams) GetMethod() string {
	if e.Method != "" {
		return e.Method
	}
	return "GET"
}

func (e *ErrMethodNotAllowed) GetMethod() string {
	if e.Method != "" {
		return e.Method
	}
	return "GET"
}
