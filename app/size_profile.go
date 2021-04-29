package app

import (
	"context"
	"go-app/model"
	"go-app/schema"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SizeProfile interface {
	CreateSizeProfile(opts *schema.CreateSizeProfileOpts) (primitive.ObjectID, error)
	GetSizeProfile(id primitive.ObjectID) (*schema.GetSizeProfileResp, error)
	GetAllSizeProfiles() ([]schema.GetAllSizeProfilesResp, error)
	AddBrandToSizeProfile(*schema.AddBrandToSizeProfileOpts) error
	GetSizeProfilesForBrand(brandID primitive.ObjectID) ([]schema.GetSizeProfileForBrandResp, error)
}

// SizeProfileImpl implements SizeProfile interface methods
type SizeProfileImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// SizeProfileImplOpts contains args required to create
type SizeProfileImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitSizeProfile returns new instance of SizeProfile implementation
func InitSizeProfile(opts *SizeProfileImplOpts) SizeProfile {
	ui := SizeProfileImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ui
}

func (sp *SizeProfileImpl) CreateSizeProfile(opts *schema.CreateSizeProfileOpts) (primitive.ObjectID, error) {

	ctx := context.TODO()
	sizeProfile := model.SizeProfile{
		Name:  opts.Name,
		Specs: opts.Specs,
	}
	res, err := sp.DB.Collection(model.SizeProfileColl).InsertOne(ctx, sizeProfile)
	if err != nil {
		return primitive.NilObjectID, errors.Wrapf(err, "unable to create new Size profile")
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (sp *SizeProfileImpl) GetSizeProfile(id primitive.ObjectID) (*schema.GetSizeProfileResp, error) {

	ctx := context.TODO()
	findQuery := bson.M{
		"_id": id,
	}
	var sizeProfile schema.GetSizeProfileResp

	err := sp.DB.Collection(model.SizeProfileColl).FindOne(ctx, findQuery).Decode(&sizeProfile)
	if err != nil {
		sp.Logger.Err(err).Interface("size profile id", id).Msg("failed to get size profile  with id")
		return nil, errors.Wrapf(err, "failed to get size profile for brand with id %s", id.Hex())
	}
	return &sizeProfile, nil
}

func (sp *SizeProfileImpl) GetSizeProfileKeeper(id primitive.ObjectID) (*schema.GetSizeProfileResp, error) {

	ctx := context.TODO()
	findQuery := bson.M{
		"_id": id,
	}
	var sizeProfile schema.GetSizeProfileResp

	err := sp.DB.Collection(model.SizeProfileColl).FindOne(ctx, findQuery).Decode(&sizeProfile)
	if err != nil {
		sp.Logger.Err(err).Interface("size profile id", id).Msg("failed to get size profile  with id")
		return nil, errors.Wrapf(err, "failed to get size profile for brand with id %s", id.Hex())
	}
	return &sizeProfile, nil
}

func (sp *SizeProfileImpl) GetAllSizeProfiles() ([]schema.GetAllSizeProfilesResp, error) {

	ctx := context.TODO()
	findQuery := bson.M{}

	cur, err := sp.DB.Collection(model.SizeProfileColl).Find(ctx, findQuery)
	if err != nil {
		sp.Logger.Err(err).Msg("failed to get size profiles")
		return nil, errors.Wrapf(err, "failed to get size profiles")
	}
	var sizeProfiles []schema.GetAllSizeProfilesResp
	if err := cur.All(ctx, &sizeProfiles); err != nil {
		return nil, errors.Wrap(err, "failed to find brands")
	}
	return sizeProfiles, nil
}

func (sp *SizeProfileImpl) AddBrandToSizeProfile(opts *schema.AddBrandToSizeProfileOpts) error {
	findQuery := bson.M{
		"_id": bson.M{
			"$in": opts.IDs,
		},
	}
	updateQuery := bson.M{
		"$addToSet": bson.M{
			"brand_ids": opts.BrandID,
		},
	}
	res, err := sp.DB.Collection(model.SizeProfileColl).UpdateMany(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "error linking Size Profile to Brand")
	}
	if int(res.MatchedCount) != len(opts.IDs) {
		return errors.Errorf("error linking Size Profile to Brand for %d brands ", len(opts.IDs)-int(res.MatchedCount))
	}
	return nil
}

func (sp *SizeProfileImpl) GetSizeProfilesForBrand(brandID primitive.ObjectID) ([]schema.GetSizeProfileForBrandResp, error) {

	ctx := context.TODO()
	findQuery := bson.M{
		"brand_ids": brandID,
	}

	cur, err := sp.DB.Collection(model.SizeProfileColl).Find(ctx, findQuery)
	if err != nil {
		sp.Logger.Err(err).Interface("brand id", brandID).Msg("failed to get size profile for brand with id")
		return nil, errors.Wrapf(err, "failed to get size profile for brand with id %s", brandID.Hex())
	}
	var sizeProfiles []schema.GetSizeProfileForBrandResp
	if err := cur.All(ctx, &sizeProfiles); err != nil {
		return nil, errors.Wrap(err, "failed to find brands")
	}
	return sizeProfiles, nil
}
