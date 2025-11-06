package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"errors"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Age       int       `gorm:"not null" json:"age"`
	IsActive  bool      `gorm:"not null" json:"is_active"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// Подключение (SQLite)
func OpenGormSQLite(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

// (Опционально) Подключение Postgres
//  func OpenGormPostgres(dsn string) (*gorm.DB, error) {
// 	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
// }

// Миграция
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

// CRUD — оставляю заготовки, ты допиши/улучши по месту

func CreateUser(db *gorm.DB, u *User) error {
	// TODO: добавить валидацию (имя/email/возраст), дефолт для CreatedAt, затем db.Create(u).Error
	if u.Name == "" || u.Email == "" || u.Age <= 0 {
		return errors.New("invalid user data")
	}
	return db.Create(u).Error
}

func GetUser(db *gorm.DB, id int) (User, error) {
	var u User
	err := db.First(&u, id).Error // TODO: обработать gorm.ErrRecordNotFound → вернуть понятную ошибку
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errors.New("user not found")
	}	
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func ListUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Order("id").Find(&users).Error // TODO: пагинация (Limit/Offset), сортировка
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []User{}, errors.New("users not found")
	}
    if err != nil { return nil, err }
    return users, nil
}

func UpdateUser(db *gorm.DB, u *User) error {
	// TODO: проверить, что запись существует; затем db.Save(u).Error или Updates(map[string]any{...})
    if u.ID == 0 { return errors.New("user id required") }
    return db.Save(u).Error
}

func DeleteUser(db *gorm.DB, id int) error {
	// TODO: м.б. проверить зависимые сущности или «мягкое удаление»
	tx := db.Delete(&User{}, id)
    if tx.Error != nil { return tx.Error }
    if tx.RowsAffected == 0 { return errors.New("user not found") }
    return nil
}