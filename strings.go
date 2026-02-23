package validator

import (
    "fmt"
    "regexp"
)

func Required() Validator[string] {
    return func(value string) error {
        if value == "" {
            return ValidationError{FieldErrors: []FieldError{
                {RuleName: "Required", Message: "value is required"},
            }}
        }
        return nil
    }
}

func MinLength(min int) Validator[string] {
    return func(value string) error {
        if len(value) < min {
            return ValidationError{FieldErrors: []FieldError{
                {RuleName: "MinLength", Message: fmt.Sprintf("must be at least %d characters", min)},
            }}
        }
        return nil
    }
}

func MaxLength(max int) Validator[string] {
    return func(value string) error {
        if len(value) > max {
            return ValidationError{FieldErrors: []FieldError{
                {RuleName: "MaxLength", Message: fmt.Sprintf("must be at most %d characters", max)},
            }}
        }
        return nil
    }
}

func Matches(pattern string) Validator[string] {
    re := regexp.MustCompile(pattern) // compile once, not on every call
    return func(value string) error {
        if !re.MatchString(value) {
            return ValidationError{FieldErrors: []FieldError{
                {RuleName: "Matches", Message: "value does not match required pattern"},
            }}
        }
        return nil
    }
}