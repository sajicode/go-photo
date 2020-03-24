package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// we want to keep the postgres dialect even though we are not using it directly
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the DB
	ErrNotFound = errors.New("models: resource not found")
)

// NewUserService DB connection
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

// UserService struct handles actions with the user model
type UserService struct {
	db *gorm.DB
}

// ByID will look up a user by the ID provided
// 1 - user, nil, 2 - nil, ErrNotFound, 3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create the provided user data
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Close DB connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset func to drop and recreate tables
//! for dev purposes only
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// User struct
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
