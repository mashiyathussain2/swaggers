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
