package errors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	msg := "test error"

	err := New(msg)

	assert.Equal(t, msg, err.err)
	assert.Equal(t, "InternalServerError", err.typ)
	assert.Equal(t, http.StatusInternalServerError, err.status)
}

func TestNewValidation(t *testing.T) {
	msg := "test error"

	err := NewValidation(msg)

	assert.Equal(t, msg, err.err)
	assert.Equal(t, "Validation", err.typ)
	assert.Equal(t, http.StatusBadRequest, err.status)
}

func TestError_SetType(t *testing.T) {
	msg := "test error"
	typ := "Test"

	err := New(msg).SetType(typ)

	assert.Equal(t, msg, err.err)
	assert.Equal(t, typ, err.typ)
	assert.Equal(t, http.StatusInternalServerError, err.status)
}

func TestError_SetStatus(t *testing.T) {
	msg := "test error"
	code := 400

	err := New(msg).SetStatus(code)

	assert.Equal(t, msg, err.err)
	assert.Equal(t, "InternalServerError", err.typ)
	assert.Equal(t, 400, err.status)
}

func TestError_SetParamName(t *testing.T) {
	paramName := "name"

	err := New("test error").SetParamName(paramName)

	assert.Equal(t, paramName, *err.paramName)
	assert.Equal(t, "InternalServerError", err.typ)
	assert.Equal(t, 500, err.status)
}

func TestError_Error(t *testing.T) {
	msg := "test error"
	err := &Error{err: msg}

	assert.Equal(t, msg, err.Error())
}

func TestError_Status(t *testing.T) {
	code := 405
	err := &Error{status: code}

	assert.Equal(t, code, err.Status())
}

func TestError_ParamName(t *testing.T) {
	name := "Name"
	err := &Error{paramName: &name}

	assert.Equal(t, &name, err.ParamName())
}

func TestError_MarshalJSON(t *testing.T) {
	msg := "test error"
	typ := "Test"
	code := 500

	e := New(msg).SetType(typ).SetStatus(code).SetParamName("name")

	bytes, err := json.Marshal(e)

	assert.NoError(t, err)
	assert.Equal(t, "{\"error\":\"Test\",\"message\":\"test error\",\"paramName\":\"name\"}", string(bytes))
}

func TestIsValidation_GivenValidationErrorType_ReturnsTrue(t *testing.T) {
	err := New("my error").SetType("Validation")

	assert.True(t, IsValidation(err))
}

func TestIsValidation_GivenErrorWithBadRequestCode_ReturnsTrue(t *testing.T) {
	err := New("my error").SetStatus(http.StatusBadRequest)

	assert.True(t, IsValidation(err))
}

func TestIsValidation_GivenNonValidationError_ReturnsFalse(t *testing.T) {
	err := New("my error").SetStatus(http.StatusInternalServerError)

	assert.False(t, IsValidation(err))
}

func TestIsValidation_GivenStandardError_ReturnsFalse(t *testing.T) {
	err := fmt.Errorf("my error")

	assert.False(t, IsValidation(err))
}

func TestWriteResponse_GivenError_WritesJSON(t *testing.T) {
	err := New("Oops").SetType("Error").SetParamName("Name").SetStatus(400)
	rr := httptest.NewRecorder()

	WriteResponse(rr, err)

	resp := rr.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	bytes, _ := ioutil.ReadAll(resp.Body)

	var data map[string]*string
	_ = json.Unmarshal(bytes, &data)

	assert.Equal(t, "Error", *data["error"])
	assert.Equal(t, "Oops", *data["message"])
	assert.Equal(t, "Name", *data["paramName"])
}
