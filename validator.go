package validator

import "fmt"

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

func Each[T any](validator Validator[T]) Validator[[]T] {
	return func(values []T) error {
		var validationError ValidationError
		for i, v := range values {
			if err := validator.Validate(v); err != nil {
				if ve, ok := err.(ValidationError); ok {
					for j := range ve.FieldErrors {
						ve.FieldErrors[j].FieldName = fmt.Sprintf("[%d]", i)
					}
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
