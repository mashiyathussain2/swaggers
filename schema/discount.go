package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateDiscountOpts serialize & validates schema to create a new catalog discount
type CreateDiscountOpts struct {
	CatalogID   primitive.ObjectID   `json:"catalog_id" validate:"required"`
	VariantsID  []primitive.ObjectID `json:"variants_id" validate:"required,gt=0"`
	SaleID      primitive.ObjectID   `json:"sale_id"`
	Type        string               `json:"type" validate:"required,oneof=flat_off percent_off"`
	Value       uint                 `json:"value" validate:"required,gt=0"`
	MaxValue    uint                 `json:"max_value" validate:"required_if=Type percent_off"`
	ValidAfter  time.Time            `json:"valid_after" validate:"required_without=SaleID"`
	ValidBefore time.Time            `json:"valid_before" validate:"required_without=SaleID,isdefault|gtfield=ValidAfter"`
}

// CreateDiscountResp contains fields to return for create discount response
type CreateDiscountResp struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID  primitive.ObjectID   `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	VariantsID []primitive.ObjectID `json:"variants_id,omitempty" bson:"variants_id,omitempty"`
	SaleID     primitive.ObjectID   `json:"sale_id,omitempty" bson:"sale_id,omitempty"`

	IsActive bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	Type     model.DiscountType `json:"type,omitempty" bson:"type,omitempty"`

	Value uint `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`

	// If discount is part of sale then ValidAfter & ValidBefore values will be inherited from sale only.
	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// CreateSaleOpts validates schema for creating a new sale
type CreateSaleOpts struct {
	Name        string    `json:"name" validate:"required"`
	Banner      Img       `json:"banner" validate:"required"`
	ValidAfter  time.Time `json:"valid_after" validate:"required"`
	ValidBefore time.Time `json:"valid_before" validate:"required,gtfield=ValidAfter"`
}

// CreateSaleResp contains fields to be returned as response when a new sale is created
type CreateSaleResp struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`

	Banner model.IMG `json:"banner,omitempty" bson:"banner,omitempty"`

	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
