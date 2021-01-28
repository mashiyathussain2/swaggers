package schema

import (
	"go-app/model"
)

// CreateCatalogOpts serialize the create catalog api arguments
type CreateCatalogOpts struct {
	Name        string   `json:"name" bson:"name" validate:"required"`
	Description string   `json:"description" bson:"description" validate:"required"`
	Keywords    []string `json:"keywords" bson:"keywords" validate:"required,gt=0"`

	HSNCode string `json:"hsn_code" bson:"hsn_code" validate:"required,gt=0"`

	VariantType model.VariantType   `json:"variant_type" bson:"hsn_code" validate:"required_with_field=Variants"`
	Variants    []CreateVariantOpts `json:"variants" bson:"variants" validate:"required_with_field=VariantType,dive"`
}

// CreateVariantOpts serialize create variant arguments
type CreateVariantOpts struct {
	SKU         string `json:"sku" bson:"sku" validate:"required"`
	BasePrice   uint32 `json:"base_price" bson:"base_price" validate:"required,gt=0"`
	RetailPrice uint32 `json:"retail_price" bson:"retail_price" validate:"required,gt=0"`
}
