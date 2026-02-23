package validator_test

import (
    "testing"
    "github.com/Brownie44l1/validator"
    govalidator "github.com/go-playground/validator/v10"
)

// --- your library ---

type BenchUser struct {
    Name  string
    Email string
    Age   int
}

func BenchmarkYours(b *testing.B) {
    user := BenchUser{Name: "Alice", Email: "alice@example.com", Age: 25}
    v := []validator.FieldValidator[BenchUser]{
        validator.Field("Name", func(u BenchUser) string { return u.Name },
            validator.And(validator.Required(), validator.MinLength(2))),
        validator.Field("Email", func(u BenchUser) string { return u.Email },
            validator.And(validator.Required(), validator.Matches(`^[^@]+@[^@]+\.[^@]+$`))),
        validator.Field("Age", func(u BenchUser) int { return u.Age },
            validator.And(validator.Min(18), validator.Max(120))),
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        validator.Validate(user, v...)
    }
}

// --- go-playground/validator ---

type BenchUserTag struct {
    Name  string `validate:"required,min=2"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=120"`
}

func BenchmarkGoPlayground(b *testing.B) {
    v := govalidator.New()
    user := BenchUserTag{Name: "Alice", Email: "alice@example.com", Age: 25}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        v.Struct(user)
    }
}