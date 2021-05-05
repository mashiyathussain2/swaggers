package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collections
const (
	DiscountColl string = "discount"
	SaleColl     string = "sale"
)

// DiscountType is type of discount applicable on catalog
/* Currently 2 discount types are supported:
1. Fixed - Flat X amount off catalog
2. Percent - x-% off catalog amount
*/
type DiscountType = string

// Contains list of discount types
const (
	FlatOffType    DiscountType = "flat_off"
	PercentOffType DiscountType = "percent_off"
)

// SaleStatusType is type of status applicable on a sale
type SaleStatusType string

//Contains list of sale status
const (
	Live     SaleStatusType = "live"
	Schedule SaleStatusType = "schedule"
	Disable  SaleStatusType = "disable"
)

// Discount contains catalog level discounts.
// 1 unique catalog can only have document document/instance
type Discount struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID  primitive.ObjectID   `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	VariantsID []primitive.ObjectID `json:"variants_id,omitempty" bson:"variants_id,omitempty"`
	SaleID     primitive.ObjectID   `json:"sale_id,omitempty" bson:"sale_id,omitempty"`

	IsActive   bool         `json:"is_active" bson:"is_active"`
	IsDisabled bool         `json:"is_disabled" bson:"is_disabled"`
	Type       DiscountType `json:"type,omitempty" bson:"type,omitempty"`

	Value uint `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`

	// If discount is part of sale then ValidAfter & ValidBefore values will be inherited from sale only.
	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Sale contains grouping of various catalog discounts.
type Sale struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug      string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Status    SaleStatusType     `json:"status" bson:"status"`
	Genders   []string           `json:"genders" bson:"genders"`
	Banner    *IMG               `json:"banner,omitempty" bson:"banner,omitempty"`
	WebBanner *IMG               `json:"web_banner,omitempty" bson:"web_banner,omitempty"`

	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
