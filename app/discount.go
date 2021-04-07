package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Discount service contains methods that handles catalog level discount
type Discount interface {
	CreateDiscount(*schema.CreateDiscountOpts) (*schema.CreateDiscountResp, error)
	DeactivateDiscount(id primitive.ObjectID) error
	CreateSale(*schema.CreateSaleOpts) (*schema.CreateSaleResp, error)
	EditSale(*schema.EditSaleOpts) (*schema.EditSaleResp, error)
	EditSaleStatus(*schema.EditSaleStatusOpts) error

	GetActiveDiscountByCatalogID(catalogID primitive.ObjectID) (*schema.DiscountInfoResp, error)

	GetSales(*schema.GetSalesOpts) ([]schema.GetSalesResp, error)
	GetDiscountAndCatalogInfoBySaleID(primitive.ObjectID) ([]schema.DiscountInfoWithCatalogInfoResp, error)
}

// DiscountImpl implements Discount service methods
type DiscountImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// DiscountOpts contains args required to create a new instance of discount service
type DiscountOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitDiscount returns a new instance of discount service implementation
func InitDiscount(opts *DiscountOpts) Discount {
	return &DiscountImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

// CreateDiscount creates a new discount for a catalog.
// Note only 1 discount can be active for a catalog at a given time
func (di *DiscountImpl) CreateDiscount(opts *schema.CreateDiscountOpts) (*schema.CreateDiscountResp, error) {
	ctx := context.TODO()
	if err := di.validateCreateDiscount(ctx, opts); err != nil {
		return nil, err
	}
	d := model.Discount{
		Type:        opts.Type,
		CatalogID:   opts.CatalogID,
		VariantsID:  opts.VariantsID,
		SaleID:      opts.SaleID,
		IsActive:    false,
		Value:       opts.Value,
		MaxValue:    opts.MaxValue,
		ValidAfter:  opts.ValidAfter,
		ValidBefore: opts.ValidBefore,
		CreatedAt:   time.Now().UTC(),
	}

	res, err := di.DB.Collection(model.DiscountColl).InsertOne(ctx, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert discount")
	}
	return &schema.CreateDiscountResp{
		ID:          res.InsertedID.(primitive.ObjectID),
		CatalogID:   d.CatalogID,
		VariantsID:  d.VariantsID,
		SaleID:      d.SaleID,
		IsActive:    d.IsActive,
		Type:        d.Type,
		Value:       d.Value,
		MaxValue:    d.MaxValue,
		ValidAfter:  d.ValidAfter,
		ValidBefore: d.ValidBefore,
		CreatedAt:   d.CreatedAt,
	}, nil
}

func (di *DiscountImpl) validateCreateDiscount(ctx context.Context, opts *schema.CreateDiscountOpts) error {
	var sale model.Sale

	if !opts.SaleID.IsZero() {
		filter := bson.M{"_id": opts.SaleID}
		if err := di.DB.Collection(model.SaleColl).FindOne(ctx, filter).Decode(&sale); err != nil {
			if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
				return errors.Errorf("sale with id:%s not found", opts.SaleID.Hex())
			}
			return errors.Wrap(err, "failed to query for sale")
		}

		// Setting up sale validAfter & validBefore date
		opts.ValidAfter = sale.ValidAfter
		opts.ValidBefore = sale.ValidBefore
	}

	filter := bson.M{
		"catalog_id": opts.CatalogID,
		"is_active":  true,
		"$or": bson.A{
			bson.M{
				"valid_after": bson.M{
					"$lte": opts.ValidAfter,
				},
				"valid_before": bson.M{
					"$gte": opts.ValidAfter,
				},
			},
			bson.M{
				"valid_after": bson.M{
					"$gte": opts.ValidAfter,
				},
				"valid_before": bson.M{
					"$lte": opts.ValidBefore,
				},
			},
			bson.M{
				"valid_after": bson.M{
					"$lte": opts.ValidBefore,
				},
				"valid_before": bson.M{
					"$gte": opts.ValidBefore,
				},
			},
		},
	}
	count, err := di.DB.Collection(model.DiscountColl).CountDocuments(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "failed to query discounts in catalog")
	}
	if count != 0 {
		return errors.Errorf("discount from %s to %s already exists for catalog id: %s", opts.ValidAfter.String(), opts.ValidBefore.String(), opts.CatalogID.Hex())
	}

	// TODO: Add validation to restrict catalog price (after discount applied) should not go below zero

	return nil
}

// DeactivateDiscount sets IsActive field to false
func (di *DiscountImpl) DeactivateDiscount(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_active": false}}

	res, err := di.DB.Collection(model.DiscountColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to update discount")
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("discount id: %s not found", id.Hex())
	}
	if res.ModifiedCount == 0 {
		return errors.Errorf("failed to update discount id: %s", id.Hex())
	}
	return nil
}

// CreateSale creates a new sale in db
func (di *DiscountImpl) CreateSale(opts *schema.CreateSaleOpts) (*schema.CreateSaleResp, error) {
	//TODO: Add status
	s := model.Sale{
		Name: opts.Name,
		Slug: UniqueSlug(opts.Name),
		Banner: &model.IMG{
			SRC: opts.Banner.SRC,
		},
		Genders:     opts.Genders,
		Status:      model.Schedule,
		ValidAfter:  opts.ValidAfter,
		ValidBefore: opts.ValidBefore,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.Banner.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "failed to load banner image")
	}

	res, err := di.DB.Collection(model.SaleColl).InsertOne(context.TODO(), s)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert sale object")
	}

	return &schema.CreateSaleResp{
		ID:          res.InsertedID.(primitive.ObjectID),
		Name:        s.Name,
		Slug:        s.Slug,
		Banner:      s.Banner,
		Genders:     s.Genders,
		ValidAfter:  s.ValidAfter,
		ValidBefore: s.ValidBefore,
		CreatedAt:   s.CreatedAt,
	}, nil
}

//EditSale edits Sale info such as Name and Banner Image
func (di *DiscountImpl) EditSale(opts *schema.EditSaleOpts) (*schema.EditSaleResp, error) {

	ctx := context.TODO()
	sale := model.Sale{}
	findQuery := bson.M{"_id": opts.ID}
	t := time.Now().UTC()

	updateData := bson.D{}

	err := di.DB.Collection(model.SaleColl).FindOne(ctx, findQuery).Decode(&sale)
	if err != nil {
		return nil, errors.Errorf("unable to find the sale with id: %s", opts.ID.Hex())
	}

	if t.After(sale.ValidAfter) && t.Before(sale.ValidBefore) {
		return nil, errors.Errorf("cannot edit the sale, since sale is already live")
	}

	if t.After(sale.ValidBefore) {
		return nil, errors.Errorf("cannot edit the sale, since sale is already finished")
	}

	if opts.Name != "" {
		sale.Name = opts.Name
		updateData = append(updateData, bson.E{
			Key: "name", Value: opts.Name,
		})
	}
	if opts.Banner != nil {
		banner := model.IMG{SRC: opts.Banner.SRC}
		err := banner.LoadFromURL()
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load banner image")
		}
		updateData = append(updateData, bson.E{
			Key: "banner", Value: banner,
		})
		sale.Banner = &banner
	}

	if len(opts.Genders) != 0 {
		sale.Genders = opts.Genders
		updateData = append(updateData, bson.E{
			Key: "genders", Value: opts.Genders,
		})
	}

	updateData = append(updateData, bson.E{
		Key: "updated_at", Value: t,
	})
	sale.UpdatedAt = t
	updateQuery := bson.M{"$set": updateData}

	res, err := di.DB.Collection(model.SaleColl).UpdateOne(ctx, findQuery, updateQuery)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to update the sale with id: %s", opts.ID.Hex())
	}

	if res.ModifiedCount == 0 {
		return nil, errors.Errorf("unable to update the sale with id: %s", opts.ID.Hex())
	}
	return &schema.EditSaleResp{
		ID:          sale.ID,
		Name:        sale.Name,
		Slug:        sale.Slug,
		Banner:      sale.Banner,
		Genders:     sale.Genders,
		ValidAfter:  sale.ValidAfter,
		ValidBefore: sale.ValidBefore,
		CreatedAt:   sale.CreatedAt,
		UpdatedAt:   sale.UpdatedAt,
	}, nil
}

//EditSaleStatus deletes the scheduled sale from the database
func (di *DiscountImpl) EditSaleStatus(opts *schema.EditSaleStatusOpts) error {
	filter := bson.M{
		"_id": opts.ID,
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"status": opts.Status,
		},
	}
	res, err := di.DB.Collection(model.SaleColl).UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to update sale status with id %s", opts.ID.Hex())
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find the sale with id: %s", opts.ID.Hex())
	}
	if res.ModifiedCount == 0 {
		return errors.Wrapf(err, "unable to update sale status with id %s", opts.ID.Hex())
	}
	return nil
}

//CheckAndUpdateStatus checks discount to be activated/deactivated and updates status
func (di *DiscountImpl) CheckAndUpdateStatus() error {

	ctx := context.TODO()
	t := time.Now().UTC()
	fmt.Println(t)
	err := di.activateDiscount(ctx, t)
	if err != nil {
		return err
	}
	err = di.deActivateDiscount(ctx, t)
	if err != nil {
		return err
	}
	return nil
}

func (di *DiscountImpl) activateDiscount(ctx context.Context, t time.Time) error {
	filter := bson.M{
		"valid_after": bson.M{
			"$lte": t,
		},
		"valid_before": bson.M{
			"$gte": t,
		},
		"is_active": false,
	}
	var discounts []model.Discount

	cur, err := di.DB.Collection(model.DiscountColl).Find(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "unable to query for discounts")
	}
	if err := cur.All(ctx, &discounts); err != nil {
		return errors.Wrap(err, "error decoding Catalogs")
	}

	var updateDiscounts []primitive.ObjectID

	var operations []mongo.WriteModel

	for _, discount := range discounts {

		discount.IsActive = true
		updateDiscounts = append(updateDiscounts, discount.ID)
		update := bson.M{
			"$set": bson.M{
				"discount_id": discount.ID,
			},
		}
		filterCat := bson.M{"_id": discount.CatalogID}

		operation := mongo.NewUpdateOneModel()
		operation.SetUpdate(update)
		operation.SetFilter(filterCat)

		operations = append(operations, operation)

		// res, err := di.DB.Collection(model.CatalogColl).UpdateOne(ctx, bson.M{"_id": discount.CatalogID}, update)
		// if err != nil {
		// 	di.Logger.Err(err)
		// 	continue
		// }
		// if res.MatchedCount == 0 {
		// 	di.Logger.Err(errors.Errorf("catalog with id: %s not found", discount.CatalogID))
		// }

	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)

	_, err = di.DB.Collection(model.CatalogColl).BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		log.Fatal(err)
		return err
	}
	filterQuery := bson.M{
		"_id": bson.M{
			"$in": updateDiscounts,
		},
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"is_active": true,
		},
	}

	if len(updateDiscounts) == 0 {
		return nil
	}

	res, err := di.DB.Collection(model.DiscountColl).UpdateMany(ctx, filterQuery, updateQuery)
	if err != nil {
		di.Logger.Log().Err(err)
		return err
	}
	if res.MatchedCount != int64(len(updateDiscounts)) {
		err := errors.Errorf("%d discount ids did not match for activating", int64(len(updateDiscounts))-res.MatchedCount)
		di.Logger.Log().Err(err)
		return err
	}
	return nil
}

func (di *DiscountImpl) deActivateDiscount(ctx context.Context, t time.Time) error {
	filter := bson.M{
		"valid_before": bson.M{
			"$lte": t,
		},
		"is_active": true,
	}

	var discounts []model.Discount
	cur, err := di.DB.Collection(model.DiscountColl).Find(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "unable to query for discounts")
	}
	if err := cur.All(ctx, &discounts); err != nil {
		return errors.Wrap(err, "error decoding Catalogs")
	}
	var updateDiscountIDs []primitive.ObjectID
	var updateCatalogIds []primitive.ObjectID
	for _, discount := range discounts {
		updateDiscountIDs = append(updateDiscountIDs, discount.ID)
		updateCatalogIds = append(updateCatalogIds, discount.CatalogID)
	}

	catlogFilterQuery := bson.M{
		"_id": bson.M{
			"$in": updateCatalogIds,
		},
	}
	catalogUpdateQuery := bson.M{
		"$unset": bson.M{
			"discount_id": 1,
		},
	}

	res, err := di.DB.Collection(model.CatalogColl).UpdateMany(ctx, catlogFilterQuery, catalogUpdateQuery)
	if err != nil {
		di.Logger.Log().Err(err)
		return err
	}
	if res.MatchedCount != int64(len(updateCatalogIds)) {
		err := errors.Errorf("%d catalog ids did not match", int64(len(updateCatalogIds))-res.MatchedCount)
		di.Logger.Log().Err(err)
		return err
	}

	filterQuery := bson.M{
		"_id": bson.M{
			"$in": updateDiscountIDs,
		},
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}

	if len(updateDiscountIDs) == 0 {
		return nil
	}
	res, err = di.DB.Collection(model.DiscountColl).UpdateMany(ctx, filterQuery, updateQuery)
	if err != nil {
		di.Logger.Log().Err(err)
		return err
	}
	if res.MatchedCount != int64(len(updateDiscountIDs)) {
		err := errors.Errorf("%d discount ids did not match in deactivate account", int64(len(updateDiscountIDs))-res.MatchedCount)
		di.Logger.Log().Err(err)
		return err
	}
	return nil
}

func (di *DiscountImpl) GetActiveDiscountByCatalogID(catalogID primitive.ObjectID) (*schema.DiscountInfoResp, error) {
	var s schema.DiscountInfoResp

	filter := bson.M{
		"catalog_id": catalogID,
		"is_active":  true,
	}
	if err := di.DB.Collection(model.CollectionColl).FindOne(context.TODO(), filter).Decode(&s); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "active discount with catalogID:%s not found", catalogID.Hex())
		}
		return nil, errors.Wrapf(err, "failed to find active discount for catalogID:%s", catalogID.Hex())
	}

	return &s, nil
}

func (di *DiscountImpl) GetSales(opts *schema.GetSalesOpts) ([]schema.GetSalesResp, error) {
	var resp []schema.GetSalesResp
	ctx := context.TODO()
	queryOpts := options.Find().SetSkip(int64(opts.Page) * 20).SetLimit(20)
	cur, err := di.DB.Collection(model.SaleColl).Find(ctx, bson.M{}, queryOpts)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to find sales")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to find sales")
	}
	return resp, nil
}

func (di *DiscountImpl) GetDiscountAndCatalogInfoBySaleID(id primitive.ObjectID) ([]schema.DiscountInfoWithCatalogInfoResp, error) {
	var resp []schema.DiscountInfoWithCatalogInfoResp
	ctx := context.TODO()
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"sale_id": id,
			},
		},
	}
	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         model.CatalogColl,
				"localField":   "catalog_id",
				"foreignField": "_id",
				"as":           "catalog_info",
			},
		},
	}

	setStage := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"$arrayElemAt": bson.A{
					"$catalog_info",
					0,
				},
			},
		},
	}

	cur, err := di.DB.Collection(model.DiscountColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, setStage})
	if err != nil {
		return nil, errors.Wrap(err, "query failed to find discount items")
	}

	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to find discount items")
	}

	return resp, nil
}
