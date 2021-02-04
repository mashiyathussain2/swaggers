package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ETAOpts serialize catalog estimated time of arrival
type etaOpts struct {
	Min  uint   `json:"min" validate:"required"`
	Max  uint   `json:"max" validate:"required"`
	Unit string `json:"unit" validate:"required,oneof=hour day month"`
}

type specsOpts struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type filterAttribute struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// CreateCatalogOpts serialize the create catalog api arguments
type CreateCatalogOpts struct {
	Name        string               `json:"name" validate:"required"`
	CategoryID  []primitive.ObjectID `json:"category_id" validate:"required,gt=0"`
	BrandID     primitive.ObjectID   `json:"brand_id" validate:"required"`
	Description string               `json:"description" validate:"required"`
	Keywords    []string             `json:"keywords" validate:"required,gt=0,unique"`

	ETA             *etaOpts          `json:"eta"`
	Specifications  []specsOpts       `json:"specifications" validate:"dive"`
	FilterAttribute []filterAttribute `json:"filter_attr" validate:"dive"`

	HSNCode string `json:"hsn_code" validate:"required,gt=0"`

	VariantType model.VariantType   `json:"variant_type" validate:"required_with_field=Variants"`
	Variants    []CreateVariantOpts `json:"variants" validate:"dive"`

	BasePrice   uint32 `json:"base_price" validate:"gt=0,gtefield=RetailPrice"`
	RetailPrice uint32 `json:"retail_price" validate:"gt=0"`
}

// CreateVariantOpts contains create variant arguments
type CreateVariantOpts struct {
	SKU       string `json:"sku" validate:"required"`
	Attribute string `json:"attribute"`
}

// CreateCatalogResp response
type CreateCatalogResp struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`

	Paths []model.Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	// Slug          string                      `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string                      `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string                    `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *model.CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []model.Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []model.Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType model.VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []model.Variant   `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string            `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	BasePrice   model.Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice model.Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`

	ETA    *model.ETA    `json:"eta,omitempty" bson:"eta,omitempty"`
	Status *model.Status `json:"status,omitempty" bson:"status,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// EditCatalogOpts contains fields which are passed in edit catalog func as args
type EditCatalogOpts struct {
	ID              primitive.ObjectID   `json:"id" validate:"required"`
	Name            string               `json:"name"`
	CategoryID      []primitive.ObjectID `json:"category_id"`
	Keywords        []string             `json:"keywords" validate:"unique"`
	ETA             *etaOpts             `json:"eta"`
	Specifications  []specsOpts          `json:"specifications" validate:"dive"`
	FilterAttribute []filterAttribute    `json:"filter_attr" validate:"dive"`

	HSNCode     string `json:"hsn_code"`
	BasePrice   uint32 `json:"base_price" validate:"gtefield=RetailPrice"`
	RetailPrice uint32 `json:"retail_price"`
}

// EditCatalogResp contains fields which are returned when a catalog is edited
type EditCatalogResp = CreateCatalogResp
