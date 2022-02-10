package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

//CreateInventoryOpts serializes the input for create inventory api
type CreateInventoryOpts struct {
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID `json:"variant_id" validate:"required"`
	SKU       string             `json:"sku,omitempty"`
	Unit      int                `json:"unit" validate:"required"`
}

//UpdateInventoryOperation defines struct for operation on Inventory
type UpdateInventoryOperation struct {
	Operator string `json:"operator" validate:"required,oneof=set add subtract"`
	Unit     int    `json:"unit" validate:"required"`
}

//UpdateInventoryOpts serializes the input for update inventory api
type UpdateInventoryOpts struct {
	ID        primitive.ObjectID        `json:"id" validate:"required"`
	Operation *UpdateInventoryOperation `json:"operation" validate:"required"`
}

//UpdateInventoryResp defines the response for update inventory api
type UpdateInventoryResp struct {
	ID   primitive.ObjectID `json:"inventory_id"`
	Unit int                `json:"unit"`
}

//UpdateInventoryOpts serializes the input for update inventory internal api with catalog and variant id

type UpdateInventoryCVOpts struct {
	CatalogID primitive.ObjectID        `json:"catalog_id" validate:"required"`
	VariantID primitive.ObjectID        `json:"variant_id" validate:"required"`
	Operation *UpdateInventoryOperation `json:"operation" validate:"required"`
}
type UpdateInventoryInternalOpts struct {
	Updates []UpdateInventoryCVOpts `json:"updates" validate:"required"`
}

//UpdateInventoryBySKUOpt defines struct for operation on Inventory by SKUs
type UpdateInventoryBySKUOpt struct {
	// BrandID primitive.ObjectID `json:"brand_id"`
	SKU  string `json:"sku" validate:"required"`
	Unit int    `json:"unit"`
}

//UpdateInventoryBySKUResp defines struct for response on Inventory by SKUs
type UpdateInventoryBySKUResp struct {
	DuplicateSKUs []string `json:"duplicate_skus"`
	InvalidSKUs   []string `json:"invalid_skus"`
}

type UnicommerceUpdateInventoryByInventoryIDsVariantOpts struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
	Unit      int    `json:"inventory"`
	HSNCode   string `json:"hsnCode"`
}

type UnicommerceUpdateInventoryByInventoryIDsOpts struct {
	InventoryList []UnicommerceUpdateInventoryByInventoryIDsVariantOpts `json:"inventoryList"`
}
