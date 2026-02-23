package main

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/Brownie44l1/validator"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

var userValidator = []validator.FieldValidator[User]{
    validator.Field("Name", func(u User) string { return u.Name },
        validator.And(validator.Required(), validator.MinLength(2))),
    validator.Field("Email", func(u User) string { return u.Email },
        validator.And(validator.Required(), validator.Matches(`^[^@]+@[^@]+\.[^@]+$`))),
    validator.Field("Age", func(u User) int { return u.Age },
        validator.And(validator.Min(18), validator.Max(120))),
}

func main() {
    data, err := os.ReadFile("example/users.json")
    if err != nil {
        fmt.Println("failed to read file:", err)
        return
    }

    var users []User
    if err := json.Unmarshal(data, &users); err != nil {
        fmt.Println("failed to parse json:", err)
        return
    }

    fmt.Printf("validating %d users...\n\n", len(users))

    passed := 0
    failed := 0

    for i, user := range users {
        err := validator.Validate(user, userValidator...)
        if err != nil {
            failed++
            ve := err.(validator.ValidationError)
			fmt.Printf("user[%d] %q passed\n", i, user.Name)
            for _, fe := range ve.FieldErrors {
                fmt.Printf("  %-10s %-15s %s\n", fe.FieldName, fe.RuleName, fe.Message)
            }
            fmt.Println()
        } else {
            passed++
            fmt.Printf("user[%d] %q passed\n", user.Name)
        }
    }

    fmt.Printf("\nresults: %d passed, %d failed out of %d\n", passed, failed, len(users))
}