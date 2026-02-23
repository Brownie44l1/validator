package validator

import "fmt"

type Number interface {
    int | int8 | int16 | int32 | int64 |
        uint | uint8 | uint16 | uint32 | uint64 |
        float32 | float64
}

func Min[T Number](n T) Validator[T] {
    return func(value T) error {
        if value < n {
            return ValidationError{FieldErrors: []FieldError{
                {RuleName: "Min", Message: fmt.Sprintf("must be at least %v", n)},
            }}
        }
        return nil
    }
}

func Max[T Number](n T) Validator[T] {
    return func(value T) error {
        if value > n {
            return ValidationError{FieldErrors: []FieldError{
                {RuleName: "Max", Message: fmt.Sprintf("must be at most %v", n)},
            }}
        }
        return nil
    }
}