package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"main-service/internal/apperror"
	"net/http"
	"strings"
)

const (
	ErrorMessageFieldRequired            = "field %s is required"
	ErrorMessageLessMinCharacters        = "minimum number of characters for the %s field - %v"
	ErrorMessageLessOrEqualThanNeeded    = "%s field must be greater than or equal to %v characters"
	ErrorMessageGreaterOrEqualThanNeeded = "%s field must be less than or equal to %v characters"
	ErrorMessageValueInFieldInvalid      = "value in field %s invalid"
	ErrorMessageInternalServerError      = "Internal server error"
	ErrorMessageInvalidRequestBody       = "invalid request body"
	ErrorMessageInvalidHeaderAuth        = "Must provide Authorization header with format 'Bearer {token}'"
)

func sendError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		err = apperror.BadRequest.New(getValidationErrorMessage(validationErrors))
	}
	if errors.Is(err, io.EOF) {
		err = apperror.BadRequest.New(ErrorMessageInvalidRequestBody)
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		err = apperror.BadRequest.New(ErrorMessageInvalidRequestBody)
	}
	var appError *apperror.AppError
	if errors.As(err, &appError) {
		switch appError.Type {
		case apperror.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, appError)
		case apperror.Authorization:
			c.AbortWithStatusJSON(http.StatusUnauthorized, appError)
		case apperror.Conflict:
			c.AbortWithStatusJSON(http.StatusConflict, appError)
		case apperror.PaymentRequired:
			c.AbortWithStatusJSON(http.StatusPaymentRequired, appError)
		case apperror.BadRequest:
			c.AbortWithStatusJSON(http.StatusBadRequest, appError)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, appError)
		}
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.Internal.New(err.Error()))
}

func getValidationErrorMessage(errors validator.ValidationErrors) string {
	var sb strings.Builder

	for _, f := range errors {
		tag := f.Tag()
		switch tag {
		case "required":
			sb.WriteString(fmt.Sprintf(ErrorMessageFieldRequired, f.Field()))
		case "min":
			sb.WriteString(fmt.Sprintf(ErrorMessageLessMinCharacters, f.Field(), f.Param()))
		case "gte":
			sb.WriteString(fmt.Sprintf(ErrorMessageLessOrEqualThanNeeded, f.Field(), f.Param()))
		case "lte":
			sb.WriteString(fmt.Sprintf(ErrorMessageGreaterOrEqualThanNeeded, f.Field(), f.Param()))
		default:
			sb.WriteString(fmt.Sprintf(ErrorMessageValueInFieldInvalid, f.Field()))
		}
		sb.WriteString("; ")
	}
	result := sb.String()
	return result[:len(result)-2]
}
