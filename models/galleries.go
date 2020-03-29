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

type galleryValFunc func(*Gallery) error

func runGalleryValFuncs(gallery *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

// * validators
type galleryValidator struct {
	GalleryDB
}

// Create validator for gallery
func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

// userIDRequired makes sure a userid is available while creating a gallery
func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

// titleRequired makes sure a title is available while creating a gallery
func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

// Create func creates a new gallery in the database
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
