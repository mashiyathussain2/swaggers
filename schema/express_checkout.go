package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemExpress struct {
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID `json:"variant_id" validate:"required"`
	Quantity  int                `json:"quantity" validate:"required"`
}

type ExpressCheckoutOpts struct {
	UserID   primitive.ObjectID `json:"user_id" validate:"required"`
	Address  *OrderAddressOpts  `json:"address" validate:"required"`
	Items    []ItemExpress      `json:"items" validate:"required"`
	Source   string             `json:"source" validate:"required"`
	SourceID primitive.ObjectID `json:"source_id" validate:"required"`
}

type ExpressCheckoutWebOpts struct {
	UserID  primitive.ObjectID `json:"user_id" validate:"required"`
	Address *OrderAddressOpts  `json:"address" validate:"required"`
	Items   []ItemExpress      `json:"items" validate:"required"`
}
