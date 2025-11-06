package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Age       int
	IsActive  bool
	CreatedAt time.Time
}

// NewUser создаёт пользователя с базовой валидацией и текущей датой
func NewUser(name, email, password string, age int) (*User, error) {
	if name == "" {
		return nil, errors.New("имя не должно быть пустым")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("email должен содержать @")
	}
	if age < 0 || age > 150 {
		return nil, errors.New("некорректный возраст")
	}
	if password != "" && len(password) < 8 {
		return nil, errors.New("пароль должен содержать минимум 8 символов")
	}

	user := &User{
		Name:      name,
		Email:     email,
		Password:  password,
		Age:       age,
		IsActive:  false,
		CreatedAt: time.Now(),
	}
	return user, nil
}

// Introduce возвращает краткое представление пользователя
func (u *User) Introduce() string {
	return fmt.Sprintf("Привет, я %s, мне %d лет", u.Name, u.Age)
}

// IsAdult возвращает true, если пользователю 18 или больше
func (u *User) IsAdult() bool {
	return u.Age >= 18
}

// GetInfo возвращает подробную информацию о пользователе
func (u *User) GetInfo() string {
	return fmt.Sprintf(
		"ID: %d, Имя: %s, Email: %s, Возраст: %d, Активен: %t, Дата создания: %s",
		u.ID, u.Name, u.Email, u.Age, u.IsActive, u.CreatedAt.Format(time.RFC3339),
	)
}

// Activate помечает пользователя активным
func (u *User) Activate() {
	u.IsActive = true
}

// Deactivate помечает пользователя неактивным
func (u *User) Deactivate() {
	u.IsActive = false
}

// ValidateEmail проверяет корректность email
func (u *User) ValidateEmail() error {
	if !strings.Contains(u.Email, "@") {
		return errors.New("неверный email: отсутствует @")
	}
	return nil
}

// CreateUser создаёт пользователя без пароля (для простых сценариев)
func CreateUser(name, email string, age int) (*User, error) {
	return NewUser(name, email, "", age)
}
