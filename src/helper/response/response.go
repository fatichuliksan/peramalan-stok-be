package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type responseHelper struct {
}

type responseFormat struct {
	C       echo.Context
	Code    int
	Status  string
	Message string
	Data    interface{}
}

// Interface ...
type Interface interface {
	SetResponse(c echo.Context, code int, status string, message string, data interface{}) responseFormat
	SendResponse(res responseFormat) error
	EmptyJSONMap() map[string]interface{}
	SendSuccess(c echo.Context, message string, data interface{}) error
	SendBadRequest(c echo.Context, message string, data interface{}) error
	SendError(c echo.Context, message string, data interface{}) error
	SendUnauthorized(c echo.Context, message string, data interface{}) error
	SendValidationError(c echo.Context, validationErrors validator.ValidationErrors) error
	SendNotFound(c echo.Context, message string, data interface{}) error
	SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error
	SendResponsByCode(c echo.Context, code int, message string, data interface{}, err error) error
}

// NewResponse ...
func NewResponse() Interface {
	return &responseHelper{}
}

// SetResponse ...
func (r *responseHelper) SetResponse(c echo.Context, code int, status string, message string, data interface{}) responseFormat {
	return responseFormat{c, code, status, message, data}
}

// SendResponse ...
func (r *responseHelper) SendResponse(res responseFormat) error {
	if len(res.Message) == 0 {
		res.Message = http.StatusText(res.Code)
	}

	if res.Data != nil {
		return res.C.JSON(res.Code, map[string]interface{}{
			"code":    res.Code,
			"status":  res.Status,
			"message": res.Message,
			"data":    res.Data,
		})
	} else {
		return res.C.JSON(res.Code, map[string]interface{}{
			"code":    res.Code,
			"status":  res.Status,
			"message": res.Message,
		})
	}
}

func (r *responseHelper) SendResponsByCode(c echo.Context, code int, message string, data interface{}, err error) error {
	if err != nil {
		message = err.Error()
	}

	res := r.SetResponse(c, code, http.StatusText(code), message, data)
	return r.SendResponse(res)
}

// EmptyJSONMap : set empty data.
func (r *responseHelper) EmptyJSONMap() map[string]interface{} {
	return make(map[string]interface{})
}

// SendSuccess : Send success response to consumers.
func (r *responseHelper) SendSuccess(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusOK, http.StatusText(http.StatusOK), message, data)
	return r.SendResponse(res)
}

// SendBadRequest : Send bad request response to consumers.
func (r *responseHelper) SendBadRequest(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), message, data)
	return r.SendResponse(res)
}

// SendError : Send error request response to consumers.
func (r *responseHelper) SendError(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), message, data)
	return r.SendResponse(res)
}

// SendUnauthorized : Send error request response to consumers.
func (r *responseHelper) SendUnauthorized(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), message, data)
	return r.SendResponse(res)
}

// SendValidationError : Send validation error request response to consumers.
func (r *responseHelper) SendValidationError(c echo.Context, validationErrors validator.ValidationErrors) error {
	errorResponse := []string{}
	for _, err := range validationErrors {
		errorResponse = append(errorResponse, strings.Trim(fmt.Sprint(err), "[]")+".")
	}
	res := r.SetResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), strings.Trim(fmt.Sprint(errorResponse), "[]"), r.EmptyJSONMap())
	return r.SendResponse(res)
}

// SendNotFound : Send error request response to consumers.
func (r *responseHelper) SendNotFound(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusNotFound, http.StatusText(http.StatusNotFound), message, data)
	return r.SendResponse(res)
}

// SendCustomResponse ...
func (r *responseHelper) SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error {
	res := r.SetResponse(c, httpCode, http.StatusText(httpCode), message, data)
	return r.SendResponse(res)
}
