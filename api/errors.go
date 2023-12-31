package api

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

// ApiError holds possible http api error fields
type ApiError struct {
	Err error `json:"-"`

	StatusCode int    `json:"-"`
	StatusText string `json:"status"`

	Location  string      `json:"location,omitempty"`
	AppCode   int64       `json:"code,omitempty"`
	ErrorText string      `json:"error,omitempty"`
	Cause     string      `json:"cause,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// Error return an error text
func (e *ApiError) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("%s: %s", e.ErrorText, e.Cause)
	}
	return e.ErrorText
}

func (a *ApiError) Wrap(cause error) *ApiError {
	ret := *a
	ret.Err = fmt.Errorf("%s: caused by '%w'", ret.ErrorText, cause)
	ret.ErrorText = ret.Err.Error()
	return &ret
}

// Render sends error message to the client
func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
	pc := make([]uintptr, 5) // maximum 5 levels to go
	runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc)
	next := false
	for {
		frame, more := frames.Next()
		if next {
			e.Location = fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if strings.Contains(frame.File, "api/renderer.go") {
			next = true
		}
		if !more {
			break
		}
	}
	return nil
}

// ErrUnauthorized is error message for Unauthorized
func ErrUnauthorized(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		StatusText: "Unauthorized",
		ErrorText:  err.Error(),
	}
}

// ErrPermission is error message for Unauthorized
func ErrPermission(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		StatusText: "Permission denied.",
		ErrorText:  err.Error(),
	}
}

// ErrInvalidRequest is error message for Unauthorized
func ErrInvalidRequest(err error, data ...interface{}) *ApiError {
	v := &ApiError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request.",
		ErrorText:  err.Error(),
	}
	if len(data) > 0 {
		if errText, ok := data[0].(string); ok {
			v.ErrorText = fmt.Sprintf("%s: %s", v.ErrorText, errText)
		}
	}
	return v
}

// ErrServiceUnavailable is error message for Service Unavailable
func ErrServiceUnavailable(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusServiceUnavailable,
		StatusText: "Service Unavailable.",
		ErrorText:  err.Error(),
	}
}

// ErrInternalServerError is error message for Internal Server.
func ErrInternalServerError(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		StatusText: "Internal Server.",
		ErrorText:  err.Error(),
	}
}

// ErrRequestEntityTooLarge is error message for Request Entity Too Large
func ErrRequestEntityTooLarge(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusRequestEntityTooLarge,
		StatusText: "Request Entity Too Large",
		ErrorText:  err.Error(),
	}
}

// ErrResourceNotFound is error message for requested resource not found
func ErrResourceNotFound(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusNotFound,
		StatusText: "Not Found",
		ErrorText:  err.Error(),
	}
}

// ErrInvalidEmailSignup shapes error message if the users email address is invalid
func ErrInvalidEmailSignup(cause error) *ApiError {
	return &ApiError{
		Err:        cause,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid email address",
		ErrorText:  fmt.Sprintf("Invalid email: %s", cause.Error()),
		Data:       map[string]string{"field": "email"},
	}
}

// IgnoreError ignores error
func IgnoreError(v ...interface{}) {}

func CatchError(v ...interface{}) {
	for _, vv := range v {
		if e, ok := vv.(error); ok {
			log.Debug().Err(e).Msg("app error caught")
		}
	}
}
