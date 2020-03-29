package models

import "github.com/jinzhu/gorm"

// NewServices func is responsible for making a connection to the database
func NewServices(dbDriver, connectionInfo string) (*Services, error) {
	db, err := gorm.Open(dbDriver, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User: NewUserService(db),
		db:   db,
	}, nil
}

// Services struct that encompasses all our services
type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops the tables and rebuilds it
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
