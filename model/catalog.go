package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CatalogFeaturedImage has one featured image for a catalog in landscape, portrait, and square.
type CatalogFeaturedImage struct {
	IMG
}

// Catalog contains catalog specific data such as name, description, linked content, brand info, keywords, specifications, variant info etc
type Catalog struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	BrandID primitive.ObjectID `json:"brand_id" bson:"brand_id"`

	Name  string `json:"name" bson:"name"`
	LName string `json:"lname" bson:"lname"`
	// slug is used in setting up catalog thumbnail image name and
	// when sharing a catalog sharing link is generated through slug
	Slug          string               `json:"slug" bson:"slug"`
	Description   string               `json:"description" bson:"description"`
	Keywords      []string             `json:"keywords" bson:"keywords"`
	FeaturedImage CatalogFeaturedImage `json:"featured_image" bson:"featured_image"`

	VariantType VariantType `json:"variant_type" bson:"variant_type"`
	Variants    []Variant   `json:"variants" bson:"variants"`
	HSNCode     string      `json:"hsn_code" bson:"hsn_code"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// VariantType is a paramater which defines the variant classification for a particular catalog such as size or color or design etc.
type VariantType string

// Defining the type of variants for variant creation in catalog
const (
	SizeType  VariantType = "size"
	ColorType VariantType = "color"
)

// Variant contains variants based on one property (size, color)
type Variant struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	InventoryID primitive.ObjectID `json:"inventory_id" bson:"inventory_id"`
	SKU         string             `json:"sku" bson:"sku"`
	BasePrice   Price              `json:"base_price" bson:"base_price"`
	RetailPrice Price              `json:"retail_price" bson:"retail_price"`
}
