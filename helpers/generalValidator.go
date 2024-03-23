package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GeneralValidator(payloadValidationError error) []string {
	errors := make([]string, 0)

	switch err := payloadValidationError.(type) {
		case *json.UnmarshalTypeError:
			errors = append(errors, fmt.Sprintf("Invalid type for field %s: expected %s", err.Field, err.Type.String()))
		case *json.InvalidUnmarshalError:
			errors = append(errors, "Invalid JSON payload")
		case *json.SyntaxError:
			errors = append(errors, fmt.Sprintf("Invalid JSON syntax at byte offset %d", err.Offset))
		case validator.ValidationErrors:
			for _, fieldErr := range err {
				var message string
				switch fieldErr.Tag() {
					case "required":
						message = fmt.Sprintf("%s is required", fieldErr.Field())
					case "min":
						message = fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
					case "max":
						message = fmt.Sprintf("%s must be at most %s characters long", fieldErr.Field(), fieldErr.Param())
					case "email":
						message = fmt.Sprintf("%s must be a valid email address", fieldErr.Field())
					case "e164":
						message = fmt.Sprintf("%s must be a valid phone number", fieldErr.Field())
				}
		
				errors = append(errors, message)
			}
		default:
			errors = append(errors, fmt.Sprintf("Validation error: %s", payloadValidationError.Error()))
	}
	
	return errors
}