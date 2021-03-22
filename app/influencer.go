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

// Influencer contains methods for influencer specific operations
type Influencer interface {
	CreateInfluencer(*schema.CreateInfluencerOpts) (*schema.CreateInfluencerResp, error)
	EditInfluencer(*schema.EditInfluencerOpts) (*schema.EditInfluencerResp, error)

	GetInfluencersByID([]primitive.ObjectID) ([]schema.GetInfluencerResp, error)
}

// InfluencerImpl implements influencer interface methods

type InfluencerImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InfluencerImplOpts contains args required to create
type InfluencerImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitInfluencer returns new instance of influencer implementation
func InitInfluencer(opts *InfluencerImplOpts) Influencer {
	ii := InfluencerImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ii
}

// CreateInfluencer create a new influencer profile
// Note: for now only creating influencer is supported. Code to link influencer profile to user
// needs to be implemented separately
func (ii *InfluencerImpl) CreateInfluencer(opts *schema.CreateInfluencerOpts) (*schema.CreateInfluencerResp, error) {
	i := model.Influencer{
		Name: opts.Name,
		Bio:  opts.Bio,
		CoverImg: &model.IMG{
			SRC: opts.CoverImg.SRC,
		},
		ExternalLinks: opts.ExternalLinks,
		CreatedAt:     time.Now().UTC(),
	}
	if err := i.CoverImg.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "invalid image for influencer cover")
	}
	if opts.SocialAccount != nil {
		i.SocialAccount = &model.SocialAccount{}
		if opts.SocialAccount.Facebook != nil {
			i.SocialAccount.Facebook = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Facebook.FollowersCount)}
		}
		if opts.SocialAccount.Instagram != nil {
			i.SocialAccount.Instagram = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Instagram.FollowersCount)}
		}
		if opts.SocialAccount.Youtube != nil {
			i.SocialAccount.Youtube = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Youtube.FollowersCount)}
		}
		if opts.SocialAccount.Twitter != nil {
			i.SocialAccount.Twitter = &model.SocialMedia{FollowersCount: uint(opts.SocialAccount.Twitter.FollowersCount)}
		}
	}

	res, err := ii.DB.Collection(model.InfluencerColl).InsertOne(context.TODO(), i)
	if err != nil {
		ii.Logger.Err(err).Interface("opts", opts).Msg("failed to insert influencer")
		return nil, errors.Wrap(err, "failed to create influencer")
	}

	return &schema.CreateInfluencerResp{
		ID:            res.InsertedID.(primitive.ObjectID),
		Name:          i.Name,
		Bio:           i.Bio,
		ExternalLinks: i.ExternalLinks,
		CoverImg:      i.CoverImg,
		SocialAccount: i.SocialAccount,
		CreatedAt:     i.CreatedAt,
	}, nil
}

// EditInfluencer updates existing influencer details
func (ii *InfluencerImpl) EditInfluencer(opts *schema.EditInfluencerOpts) (*schema.EditInfluencerResp, error) {
	var update bson.D

	if opts.Name != "" {
		update = append(update, bson.E{Key: "name", Value: opts.Name})
	}
	if opts.CoverImg != nil {
		img := model.IMG{SRC: opts.CoverImg.SRC}
		if err := img.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid image for brand cover")
		}
		update = append(update, bson.E{Key: "cover_img", Value: img})
	}
	if opts.Bio != "" {
		update = append(update, bson.E{Key: "bio", Value: opts.Bio})
	}
	if len(opts.ExternalLinks) > 0 {
		update = append(update, bson.E{Key: "external_links", Value: opts.ExternalLinks})
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

	if update == nil {
		return nil, errors.New("no fields found to update")
	}

	update = append(update, bson.E{Key: "updated_at", Value: time.Now().UTC()})

	filterQuery := bson.M{"_id": opts.ID}
	updateQuery := bson.M{"$set": update}

	var influencer model.Influencer
	queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := ii.DB.Collection(model.InfluencerColl).FindOneAndUpdate(context.TODO(), filterQuery, updateQuery, queryOpts).Decode(&influencer); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "influencer with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrapf(err, "failed to update influencer with id:%s", opts.ID.Hex())
	}

	resp := schema.EditInfluencerResp{
		ID:            influencer.ID,
		Name:          influencer.Name,
		ExternalLinks: influencer.ExternalLinks,
		CoverImg:      influencer.CoverImg,
		SocialAccount: influencer.SocialAccount,
		Bio:           influencer.Bio,
		CreatedAt:     influencer.CreatedAt,
		UpdatedAt:     influencer.UpdatedAt,
	}
	return &resp, nil
}

// GetInfluencersByID returns influencer info with matching id
func (ii *InfluencerImpl) GetInfluencersByID(ids []primitive.ObjectID) ([]schema.GetInfluencerResp, error) {
	var resp []schema.GetInfluencerResp
	ctx := context.TODO()
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	cur, err := ii.DB.Collection(model.InfluencerColl).Find(ctx, filter)
	if err != nil {
		ii.Logger.Err(err).Interface("ids", ids).Msg("failed to get influencer with ids")
		return nil, errors.Wrap(err, "failed to get influencer with id")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to find influencer")
	}
	return resp, nil
}
