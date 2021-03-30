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

// Discount service contains methods that handles catalog level discount
type Discount interface {
	CreateDiscount(*schema.CreateDiscountOpts) (*schema.CreateDiscountResp, error)
	DeactivateDiscount(id primitive.ObjectID) error
	CreateSale(*schema.CreateSaleOpts) (*schema.CreateSaleResp, error)
	EditSale(*schema.EditSaleOpts) (*schema.EditSaleResp, error)
	EditSaleStatus(*schema.EditSaleStatusOpts) error
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
		IsActive:    true,
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
	t_low := t.Add(time.Second * time.Duration(-1))
	t_high := t.Add(time.Second * time.Duration(1))

	err := di.activateDiscount(ctx, t_low, t_high)
	if err != nil {
		return err
	}
	err = di.deActivateDiscount(ctx, t_low, t_high)
	if err != nil {
		return err
	}
	return nil
}

func (di *DiscountImpl) activateDiscount(ctx context.Context, t_low, t_high time.Time) error {
	filter := bson.M{
		"valid_after": bson.M{
			"lte": t_high,
			"gte": t_low,
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

	for _, discount := range discounts {

		if discount.IsActive {
			continue
		}
		discount.IsActive = true
		updateDiscounts = append(updateDiscounts, discount.ID)
		update := bson.M{
			"$set": bson.M{
				"discount_id": discount.ID,
			},
		}
		res, err := di.DB.Collection(model.CatalogColl).UpdateOne(ctx, bson.M{"_id": discount.CatalogID}, update)
		if err != nil {
			di.Logger.Err(err)
			continue
		}
		if res.MatchedCount == 0 {
			di.Logger.Err(errors.Errorf("catalog with id: %s not found", discount.CatalogID))
		}

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
	res, err := di.DB.Collection(model.CatalogColl).UpdateMany(ctx, filterQuery, updateQuery)
	if err != nil {
		di.Logger.Log().Err(err)
		return err
	}
	if res.MatchedCount != int64(len(updateDiscounts)) {
		err := errors.Errorf("%d discount ids did not match", int64(len(updateDiscounts))-res.MatchedCount)
		di.Logger.Log().Err(err)
		return err
	}
	return nil
}

func (di *DiscountImpl) deActivateDiscount(ctx context.Context, t_low, t_high time.Time) error {
	filter := bson.M{
		"valid_after": bson.M{
			"lte": t_high,
			"gte": t_low,
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
	var updateDiscounts []primitive.ObjectID

	for _, discount := range discounts {

		if discount.IsActive {
			continue
		}
		discount.IsActive = true
		updateDiscounts = append(updateDiscounts, discount.ID)
		update := bson.M{
			"$set": bson.M{
				"discount_id": primitive.NilObjectID,
			},
		}
		res, err := di.DB.Collection(model.CatalogColl).UpdateOne(ctx, bson.M{"_id": discount.CatalogID}, update)
		if err != nil {
			di.Logger.Err(err)
			continue
		}
		if res.MatchedCount == 0 {
			di.Logger.Err(errors.Errorf("catalog with id: %s not found", discount.CatalogID))
		}

	}
	filterQuery := bson.M{
		"_id": bson.M{
			"$in": updateDiscounts,
		},
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}
	if len(updateDiscounts) == 0 {
		return nil
	}
	res, err := di.DB.Collection(model.CatalogColl).UpdateMany(ctx, filterQuery, updateQuery)
	if err != nil {
		di.Logger.Log().Err(err)
		return err
	}
	if res.MatchedCount != int64(len(updateDiscounts)) {
		err := errors.Errorf("%d discount ids did not match", int64(len(updateDiscounts))-res.MatchedCount)
		di.Logger.Log().Err(err)
		return err
	}
	return nil
}
