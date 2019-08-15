package models

import (
	"github.com/disintegration/imaging"
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const imageIDLength = 10

type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time
	Description string
}

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offer int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

func NewImage(user *User) *Image {
	return &Image{
		ID:        GenerateID("img", imageIDLength),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
}

var mimeExtensions = map[string]string{
	"image/png": ".png",
	"image/jpeg": ".jpg",
	"image/gif": ".gif",
}

var widthThumbnail = 400
var widthPreview = 800

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (image *Image) CreateFromURL(imageURL string) error {
	response, err := http.Get(imageURL)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return ErrImageURLInvalid
	}

	defer response.Body.Close()

	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))

	if err != nil {
		return ErrInvalidImageType
	}

	ext, valid := mimeExtensions[mimeType]

	if !valid {
		return ErrInvalidImageType
	}

	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext

	savedFile, err := os.Create("./data/images/" + image.Location)

	if err != nil {
		return err
	}

	defer savedFile.Close()

	size, err := io.Copy(savedFile, response.Body)

	if err != nil {
		return err
	}

	image.Size = size

	err = image.CreateResizedImages()

	if err != nil {
		return err
	}

	return GlobalImageStore.Save(image)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	savedFile, err := os.Create("./data/images/" + image.Location)

	if err != nil {
		return err
	}

	defer savedFile.Close()

	size, err := io.Copy(savedFile, file)

	if err != nil {
		return err
	}

	image.Size = size

	err = image.CreateResizedImages()

	if err != nil {
		return err
	}

	return GlobalImageStore.Save(image)
}

func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}

func (image *Image) CreateResizedImages() error {
	srcImage, err := imaging.Open("./data/images/" + image.Location)

	if err != nil {
		return err
	}

	errorChan := make(chan error)

	go image.resizePreview(errorChan, srcImage)
	go image.resizeThumbnail(errorChan, srcImage)

	var err1 error

	for i := 0; i < 2; i++ {
		err1 := <-errorChan

		if err1 == nil {
			err = err1
		}
	}

	return err1
}

func (image *Image) resizeThumbnail(errorChan chan error, srcImage image.Image) {
	dstImage := imaging.Thumbnail(srcImage, widthThumbnail, widthThumbnail, imaging.Lanczos)

	destination := "./data/images/thumbnail/" + image.Location
	errorChan <- imaging.Save(dstImage, destination)
}

func (image *Image) resizePreview(errorChan chan error, srcImage image.Image) {

	size := srcImage.Bounds().Size()

	ratio := float64(size.Y) / float64(size.X)

	targetHeigth := int(float64(widthPreview) * ratio)

	dstImage := imaging.Resize(srcImage, widthPreview, targetHeigth, imaging.Lanczos)

	destination := "./data/images/preview/" + image.Location

	errorChan <- imaging.Save(dstImage, destination)
}

func (image *Image) StaticThumbnailRoute() string  {
	return "/im/thumbnail/" + image.Location
}

func (image *Image) StaticPreviewRoute() string  {
	return "/im/preview/" + image.Location
}
