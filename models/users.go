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

	// ErrInvalidID is returned when an Invalid ID is provided to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")
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
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
}

// ByEmail returns a user object based on email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user data
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update will update the provided user with all of the data in the provided user object
func (us *UserService) Update(user *User) {
	return us.db.Save(user)
}

// Delete a database user
func (us *UserService) Delete(id uint) error {
	// we are preventing data with id of 0 bcos gorm will delete all records if allowed
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
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
