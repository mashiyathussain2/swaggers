package app

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/avelino/slugify"
	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
	"github.com/vasupal1996/goerror"
)

// IMG represents image config, object and type
type IMG struct {
	Type string
	Conf *image.Config
	Img  *image.Image
}

// DecodeBase64StrToIMG converts base64 image string to IMG object
// b64Str -> "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQE...."
// image meta data -> data:image/jpeg;base64
// actual image data -> /9j/4AAQSkZJRgABAQE....
func (i *IMG) DecodeBase64StrToIMG(b64Str string) error {
	// finding the index where base64Image data starts
	coI := strings.Index(string(b64Str), ",")
	if coI == -1 {
		return goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil)
	}

	b64ImgData := string(b64Str)[coI+1:]
	if b64ImgData == "" {
		return goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil)
	}
	imgType := strings.TrimSuffix(b64Str[5:coI], ";base64")
	// getting the image type

	// getting image meta data by reading image
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64ImgData))
	imgConf, _, err := image.DecodeConfig(r)
	if err != nil {
		return goerror.New(fmt.Sprintf("failed to decode image config: %s", err), nil)
	}

	// un-basing image and getting image bytes reader
	unBaseImg, err := base64.StdEncoding.DecodeString(string(b64ImgData))
	if err != nil {
		return goerror.New(fmt.Sprintf("failed to un-base image string: %s", err), nil)
	}
	imgByteReader := bytes.NewReader(unBaseImg)

	// decoding image
	var img image.Image
	switch imgType {
	case "image/png":
		img, err = png.Decode(imgByteReader)
		if err != nil {
			return goerror.New(fmt.Sprintf("failed to decode png image: %s", err), nil)
		}
	case "image/jpeg":
		img, err = jpeg.Decode(imgByteReader)
		if err != nil {
			return goerror.New(fmt.Sprintf("failed to decode jpeg/jpg image: %s", err), nil)
		}
	default:
		return goerror.New("invalid image type [only jpeg and png are allowed]", nil)
	}

	i.Type = imgType
	i.Conf = &imgConf
	i.Img = &img

	return nil
}

// Resize resizes image to passed width and height and returns the resized image
// If one of the parameters width or height is set to 0, its size will be calculated so that the aspect ratio is that of the originating image.
func (i *IMG) Resize(width, height uint) *image.Image {
	m := resize.Resize(width, height, *i.Img, resize.NearestNeighbor)
	return &m
}

// UniqueSlug converts a string into unique lowercase slug string
func UniqueSlug(s string) string {
	return slugify.Slugify(fmt.Sprintf("%s-%s", strings.ToLower(s), uuid.NewV1().String()[:4]))
}
