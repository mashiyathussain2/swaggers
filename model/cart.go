package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Define name of the cart collection
const (
	CartColl string = "cart"
)

//CouponType
const (
	FreeDelivery string = "free_delivery"
)

//CheckoutType
const (
	CartCheckout    string = "cart_checkout"
	ExpressCheckout string = "express_checkout"
)

//Cart contains the users cart details
type Cart struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ShippingAddress *Address           `json:"shipping_address,omitempty" bson:"shipping_address,omitempty"`
	BillingAddress  *Address           `json:"billing_address,omitempty" bson:"billing_address,omitempty"`
	Items           []Item             `json:"items,omitempty" bson:"items,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Coupon          *Coupon            `json:"coupon,omitempty" bson:"coupon,omitempty"`
}

//Item is a unique catalogs data inside the cart
type Item struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID     primitive.ObjectID `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	BrandID       primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	VariantID     primitive.ObjectID `json:"variant_id,omitempty" bson:"variant_id,omitempty"`
	CatalogInfo   *CatalogInfo       `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	DiscountID    primitive.ObjectID `json:"discount_id,omitempty" bson:"discount_id,omitempty"`
	DiscountInfo  *DiscountInfo      `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	BasePrice     *Price             `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   *Price             `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice *Price             `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	Quantity      uint               `json:"quantity,omitempty" bson:"quantity,omitempty"`
	BrandInfo     *BrandInfoResp     `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	InStock       *bool              `json:"in_stock,omitempty" bson:"in_stock,omitempty"`
	Source        *Source            `json:"source,omitempty" bson:"source,omitempty"`
}
type Coupon struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Code             string             `json:"code,omitempty" bson:"code,omitempty"`
	Description      string             `json:"description,omitempty" bson:"description"`
	Type             DiscountType       `json:"type,omitempty" bson:"type,omitempty"`
	Value            int                `json:"value" bson:"value"`
	ApplicableON     *ApplicableON      `json:"applicable_on,omitempty" bson:"applicable_on,omitempty"`
	MaxDiscount      *Price             `json:"max_discount,omitempty" bson:"max_discount,omitempty"`
	MinPurchaseValue *Price             `json:"min_purchase_value,omitempty" bson:"min_purchase_value,omitempty"`
	ValidAfter       time.Time          `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore      time.Time          `json:"valid_before,omitempty" bson:"valid_before,omitempty"`
	Status           string             `json:"status,omitempty" bson:"status,omitempty"`
}
type ApplicableON struct {
	Name string               `json:"name,omitempty" bson:"name,omitempty"`
	IDs  []primitive.ObjectID `json:"ids,omitempty" bson:"ids,omitempty"`
}

type Source struct {
	ID   string `json:"id,omitempty" bson:"id,omitempty"`
	Type string `json:"type,omitempty" bson:"type,omitempty"`
}
