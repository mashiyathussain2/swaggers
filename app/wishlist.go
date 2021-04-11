package app

import (
	"context"
	"encoding/json"
	"go-app/model"
	"go-app/schema"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Wishlist contains methods for Wishlist service functionality
type Wishlist interface {
	AddToWishlist(*schema.AddToWishlistOpts) error
	RemoveFromWishlist(*schema.RemoveFromWishlistOpts) error
	GetWishlist(primitive.ObjectID) ([]schema.GetWishlistResp, error)
	GetWishlistMap(primitive.ObjectID) (map[string]schema.CatalogWishListinfo, error)
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

func (wi *WishlistImpl) GetWishlist(id primitive.ObjectID) ([]schema.GetWishlistResp, error) {

	ctx := context.TODO()
	wishlistResp := []schema.GetWishlistResp{}
	var wishlist model.Wishlist

	err := wi.DB.Collection(model.WishlistColl).FindOne(ctx, bson.M{"user_id": id}).Decode(&wishlist)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Errorf("unable to find wishlist for user with id: %s", id.Hex())
		}
		return nil, errors.Wrapf(err, "unable to query for document")
	}

	for _, cat := range wishlist.CatalogIDS {
		var s model.GetAllCatalogInfoResp

		url := wi.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + cat.Hex()
		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to fetch catlog data")
		}
		defer resp.Body.Close()

		//Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			wi.Logger.Err(err).Msgf("failed to read response from api %s", url)
			return nil, errors.Wrap(err, "failed to get catalog info")
		}
		if err := json.Unmarshal(body, &s); err != nil {
			wi.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
			return nil, errors.Wrap(err, "failed to decode body into struct")
		}
		if !s.Success {
			wi.Logger.Err(errors.New("success false from catalog")).Str("body", string(body)).Msg("got success false response from catalog")
			return nil, errors.New("got success false response from catalog")
		}

		wishlistItem := schema.GetWishlistResp{
			CatalogID: cat,
			CatalogInfo: schema.CatalogWishListinfo{
				ID:            cat,
				Name:          s.Payload.Name,
				FeaturedImage: s.Payload.FeaturedImage,

				BasePrice:   s.Payload.BasePrice,
				RetailPrice: s.Payload.RetailPrice,

				Status: s.Payload.Status,

				DiscountInfo: s.Payload.DiscountInfo,
				BrandInfo:    s.Payload.BrandInfo,
			},
		}
		wishlistResp = append(wishlistResp, wishlistItem)
	}

	return wishlistResp, nil
}

func (wi *WishlistImpl) GetWishlistMap(id primitive.ObjectID) (map[string]bool, error) {

	ctx := context.TODO()
	var wishlist model.Wishlist
	mapWish := make(map[string]bool)

	err := wi.DB.Collection(model.WishlistColl).FindOne(ctx, bson.M{"user_id": id}).Decode(&wishlist)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Errorf("unable to find wishlist for user with id: %s", id.Hex())
		}
		return nil, errors.Wrapf(err, "unable to query for document")
	}

	for _, cat := range wishlist.CatalogIDS {
		mapWish[cat.Hex()] = true
	}

	return mapWish, nil
}
