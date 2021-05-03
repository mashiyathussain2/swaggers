package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Brand contains methods for brand service functionality
type Brand interface {
	CreateBrand(*schema.CreateBrandOpts) (*schema.CreateBrandResp, error)
	EditBrand(*schema.EditBrandOpts) (*schema.EditBrandResp, error)
	GetBrandByID(primitive.ObjectID) (*schema.GetBrandResp, error)
	CheckBrandByID(primitive.ObjectID) (bool, error)
	GetBrandsByID([]primitive.ObjectID) ([]schema.GetBrandResp, error)
	GetBrands() ([]schema.GetBrandResp, error)

	AddFollower(opts *schema.AddBrandFollowerOpts) (bool, error)
	RemoveFollower(opts *schema.AddBrandFollowerOpts) (bool, error)
}

// BrandImpl implements brand interface methods
type BrandImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// BrandImplOpts contains args required to create
type BrandImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitBrand returns new instance of brand implementation
func InitBrand(opts *BrandImplOpts) Brand {
	ui := BrandImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ui
}

// CreateBrand inserts a new brand document into collection
// Note: this method only creates brand profile not brand authentication
func (bi *BrandImpl) CreateBrand(opts *schema.CreateBrandOpts) (*schema.CreateBrandResp, error) {
	b := model.Brand{
		Name:               opts.Name,
		LName:              strings.ToLower(opts.Name),
		RegisteredName:     opts.RegisteredName,
		FulfillmentEmail:   opts.FulfillmentEmail,
		FulfillmentCCEmail: opts.FulfillmentCCEmail,
		Domain:             opts.Domain,
		Website:            opts.Website,
		Logo: &model.IMG{
			SRC: opts.Logo.SRC,
		},
		Bio: opts.Bio,
		CoverImg: &model.IMG{
			SRC: opts.CoverImg.SRC,
		},
		CreatedAt: time.Now().UTC(),
	}
	if err := b.Logo.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "invalid image for brand logo")
	}
	if err := b.CoverImg.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "invalid image for brand cover")
	}
	if opts.SocialAccount != nil {
		b.SocialAccount = &model.SocialAccount{}
		if opts.SocialAccount.Facebook != nil {
			b.SocialAccount.Facebook = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Facebook.FollowersCount)}
		}
		if opts.SocialAccount.Instagram != nil {
			b.SocialAccount.Instagram = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Instagram.FollowersCount)}
		}
		if opts.SocialAccount.Youtube != nil {
			b.SocialAccount.Youtube = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Youtube.FollowersCount)}
		}
		if opts.SocialAccount.Twitter != nil {
			b.SocialAccount.Twitter = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Twitter.FollowersCount)}
		}
	}
	if len(opts.SizeProfiles) > 0 {
		b.SizeProfiles = opts.SizeProfiles
	}

	res, err := bi.DB.Collection(model.BrandColl).InsertOne(context.TODO(), b)
	if err != nil {
		bi.Logger.Err(err).Interface("opts", opts).Msg("failed to insert brand")
		return nil, errors.Wrap(err, "failed to create brand")
	}

	err = bi.App.SizeProfile.AddBrandToSizeProfile(&schema.AddBrandToSizeProfileOpts{
		IDs:     opts.SizeProfiles,
		BrandID: res.InsertedID.(primitive.ObjectID),
	})
	if err != nil {
		bi.Logger.Err(err).Interface("opts", opts).Msg("failed to link size profiles with brand id")
	}
	var sp []schema.GetSizeProfileForBrandResp
	sp, err = bi.App.SizeProfile.GetSizeProfilesForBrand(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		bi.Logger.Err(err).Interface("opts", opts).Msg("failed to get size profiles for brand with id")
	}

	resp := schema.CreateBrandResp{
		ID:                 res.InsertedID.(primitive.ObjectID),
		Name:               b.Name,
		RegisteredName:     b.RegisteredName,
		FulfillmentEmail:   b.FulfillmentEmail,
		FulfillmentCCEmail: b.FulfillmentCCEmail,
		Domain:             b.Domain,
		Website:            b.Website,
		Logo:               b.Logo,
		CoverImg:           b.CoverImg,
		SocialAccount:      b.SocialAccount,
		Bio:                b.Bio,
		CreatedAt:          b.CreatedAt,
		SizeProfiles:       sp,
	}
	return &resp, nil
}

func (bi *BrandImpl) EditBrand(opts *schema.EditBrandOpts) (*schema.EditBrandResp, error) {
	var update bson.D
	if opts.Name != "" {
		update = append(update, bson.E{Key: "name", Value: opts.Name})
		update = append(update, bson.E{Key: "lname", Value: string(opts.Name)})
	}
	if opts.Domain != "" {
		update = append(update, bson.E{Key: "domain", Value: opts.Domain})
	}
	if opts.RegisteredName != "" {
		update = append(update, bson.E{Key: "registered_name", Value: opts.RegisteredName})
	}
	if opts.FulfillmentEmail != "" {
		update = append(update, bson.E{Key: "fulfillment_email", Value: opts.FulfillmentEmail})
	}
	if len(opts.FulfillmentCCEmail) != 0 {
		update = append(update, bson.E{Key: "fulfillment_cc_email", Value: opts.FulfillmentCCEmail})
	}
	if opts.Website != "" {
		update = append(update, bson.E{Key: "website", Value: opts.Website})
	}
	if opts.CoverImg != nil {
		img := model.IMG{SRC: opts.CoverImg.SRC}
		if err := img.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid image for brand cover")
		}
		update = append(update, bson.E{Key: "cover_img", Value: img})
	}
	if opts.Logo != nil {
		img := model.IMG{SRC: opts.Logo.SRC}
		if err := img.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid image for brand logo")
		}
		update = append(update, bson.E{Key: "logo", Value: img})
	}
	if opts.SocialAccount != nil {
		if opts.SocialAccount.Facebook != nil {
			update = append(update, bson.E{Key: "social_account.facebook.followers_count", Value: opts.SocialAccount.Facebook.FollowersCount})
		}
		if opts.SocialAccount.Instagram != nil {
			update = append(update, bson.E{Key: "social_account.instagram.followers_count", Value: opts.SocialAccount.Instagram.FollowersCount})
		}
		if opts.SocialAccount.Youtube != nil {
			update = append(update, bson.E{Key: "social_account.youtube.followers_count", Value: opts.SocialAccount.Youtube.FollowersCount})
		}
		if opts.SocialAccount.Twitter != nil {
			update = append(update, bson.E{Key: "social_account.twitter.followers_count", Value: opts.SocialAccount.Twitter.FollowersCount})
		}
	}
	if opts.Bio != "" {
		update = append(update, bson.E{Key: "bio", Value: opts.Bio})
	}
	if len(opts.SizeProfiles) != 0 {
		update = append(update, bson.E{Key: "size_profiles", Value: opts.SizeProfiles})
		bi.App.SizeProfile.AddBrandToSizeProfile(&schema.AddBrandToSizeProfileOpts{IDs: opts.SizeProfiles, BrandID: opts.ID})
	}
	update = append(update, bson.E{Key: "updated_at", Value: time.Now().UTC()})

	filterQuery := bson.M{"_id": opts.ID}
	updateQuery := bson.M{"$set": update}

	var brand model.Brand
	queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := bi.DB.Collection(model.BrandColl).FindOneAndUpdate(context.TODO(), filterQuery, updateQuery, queryOpts).Decode(&brand); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "brand with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrapf(err, "failed to update brand with id:%s", opts.ID.Hex())
	}
	var sp []schema.GetSizeProfileForBrandResp
	sp, err := bi.App.SizeProfile.GetSizeProfilesForBrand(brand.ID)
	if err != nil {
		bi.Logger.Err(err).Interface("opts", opts).Msg("failed to get size profiles for brand with id")
	}

	resp := schema.EditBrandResp{
		ID:                 brand.ID,
		Name:               brand.Name,
		RegisteredName:     brand.RegisteredName,
		FulfillmentEmail:   brand.FulfillmentEmail,
		FulfillmentCCEmail: brand.FulfillmentCCEmail,
		Domain:             brand.Domain,
		Website:            brand.Website,
		Logo:               brand.Logo,
		CoverImg:           brand.CoverImg,
		SocialAccount:      brand.SocialAccount,
		Bio:                brand.Bio,
		SizeProfiles:       sp,
		CreatedAt:          brand.CreatedAt,
		UpdatedAt:          brand.UpdatedAt,
	}
	return &resp, nil
}

// GetBrandByID returns brand info with matching id
func (bi *BrandImpl) GetBrandByID(id primitive.ObjectID) (*schema.GetBrandResp, error) {

	var resp []schema.GetBrandResp
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": id,
		},
	}}
	lookupStage := bson.D{{

		Key: "$lookup", Value: bson.M{
			"from":         model.SizeProfileColl,
			"localField":   "size_profiles",
			"foreignField": "_id",
			"as":           "size_profiles",
		},
	}}

	cursor, err := bi.DB.Collection(model.BrandColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get brand data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return &resp[0], nil
}

// CheckBrandByID check if brand exists with matching id
func (bi *BrandImpl) CheckBrandByID(id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}
	count, err := bi.DB.Collection(model.BrandColl).CountDocuments(context.TODO(), filter)
	if err != nil {
		bi.Logger.Err(err).Msgf("failed to check brand with id:%s", id.Hex())
		return false, errors.Wrapf(err, "failed to check brand with id:%s", id.Hex())
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

// GetBrandByID returns brand info with matching id
func (bi *BrandImpl) GetBrandsByID(ids []primitive.ObjectID) ([]schema.GetBrandResp, error) {
	var resp []schema.GetBrandResp
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": bson.M{"$in": ids},
		},
	}}
	lookupStage := bson.D{{

		Key: "$lookup", Value: bson.M{
			"from":         model.SizeProfileColl,
			"localField":   "size_profiles",
			"foreignField": "_id",
			"as":           "size_profiles",
		},
	}}

	cursor, err := bi.DB.Collection(model.BrandColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get brand data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return resp, nil
}

func (bi *BrandImpl) GetBrands() ([]schema.GetBrandResp, error) {
	var resp []schema.GetBrandResp
	ctx := context.TODO()

	lookupStage := bson.D{{

		Key: "$lookup", Value: bson.M{
			"from":         model.SizeProfileColl,
			"localField":   "size_profiles",
			"foreignField": "_id",
			"as":           "size_profiles",
		},
	}}

	cursor, err := bi.DB.Collection(model.BrandColl).Aggregate(ctx, mongo.Pipeline{
		lookupStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get brands data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return resp, nil
}

func (bi *BrandImpl) AddFollower(opts *schema.AddBrandFollowerOpts) (bool, error) {
	ctx := context.TODO()
	session, err := bi.DB.Client().StartSession()
	if err != nil {
		bi.Logger.Err(err).Msg("unable to create db session")
		return false, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		bi.Logger.Err(err).Msg("unable to start transaction")
		return false, errors.Wrap(err, "failed to add follower")
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		isFollowing, err := bi.DB.Collection(model.BrandColl).CountDocuments(sc, bson.M{"followers_id": opts.CustomerID})
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to check is user already follow brand")
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to follow brand")
		}

		if isFollowing != 0 {
			session.AbortTransaction(sc)
			return errors.New("user already follow the brand")
		}

		filter := bson.M{
			"_id": opts.BrandID,
		}
		update := bson.M{
			"$addToSet": bson.M{
				"followers_id": opts.CustomerID,
			},
			"$inc": bson.M{
				"followers_count": 1,
			},
		}

		res, err := bi.DB.Collection(model.BrandColl).UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed add follower")
			return errors.Wrap(err, "failed to add follower")
		}

		if res.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return errors.New("brand not found")
		}

		if err := bi.App.Customer.AddBrandFollowing(sc, opts); err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed add brand_id in customer following")
			return errors.Wrap(err, "failed to add follower")
		}

		if err := session.CommitTransaction(sc); err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to add follower")
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (bi *BrandImpl) RemoveFollower(opts *schema.AddBrandFollowerOpts) (bool, error) {
	ctx := context.TODO()
	session, err := bi.DB.Client().StartSession()
	if err != nil {
		bi.Logger.Err(err).Msg("unable to create db session")
		return false, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		bi.Logger.Err(err).Msg("unable to start transaction")
		return false, errors.Wrap(err, "failed to add follower")
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		isFollowing, err := bi.DB.Collection(model.BrandColl).CountDocuments(sc, bson.M{"followers_id": opts.CustomerID})
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to check is user already follow brand")
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to unfollow brand")
		}

		if isFollowing == 0 {
			session.AbortTransaction(sc)
			return errors.New("user does not follow the brand")
		}

		filter := bson.M{
			"_id": opts.BrandID,
		}
		update := bson.M{
			"$pull": bson.M{
				"followers_id": opts.CustomerID,
			},
			"$inc": bson.M{
				"followers_count": -1,
			},
		}

		res, err := bi.DB.Collection(model.BrandColl).UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed remove follower")
			return errors.Wrap(err, "failed to remove follower")
		}

		if res.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return errors.New("brand not found")
		}

		if err := bi.App.Customer.RemoveBrandFollowing(sc, opts); err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed remove brand_id in customer following")
			return errors.Wrap(err, "failed to remove follower")
		}

		if err := session.CommitTransaction(sc); err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to remove follower")
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}
