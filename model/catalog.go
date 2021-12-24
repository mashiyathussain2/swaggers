package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Path contains entire category hierarchy
type Path = string

// CatalogFeaturedImage has one featured image for a catalog in landscape, portrait, and square.
type CatalogFeaturedImage = IMG

// Specification contains catalog specification in key:value format
type Specification struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}

// Attribute define key value pair that defines catalog properties
type Attribute struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}

// Variant contains variants based on one property (size, color)
type Variant struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Attribute   string             `json:"attribute,omitempty" bson:"attribute,omitempty"`
	InventoryID primitive.ObjectID `json:"inventory_id,omitempty" bson:"inventory_id,omitempty"`
	SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	IsDeleted   bool               `json:"is_deleted" bson:"is_deleted"`
	Inventory   *Inventory         `json:"inventory_info" bson:"inventory_info"`
}

// VariantType is a paramater which defines the variant classification for a particular catalog such as size or color or design etc.
type VariantType = string

// Defining the type of variants for variant creation in catalog
const (
	DefaultType VariantType = ""
	SizeType    VariantType = "size"
	ColorType   VariantType = "color"
)

// ETA contains maximum and minimum delivery time of a catalog
type ETA struct {
	Min  int    `json:"min,omitempty" bson:"min,omitempty"`
	Max  int    `json:"max,omitempty" bson:"max,omitempty"`
	Unit string `json:"unit,omitempty" bson:"unit,omitempty"`
}

/*Status stores catalog status such as unlisted (default), published, archive

Unlist: status is set by default when a new catalog instance is created.
		catalog with this status are now shown to the customer as this represents WIP catalog.
		unlist is one time status only once unlist status is changed it cannot be reverted.

		allowed status transitions 	-> publish
									-> discard

Publish: 	status is set thorugh admin/keeper dashboard to allow visibility of a catalog to customer.
			publish status can only be changed to discard

			allowed status transitions 	-> discard

Discard:	discard is an alias to delete a catalog without actually deleting it from the database to avoid NOT FOUND for
			other services while searching for catalog
*/
type Status struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Value     string    `json:"value,omitempty" bson:"value,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type CatalogInfo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID   primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName string             `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	// Paths []Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	// Slug        string `json:"slug,omitempty" bson:"slug,omitempty"`
	// Description string `json:"description,omitempty" bson:"description,omitempty"`
	// Keywords      []string              `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	// Specifications  []Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	// FilterAttribute []Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []Variant   `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string      `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	// BasePrice     Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	// RetailPrice   Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice Price `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`

	ETA            *ETA              `json:"eta,omitempty" bson:"eta,omitempty"`
	Status         *Status           `json:"status,omitempty" bson:"status,omitempty"`
	DiscountInfo   *DiscountInfoResp `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	Tax            *Tax              `json:"tax,omitempty" bson:"tax,omitempty"`
	CommissionRate uint              `json:"commission_rate,omitempty" bson:"commission_rate,omitempty"`
}

//Defined Multiple Status for Inventory
const (
	InStockStatus    string = "in_stock"
	OutOfStockStatus string = "out_of_stock"
)

//InventoryStatus stores catalog status such as out_of_stock, in_stock, inactive
type InventoryStatus struct {
	Value     string    `json:"value,omitempty" bson:"value,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//Inventory contains inventory specific data
type Inventory struct {
	// ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// CatalogID   primitive.ObjectID `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	// VariantID   primitive.ObjectID `json:"variant_id,omitempty" bson:"variant_id,omitempty"`
	// SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	Status      *InventoryStatus `json:"status,omitempty" bson:"status,omitempty"`
	UnitInStock int              `json:"unit_in_stock,omitempty" bson:"unit_in_stock,omitempty"`
	// CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	// UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

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

	IsActive bool         `json:"is_active,omitempty" bson:"is_active,omitempty"`
	Type     DiscountType `json:"type,omitempty" bson:"type,omitempty"`

	Value uint `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`

	// If discount is part of sale then ValidAfter & ValidBefore values will be inherited from sale only.
	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//CatalogVariant contains fields which are returned to get variant
type CatalogVariant struct {
	ID            primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string                `json:"name,omitempty" bson:"name,omitempty"`
	BasePrice     Price                 `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   Price                 `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice Price                 `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	BrandID       primitive.ObjectID    `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	VariantType   VariantType           `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variant       Variant               `json:"variant,omitempty" bson:"variant,omitempty"`
	DiscountInfo  *DiscountInfo         `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	FeaturedImage *CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	InventoryInfo Inventory             `json:"inventory_info,omitempty" bson:"inventory_info,omitempty"`
}

type GetCatalogVariant struct {
	Success bool           `json:"success"`
	Payload CatalogVariant `json:"payload"`
}

//DiscountInfo contains discount data for particular variant
type DiscountInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type     DiscountType       `json:"type,omitempty" bson:"type,omitempty"`
	Value    uint               `json:"value,omitempty" bson:"value,omitempty"`
	MaxValue uint               `json:"max_value,omitempty" bson:"max_value,omitempty"`
}

//GetAllCatalogInfoResp contains fields which are returned on calling api getAllcatalogInfo
type GetAllCatalogInfoResp struct {
	Success bool               `json:"success"`
	Payload AllCatalogInfoResp `json:"payload"`
}

// VariantAllInfo contains all variant data (size, color)
type VariantAllInfo struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Attribute     string             `json:"attribute,omitempty" bson:"attribute,omitempty"`
	InventoryID   primitive.ObjectID `json:"inventory_id,omitempty" bson:"inventory_id,omitempty"`
	SKU           string             `json:"sku,omitempty" bson:"sku,omitempty"`
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`
	InventoryInfo Inventory          `json:"inventory_info" bson:"inventory_info"`
}

type AllCatalogInfoResp struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`

	Paths []Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	Slug          string                `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string                `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string              `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []Variant   `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string      `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	BasePrice     Price   `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   Price   `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice Price   `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	ETA           *ETA    `json:"eta,omitempty" bson:"eta,omitempty"`
	Status        *Status `json:"status,omitempty" bson:"status,omitempty"`

	CreatedAt      time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DiscountInfo   *DiscountInfoResp `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	BrandInfo      *BrandInfoResp    `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	Tax            *Tax              `json:"tax,omitempty" bson:"tax,omitempty"`
	CommissionRate uint              `json:"commission_rate,omitempty" bson:"commission_rate,omitempty"`
}
type Tax struct {
	Type      string     `json:"type,omitempty" bson:"type,omitempty"`
	Rate      float32    `json:"rate,omitempty" bson:"rate,omitempty"`
	TaxRanges []TaxRange `json:"tax_ranges,omitempty" bson:"tax_ranges,omitempty"`
}
type TaxRange struct {
	MinValue int     `json:"min_value" bson:"min_value"`
	MaxValue int     `json:"max_value,omitempty" bson:"max_value,omitempty"`
	Rate     float32 `json:"rate,omitempty" bson:"rate,omitempty"`
}
type BrandInfoResp struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Slug        string             `json:"slug,omitempty"`
	Description string             `json:"description,omitempty"`
	Logo        *IMG               `json:"logo"`
}

type DiscountInfoResp struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID  primitive.ObjectID   `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	VariantsID []primitive.ObjectID `json:"variants_id,omitempty" bson:"variants_id,omitempty"`
	SaleID     primitive.ObjectID   `json:"sale_id,omitempty" bson:"sale_id,omitempty"`

	IsActive bool   `json:"is_active,omitempty" bson:"is_active,omitempty"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`

	Value uint `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`

	// If discount is part of sale then ValidAfter & ValidBefore values will be inherited from sale only.
	ValidAfter  time.Time `json:"valid_after,omitempty" bson:"valid_after,omitempty"`
	ValidBefore time.Time `json:"valid_before,omitempty" bson:"valid_before,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
