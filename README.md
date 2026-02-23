# validator

A type-safe, zero-reflection data validation library for Go, built with generics.

## Why

Most Go validation libraries use reflection. Reflection is slow and pushes errors to runtime that should be caught at compile time. This library uses generics instead — the compiler catches type mismatches, and you pay no runtime cost for it.

## Benchmarks

Benchmarked against [go-playground/validator](https://github.com/go-playground/validator) on an Intel Pentium N4200:

```
BenchmarkYours-4          720924    1645 ns/op      0 B/op    0 allocs/op
BenchmarkGoPlayground-4   229621    5131 ns/op    145 B/op    6 allocs/op
```

**3x faster. Zero allocations.**

## Install

```bash
go get github.com/Brownie44l1/validator
```

## Usage

### Validating a single value

```go
err := validator.And(
    validator.Required(),
    validator.MinLength(2),
    validator.MaxLength(50),
).Validate("Alice")
```

### Validating a struct

```go
type User struct {
    Name  string
    Email string
    Age   int
}

user := User{Name: "A", Email: "notanemail", Age: 15}

err := validator.Validate(user,
    validator.Field("Name", func(u User) string { return u.Name },
        validator.And(validator.Required(), validator.MinLength(2))),

    validator.Field("Email", func(u User) string { return u.Email },
        validator.And(validator.Required(), validator.Matches(`^[^@]+@[^@]+\.[^@]+$`))),

    validator.Field("Age", func(u User) int { return u.Age },
        validator.And(validator.Min(18), validator.Max(120))),
)

if err != nil {
    ve := err.(validator.ValidationError)
    for _, fe := range ve.FieldErrors {
        fmt.Printf("field: %s, rule: %s, message: %s\n", fe.FieldName, fe.RuleName, fe.Message)
    }
}
// field: Name,  rule: MinLength, message: must be at least 2 characters
// field: Email, rule: Matches,   message: value does not match required pattern
// field: Age,   rule: Min,       message: must be at least 18
```

### Nested structs

```go
type Address struct {
    Street string
    City   string
}

type Order struct {
    Name    string
    Address Address
}

err := validator.Validate(order,
    validator.Field("Name", func(o Order) string { return o.Name }, validator.Required()),
    validator.Field("Address", func(o Order) Address { return o.Address },
        validator.ValidateNested(
            validator.Field("Street", func(a Address) string { return a.Street }, validator.Required()),
            validator.Field("City", func(a Address) string { return a.City }, validator.Required()),
        )),
)
```

### Slice validation

```go
err := validator.Each(validator.MinLength(2)).Validate([]string{"go", "x", "hi"})
// validates every element, collects all errors with their index
// field: [1], rule: MinLength, message: must be at least 2 characters
```

### Composing validators

```go
// And — runs all validators, collects every error
validator.And(validator.Required(), validator.MinLength(2), validator.MaxLength(50))

// Not — passes if the validator fails
validator.Not(validator.Matches(`[0-9]`)) // must not contain digits
```

## Error structure

Every error is a `ValidationError` containing a slice of `FieldError`:

```go
type FieldError struct {
    FieldName string // which field failed
    RuleName  string // which rule failed
    Message   string // human readable message
}
```

This makes it trivial to serialize errors to JSON for APIs:

```go
if err != nil {
    ve := err.(validator.ValidationError)
    json.NewEncoder(w).Encode(ve.FieldErrors)
}
// [{"FieldName":"Email","RuleName":"Matches","Message":"value does not match required pattern"}]
```

## Built-in validators

### Strings
| Validator | Description |
|---|---|
| `Required()` | value must not be empty |
| `MinLength(n)` | minimum character length |
| `MaxLength(n)` | maximum character length |
| `Matches(pattern)` | must match regex pattern |

### Numbers
| Validator | Description |
|---|---|
| `Min(n)` | minimum value |
| `Max(n)` | maximum value |

Supports all Go number types: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`.

## Design decisions

**No reflection.** Fields are accessed via plain getter functions. The compiler enforces types — if you pass a `Validator[string]` to a `Field` that returns an `int`, it won't compile.

**All errors, not just the first.** `And` and `Validate` collect every failure before returning. Callers get a complete picture of what's wrong in a single call.

**Composable by design.** `Validator[T]` is just a function. `And`, `Not`, `Each`, and `ValidateNested` are all combinators that return validators — they can be nested and reused freely.

## License

MIT