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
	}, nil
}

// Services struct that encompasses all our services
type Services struct {
	Gallery GalleryService
	User    UserService
}
