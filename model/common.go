package model

import (
	"image"
	"net/http"

	"github.com/pkg/errors"
)

// PhoneNumber represents a contact number contains prefix (country code) and phone number
type PhoneNumber struct {
	Prefix string `json:"prefix,omitempty" bson:"prefix,omitempty"`
	Number string `json:"number,omitempty" bson:"number,omitempty"`
}

// IMG contains image url, src, height and id
type IMG struct {
	SRC    string `json:"src" bson:"src"`
	Height int    `json:"height" bson:"height"`
	Width  int    `json:"width" bson:"width"`
}

// LoadFromURL loads the image from url and sets the width and height
func (i *IMG) LoadFromURL() error {
	resp, err := http.Get(i.SRC)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	m, _, err := image.Decode(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to decode img src: %s", i.SRC)
	}

	i.Height = m.Bounds().Dy()
	i.Width = m.Bounds().Dx()

	return nil
}

// CurrencyISO iso representation of currency
type CurrencyISO string

// Types of supported currency iso
const (
	INR CurrencyISO = "inr"
)

// Price represents cost of an entity
type Price struct {
	CurrencyISO CurrencyISO `json:"iso" bson:"iso"`
	Value       float32     `json:"value" bson:"value"`
}

// SetINRPrice sets INR and passed value returns the price struct
func SetINRPrice(v float32) *Price {
	return &Price{CurrencyISO: INR, Value: v}
}
