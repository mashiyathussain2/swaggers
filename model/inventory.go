package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//InventoryColl defines the Collection for storing Inventory
const (
	InventoryColl string = "inventory"
)

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
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID   primitive.ObjectID `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	VariantID   primitive.ObjectID `json:"variant_id,omitempty" bson:"variant_id,omitempty"`
	SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	Status      *InventoryStatus   `json:"status,omitempty" bson:"status,omitempty"`
	UnitInStock int                `json:"unit_in_stock,omitempty" bson:"unit_in_stock,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
