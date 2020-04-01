package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ImageService interface describes methods present on this service
type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

// Image is not stored in the database
type Image struct {
	GalleryID uint
	Filename  string
}

// Path returns an image path as a string
func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

// RelativePath returns a path to delete an image
func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

// NewImageService function ?
func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

// Create initiates image upload
func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	// Create a destination file
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

// ByGalleryID fetches images linked to a gallery
func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	imgStrings, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	ret := make([]Image, len(imgStrings))
	for i := range imgStrings {
		imgStrings[i] = strings.Replace(imgStrings[i], path, "", 1)
		ret[i] = Image{
			Filename:  imgStrings[i],
			GalleryID: galleryID,
		}

	}
	return ret, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

// Delete removes an image from the filesystem
func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}
