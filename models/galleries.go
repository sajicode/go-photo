package models

import "github.com/jinzhu/gorm"

// Gallery is our image container resources that visitors view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

// GalleryService interface communicates with the DB
type GalleryService interface {
	GalleryDB
}

// GalleryDB interface
type GalleryDB interface {
	Create(gallery *Gallery) error
}

// NewGalleryService tells the db to create a new gallery
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

// GalleryService implementation
type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

// Create func creates a new gallery in the database
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
