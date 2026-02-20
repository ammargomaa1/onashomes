package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IResource interface {
	GetStatusCode() int
	GetMessage() string
	GetData() interface{}
	GetErrors() interface{}
	GetMeta() interface{}
}

type baseResource struct {
	statusCode int
	message    string
	data       interface{}
	errors     interface{}
	meta       interface{}
}

func (r *baseResource) GetStatusCode() int     { return r.statusCode }
func (r *baseResource) GetMessage() string     { return r.message }
func (r *baseResource) GetData() interface{}   { return r.data }
func (r *baseResource) GetErrors() interface{} { return r.errors }
func (r *baseResource) GetMeta() interface{}   { return r.meta }

func NewOKResource(message string, data interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusOK,
		message:    message,
		data:       data,
	}
}

func NewCreatedResource(message string, data interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusCreated,
		message:    message,
		data:       data,
	}
}

func NewNoContentResource() IResource {
	return &baseResource{
		statusCode: http.StatusNoContent,
	}
}

func NewPaginatedOKResource(message string, data interface{}, meta interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusOK,
		message:    message,
		data:       data,
		meta:       meta,
	}
}

func NewBadRequestResource(message string, errors interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusBadRequest,
		message:    message,
		errors:     errors,
	}
}

func NewBadRequestWithBodyResource(message string, data interface{}, errors interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusBadRequest,
		message:    message,
		data:       data,
		errors:     errors,
	}
}

func NewUnauthorizedResource(message string, errors interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusUnauthorized,
		message:    message,
		errors:     errors,
	}
}

func NewForbiddenResource(message string, errors interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusForbidden,
		message:    message,
		errors:     errors,
	}
}

func NewNotFoundResource(message string, errors interface{}) IResource {
	return &baseResource{
		statusCode: http.StatusNotFound,
		message:    message,
		errors:     errors,
	}
}

func NewInternalErrorResource(message string, errors error) IResource {
	fmt.Printf("Internal server error: %s, details: %v", message, errors)
	return &baseResource{
		statusCode: http.StatusInternalServerError,
		message:    message,
	}
}

func WriteResource(c *gin.Context, res IResource) {
	statusCode := res.GetStatusCode()

	if statusCode == http.StatusNoContent {
		c.Status(http.StatusNoContent)
		return
	}

	if statusCode >= 200 && statusCode < 300 {
		if res.GetMeta() != nil {
			SuccessResponseWithMeta(c, statusCode, res.GetMessage(), res.GetData(), res.GetMeta())
		} else {
			SuccessResponse(c, statusCode, res.GetMessage(), res.GetData())
		}
		return
	}

	if statusCode == http.StatusBadRequest && res.GetMessage() == "Validation failed" {
		ValidationErrorResponse(c, res.GetErrors())
		return
	}

	ErrorResponse(c, statusCode, res.GetMessage(), res.GetErrors())
}
