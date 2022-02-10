//go:generate $GOBIN/mockgen -destination=./../mock/mock_inventory.go -package=mock go-app/app Inventory

package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inventory service allows `Inventory` to execute admin operations.
type Inventory interface {
	CreateInventory(*schema.CreateInventoryOpts) (primitive.ObjectID, error)
	UpdateInventory(*schema.UpdateInventoryOpts) error
	SetOutOfStock(primitive.ObjectID) error
	CheckInventoryExists(primitive.ObjectID, primitive.ObjectID, int) (bool, error)
	UpdateInventoryInternal([]schema.UpdateInventoryCVOpts) error
	UpdateInventorybySKUs([]schema.UpdateInventoryBySKUOpt) (*schema.UpdateInventoryBySKUResp, error)

	UnicommerceUpdateInventoryByVariantIDs(*schema.UnicommerceUpdateInventoryByInventoryIDsOpts) error
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

	if err != nil {
		return errors.Wrap(err, "unable to query for inventory")
	}
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

func (ii *InventoryImpl) CheckInventoryExists(cat_id, var_id primitive.ObjectID, qty int) (bool, error) {
	ctx := context.TODO()

	if qty <= 0 {
		return false, errors.Errorf("quantity must be greater than 0")
	}

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
	if inventory.UnitInStock < qty || inventory.Status.Value == model.OutOfStockStatus {
		return false, nil
	}
	return true, nil
}

//UpdateInventory updates inventory in 3 ways - set(replace), add and remove
func (ii *InventoryImpl) UpdateInventoryInternal(opts []schema.UpdateInventoryCVOpts) error {

	ctx := context.TODO()

	var operations []mongo.WriteModel

	for i := range opts {

		operation := mongo.NewUpdateOneModel()

		findQuery := bson.M{
			"catalog_id": opts[i].CatalogID,
			"variant_id": opts[i].VariantID,
		}

		updateQuery := bson.D{}
		var inventory model.Inventory
		err := ii.DB.Collection(model.InventoryColl).FindOne(ctx, findQuery).Decode(&inventory)

		if err != nil {
			return errors.Wrapf(err, "unable to query for inventory with catalog id %s", opts[i].CatalogID.Hex())
		}
		switch opts[i].Operation.Operator {
		case "set":
			updateQuery = append(updateQuery, bson.E{
				Key: "$set", Value: bson.M{
					"unit_in_stock": opts[i].Operation.Unit,
				},
			})
			if opts[i].Operation.Unit == 0 {
				updateQuery = append(updateQuery, bson.E{
					Key: "$set", Value: bson.M{
						"status": model.InventoryStatus{
							Value:     model.OutOfStockStatus,
							CreatedAt: time.Now().UTC(),
						},
					},
				})
			}
			if inventory.UnitInStock == 0 && opts[i].Operation.Unit != 0 {
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
					"unit_in_stock": opts[i].Operation.Unit,
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

			if inventory.UnitInStock-opts[i].Operation.Unit < 0 {
				return errors.Errorf("inventory for catalog id: %s, cannot be negative", opts[i].CatalogID)
			}

			updateQuery = append(updateQuery, bson.E{
				Key: "$inc", Value: bson.M{
					"unit_in_stock": -opts[i].Operation.Unit,
				},
			})

			if inventory.UnitInStock-opts[i].Operation.Unit == 0 {
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

		operation.SetFilter(findQuery)
		operation.SetUpdate(updateQuery)
		operations = append(operations, operation)

	}

	if len(operations) == 0 {
		ii.Logger.Info().Msgf("no operations")
		return nil
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	_, err := ii.DB.Collection(model.InventoryColl).BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		ii.Logger.Err(err).Msgf("failed to update inventory")
		return errors.Wrap(err, "failed to update inventory")
	}
	return nil
}

func (ii *InventoryImpl) UnicommerceUpdateInventoryByVariantIDs(opts *schema.UnicommerceUpdateInventoryByInventoryIDsOpts) error {
	var operations []mongo.WriteModel
	for _, opts := range opts.InventoryList {

		var catalogId primitive.ObjectID
		var variantId primitive.ObjectID
		var err error
		operation := mongo.NewUpdateOneModel()

		if catalogId, err = primitive.ObjectIDFromHex(opts.ProductID); err != nil {
			ii.Logger.Err(err).Msgf("invalid catalog id for: %s", opts.ProductID)
			continue
		}

		if variantId, err = primitive.ObjectIDFromHex(opts.VariantID); err != nil {
			ii.Logger.Err(err).Msgf("invalid variant id for: %s", opts.VariantID)
			continue
		}

		findQuery := bson.M{
			"catalog_id": catalogId,
			"variant_id": variantId,
		}

		var updateQuery bson.D
		if opts.Unit != 0 {
			updateQuery = append(updateQuery, bson.E{
				Key: "$set",
				Value: bson.M{
					"status": model.InventoryStatus{
						Value:     model.InStockStatus,
						CreatedAt: time.Now().UTC(),
					},
					"unit_in_stock": opts.Unit,
				},
			})
		} else {
			updateQuery = append(updateQuery, bson.E{
				Key: "$set",
				Value: bson.M{
					"status": model.InventoryStatus{
						Value:     model.OutOfStockStatus,
						CreatedAt: time.Now().UTC(),
					},
					"unit_in_stock": 0,
				},
			})
		}
		operation.SetFilter(findQuery)
		operation.SetUpdate(updateQuery)
		operations = append(operations, operation)
	}

	if len(operations) == 0 {
		ii.Logger.Info().Msgf("no operations")
		return errors.Errorf("no products in the request to update")
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	resp, err := ii.DB.Collection(model.InventoryColl).BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		ii.Logger.Err(err).Msgf("failed to update inventory")
		return errors.Wrap(err, "failed to update inventory")
	}
	fmt.Printf("RESP: %+v\n", resp)
	return nil
}

func (ii *InventoryImpl) UpdateInventorybySKUs(opts []schema.UpdateInventoryBySKUOpt) (*schema.UpdateInventoryBySKUResp, error) {
	ctx := context.TODO()
	inStockStatus := model.InventoryStatus{
		Value:     model.InStockStatus,
		CreatedAt: time.Now(),
	}
	outOfStockStatus := model.InventoryStatus{
		Value:     model.OutOfStockStatus,
		CreatedAt: time.Now(),
	}
	skuUnitMap := make(map[string]int)

	var skus []string
	// var operations []mongo.WriteModel
	for _, uis := range opts {
		skus = append(skus, uis.SKU)
		skuUnitMap[uis.SKU] = uis.Unit
	}

	//pipeline 1 to find unique skus

	matchQuery := bson.D{{
		Key: "$match", Value: bson.M{
			"sku": bson.M{
				"$in": skus,
			},
		},
	}}

	groupQuery := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$sku",
			"inventory_id": bson.M{
				"$addToSet": "$_id",
			},
		},
	}}

	projectQuery := bson.D{{
		Key: "$project", Value: bson.M{

			"inventory_id": bson.M{"$first": "$inventory_id"},
			"is_unique": bson.M{"$eq": bson.A{
				bson.M{
					"$size": "$inventory_id",
				}, 1}},
		},
	}}

	pipeline := mongo.Pipeline{
		matchQuery,
		groupQuery,
		projectQuery,
	}
	cur, err := ii.DB.Collection(model.InventoryColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrapf(err, "error in aggregation")
	}
	var variantResp []model.InventorySearchBySKU
	if err := cur.All(ctx, &variantResp); err != nil {
		return nil, err
	}
	resp := schema.UpdateInventoryBySKUResp{}

	//finding invalid SKUs
	inValidSkuMap := make(map[string]bool)
	for _, s := range variantResp {
		inValidSkuMap[s.SKU] = true
	}

	for _, s := range opts {
		if !inValidSkuMap[s.SKU] {
			resp.InvalidSKUs = append(resp.InvalidSKUs, s.SKU)
		}
	}

	var operations []mongo.WriteModel

	for _, vr := range variantResp {
		if !vr.IsUnique {
			resp.DuplicateSKUs = append(resp.DuplicateSKUs, vr.SKU)
			continue
		}
		operation := mongo.NewUpdateOneModel()
		var status model.InventoryStatus
		operation.SetFilter(bson.M{"_id": vr.InventoryID})
		if skuUnitMap[vr.SKU] > 0 {
			status = inStockStatus
		} else {
			status = outOfStockStatus
			skuUnitMap[vr.SKU] = 0
		}
		operation.SetUpdate(bson.M{
			"$set": bson.M{
				"unit_in_stock": skuUnitMap[vr.SKU],
				"status":        status,
			},
		})

		operations = append(operations, operation)

	}
	if len(operations) == 0 {
		ii.Logger.Info().Msgf("no operations for skus")
		return &resp, nil
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	_, err = ii.DB.Collection(model.InventoryColl).BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		ii.Logger.Err(err).Msgf("failed to update inventory")
		return nil, errors.Wrapf(err, "failed to update inventory")
	}

	return &resp, nil
}
