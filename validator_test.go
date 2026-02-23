package validator

import "testing"

// helper — stops the test and prints a message if condition is false
func assert(t *testing.T, condition bool, msg string) {
    t.Helper()
    if !condition {
        t.Fatal(msg)
    }
}

// --- Required ---

func TestRequired_EmptyString_Fails(t *testing.T) {
    err := Required().Validate("")
    assert(t, err != nil, "expected error for empty string")
}

func TestRequired_NonEmptyString_Passes(t *testing.T) {
    err := Required().Validate("hello")
    assert(t, err == nil, "expected no error for non-empty string")
}

// --- MinLength ---

func TestMinLength_TooShort_Fails(t *testing.T) {
    err := MinLength(5).Validate("hi")
    assert(t, err != nil, "expected error for string shorter than min")
}

func TestMinLength_ExactLength_Passes(t *testing.T) {
    err := MinLength(5).Validate("hello")
    assert(t, err == nil, "expected no error for string at exact min length")
}

func TestMinLength_LongerThanMin_Passes(t *testing.T) {
    err := MinLength(3).Validate("hello")
    assert(t, err == nil, "expected no error for string longer than min")
}

// --- MaxLength ---

func TestMaxLength_TooLong_Fails(t *testing.T) {
    err := MaxLength(3).Validate("hello")
    assert(t, err != nil, "expected error for string longer than max")
}

func TestMaxLength_ExactLength_Passes(t *testing.T) {
    err := MaxLength(5).Validate("hello")
    assert(t, err == nil, "expected no error for string at exact max length")
}

// --- Matches ---

func TestMatches_ValidEmail_Passes(t *testing.T) {
    err := Matches(`^[^@]+@[^@]+\.[^@]+$`).Validate("user@example.com")
    assert(t, err == nil, "expected no error for valid email")
}

func TestMatches_InvalidEmail_Fails(t *testing.T) {
    err := Matches(`^[^@]+@[^@]+\.[^@]+$`).Validate("notanemail")
    assert(t, err != nil, "expected error for invalid email")
}

// --- Min / Max ---

func TestMin_BelowMin_Fails(t *testing.T) {
    err := Min(10).Validate(5)
    assert(t, err != nil, "expected error for value below min")
}

func TestMin_AboveMin_Passes(t *testing.T) {
    err := Min(10).Validate(15)
    assert(t, err == nil, "expected no error for value above min")
}

func TestMax_AboveMax_Fails(t *testing.T) {
    err := Max(10).Validate(15)
    assert(t, err != nil, "expected error for value above max")
}

func TestMax_BelowMax_Passes(t *testing.T) {
    err := Max(10).Validate(5)
    assert(t, err == nil, "expected no error for value below max")
}

// --- And ---

func TestAnd_AllPass(t *testing.T) {
    v := And(Required(), MinLength(3), MaxLength(10))
    err := v.Validate("hello")
    assert(t, err == nil, "expected no error when all validators pass")
}

func TestAnd_CollectsAllErrors(t *testing.T) {
    v := And(MinLength(10), MaxLength(3)) // both will fail on "hello"
    err := v.Validate("hello")
    assert(t, err != nil, "expected errors when validators fail")

    ve, ok := err.(ValidationError)
    assert(t, ok, "expected error to be a ValidationError")
    assert(t, len(ve.FieldErrors) == 2, "expected 2 field errors, got a different number")
}

type TestUser struct {
    Name string
    Age  int
}

func TestValidate_AllPass(t *testing.T) {
    user := TestUser{Name: "Alice", Age: 25}
    err := Validate(user,
        Field("Name", func(u TestUser) string { return u.Name }, And(Required(), MinLength(2))),
        Field("Age", func(u TestUser) int { return u.Age }, And(Min(18), Max(120))),
    )
    assert(t, err == nil, "expected no error when all fields are valid")
}

func TestValidate_CollectsAllFieldErrors(t *testing.T) {
    user := TestUser{Name: "", Age: 15}
    err := Validate(user,
        Field("Name", func(u TestUser) string { return u.Name }, Required()),
        Field("Age", func(u TestUser) int { return u.Age }, Min(18)),
    )
    assert(t, err != nil, "expected errors for invalid fields")

    ve := err.(ValidationError)
    assert(t, len(ve.FieldErrors) == 2, "expected 2 field errors")
    assert(t, ve.FieldErrors[0].FieldName == "Name", "expected first error to be on Name field")
    assert(t, ve.FieldErrors[1].FieldName == "Age", "expected second error to be on Age field")
}

func TestValidate_FieldNameStampedOnError(t *testing.T) {
    user := TestUser{Name: "", Age: 25}
    err := Validate(user,
        Field("Name", func(u TestUser) string { return u.Name }, Required()),
    )
    ve := err.(ValidationError)
    assert(t, ve.FieldErrors[0].FieldName == "Name", "expected field name to be stamped on error")
}