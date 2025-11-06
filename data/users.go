package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, Age: %d, IsActive: %t, CreatedAt: %s}", u.ID, u.Name, u.Email, u.Age, u.IsActive, u.CreatedAt.Format(time.RFC3339))
}

func LoadUsers(filePath string) ([]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Проверка на пустой файл (ПЕРЕД чтением)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return []User{}, nil // возвращаем пустой массив
	}

	users := []User{}
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func SaveUsers(filePath string, users []User) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // красивое форматирование
	return encoder.Encode(users)
}

func AddUser(users []User, newUser User) []User {
	return append(users, newUser)
}

func DeleteUser(users []User, id int) []User {
	for i, user := range users {
		if user.ID == id {
			return append(users[:i], users[i+1:]...)
		}
	}
	return users
}

func UpdateUser(users []User, id int, updatedUser User) []User {
	updatedUser.ID = id // убедись что ID правильный
	for i, user := range users {
		if user.ID == id {
			users[i] = updatedUser
			break
		}
	}
	return users
}

func FindUserByID(users []User, id int) (*User, error) {
	for i := range users {
		if users[i].ID == id {
			return &users[i], nil // возвращаем указатель на элемент массива
		}
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}
