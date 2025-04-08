package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// GetValidationErrorMessage returns a user-friendly error message for validation errors
func GetValidationErrorMessage(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid email address", e.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param()))
			case "oneof":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", e.Field()))
			}
		}
		return strings.Join(errorMessages, "; ")
	}
	return err.Error()
}
