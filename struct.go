package validator

type FieldValidator[S any] struct {
	name string
	run  func(s S) []FieldError
}

func Field[S any, T any](name string, getter func(S) T, validator Validator[T]) FieldValidator[S] {
	return FieldValidator[S]{
		name: name,
		run: func(s S) []FieldError {
			err := validator.Validate(getter(s))
			if err == nil {
				return nil
			}
			ve, ok := err.(ValidationError)
			if !ok {
				return []FieldError{{FieldName: name, RuleName: "", Message: err.Error()}}
			}
			for i := range ve.FieldErrors {
				ve.FieldErrors[i].FieldName = name
			}
			return ve.FieldErrors
		},
	}
}

func Validate[S any](s S, fields ...FieldValidator[S]) error {
	var result ValidationError
	for _, field := range fields {
		result.FieldErrors = append(result.FieldErrors, field.run(s)...)
	}
	if result.IsEmpty() {
		return nil
	}
	return result
}

func ValidateNested[S any](fields ...FieldValidator[S]) Validator[S] {
	return func(s S) error {
		return Validate(s, fields...)
	}
}
