package errors

import (
	"encoding/json"
	"net/http"
)

// Error is a custom implementation of error, providing
// an error type and HTTP status code.
type Error struct {
	err       string
	typ       string
	status    int
	paramName *string
}

// New returns a new instance of Error, with the given error message.
func New(err string) *Error {
	return &Error{
		err:    err,
		typ:    "InternalServerError",
		status: http.StatusInternalServerError,
	}
}

// NewValidation returns a new instance of Error, with the typ "Validation".
func NewValidation(err string) *Error {
	return New(err).
		SetType("Validation").
		SetStatus(http.StatusBadRequest)
}

// SetType sets the type of the error - this is used for the HTTP response.
func (err *Error) SetType(typ string) *Error {
	err.typ = typ

	return err
}

// SetStatus sets the status code of the error - this is used
// for the HTTP response.
func (err *Error) SetStatus(statusCode int) *Error {
	err.status = statusCode

	return err
}

// SetParamName sets the error's paramName property.
func (err *Error) SetParamName(name string) *Error {
	err.paramName = &name

	return err
}

// Error returns the error message.
func (err *Error) Error() string {
	return err.err
}

// Status returns the error status code.
func (err *Error) Status() int {
	return err.status
}

// ParamName returns the error's paramName.
func (err *Error) ParamName() *string {
	return err.paramName
}

// MarshalJSON is a custom JSON marshal func used
// when writing the error to a http.ResponseWriter.
func (err *Error) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"error":   err.typ,
		"message": err.err,
	}

	if err.paramName != nil {
		data["paramName"] = err.paramName
	}

	return json.Marshal(data)
}

// IsValidation returns true if the error is a validation error.
func IsValidation(err error) bool {
	e, ok := err.(*Error)

	return ok && (e.status == http.StatusBadRequest ||
		e.typ == "Validation")
}

// WriteResponse writes the instance of error to the response. This
// results in all errors being returned in a standard structure.
func WriteResponse(w http.ResponseWriter, err error) {
	var e *Error

	switch t := err.(type) {
	case *Error:
		e = t
	default:
		e = New(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status())

	json.NewEncoder(w).Encode(e)
}
