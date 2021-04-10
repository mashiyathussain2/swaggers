package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Define name of the cart collection
const (
	CartColl string = "cart"
)

//Cart contains the users cart details
type Cart struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ShippingAddress *Address           `json:"shipping_address,omitempty" bson:"shipping_address,omitempty"`
	BillingAddress  *Address           `json:"billing_address,omitempty" bson:"billing_address,omitempty"`
	Items           []Item             `json:"items,omitempty" bson:"items,omitempty"`
	TotalPrice      *Price             `json:"total_price,omitempty" bson:"total_price,omitempty"`
	TotalDiscount   *Price             `json:"total_discount,omitempty" bson:"total_discount,omitempty"`
	GrandTotal      *Price             `json:"grand_total,omitempty" bson:"grand_total,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

//Item is a unique catalogs data inside the cart
type Item struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID       primitive.ObjectID `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	BrandID         primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	VariantID       primitive.ObjectID `json:"variant_id,omitempty" bson:"variant_id,omitempty"`
	CatalogInfo     CatalogInfo        `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	DiscountID      primitive.ObjectID `json:"discount_id,omitempty" bson:"discount_id,omitempty"`
	DiscountInfo    *DiscountInfo      `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	BasePrice       *Price             `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice     *Price             `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice   *Price             `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	DiscountedPrice *Price             `json:"discounted_price,omitempty" bson:"discounted_price,omitempty"`
	Quantity        uint               `json:"quantity,omitempty" bson:"quantity,omitempty"`
	BrandInfo       *BrandInfoResp     `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
}
