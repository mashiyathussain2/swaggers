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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InfluencerProducts service allows `InfluencerProducts` to execute admin operations.
type InfluencerProducts interface {
	AddInfluencerProductsOpts(opts *schema.AddInfluencerProductsOpts) error
	RemoveInfluencerProductsOpts(opts *schema.RemoveInfluencerProductsOpts) error
	GetInfluencerProductsOpts(id string) (*schema.GetInfluencerProductESResp, error)
}

// InfluencerProductsImpl implements InfluencerProducts related operations
type InfluencerProductsImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InfluencerProductsOpts contains arguments required to create a new instance of InfluencerProducts
type InfluencerProductsOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

//InitInfluencerProducts returns InfluencerProducts instance
func InitInfluencerProducts(opts *InfluencerProductsOpts) InfluencerProducts {
	return &InfluencerProductsImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

func (ip *InfluencerProductsImpl) AddInfluencerProductsOpts(opts *schema.AddInfluencerProductsOpts) error {
	ctx := context.TODO()
	filter := bson.M{
		"influencer_id": opts.InfluencerID,
	}
	update := bson.M{
		"$set": bson.M{
			"influencer_id": opts.InfluencerID,
			"updated_at":    time.Now(),
		},
		"$addToSet": bson.M{
			"catalog_ids": bson.M{
				"$each": opts.CatalogIDs,
			},
		},
	}
	option := options.Update().SetUpsert(true)

	_, err := ip.DB.Collection(model.InfluencerProductColl).UpdateOne(ctx, filter, update, option)
	if err != nil {
		return errors.Wrapf(err, "error updating influencer product")
	}
	return nil
}

func (ip *InfluencerProductsImpl) RemoveInfluencerProductsOpts(opts *schema.RemoveInfluencerProductsOpts) error {
	ctx := context.TODO()
	filter := bson.M{
		"influencer_id": opts.InfluencerID,
	}
	update := bson.M{
		"$pull": bson.M{
			"catalog_ids": bson.M{
				"$in": opts.CatalogIDs,
			},
		},
	}
	_, err := ip.DB.Collection(model.InfluencerProductColl).UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrapf(err, "error updating influencer product")
	}
	return nil
}

func (ip *InfluencerProductsImpl) GetInfluencerProductsOpts(id string) (*schema.GetInfluencerProductESResp, error) {
	ctx := context.TODO()
	iid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrapf(err, "error converting influncer ids")
	}
	filter := bson.M{
		"influencer_id": iid,
	}
	var res schema.GetInfluencerProductESResp
	err = ip.DB.Collection(model.InfluencerProductColl).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, errors.Wrapf(err, "error updating influencer product")
	}
	return &res, err
}
