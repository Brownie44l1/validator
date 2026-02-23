package validator

type FieldError struct {
	FieldName string
	RuleName  string
	Message  string
}

type ValidationError struct {
	FieldErrors []FieldError
}

func (ve ValidationError) Error() string {
	return "validation failed"
}

func (ve ValidationError) IsEmpty() bool {
	return len(ve.FieldErrors) == 0
}

func (ve *ValidationError) AddFieldError(fieldName, ruleName, message string) {
	ve.FieldErrors = append(ve.FieldErrors, FieldError{
		FieldName: fieldName,
		RuleName:  ruleName,
		Message:  message,
	})
}

