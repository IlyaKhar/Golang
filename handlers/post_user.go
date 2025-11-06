package handlers

import (
	"day2/data"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ValidateUser проверяет данные пользователя
func ValidateUser(user data.User, existingUsers []data.User) error {
	log.Println("ValidateUser")
	// 1. Проверка имени
	if user.Name == "" {
		return errors.New("name is required")
	}
	if len(user.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}

	// 2. Проверка email
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("email must contain @")
	}

	// 3. Проверка возраста
	if user.Age < 0 || user.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	// 4. Проверка уникальности email
	for _, u := range existingUsers {
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	// 5. Проверка уникальности ID (если ID задан вручную)
	for _, u := range existingUsers {
		if u.ID == user.ID {
			return errors.New("user ID already exists")
		}
	}

	return nil // всё ок
}

// PostUserHandler обрабатывает POST /users
func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Читаем тело
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// 2. Парсим JSON
	var newUser data.User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 3. Загружаем существующих (для проверки уникальности)
	users, err := data.LoadUsers("data/users.json")
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	// 4. Устанавливаем CreatedAt если не задан
	if newUser.CreatedAt.IsZero() {
		newUser.CreatedAt = time.Now()
	}

	// 5. Валидируем
	err = ValidateUser(newUser, users)
	if err != nil {
		// Отправляем JSON с ошибкой
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 6. Добавляем и сохраняем
	users = data.AddUser(users, newUser)
	err = data.SaveUsers("data/users.json", users)
	if err != nil {
		http.Error(w, "Failed to save users", http.StatusInternalServerError)
		return
	}

	// 7. Отправляем ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
