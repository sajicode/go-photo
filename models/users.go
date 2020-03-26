package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sajicode/go-photo/hash"
	"github.com/sajicode/go-photo/rand"

	// we want to keep the postgres dialect even though we are not using it directly
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the DB
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an Invalid ID is provided to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned whenever a user passes a wrong password
	ErrInvalidPassword = errors.New("models: incorrect passsword provided")
)

const userPwPepper = "secret-random-string"
const hmacSecretKey = "secret-hmac-key"

// User struct represents the user model in our DB
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB interface talks directly to the DB
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// newUserGorm responsible for making a connection to the database
func newUserGorm(dbDriver, connectionInfo string) (*UserGorm, error) {
	db, err := gorm.Open(dbDriver, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	// Authenticate will verify the provided email address & password are correct
	// If correct, the corresonnding user is returned, otherwise errors
	Authenticate(email, password string) (*User, error)
	// We need all the methods in the UserDB type
	UserDB
}

// NewUserService DB connection
func NewUserService(dbDriver, connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(dbDriver, connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

var _ UserService = &userService{}

// UserService struct handles actions with the user model
type userService struct {
	UserDB
}

// userValidator responsible for validating data before it gets to the DB
type userValidator struct {
	UserDB
}

// If the userGorm interface stops matching UserDB, compilation error
// basically, we are ensuring UserGorm implements the UserDB type
var _ UserDB = &UserGorm{}

// UserGorm struct
type UserGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// ByID will look up a user by the ID provided
// 1 - user, nil, 2 - nil, ErrNotFound, 3 - nil, otherError
func (ug *UserGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail returns a user object based on email
func (ug *UserGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token and returns that user
// This method will handle hashing the token for us
func (ug *UserGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Authenticate function returns a user or an error when verifying a user
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

// Create will create the provided user data
func (ug *UserGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		// create a remember token
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

// Update will update the provided user with all of the data in the provided user object
func (ug *UserGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

// Delete a database user
func (ug *UserGorm) Delete(id uint) error {
	// we are preventing data with id of 0 bcos gorm will delete all records if allowed
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close DB connection
func (ug *UserGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset func to drop and recreate tables
func (ug *UserGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate creates tables in the db
func (ug *UserGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// first will query using the provided gorm.DB & it will get the first item returned & place it into dst
// If  othing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
