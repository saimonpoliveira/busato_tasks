package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errors[field] = formatValidationError(e)
	}
	return errors
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("maximum length is %s", e.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	default:
		return fmt.Sprintf("failed validation: %s", e.Tag())
	}
}
