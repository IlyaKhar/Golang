// internal/domain/user.go
package domain

import (
  "errors"
  "time"
)

type User struct {
  ID int
  Name string
  Email string
  Age int
  IsActive bool
  CreatedAt time.Time
}

// Ошибки домена (для маппинга на HTTP коды)
var (
  ErrUserNotFound = errors.New("user not found")
  ErrEmailExists  = errors.New("email already exists")
  ErrInvalidInput = errors.New("invalid input")
)