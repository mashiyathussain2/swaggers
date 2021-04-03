package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name
const (
	CatalogColl = "catalog"
)

// list of tax type
const (
	SingleTax   = "single"
	MultipleTax = "multiple"
)

// CatalogFeaturedImage has one featured image for a catalog in landscape, portrait, and square.
type CatalogFeaturedImage struct {
	IMG
}

// ETA contains maximum and minimum delivery time of a catalog
type ETA struct {
	Min  int    `json:"min,omitempty" bson:"min,omitempty"`
	Max  int    `json:"max,omitempty" bson:"max,omitempty"`
	Unit string `json:"unit,omitempty" bson:"unit,omitempty"`
}

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

// Path contains entire category hierarchy
type Path = string

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

// Catalog contains catalog specific data such as name, description, linked content, brand info, keywords, specifications, variant info etc
type Catalog struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`

	/*
		Path stores entire path of category as a string of hyphen seperated ids Eg. /84700/80009/1282094266/1200003270
		Using this we can store multiple category path inside a single catalog
		category_path: [/84700/80009/1282094266, /84701/80008/1282094267]
						men/footwear/casual, women/footwear/casual

		filtering category can be done using regex such as {$regex: "^84700/$"} {$regex: "^84700/80009/$"}
	*/
	Paths []Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	LName string `json:"lname,omitempty" bson:"lname,omitempty"`
	// slug is used in setting up catalog thumbnail image name and
	// when sharing a catalog sharing link is generated through slug
	Slug          string                `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string                `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string              `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []Attribute     `json:"filter_attrs,omitempty" bson:"filter_attrs,omitempty"`

	VariantType VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []Variant   `json:"variants,omitempty" bson:"variants,omitempty"`

	ETA           *ETA     `json:"eta,omitempty" bson:"eta,omitempty"`
	Status        *Status  `json:"status,omitempty" bson:"status,omitempty"`
	StatusHistory []Status `json:"status_history,omitempty" bson:"status_history,omitempty"`

	HSNCode       string `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`
	BasePrice     *Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   *Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice *Price `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`

	Tax *Tax `json:"tax,omitempty" bson:"tax,omitempty"`

	CatalogContent []primitive.ObjectID `json:"catalog_content,omitempty" bson:"catalog_content,omitempty"`

	// CatalogContentInfo []C

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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

// VariantType is a paramater which defines the variant classification for a particular catalog such as size or color or design etc.
type VariantType = string

// Defining the type of variants for variant creation in catalog
const (
	DefaultType VariantType = ""
	SizeType    VariantType = "size"
	ColorType   VariantType = "color"
)

// Variant contains variants based on one property (size, color)
type Variant struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Attribute   string             `json:"attribute,omitempty" bson:"attribute,omitempty"`
	InventoryID primitive.ObjectID `json:"inventory_id,omitempty" bson:"inventory_id,omitempty"`
	SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	IsDeleted   bool               `json:"is_deleted" bson:"is_deleted"`
}

// Defining the type of Catalog Status
const (
	Draft   string = "draft"
	Unlist  string = "unlist"
	Archive string = "archive"
	Publish string = "publish"
)
