package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderCatalogInfo struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`

	FeaturedImage Img `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	VariantType model.VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variant     OrderVariant      `json:"variant,omitempty" bson:"variant,omitempty"`
	HSNCode     string            `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	TransferPrice model.Price `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	ETA           *model.ETA  `json:"eta,omitempty" bson:"eta,omitempty"`
	Tax           *model.Tax  `json:"tax,omitempty" bson:"tax,omitempty"`
}

type OrderVariant struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Attribute string             `json:"attribute,omitempty" bson:"attribute,omitempty"`
	SKU       string             `json:"sku,omitempty" bson:"sku,omitempty"`
}
