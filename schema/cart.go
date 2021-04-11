package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AddToCartOpts contains field required to add item into Cart
type AddToCartOpts struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID `json:"variant_id" validate:"required"`
	Quantity  uint               `json:"quantity" validate:"required,gt=0"`
}

//UpdateItemQtyOpts contains field required to update the quantity of a item already in the user's cart
type UpdateItemQtyOpts struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID `json:"variant_id" validate:"required"`
	Quantity  int                `json:"quantity" validate:"oneof=-1 0 1"`
}

//AddressOpts contains field required to add/edit the address of the user's cart
type AddressOpts struct {
	ID                primitive.ObjectID `json:"id" validate:"required"`
	AddressID         primitive.ObjectID `json:"address_id" validate:"required"`
	DisplayName       string             `json:"display_name"`
	Line1             string             `json:"line1" validate:"required"`
	Line2             string             `json:"line2"`
	District          string             `json:"district"`
	City              string             `json:"city" validate:"required"`
	State             *model.State       `json:"state" validate:"required"`
	PostalCode        string             `json:"postal_code" validate:"required"`
	Country           *model.Country     `json:"country" validate:"required"`
	PlainAddress      string             `json:"plain_address" validate:"required"`
	IsBillingAddress  bool               `json:"is_billing_address" validate:"required"`
	IsShippingAddress bool               `json:"is_shipping_address" validate:"required"`
	IsDefaultAddress  bool               `json:"is_default_address" validate:"required"`
	ContactNumber     *model.PhoneNumber `json:"contact_number" validate:"required"`
}

type CheckInventoryResp struct {
	Success bool `json:"success"`
	Payload bool `json:"payload"`
}

type CartUnwindBrand struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	BrandID         primitive.ObjectID `json:"brand_id" bson:"brand_id"`
	Items           []model.Item       `json:"items" bson:"items"`
	ShippingAddress *OrderAddressOpts  `json:"shipping_address" bson:"shipping_address"`
	BillingAddress  *OrderAddressOpts  `json:"billing_address" bson:"billing_address"`
}

type OrderOpts struct {
	UserID          primitive.ObjectID `json:"user_id"`
	BrandID         primitive.ObjectID `json:"brand_id"`
	ShippingAddress *OrderAddressOpts  `json:"shipping_address"`
	BillingAddress  *OrderAddressOpts  `json:"billing_address"`
	Source          string             `json:"source"`
	// SourceID        primitive.ObjectID `json:"source_id,omitempty"`
	OrderItems []OrderItem `json:"order_items" bson:"order_items"`
}

type OrderResp struct {
	Success bool      `json:"success"`
	Payload OrderInfo `json:"payload"`
}

type OrderInfo struct {
	OrderID    string  `json:"order_id" bson:"order_id"`
	RazorpayID string  `json:"razorpay_id" bson:"razorpay_id"`
	Amount     float32 `json:"amount" bson:"amount"`
}

//OrderItem is a unique catalogs data inside the cart
type OrderItem struct {
	CatalogID       primitive.ObjectID  `json:"catalog_id" bson:"catalog_id"`
	VariantID       primitive.ObjectID  `json:"variant_id" bson:"variant_id"`
	CatalogInfo     OrderCatalogInfo    `json:"catalog_info" bson:"catalog_info"`
	DiscountID      primitive.ObjectID  `json:"discount_id" bson:"discount_id"`
	DiscountInfo    *model.DiscountInfo `json:"discount_info" bson:"discount_info"`
	BasePrice       *model.Price        `json:"base_price" bson:"base_price"`
	RetailPrice     *model.Price        `json:"retail_price" bson:"retail_price"`
	DiscountedPrice *model.Price        `json:"discounted_price" bson:"discounted_price"`
	Quantity        uint                `json:"quantity" bson:"quantity"`
}

//OrderAddressOpts contains field required to add/edit the address of the user's cart
type OrderAddressOpts struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DisplayName   string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
	ContactNumber *model.PhoneNumber `json:"phone_no,omitempty" bson:"contact_number,omitempty"`
	Line1         string             `json:"line1,omitempty" bson:"line1,omitempty"`
	Line2         string             `json:"line2,omitempty" bson:"line2,omitempty"`
	District      string             `json:"district,omitempty" bson:"district,omitempty"`
	City          string             `json:"city,omitempty" bson:"city,omitempty"`
	State         *model.State       `json:"state,omitempty" bson:"state,omitempty"`
	PostalCode    string             `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	PlainAddress  string             `json:"plain_address,omitempty" bson:"plain_address,omitempty"`
}

type DiscountKafkaMessage struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	CatalogID  primitive.ObjectID   `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	VariantsID []primitive.ObjectID `json:"variants_id,omitempty" bson:"variants_id,omitempty"`
	SaleID     primitive.ObjectID   `json:"sale_id,omitempty" bson:"sale_id,omitempty"`

	IsActive   bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsDisabled bool               `json:"is_disabled,omitempty" bson:"is_disabled,omitempty"`
	Type       model.DiscountType `json:"type,omitempty" bson:"type,omitempty"`

	Value uint `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`

	// If discount is part of sale then ValidAfter & ValidBefore values will be inherited from sale only.
	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type DiscountInCartItemsOpts struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID  primitive.ObjectID   `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	VariantsID []primitive.ObjectID `json:"variants_id,omitempty" bson:"variants_id,omitempty"`
	SaleID     primitive.ObjectID   `json:"sale_id,omitempty" bson:"sale_id,omitempty"`

	IsActive   bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsDisabled bool               `json:"is_disabled,omitempty" bson:"is_disabled,omitempty"`
	Type       model.DiscountType `json:"type,omitempty" bson:"type,omitempty"`

	Value uint `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`

	// If discount is part of sale then ValidAfter & ValidBefore values will be inherited from sale only.
	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
