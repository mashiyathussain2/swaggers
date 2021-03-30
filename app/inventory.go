//go:generate $GOBIN/mockgen -destination=./../mock/mock_inventory.go -package=mock go-app/app Inventory

package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Inventory service allows `Inventory` to execute admin operations.
type Inventory interface {
	CreateInventory(*schema.CreateInventoryOpts) (primitive.ObjectID, error)
	UpdateInventory(*schema.UpdateInventoryOpts) error
	SetOutOfStock(primitive.ObjectID) error
	CheckInventoryExists(primitive.ObjectID, primitive.ObjectID) (bool, error)
}

// InventoryImpl implements Inventory related operations
type InventoryImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InventoryOpts contains arguments required to create a new instance of Inventory
type InventoryOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

//InitInventory returns Inventory instance
func InitInventory(opts *InventoryOpts) Inventory {
	return &InventoryImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

//CreateInventory creates and set inventory for Variant in DB
func (ii *InventoryImpl) CreateInventory(opts *schema.CreateInventoryOpts) (primitive.ObjectID, error) {

	ctx := context.TODO()

	findFilter := bson.M{"catalog_id": opts.CatalogID, "variant_id": opts.VariantID}
	count, err := ii.DB.Collection(model.InventoryColl).CountDocuments(ctx, findFilter)
	if err != nil {
		return primitive.NilObjectID, errors.Wrapf(err, "unable to create inventory ")
	}
	if count != 0 {
		return primitive.NilObjectID, errors.Errorf("inventory already exist")
	}

	ti := time.Now().UTC()
	inventory := model.Inventory{
		CatalogID: opts.CatalogID,
		VariantID: opts.VariantID,
		SKU:       opts.SKU,
		Status: &model.InventoryStatus{
			Value:     model.OutOfStockStatus,
			CreatedAt: ti,
		},
		UnitInStock: opts.Unit,
		CreatedAt:   ti,
	}

	if opts.Unit > 0 {
		inventory.Status.Value = model.InStockStatus
	}

	res, err := ii.DB.Collection(model.InventoryColl).InsertOne(ctx, inventory)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

//UpdateInventory updates inventory in 3 ways - set(replace), add and remove
func (ii *InventoryImpl) UpdateInventory(opts *schema.UpdateInventoryOpts) error {

	ctx := context.TODO()
	findQuery := bson.M{
		"_id": opts.ID,
	}
	updateQuery := bson.D{}
	var inventory model.Inventory
	err := ii.DB.Collection(model.InventoryColl).FindOne(ctx, findQuery).Decode(&inventory)

	switch opts.Operation.Operator {
	case "set":
		updateQuery = append(updateQuery, bson.E{
			Key: "$set", Value: bson.M{
				"unit_in_stock": opts.Operation.Unit,
			},
		})
		if opts.Operation.Unit == 0 {
			updateQuery = append(updateQuery, bson.E{
				Key: "$set", Value: bson.M{
					"status": model.InventoryStatus{
						Value:     model.OutOfStockStatus,
						CreatedAt: time.Now().UTC(),
					},
				},
			})
		}
		if inventory.UnitInStock == 0 && opts.Operation.Unit != 0 {
			updateQuery = append(updateQuery, bson.E{
				Key: "$set", Value: bson.M{
					"status": model.InventoryStatus{
						Value:     model.InStockStatus,
						CreatedAt: time.Now().UTC(),
					},
				},
			})
		}
	case "add":
		updateQuery = append(updateQuery, bson.E{
			Key: "$inc", Value: bson.M{
				"unit_in_stock": opts.Operation.Unit,
			},
		})
		if inventory.Status.Value == model.OutOfStockStatus {
			updateQuery = append(updateQuery, bson.E{
				Key: "$set", Value: bson.M{
					"status": model.InventoryStatus{
						Value:     model.InStockStatus,
						CreatedAt: time.Now().UTC(),
					},
				},
			})
		}

	case "subtract":

		if inventory.UnitInStock-opts.Operation.Unit < 0 {
			return errors.Errorf("inventory for id: %s, cannot be negative", opts.ID)
		}

		updateQuery = append(updateQuery, bson.E{
			Key: "$inc", Value: bson.M{
				"unit_in_stock": -opts.Operation.Unit,
			},
		})

		if inventory.UnitInStock-opts.Operation.Unit == 0 {
			updateQuery = append(updateQuery, bson.E{
				Key: "$set", Value: bson.M{
					"status": model.InventoryStatus{
						Value:     model.OutOfStockStatus,
						CreatedAt: time.Now().UTC(),
					},
				},
			})
		}
	}
	// updateQuery = append(updateQuery, bson.E{
	// 	Key: "$set", Value: bson.M{
	// 		"updated_at": time.Now(),
	// 	},
	// })

	res, err := ii.DB.Collection(model.InventoryColl).UpdateOne(ctx, findQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to update inventory with id: %s", opts.ID.Hex())
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find the inventory with id: %s", opts.ID.Hex())
	}
	if res.ModifiedCount == 0 {
		return errors.Errorf("unable to update the inventory with id: %s", opts.ID.Hex())
	}
	return nil
}

//SetOutOfStock sets the status out of stock for the inventory with given id
func (ii *InventoryImpl) SetOutOfStock(id primitive.ObjectID) error {
	filterQuery := bson.M{"_id": id}

	status := model.InventoryStatus{
		Value:     model.OutOfStockStatus,
		CreatedAt: time.Now().UTC(),
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"status":        status,
			"unit_in_stock": 0,
		},
	}
	res, err := ii.DB.Collection(model.InventoryColl).UpdateOne(context.TODO(), filterQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "error updating inventory with id: %s", id.Hex())
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find inventory with id: %s", id.Hex())
	}
	return nil
}

func (ii *InventoryImpl) CheckInventoryExists(cat_id, var_id primitive.ObjectID) (bool, error) {
	ctx := context.TODO()

	filter := bson.M{
		"catalog_id": cat_id,
		"variant_id": var_id,
	}
	var inventory model.Inventory
	err := ii.DB.Collection(model.InventoryColl).FindOne(ctx, filter).Decode(&inventory)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, errors.Errorf("inventory not found")
		}
		return false, errors.Errorf("unable to query for document")
	}
	if inventory.UnitInStock == 0 || inventory.Status.Value == model.OutOfStockStatus {
		return false, nil
	}
	return true, nil
}
