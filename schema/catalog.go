package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ETAOpts serialize catalog estimated time of arrival
type etaOpts struct {
	Min  uint   `json:"min" validate:"required"`
	Max  uint   `json:"max" validate:"required"`
	Unit string `json:"unit" validate:"required,oneof=hour day month"`
}

type specsOpts struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type FilterAttribute struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type TaxOpts struct {
	Type      string           `json:"type,omitempty" validate:"required,oneof=single multiple"`
	Rate      float32          `json:"rate,omitempty"`
	TaxRanges []model.TaxRange `json:"tax_ranges,omitempty" validate:"required_without=Rate" `
}

// CreateCatalogOpts serialize the create catalog api arguments
type CreateCatalogOpts struct {
	Name            string               `json:"name" validate:"required"`
	CategoryID      []primitive.ObjectID `json:"category_id" validate:"required,gt=0"`
	BrandID         primitive.ObjectID   `json:"brand_id" validate:"required"`
	Description     string               `json:"description" validate:"required"`
	Keywords        []string             `json:"keywords" validate:"required,gt=0,unique"`
	FeaturedImage   *Img                 `json:"featured_image" validate:"required"`
	ETA             *etaOpts             `json:"eta"`
	Specifications  []specsOpts          `json:"specifications" validate:"dive"`
	FilterAttribute []FilterAttribute    `json:"filter_attr" validate:"dive"`

	HSNCode string `json:"hsn_code" validate:"required,gt=0"`

	VariantType model.VariantType   `json:"variant_type" validate:"required_with_field=Variants"`
	Variants    []CreateVariantOpts `json:"variants" validate:"dive"`

	BasePrice     uint32             `json:"base_price" validate:"gt=0,gtefield=RetailPrice"`
	RetailPrice   uint32             `json:"retail_price" validate:"gt=0"`
	TransferPrice uint32             `json:"transfer_price" validate:"gt=0"`
	SizeProfile   *model.SizeProfile `json:"size_profile"`
	Tax           *TaxOpts           `json:"tax" validate:"required"`
}

// CreateCatalogResp response
type CreateCatalogResp struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`

	Paths []model.Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	Slug          string     `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string     `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string   `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *model.IMG `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []model.Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []model.Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType model.VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []model.Variant   `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string            `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	BasePrice     model.Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice model.Price `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`

	ETA         *model.ETA         `json:"eta,omitempty" bson:"eta,omitempty"`
	Status      *model.Status      `json:"status,omitempty" bson:"status,omitempty"`
	SizeProfile *model.SizeProfile `json:"size_profile,omitempty" bson:"size_profile,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`

	Tax *model.Tax `json:"tax,omitempty" bson:"tax,omitempty"`
}

// CreateVariantOpts contains create variant arguments
type CreateVariantOpts struct {
	SKU       string `json:"sku" validate:"required"`
	Attribute string `json:"attribute" validate:"required"`
	Unit      int    `json:"unit" validate:"required"`
}

// AddVariantOpts contains fields to add a new variant into existing catalog
type AddVariantOpts struct {
	ID          primitive.ObjectID `json:"id" validate:"required"`
	VariantType string             `json:"variant_type" validate:"required"`
	SKU         string             `json:"sku" validate:"required"`
	Attribute   string             `json:"attribute" validate:"required"`
	Unit        int                `json:"unit" validate:"required"`
}

// CreateVariantResp contains response fields when a new variant is created
type CreateVariantResp struct {
	ID        primitive.ObjectID `json:"id"`
	SKU       string             `json:"sku"`
	Attribute string             `json:"attribute"`
	Unit      int                `json:"unit"`
}

// AddVariantResp contains response fields when a new variant is added into existing catalog
type AddVariantResp = CreateVariantResp

// EditCatalogOpts contains fields which are passed in edit catalog func as args
type EditCatalogOpts struct {
	ID              primitive.ObjectID   `json:"id" validate:"required"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	FeaturedImage   *Img                 `json:"featured_image"`
	CategoryID      []primitive.ObjectID `json:"category_id"`
	Keywords        []string             `json:"keywords" validate:"unique"`
	ETA             *etaOpts             `json:"eta"`
	Specifications  []specsOpts          `json:"specifications" validate:"dive"`
	FilterAttribute []FilterAttribute    `json:"filter_attr" validate:"dive"`
	HSNCode         string               `json:"hsn_code"`
	VariantType     string               `json:"variant_type"`
	BasePrice       uint32               `json:"base_price" validate:"isdefault|gtefield=RetailPrice"`
	RetailPrice     uint32               `json:"retail_price" validate:"isdefault|gt=0"`
	TransferPrice   uint32               `json:"transfer_price" validate:"isdefault|gt=0"`
	Tax             *TaxOpts             `json:"tax"`
	SizeProfile     *model.SizeProfile   `json:"size_profile"`
}

// EditCatalogResp contains fields which are returned when a catalog is edited
type EditCatalogResp struct {
	ID              primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Paths           []model.Path          `json:"category_path,omitempty" bson:"category_path,omitempty"`
	Name            string                `json:"name,omitempty" bson:"name,omitempty"`
	Description     string                `json:"description,omitempty" bson:"description,omitempty"`
	Keywords        []string              `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage   *model.IMG            `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	Specifications  []model.Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []model.Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`
	HSNCode         string                `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`
	BasePrice       model.Price           `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice     model.Price           `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	VariantType     string                `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	ETA             *model.ETA            `json:"eta,omitempty" bson:"eta,omitempty"`
	UpdatedAt       time.Time             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	TransferPrice   *model.Price          `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	Tax             model.Tax             `json:"tax,omitempty" bson:"tax,omitempty"`
	SizeProfile     *model.SizeProfile    `json:"size_profile,omitempty" bson:"size_profile,omitempty"`
}

// GetBasicCatalogFilter contains filter fields for GetCatalog
type GetBasicCatalogFilter struct {
	BrandID    []primitive.ObjectID `json:"id"`
	CategoryID []primitive.ObjectID `json:"category_id"`
}

// GetBasicCatalogResp contains fields to be returned as GetCatalog response
type GetBasicCatalogResp struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Paths       []model.Path       `json:"category_path,omitempty" bson:"category_path,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	RetailPrice model.Price        `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
}

// GetCatalogFilterResp response contains filter list and their values to be returned
type GetCatalogFilterResp struct {
	Category []GetCategoriesBasicResp `json:"category"`
}

// KeeperSearchCatalogOpts contains fields which are passed on catalog search function
type KeeperSearchCatalogOpts struct {
	Name string `json:"name" validate:"required"`
	Page int64  `json:"page" validate:"gte=0"`
}

// KeeperSearchCatalogResp contains fields which are returned on catalog search to Keeper
type KeeperSearchCatalogResp struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Path          []model.Path       `json:"category_path,omitempty" bson:"category_path,omitempty"`
	BasePrice     model.Price        `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price        `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice model.Price        `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	Status        *model.Status      `json:"status,omitempty" bson:"status,omitempty"`
	Variants      []model.Variant    `json:"variants,omitempty" bson:"variants,omitempty"`
	VariantType   model.VariantType  `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
}

// DeleteVariantOpts contains fields which are passed on delete variant from catalog function
type DeleteVariantOpts struct {
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID `json:"variant_id" validate:"required"`
}

// UpdateCatalogStatusOpts contains fields which are passed on update catalog status function
type UpdateCatalogStatusOpts struct {
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	Status    string             `json:"status" validate:"required,oneof=publish unlist draft archive"`
}

//UpdateCatalogStatusResp contains fields which are returned on update catalog status function
type UpdateCatalogStatusResp struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

//GetCatalogResp alias for Get Catalog
type GetCatalogResp struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`

	Paths []model.Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	Slug          string     `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string     `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string   `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *model.IMG `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []model.Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []model.Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType model.VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []VariantInfo     `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string            `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	BasePrice     model.Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice model.Price `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`

	ETA    *model.ETA    `json:"eta,omitempty" bson:"eta,omitempty"`
	Status *model.Status `json:"status,omitempty" bson:"status,omitempty"`

	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CatalogContent []primitive.ObjectID `json:"catalog_content,omitempty" bson:"catalog_content,omitempty"`

	Tax         *model.Tax         `json:"tax,omitempty" bson:"tax,omitempty"`
	SizeProfile *model.SizeProfile `json:"size_profile,omitempty" bson:"size_profile,omitempty"`
}

//AddCatalogContentOpts contains fields which are passed on add content api
type AddCatalogContentOpts struct {
	BrandID   primitive.ObjectID `json:"brand_id"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	FileName  string             `json:"file_name" validate:"required"`
}

//AddCatalogContentImageOpts contains fields which are passed on add content api
type AddCatalogContentImageOpts struct {
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	MediaID   primitive.ObjectID `json:"media_id" validate:"required"`
}

//AddCatalogContentResp contains fields which are received from CMS and passed to Keeper
type AddCatalogContentResp struct {
	Success bool         `json:"success"`
	Payload PayloadVideo `json:"payload"`
	Error   []ErrorCMS   `json:"error"`
}

type PayloadImage struct {
	ID primitive.ObjectID `json:"id"`
}

type PayloadVideo struct {
	ID    primitive.ObjectID `json:"id"`
	Token string             `json:"token"`
}

//AddCatalogContentImageResp contains fields which are received from CMS and passed to Keeper
type AddCatalogContentImageResp struct {
	Success bool         `json:"success"`
	Payload PayloadImage `json:"payload"`
	Error   []ErrorCMS   `json:"error"`
}

type ContentLabel struct {
	Interests []string `json:"interests" validate:"required" bson:"interests"`
	AgeGroup  []string `json:"age_group" validate:"required" bson:"interests"`
	Gender    []string `json:"gender" validate:"required" bson:"gender"`
}

//GetCatalogsByFilterOpts contains fields which are used to filter catalogs
type GetCatalogsByFilterOpts struct {
	Page     int                  `json:"page" validate:"gte=0"`
	BrandIDs []primitive.ObjectID `json:"brands"`
	Status   []string             `json:"status"`
	Name     string               `json:"name"`
}

//ErrorCMS contains Error response from CMS
type ErrorCMS struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type GetAllCatalogInfoResp struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID    primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	DiscountID primitive.ObjectID `json:"discount_id,omitempty" bson:"discount_id,omitempty"`

	Paths []model.Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	Slug          string     `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string     `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string   `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *model.IMG `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []model.Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []model.Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType model.VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []VariantInfo     `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string            `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	BasePrice     model.Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice model.Price `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`

	ETA    *model.ETA    `json:"eta,omitempty" bson:"eta,omitempty"`
	Status *model.Status `json:"status,omitempty" bson:"status,omitempty"`

	CreatedAt    time.Time                `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time                `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	GroupInfo    []model.Group            `json:"group_info,omitempty" bson:"group_info,omitempty"`
	DiscountInfo *DiscountInfoResp        `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	ContentInfo  []CatalogContentInfoResp `json:"content_info,omitempty" bson:"content_info,omitempty"`
	BrandInfo    *BrandInfoResp           `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	Tax          *model.Tax               `json:"tax,omitempty" bson:"tax,omitempty"`
	SizeProfile  *model.SizeProfile       `json:"size_profile,omitempty" bson:"size_profile,omitempty"`
}

type CatalogKafkaMessage struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BrandID    primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	DiscountID primitive.ObjectID `json:"discount_id,omitempty" bson:"discount_id,omitempty"`

	Paths []model.Path `json:"category_path,omitempty" bson:"category_path,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// LName string `json:"lname,omitempty" bson:"lname,omitempty"`

	Slug          string     `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string     `json:"description,omitempty" bson:"description,omitempty"`
	Keywords      []string   `json:"keywords,omitempty" bson:"keywords,omitempty"`
	FeaturedImage *model.IMG `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	Specifications  []model.Specification `json:"specs,omitempty" bson:"specs,omitempty"`
	FilterAttribute []model.Attribute     `json:"filter_attr,omitempty" bson:"filter_attr,omitempty"`

	VariantType model.VariantType `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants    []VariantInfo     `json:"variants,omitempty" bson:"variants,omitempty"`
	HSNCode     string            `json:"hsn_code,omitempty" bson:"hsn_code,omitempty"`

	BasePrice   model.Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice model.Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`

	ETA    *model.ETA    `json:"eta,omitempty" bson:"eta,omitempty"`
	Status *model.Status `json:"status,omitempty" bson:"status,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CatalogContentInfoResp struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Type        string             `json:"type,omitempty"`
	MediaType   string             `json:"media_type,omitempty"`
	MediaID     primitive.ObjectID `json:"media_id,omitempty"`
	MediaInfo   interface{}        `json:"media_info,omitempty"`
	IsActive    bool               `json:"is_active"`
	IsProcessed bool               `json:"is_processed"`
}

type CatalogContentKafkaMessage struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	Type      string             `json:"type,omitempty"`
	MediaType string             `json:"media_type,omitempty"`
	MediaID   primitive.ObjectID `json:"media_id,omitempty"`
	MediaInfo interface{}        `json:"media_info,omitempty"`
}

type GetCatalogContentInfoResp struct {
	Success bool                     `json:"success"`
	Payload []CatalogContentInfoResp `json:"payload"`
}

type DiscountBasicResp struct {
	ID     primitive.ObjectID `json:"id,omitempty"`
	SaleID primitive.ObjectID `json:"sale_id,omitempty"`

	IsActive bool   `json:"is_active,omitempty"`
	Type     string `json:"type,omitempty"`

	Value uint `json:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty"`
}

type GetCatalogBasicResp struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	BrandID       primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandInfo     *BrandInfoResp     `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	DiscountID    primitive.ObjectID `json:"discount_id,omitempty" bson:"discount_id,omitempty"`
	DiscountInfo  *DiscountBasicResp `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	FeaturedImage *model.IMG         `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	BasePrice     model.Price        `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price        `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
}

type GetCatalogPebbleResp struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	BrandID       primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandInfo     *BrandInfoResp     `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	DiscountID    primitive.ObjectID `json:"discount_id,omitempty" bson:"discount_id,omitempty"`
	DiscountInfo  *DiscountBasicResp `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	FeaturedImage *model.IMG         `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	BasePrice     model.Price        `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price        `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
}

type GetCatalogInfoResp struct {
	ID              primitive.ObjectID       `json:"id,omitempty"`
	Name            string                   `json:"name,omitempty"`
	BrandID         primitive.ObjectID       `json:"brand_id,omitempty"`
	BrandInfo       *BrandInfoResp           `json:"brand_info,omitempty"`
	DiscountID      primitive.ObjectID       `json:"discount_id,omitempty"`
	DiscountInfo    *DiscountBasicResp       `json:"discount_info,omitempty"`
	FeaturedImage   *model.IMG               `json:"featured_image,omitempty"`
	BasePrice       model.Price              `json:"base_price,omitempty"`
	RetailPrice     model.Price              `json:"retail_price,omitempty"`
	Slug            string                   `json:"slug,omitempty"`
	Description     string                   `json:"description,omitempty"`
	Specifications  []model.Specification    `json:"specs,omitempty"`
	FilterAttribute []model.Attribute        `json:"filter_attr,omitempty"`
	VariantType     model.VariantType        `json:"variant_type,omitempty"`
	Variants        []VariantInfo            `json:"variants,omitempty"`
	ContentInfo     []CatalogContentInfoResp `json:"content_info,omitempty" bson:"content_info,omitempty"`
	SizeProfile     *model.SizeProfile       `json:"size_profile,omitempty" bson:"size_profile,omitempty"`
}

type GetCatalogByIDFilter struct {
	IDs []string `qs:"id" json:"id"`
}

type GetPebbleCatalogInfoByIDs struct {
	IDs []primitive.ObjectID `json:"id" validate:"required,min=1"`
}

//GetCatalogVariantResp contains fields which are returned to get variant
type GetCatalogVariantResp struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	BasePrice     model.Price        `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   model.Price        `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	TransferPrice model.Price        `json:"transfer_price,omitempty" bson:"transfer_price,omitempty"`
	VariantType   model.VariantType  `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variant       model.Variant      `json:"variant,omitempty" bson:"variant,omitempty"`
	DiscountInfo  DiscountInfo       `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	FeaturedImage *model.IMG         `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	InventoryInfo *model.Inventory   `json:"inventory_info" bson:"inventory_info"`
}

type DiscountInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type     model.DiscountType `json:"type,omitempty" bson:"type,omitempty"`
	Value    uint               `json:"value,omitempty" bson:"value,omitempty"`
	MaxValue uint               `json:"max_value,omitempty" bson:"max_value,omitempty"`
}

// VariantInfo contains all variant data (size, color)
type VariantInfo struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Attribute     string             `json:"attribute,omitempty" bson:"attribute,omitempty"`
	InventoryID   primitive.ObjectID `json:"inventory_id,omitempty" bson:"inventory_id,omitempty"`
	SKU           string             `json:"sku,omitempty" bson:"sku,omitempty"`
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`
	InventoryInfo *model.Inventory   `json:"inventory_info" bson:"inventory_info"`
}

type RemoveContentOpts struct {
	CatalogID primitive.ObjectID `json:"catalog_id"`
	ContentID primitive.ObjectID `json:"content_id"`
}

type EditVariantSKU struct {
	ID  primitive.ObjectID `json:"id" validate:"required"`
	SKU string             `json:"sku" validate:"required"`
}
