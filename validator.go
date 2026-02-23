package validator

type Validator[T any] func(T) error

func (v Validator[T]) Validate(value T) error {
	return v(value)
}

func And[T any](validators ...Validator[T]) Validator[T] {
	return func(value T) error {
		var validationError ValidationError
		for _, validator := range validators {
			if err := validator.Validate(value); err != nil {
				if ve, ok := err.(ValidationError); ok {
					validationError.FieldErrors = append(validationError.FieldErrors, ve.FieldErrors...)
				} else {
					validationError.AddFieldError("", "", err.Error())
				}
			}
		}
		if !validationError.IsEmpty() {
			return validationError
		}
		return nil
	}
}

func Not[T any](validator Validator[T]) Validator[T] {
	return func(value T) error {
		if err := validator.Validate(value); err == nil {
			return ValidationError{
				FieldErrors: []FieldError{
					{
						FieldName: "",
						RuleName:  "Not",
						Message:   "validation should fail but passed",
					},
				},
			}
		}
		return nil
	}
}