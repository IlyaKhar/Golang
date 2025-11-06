package handlers

import (
    "day2/data"
    "testing"
    "strings"
)

func TestValidateUser(t *testing.T) {
    existing := []data.User{
        {ID:1, Email:"a@x.com"},
    }

    tests := []struct{
        name string
        u    data.User
        wantErr string // часть текста ошибки
    }{
        {"ok", data.User{Name:"Al", Email:"b@x.com", Age:20}, ""},
        {"no_name", data.User{Name:"", Email:"b@x.com", Age:20}, "name is required"},
        {"short_name", data.User{Name:"A", Email:"b@x.com", Age:20}, "at least 2"},
        {"no_email", data.User{Name:"Al", Email:"", Age:20}, "email is required"},
        {"bad_email", data.User{Name:"Al", Email:"bx.com", Age:20}, "must contain @"},
        {"bad_age", data.User{Name:"Al", Email:"b@x.com", Age:151}, "between 0 and 150"},
        {"dup_email", data.User{Name:"Al", Email:"a@x.com", Age:20}, "email already exists"},
        {"dup_id", data.User{ID:1, Name:"Al", Email:"b@x.com", Age:20}, "ID already exists"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUser(tt.u, existing)
            if tt.wantErr == "" {
                if err != nil { t.Fatalf("unexpected err: %v", err) }
            } else {
                if err == nil || !contains(err.Error(), tt.wantErr) {
                    t.Fatalf("want err containing %q, got %v", tt.wantErr, err)
                }
            }
        })
    }
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }