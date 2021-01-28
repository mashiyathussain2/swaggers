package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// IMG contains image url, src, height and id
type IMG struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SRC    string             `json:"src" bson:"src"`
	Height int                `json:"height" bson:"height"`
	Width  int                `json:"width" bson:"width"`
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
func SetINRPrice(v float32) Price {
	return Price{CurrencyISO: INR, Value: v}
}
