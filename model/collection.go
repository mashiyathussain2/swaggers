package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollectionColl defines the name of the collection for Collections
const (
	CollectionColl = "collection"
)

// Different types of Collections
const (
	BourbonCollection   = "bourbon"
	DialCollection      = "dial"
	ProductCollection   = "product"
	EditorialCollection = "editorial"
)

//Collection contains Collection specific data such as CollectionType, CatalogIDs, Gender
type Collection struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name,omitempty"`
	Type           string             `json:"type,omitempty" bson:"type,omitempty"`
	Genders        []string           `json:"genders,omitempty" bson:"genders,omitempty"`
	Title          string             `json:"title,omitempty" bson:"title,omitempty"`
	SubCollections []SubCollection    `json:"sub_collections,omitempty" bson:"sub_collections,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Status         string             `json:"status,omitempty" bson:"status,omitempty"`
	Order          int                `json:"order,omitempty" bson:"order,omitempty"`
}

//SubCollection contains SubCollection specific data such as Name, Image and CatalogIDs
type SubCollection struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty" bson:"name,omitempty"`
	Image      *IMG                 `json:"image,omitempty" bson:"image,omitempty"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
