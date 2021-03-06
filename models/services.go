package models

import (
	"os"

	"github.com/jinzhu/gorm"
	// we want to keep the postgres dialect even though we are not using it directly
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewServices func is responsible for making a connection to the database
func NewServices(dbDriver, connectionInfo string) (*Services, error) {
	db, err := gorm.Open(dbDriver, connectionInfo)
	if err != nil {
		return nil, err
	}
	var logDB bool
	if os.Getenv("APP_ENV") == "production" {
		logDB = false
	} else {
		logDB = true
	}
	db.LogMode(logDB)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
}

// Services struct that encompasses all our services
type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops the tables and rebuilds it
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}, &pwReset{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{}).Error
}
