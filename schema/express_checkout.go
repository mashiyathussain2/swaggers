package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemExpress struct {
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID `json:"variant_id" validate:"required"`
	Quantity  int                `json:"quantity" validate:"required"`
	Source    *model.Source      `json:"source"`
}

type ExpressCheckoutOpts struct {
	UserID   primitive.ObjectID `json:"user_id" validate:"required"`
	Address  *OrderAddressOpts  `json:"address" validate:"required"`
	Items    []ItemExpress      `json:"items" validate:"required"`
	Source   string             `json:"source" validate:"required"`
	SourceID primitive.ObjectID `json:"source_id" validate:"required"`
	Coupon   string             `json:"coupon"`
}

type ExpressCheckoutWebOpts struct {
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	Address   *OrderAddressOpts  `json:"address" validate:"required"`
	Items     []ItemExpress      `json:"items" validate:"required"`
	Coupon    string             `json:"coupon"`
	Source    string             `json:"source"`
	SourceID  primitive.ObjectID `json:"source_id"`
	IsCOD     bool               `json:"is_cod,omitempty"`
	RequestID string             `json:"request_id,omitempty"`
}

type ExpressCheckoutWebV2Opts struct {
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	Address   *OrderAddressOpts  `json:"address" validate:"required"`
	Item      ItemExpress        `json:"item" validate:"required"`
	Coupon    string             `json:"coupon"`
	Source    string             `json:"source"`
	SourceID  primitive.ObjectID `json:"source_id"`
	IsCOD     bool               `json:"is_cod,omitempty"`
	RequestID string             `json:"request_id,omitempty"`
}
