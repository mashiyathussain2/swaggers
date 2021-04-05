package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Wishlist contains methods for Wishlist service functionality
type Wishlist interface {
	AddToWishlist(*schema.AddToWishlistOpts) error
	RemoveFromWishlist(opts *schema.RemoveFromWishlistOpts) error
}

// WishlistImpl implements Wishlist interface methods
type WishlistImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// WishlistImplOpts contains args required to create
type WishlistImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitWishlist returns new instance of Wishlist implementation
func InitWishlist(opts *WishlistImplOpts) Wishlist {
	wi := WishlistImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &wi
}

func (wi *WishlistImpl) AddToWishlist(opts *schema.AddToWishlistOpts) error {

	findQuery := bson.M{
		"user_id": opts.UserID,
	}
	updateQuery := bson.M{
		"$addToSet": bson.M{
			"catalog_ids": opts.CatalogID,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}
	var wishlist model.Wishlist
	opt := options.FindOneAndUpdate().SetUpsert(true)
	opt.SetReturnDocument(options.After)
	err := wi.DB.Collection(model.WishlistColl).FindOneAndUpdate(context.TODO(), findQuery, updateQuery, opt).Decode(&wishlist)
	if err != nil {
		return errors.Wrapf(err, "unable to add catalog with id: %s to wishlist", opts.CatalogID.Hex())
	}
	return nil
}

func (wi *WishlistImpl) RemoveFromWishlist(opts *schema.RemoveFromWishlistOpts) error {

	findQuery := bson.M{
		"user_id": opts.UserID,
	}
	updateQuery := bson.M{
		"$pull": bson.M{
			"catalog_ids": opts.CatalogID,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}
	var wishlist model.Wishlist
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := wi.DB.Collection(model.WishlistColl).FindOneAndUpdate(context.TODO(), findQuery, updateQuery, opt).Decode(&wishlist)
	if err != nil {
		return errors.Wrapf(err, "unable to add catalog with id: %s to wishlist", opts.CatalogID.Hex())
	}
	return nil
}
