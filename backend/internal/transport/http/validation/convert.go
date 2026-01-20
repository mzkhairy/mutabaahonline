package validation

import (
	"errors"
	"strings"

	"mutabaahapi/internal/platform/apperror"

	"github.com/go-playground/validator/v10"
)

// FromValidator converts go-playground/validator errors to a structured ValidationErrors.
func FromValidator(err error) (*apperror.ValidationErrors, bool) {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fields := make(map[string]string, len(verrs))
		for _, fe := range verrs {
			// Field path comes like "CreateNoteReq.Title" sometimes; keep the last segment lowercased.
			parts := strings.Split(fe.Namespace(), ".")
			field := fe.Field()
			if len(parts) > 0 {
				field = parts[len(parts)-1]
			}
			fields[strings.ToLower(field)] = friendlyMsg(fe)
		}
		return apperror.NewValidationErrors(fields), true
	}
	return nil, false
}

func friendlyMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "min=" + fe.Param()
	case "max":
		return "max=" + fe.Param()
	case "len":
		return "len=" + fe.Param()
	case "uuid", "uuid4":
		return "must be a valid uuid"
	default:
		return fe.Error()
	}
}
