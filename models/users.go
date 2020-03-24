package models

import (
	"github.com/jinzhu/gorm"
	// we want to keep the postgres dialect even though we are not using it directly
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User struct
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
