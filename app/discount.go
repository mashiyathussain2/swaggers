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
	s := model.Sale{
		Name: opts.Name,
		Slug: UniqueSlug(opts.Name),
		Banner: model.IMG{
			SRC: opts.Banner.SRC,
		},
		ValidAfter:  opts.ValidAfter,
		ValidBefore: opts.ValidBefore,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.Banner.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "failed to load banner image")
	}

	res, err := di.DB.Collection(model.DiscountColl).InsertOne(context.TODO(), s)
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
