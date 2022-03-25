package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"math"
	"regexp"
	"strings"
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
	GetInfluencerByID(primitive.ObjectID) (*schema.GetInfluencerResp, error)

	GetInfluencerByName(string) ([]schema.GetInfluencerResp, error)
	AddFollower(*schema.AddInfluencerFollowerOpts) (bool, error)
	RemoveFollower(*schema.AddInfluencerFollowerOpts) (bool, error)

	InfluencerAccountRequest(opts *schema.InfluencerAccountRequestOpts) error
	GetInfluencerAccountRequestStatus(id primitive.ObjectID) (string, error)
	GetInfluencerAccountRequest() ([]schema.InfluencerAccountRequestResp, error)
	UpdateInfluencerAccountRequestStatus(opts *schema.UpdateInfluencerAccountRequestStatusOpts) error
	CheckInfluencerUsernameExists(string, *mongo.SessionContext) error
	EditInfluencerApp(*schema.EditInfluencerAppOpts) (*schema.EditInfluencerResp, error)
	AddCreditTransaction(opts *schema.CommisionOrderItem) error
	DebitRequest(opts *schema.CommissionDebitRequest) error
	UpdateDebitRequest(opts *schema.UpdateCommissionDebitRequest) error
	GetActiveDebitRequest() ([]schema.GetDebitRequestResponse, error)
	GetInfluencerDashboard(opts *schema.GetInfluencerDashboardOpts) (*schema.GetInfluencerDashboardResp, error)
	GetInfluencerLedger(opts *schema.GetInfluencerLedgerOpts) ([]schema.GetInfluencerLedgerResp, error)
	GetInfluencerPayoutInfo(id primitive.ObjectID) (*schema.GetPayoutInfoResp, error)
	GetCommissionAndRevenue(opts *schema.GetCommissionAndRevenueOpts) (*schema.GetCommissionAndRevenueResp, error)
	EditInfluencerAppV2(opts *schema.EditInfluencerAppV2Opts) (*schema.EditInfluencerResp, error)

	//v2
	InfluencerAccountRequestV2(opts *schema.InfluencerAccountRequestV2Opts) error
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

	ctx := context.TODO()
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		ii.Logger.Err(err).Msg("unable to create db session")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		ii.Logger.Err(err).Msg("unable to start transaction")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	var i model.Influencer
	var res *mongo.InsertOneResult
	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		i = model.Influencer{
			Name: opts.Name,
			Bio:  opts.Bio,
			CoverImg: &model.IMG{
				SRC: opts.CoverImg.SRC,
			},
			ProfileImage: &model.IMG{
				SRC: opts.ProfileImage.SRC,
			},
			ExternalLinks: opts.ExternalLinks,
			CreatedAt:     time.Now().UTC(),
		}
		//TODO: check if username is unique
		// isAlpha := regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
		// if !isAlpha(opts.Username) {
		// 	session.AbortTransaction(sc)
		// 	return errors.Errorf("%s is not valid", opts.Username)
		// }
		// filter := bson.M{
		// 	"username": opts.Username,
		// }
		// var influencer *model.Influencer
		// err := ii.DB.Collection(model.InfluencerColl).FindOne(sc, filter).Decode(&influencer)
		// if err != nil {
		// 	if err == mongo.ErrNilDocument || err == mongo.ErrNilValue {
		// 		i.Username = opts.Username
		// 	} else {
		// 		session.AbortTransaction(sc)
		// 		return errors.Wrapf(err, "error checking if username exists or not")
		// 	}
		// }
		// if influencer.Username == opts.Username {
		// 	session.AbortTransaction(sc)
		// 	return errors.Errorf("username: %s already exist", opts.Username)
		// }
		err = ii.CheckInfluencerUsernameExists(opts.Username, &sc)
		if err != nil {
			return err
		}
		i.Username = opts.Username

		if err := i.CoverImg.LoadFromURL(); err != nil {
			session.AbortTransaction(sc)
			return errors.Wrap(err, "invalid image for influencer cover")
		}
		if err := i.ProfileImage.LoadFromURL(); err != nil {
			session.AbortTransaction(sc)
			return errors.Wrap(err, "invalid image for profile image")
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

		res, err = ii.DB.Collection(model.InfluencerColl).InsertOne(sc, i)
		if err != nil {
			session.AbortTransaction(sc)
			ii.Logger.Err(err).Interface("opts", opts).Msg("failed to insert influencer")
			return errors.Wrap(err, "failed to create influencer")
		}
		if err := session.CommitTransaction(sc); err != nil {
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to create influencer")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &schema.CreateInfluencerResp{
		ID:            res.InsertedID.(primitive.ObjectID),
		Name:          i.Name,
		Username:      i.Username,
		Bio:           i.Bio,
		ExternalLinks: i.ExternalLinks,
		CoverImg:      i.CoverImg,
		ProfileImage:  i.ProfileImage,
		SocialAccount: i.SocialAccount,
		CreatedAt:     i.CreatedAt,
	}, nil
}

// EditInfluencer updates existing influencer details
func (ii *InfluencerImpl) EditInfluencer(opts *schema.EditInfluencerOpts) (*schema.EditInfluencerResp, error) {

	ctx := context.TODO()
	var influencer model.Influencer
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		ii.Logger.Err(err).Msg("unable to create db session")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		ii.Logger.Err(err).Msg("unable to start transaction")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		var update bson.D

		if opts.Name != "" {
			update = append(update, bson.E{Key: "name", Value: opts.Name})
		}
		if opts.Username != "" {
			err := ii.CheckInfluencerUsernameExists(opts.Username, &sc)
			if err != nil {
				return err
			}
			update = append(update, bson.E{Key: "username", Value: opts.Username})
		}

		if opts.CoverImg != nil {
			img := model.IMG{SRC: opts.CoverImg.SRC}
			if err := img.LoadFromURL(); err != nil {
				return errors.Wrap(err, "invalid image for brand cover")
			}
			update = append(update, bson.E{Key: "cover_img", Value: img})
		}

		if opts.ProfileImage != nil {
			img := model.IMG{SRC: opts.ProfileImage.SRC}
			if err := img.LoadFromURL(); err != nil {
				return errors.Wrap(err, "invalid image for profile image")
			}
			update = append(update, bson.E{Key: "profile_image", Value: img})
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
			return errors.New("no fields found to update")
		}

		update = append(update, bson.E{Key: "updated_at", Value: time.Now().UTC()})

		filterQuery := bson.M{"_id": opts.ID}
		updateQuery := bson.M{"$set": update}

		queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		if err := ii.DB.Collection(model.InfluencerColl).FindOneAndUpdate(context.TODO(), filterQuery, updateQuery, queryOpts).Decode(&influencer); err != nil {
			if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
				return errors.Wrapf(err, "influencer with id:%s not found", opts.ID.Hex())
			}
			return errors.Wrapf(err, "failed to update influencer with id:%s", opts.ID.Hex())
		}
		return nil
	}); err != nil {
		return nil, err
	}

	resp := schema.EditInfluencerResp{
		ID:            influencer.ID,
		Name:          influencer.Name,
		Username:      influencer.Username,
		ExternalLinks: influencer.ExternalLinks,
		CoverImg:      influencer.CoverImg,
		ProfileImage:  influencer.ProfileImage,
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

// GetInfluencerByID returns influencer info with matching id
func (ii *InfluencerImpl) GetInfluencerByID(id primitive.ObjectID) (*schema.GetInfluencerResp, error) {
	var resp schema.GetInfluencerResp
	ctx := context.TODO()
	filter := bson.M{
		"_id": id,
	}
	if err := ii.DB.Collection(model.InfluencerColl).FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, errors.Wrapf(err, "failed to get influencer with id: %s", id.Hex())
	}

	return &resp, nil
}

func (ii *InfluencerImpl) GetInfluencerByName(name string) ([]schema.GetInfluencerResp, error) {
	ctx := context.TODO()
	filter := bson.M{
		"$or": bson.A{
			bson.M{
				"name": primitive.Regex{
					Pattern: name,
					Options: "i",
				},
			},
			bson.M{
				"external_links": primitive.Regex{
					Pattern: name,
					Options: "i",
				},
			},
			bson.M{
				"social_account.instagram.url": primitive.Regex{
					Pattern: name,
					Options: "i",
				},
			},
		},
	}
	cur, err := ii.DB.Collection(model.InfluencerColl).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to find influencer")
	}

	var resp []schema.GetInfluencerResp
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to find influencer")
	}
	return resp, nil
}

func (ii *InfluencerImpl) AddFollower(opts *schema.AddInfluencerFollowerOpts) (bool, error) {
	ctx := context.TODO()
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		ii.Logger.Err(err).Msg("unable to create db session")
		return false, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		ii.Logger.Err(err).Msg("unable to start transaction")
		return false, errors.Wrap(err, "failed to add follower")
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		isFollowing, err := ii.DB.Collection(model.InfluencerColl).CountDocuments(sc, bson.M{"_id": opts.InfluencerID, "followers_id": opts.CustomerID})
		if err != nil {
			ii.Logger.Err(err).Interface("opts", opts).Msg("failed to check is user already follow influencer")
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to follow influencer")
		}

		if isFollowing != 0 {
			session.AbortTransaction(sc)
			return errors.New("user already follow the influencer")
		}

		filter := bson.M{
			"_id": opts.InfluencerID,
		}
		update := bson.M{
			"$addToSet": bson.M{
				"followers_id": opts.CustomerID,
			},
			"$inc": bson.M{
				"followers_count": 1,
			},
		}

		res, err := ii.DB.Collection(model.InfluencerColl).UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed add follower")
			return errors.Wrap(err, "failed to add follower")
		}
		if res.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return errors.New("influencer not found")
		}
		if err := ii.App.Customer.AddInfluencerFollowing(sc, opts); err != nil {
			session.AbortTransaction(sc)
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed add influencer_id in customer following")
			return errors.Wrap(err, "failed to add follower")
		}
		if err := session.CommitTransaction(sc); err != nil {
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to add follower")
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (ii *InfluencerImpl) RemoveFollower(opts *schema.AddInfluencerFollowerOpts) (bool, error) {
	ctx := context.TODO()
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		ii.Logger.Err(err).Msg("unable to create db session")
		return false, errors.Wrap(err, "failed to remove follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		ii.Logger.Err(err).Msg("unable to start transaction")
		return false, errors.Wrap(err, "failed to remove follower")
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		isFollowing, err := ii.DB.Collection(model.InfluencerColl).CountDocuments(sc, bson.M{"_id": opts.InfluencerID, "followers_id": opts.CustomerID})
		if err != nil {
			ii.Logger.Err(err).Interface("opts", opts).Msg("failed to check is user already follow influencer")
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to follow influencer")
		}

		if isFollowing == 0 {
			session.AbortTransaction(sc)
			return errors.New("user does not follow the influencer")
		}

		filter := bson.M{
			"_id": opts.InfluencerID,
		}
		update := bson.M{
			"$pull": bson.M{
				"followers_id": opts.CustomerID,
			},
			"$inc": bson.M{
				"followers_count": -1,
			},
		}

		res, err := ii.DB.Collection(model.InfluencerColl).UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed remove follower")
			return errors.Wrap(err, "failed to remove follower")
		}
		if res.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return errors.New("influencer not found")
		}
		if err := ii.App.Customer.RemoveInfluencerFollowing(sc, opts); err != nil {
			session.AbortTransaction(sc)
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed remove influencer_id in customer following")
			return errors.Wrap(err, "failed to remove follower")
		}
		if err := session.CommitTransaction(sc); err != nil {
			ii.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to remove follower")
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (ii *InfluencerImpl) linkUserAccount(sc mongo.SessionContext, opts *schema.LinkUserAccountOpts) error {
	// get influencer
	var influencer model.Influencer
	if err := ii.DB.Collection(model.InfluencerColl).FindOne(sc, bson.M{"_id": opts.InfluencerID}).Decode(&influencer); err != nil {
		return errors.Wrap(err, "failed to get influencer")
	}

	// get user
	var user model.User
	if err := ii.DB.Collection(model.UserColl).FindOne(sc, bson.M{"_id": opts.UserID}).Decode(&user); err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	if !user.InfluencerID.IsZero() {
		if user.InfluencerID != opts.InfluencerID {
			return errors.Errorf("user already has another influencer account attached to it")
		}
		return errors.Errorf("user already has this influencer account attached to it")
	}

	filter := bson.M{
		"_id": opts.UserID,
	}
	update := bson.M{
		"$set": bson.M{
			"influencer_id": opts.InfluencerID,
		},
	}
	if _, err := ii.DB.Collection(model.UserColl).UpdateOne(sc, filter, update); err != nil {
		return errors.Wrap(err, "failed to link influencer with user account")
	}
	return nil
}

func (ii *InfluencerImpl) InfluencerAccountRequest(opts *schema.InfluencerAccountRequestOpts) error {
	var request model.InfluencerAccountRequest
	ctx := context.TODO()
	filter := bson.M{
		"user_id":   opts.UserID,
		"is_active": true,
	}
	if err := ii.DB.Collection(model.InfluencerAccountRequestColl).FindOne(ctx, filter).Decode(&request); err != nil {
		if err != mongo.ErrNilDocument && err != mongo.ErrNoDocuments {
			return errors.Wrap(err, "failed to check for existing requests")
		}
	}
	if !request.ID.IsZero() {
		if request.Status == model.AcceptedStatus {
			return errors.Errorf("account already has influencer access")
		}
		return errors.Errorf("account upgrade request is already in active status")
	}
	if opts.Username == "" {
		opts.Username = GenerateUsernameInfluencer(opts.FullName)
	}
	err := ii.CheckInfluencerUsernameExists(opts.Username, nil)
	if err != nil {
		return err
	}
	r := model.InfluencerAccountRequest{
		UserID:     opts.UserID,
		CustomerID: opts.CustomerID,
		// InfluencerID: opts.InfluencerID,
		Name:     opts.FullName,
		Username: opts.Username,
		ProfileImage: &model.IMG{
			SRC: opts.ProfileImage.SRC,
		},
		CoverImage: &model.IMG{
			SRC: opts.CoverImage.SRC,
		},
		Bio:       opts.Bio,
		Website:   opts.Website,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		Status:    model.InReviewStatus,
	}
	if err := r.ProfileImage.LoadFromURL(); err != nil {
		return errors.Wrap(err, "invalid profile image for influencer")
	}
	if err := r.CoverImage.LoadFromURL(); err != nil {
		return errors.Wrap(err, "invalid cover image for influencer")
	}

	if opts.SocialAccount != nil {
		r.SocialAccount = &model.SocialAccount{}
		if opts.SocialAccount.Facebook != nil {
			r.SocialAccount.Facebook = &model.SocialMedia{
				URL:            opts.SocialAccount.Facebook.URL,
				FollowersCount: uint(opts.SocialAccount.Facebook.FollowersCount),
			}
		}
		if opts.SocialAccount.Instagram != nil {
			r.SocialAccount.Instagram = &model.SocialMedia{
				URL:            opts.SocialAccount.Instagram.URL,
				FollowersCount: uint(opts.SocialAccount.Instagram.FollowersCount),
			}
		}
		if opts.SocialAccount.Youtube != nil {
			r.SocialAccount.Youtube = &model.SocialMedia{
				URL:            opts.SocialAccount.Youtube.URL,
				FollowersCount: uint(opts.SocialAccount.Youtube.FollowersCount),
			}
		}
		if opts.SocialAccount.Twitter != nil {
			r.SocialAccount.Twitter = &model.SocialMedia{
				URL:            opts.SocialAccount.Twitter.URL,
				FollowersCount: uint(opts.SocialAccount.Twitter.FollowersCount),
			}
		}
	}
	if _, err := ii.DB.Collection(model.InfluencerAccountRequestColl).InsertOne(ctx, r); err != nil {
		return errors.Wrap(err, "failed to create account upgrade request")
	}
	return nil
}

func (ii *InfluencerImpl) GetInfluencerAccountRequestStatus(id primitive.ObjectID) (string, error) {
	ctx := context.TODO()
	// checking if influencer profile is already associated with user model
	var user model.User
	if err := ii.DB.Collection(model.UserColl).FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return "", errors.Wrap(err, "failed to find user")
	}
	if !user.InfluencerID.IsZero() {
		return model.AcceptedStatus, nil
	}

	// checking if influencer request is accepted or rejected
	var request model.InfluencerAccountRequest
	filter := bson.M{
		"user_id": id,
	}
	queryOpts := options.FindOne().SetSort(bson.M{"_id": -1})

	if err := ii.DB.Collection(model.InfluencerAccountRequestColl).FindOne(ctx, filter, queryOpts).Decode(&request); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return "", errors.Errorf("not account upgrade request found")
		}
		return "", errors.Wrap(err, "failed to get account upgrade request")
	}

	if !request.IsActive {
		if request.Status != "" {
			return request.Status, nil
		}
	}
	return request.Status, nil
}

func (ii *InfluencerImpl) UpdateInfluencerAccountRequestStatus(opts *schema.UpdateInfluencerAccountRequestStatusOpts) error {
	ctx := context.TODO()

	// creating session for atomic updates
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		return errors.Wrap(err, "failed to start session")
	}
	// Closing session at the end for function execution
	defer session.EndSession(ctx)

	// staring a new transaction
	if err := session.StartTransaction(); err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	if err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {
		filter := bson.M{
			"_id":       opts.ID,
			"is_active": true,
		}
		status := model.AcceptedStatus
		if !*opts.Grant {
			status = model.RejectedStatus
		}
		update := bson.M{
			"$set": bson.M{
				"grantee_id": opts.GranteeID,
				"status":     status,
				"is_active":  false,
				"granted_at": time.Now().UTC(),
			},
		}
		var request model.InfluencerAccountRequest
		queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		if err := ii.DB.Collection(model.InfluencerAccountRequestColl).FindOneAndUpdate(sc, filter, update, queryOpts).Decode(&request); err != nil {
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to update request status")
		}
		if request.ID.IsZero() {
			session.AbortTransaction(sc)
			return errors.Errorf("influencer account request failed")
		}

		if request.Status == model.AcceptedStatus {
			res, err := ii.createInfluencerFromRequest(sc, &request)
			if err != nil {
				session.AbortTransaction(sc)
				return errors.Wrap(err, "failed to create influencer")
			}
			fmt.Println(res)
			if err := ii.linkUserAccount(sc, &schema.LinkUserAccountOpts{RequestID: opts.ID, InfluencerID: res.ID, UserID: request.UserID}); err != nil {
				session.AbortTransaction(sc)
				return errors.Wrap(err, "failed to link user with influencer")
			}
		}

		if err := session.CommitTransaction(sc); err != nil {
			return errors.Wrap(err, "failed to commit transaction")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (ii *InfluencerImpl) GetInfluencerAccountRequest() ([]schema.InfluencerAccountRequestResp, error) {
	var resp []schema.InfluencerAccountRequestResp
	ctx := context.TODO()
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"is_active": true,
			},
		},
	}
	sortStage := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"_id": -1,
		},
	}}
	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         model.CustomerColl,
				"localField":   "user_id",
				"foreignField": "user_id",
				"as":           "customer_info",
			},
		},
	}
	lookupStage2 := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         model.UserColl,
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user_info",
			},
		},
	}
	setStage := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"user_info": bson.M{
					"$arrayElemAt": bson.A{
						"$user_info",
						0,
					},
				},
				"customer_info": bson.M{
					"$arrayElemAt": bson.A{
						"$customer_info",
						0,
					},
				},
			},
		},
	}

	cur, err := ii.DB.Collection(model.InfluencerAccountRequestColl).Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, lookupStage, lookupStage2, setStage})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query influencer account requests")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return resp, nil
}

// createInfluencerFromRequest create a new influencer profile
// Note: for now only creating influencer is supported. Code to link influencer profile to user
// needs to be implemented separately
func (ii *InfluencerImpl) createInfluencerFromRequest(sc mongo.SessionContext, opts *model.InfluencerAccountRequest) (*schema.CreateInfluencerResp, error) {
	i := model.Influencer{
		Name:          opts.Name,
		Username:      opts.Username,
		Bio:           opts.Bio,
		CoverImg:      opts.CoverImage,
		ProfileImage:  opts.ProfileImage,
		ExternalLinks: []string{opts.Website},
		SocialAccount: opts.SocialAccount,
		CreatedAt:     time.Now().UTC(),
	}
	err := ii.CheckInfluencerUsernameExists(opts.Username, &sc)
	if err != nil {
		return nil, err
	}
	res, err := ii.DB.Collection(model.InfluencerColl).InsertOne(sc, i)
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
		ProfileImage:  i.ProfileImage,
		SocialAccount: i.SocialAccount,
		CreatedAt:     i.CreatedAt,
	}, nil
}

func (ii *InfluencerImpl) CheckInfluencerUsernameExists(username string, sc *mongo.SessionContext) error {
	// ctx := context.TODO()

	isAlpha := regexp.MustCompile(`^[a-z0-9\_\-\.]{5,30}$`).MatchString
	if !isAlpha(username) {
		return errors.Errorf("%s is not valid", username)
	}
	filter := bson.M{
		"username": username,
	}
	var influencer *model.Influencer
	var err error
	if sc != nil {
		err = ii.DB.Collection(model.InfluencerColl).FindOne(*sc, filter).Decode(&influencer)
	} else {
		err = ii.DB.Collection(model.InfluencerColl).FindOne(context.TODO(), filter).Decode(&influencer)
	}
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNilValue || err == mongo.ErrNoDocuments {
			return nil
		}
		return errors.Wrapf(err, "error checking if username exists or not")
	}
	return errors.Errorf("username: %s already exist", username)
}

// EditInfluencerApp updates existing influencer details
func (ii *InfluencerImpl) EditInfluencerApp(opts *schema.EditInfluencerAppOpts) (*schema.EditInfluencerResp, error) {
	ctx := context.TODO()
	var influencer model.Influencer
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		ii.Logger.Err(err).Msg("unable to create db session")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		ii.Logger.Err(err).Msg("unable to start transaction")
		return nil, errors.Wrap(err, "failed to edit influencer")
	}
	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		var update bson.D
		if opts.Username != "" {
			err := ii.CheckInfluencerUsernameExists(opts.Username, &sc)
			if err != nil {
				return err
			}
			update = append(update, bson.E{Key: "username", Value: opts.Username})
		}
		if opts.PayoutInformation != nil {
			update = append(update, bson.E{Key: "payout_information", Value: model.PayoutInformation{
				UPIID:           opts.PayoutInformation.UPIID,
				BankInformation: opts.PayoutInformation.BankInformation,
				PanCard:         strings.ToUpper(opts.PayoutInformation.PanCard),
			}})
		}
		if update == nil {
			return errors.New("no fields found to update")
		}
		update = append(update, bson.E{Key: "updated_at", Value: time.Now().UTC()})

		filterQuery := bson.M{"_id": opts.ID}
		updateQuery := bson.M{"$set": update}

		queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		if err := ii.DB.Collection(model.InfluencerColl).FindOneAndUpdate(context.TODO(), filterQuery, updateQuery, queryOpts).Decode(&influencer); err != nil {
			if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
				return errors.Wrapf(err, "influencer with id:%s not found", opts.ID.Hex())
			}
			return errors.Wrapf(err, "failed to update influencer with id:%s", opts.ID.Hex())
		}
		if err := session.CommitTransaction(sc); err != nil {
			return errors.Wrapf(err, "failed to commit transaction")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	resp := schema.EditInfluencerResp{
		ID:            influencer.ID,
		Name:          influencer.Name,
		Username:      influencer.Username,
		ExternalLinks: influencer.ExternalLinks,
		CoverImg:      influencer.CoverImg,
		ProfileImage:  influencer.ProfileImage,
		SocialAccount: influencer.SocialAccount,
		Bio:           influencer.Bio,
		CreatedAt:     influencer.CreatedAt,
		UpdatedAt:     influencer.UpdatedAt,
	}
	return &resp, nil
}

func (ii *InfluencerImpl) AddCreditTransaction(opts *schema.CommisionOrderItem) error {

	ctx := context.TODO()
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		return errors.Wrap(err, "failed to start session")
	}
	// Closing session at the end for function execution
	defer session.EndSession(ctx)

	// staring a new transaction
	if err := session.StartTransaction(); err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	if err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {

		// 1. Get Commission rate from Catalog

		cr := float64(opts.CatalogInfo.CommissionRate)
		if cr == 0 {
			return nil
		}

		//2. Calculate Commission based on item total price
		commission := math.Floor(cr / 100.0 * float64(opts.TotalPrice.Value))
		//3. Add transaction to ledger collection

		iID, err := primitive.ObjectIDFromHex(opts.Source.ID)
		if err != nil {
			return err
		}

		oldBalance, err := ii.GetBalance(&sc, iID)
		if err != nil {
			return err
		}
		balance := oldBalance + commission

		brand, err := ii.App.Brand.GetBrandByID(opts.CatalogInfo.BrandID)
		if err != nil {
			return err
		}
		opts.CatalogInfo.BrandName = brand.Name
		transaction := model.Transaction{
			InfluencerID:    iID,
			Type:            model.CreditTransaction,
			ItemID:          opts.ID,
			OrderID:         opts.OrderID,
			OrderNo:         opts.OrderNo,
			OrderValue:      opts.TotalPrice,
			CatalogInfo:     opts.CatalogInfo,
			CommissionValue: commission,
			CreatedAt:       time.Now(),
			Balance:         balance,
			OrderDate:       opts.OrderDate,
		}
		_, err = ii.DB.Collection(model.CommissionLedgerColl).InsertOne(sc, transaction)
		if err != nil {
			return err
		}
		if err := session.CommitTransaction(sc); err != nil {
			return errors.Wrapf(err, "failed to commit transaction")
		}

		return nil
	}); err != nil {
		ii.Logger.Err(err).Msgf("failed to create credit transaction: %s", opts.OrderID)
		return err
	}

	return nil
}

// GetBalance returns last balance of Influencer
func (ii *InfluencerImpl) GetBalance(sc *mongo.SessionContext, id primitive.ObjectID) (float64, error) {

	var transaction model.Transaction
	filterOpts := options.FindOne().SetSort(bson.M{"_id": -1})
	var err error
	if sc != nil {
		err = ii.DB.Collection(model.CommissionLedgerColl).FindOne(*sc, bson.M{"influencer_id": id}, filterOpts).Decode(&transaction)
	} else {
		err = ii.DB.Collection(model.CommissionLedgerColl).FindOne(context.TODO(), bson.M{"influencer_id": id}, filterOpts).Decode(&transaction)
	}
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return transaction.Balance, nil
}

func (ii *InfluencerImpl) DebitRequest(opts *schema.CommissionDebitRequest) error {
	bal, err := ii.GetBalance(nil, opts.ID)
	if err != nil {
		return errors.Wrapf(err, "error getting current balance")
	}
	if bal < float64(opts.Amount) {
		return errors.New("amount requested is invalid")
	}
	ctx := context.TODO()
	//checking if payout info is available
	var influencer model.Influencer
	err = ii.DB.Collection(model.InfluencerColl).FindOne(ctx, bson.M{"_id": opts.ID}).Decode(&influencer)
	if err != nil {
		return errors.Wrapf(err, "error getting influencer info")
	}
	if influencer.PayoutInformation == nil {
		return errors.New("error: payout info missing")
	} else if influencer.PayoutInformation.UPIID == "" && influencer.PayoutInformation.BankInformation == nil {
		return errors.New("error: payout info missing")
	}
	if influencer.PayoutInformation.PanCard == "" {
		return errors.New("error: pancard info missing")
	}

	filter := bson.M{
		"influencer_id": opts.ID,
		"status":        model.InReviewStatus,
	}
	var dr *model.DebitRequest
	err = ii.DB.Collection(model.DebitRequestColl).FindOne(ctx, filter).Decode(&dr)
	if err != nil {
		if err != mongo.ErrNoDocuments && err != mongo.ErrNilDocument {
			return errors.Wrapf(err, "error checking for debit request")
		}
	}
	if dr != nil {
		return errors.New("another request is active")
	}
	dr = &model.DebitRequest{
		InfluencerID:      opts.ID,
		Amount:            float64(opts.Amount),
		Status:            model.InReviewStatus,
		PayoutInformation: (*model.PayoutInformation)(&opts.PayoutInformation),
		CreatedAt:         time.Now(),
	}
	_, err = ii.DB.Collection(model.DebitRequestColl).InsertOne(ctx, dr)
	if err != nil {
		return errors.Wrapf(err, "error creating debit request")
	}
	return nil
}

func (ii *InfluencerImpl) UpdateDebitRequest(opts *schema.UpdateCommissionDebitRequest) error {

	ctx := context.TODO()
	session, err := ii.DB.Client().StartSession()
	if err != nil {
		return errors.Wrap(err, "failed to start session")
	}
	// Closing session at the end for function execution
	defer session.EndSession(ctx)

	// staring a new transaction
	if err := session.StartTransaction(); err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	if err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {

		dr := model.DebitRequest{
			Status:    opts.Status,
			GranteeID: opts.GranteeID,
			UpdatedAt: time.Now(),
		}
		filter := bson.M{
			"_id":    opts.ID,
			"status": model.InReviewStatus,
		}
		update := bson.M{
			"$set": dr,
		}
		queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err := ii.DB.Collection(model.DebitRequestColl).FindOneAndUpdate(sc, filter, update, queryOpts).Decode(&dr)
		if err != nil {
			return errors.Wrapf(err, "error updating debit request")
		}
		if dr.Status == model.RejectedStatus {
			if err := session.CommitTransaction(sc); err != nil {
				return errors.Wrapf(err, "failed to commit transaction")
			}
			return nil
		}
		oldBal, err := ii.GetBalance(&sc, dr.InfluencerID)
		if err != nil {
			return errors.Wrapf(err, "error getting current balance")
		}
		bl := oldBal - dr.Amount
		if bl < 0 {
			return errors.New("error amount exceeding current balance")
		}
		if err != nil {
			return errors.Wrapf(err, "error getting influencer payout info")
		}
		//create debit transaction
		transaction := model.Transaction{
			InfluencerID:      dr.InfluencerID,
			Type:              model.DebitTransaction,
			DebitAmount:       dr.Amount,
			CreatedAt:         time.Now(),
			Balance:           bl,
			PayoutInformation: dr.PayoutInformation,
		}
		_, err = ii.DB.Collection(model.CommissionLedgerColl).InsertOne(sc, transaction)
		if err != nil {
			return err
		}
		if err := session.CommitTransaction(sc); err != nil {
			return errors.Wrapf(err, "failed to commit transaction")
		}
		return nil
	}); err != nil {
		ii.Logger.Err(err).Msgf("failed to create debit transaction: %s", opts.ID)
		return err
	}

	return nil
}

func (ii *InfluencerImpl) GetActiveDebitRequest() ([]schema.GetDebitRequestResponse, error) {
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"status": model.InReviewStatus,
		},
	}}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "influencer",
			"localField":   "influencer_id",
			"foreignField": "_id",
			"as":           "influencer_info",
		},
	}}

	lookupStage2 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "user",
			"localField":   "influencer_id",
			"foreignField": "influencer_id",
			"as":           "user_info",
		},
	}}

	projectStage := bson.D{{
		Key: "$project", Value: bson.M{
			"amount":     1,
			"created_at": 1,
			"influencer_info": bson.M{
				"$first": "$influencer_info",
			},
			"status": 1,
			"phone_no": bson.M{
				"$first": "$user_info.phone_no.number",
			},
			"email": bson.M{
				"$first": "$user_info.email",
			},
			"payout_information": 1,
		},
	}}

	cursor, err := ii.DB.Collection(model.DebitRequestColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		lookupStage2,
		projectStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get debit request data")
	}
	var resp []schema.GetDebitRequestResponse
	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return resp, nil
}

func (ii *InfluencerImpl) GetInfluencerDashboard(opts *schema.GetInfluencerDashboardOpts) (*schema.GetInfluencerDashboardResp, error) {
	var matchStage bson.D
	matchStage = bson.D{{
		Key: "$match", Value: bson.M{
			"influencer_id": opts.ID,
		},
	}}
	if opts.StartDate != nil && opts.EndDate != nil {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"influencer_id": opts.ID,
				"created_at": bson.M{
					"$gte": opts.StartDate,
					"$lte": opts.EndDate,
				},
			},
		}}
	}

	facetStage := bson.D{{
		Key: "$facet", Value: bson.M{
			"overall_data": bson.A{
				bson.D{{
					Key: "$group", Value: bson.M{
						"_id": "$influencer_id",
						"revenue": bson.M{
							"$sum": "$order_value.value",
						},
						"total_commission": bson.M{
							"$sum": "$commission_value",
						},
					},
				}},
			},
			"monthly_data": bson.A{
				bson.D{{
					Key: "$group", Value: bson.M{
						"_id": bson.M{
							"$month": "$created_at",
						},
						"count": bson.M{
							"$sum": 1,
						},
					},
				}},
			},
		},
	}}

	projectStage := bson.D{{
		Key: "$project", Value: bson.M{
			"overall_data": bson.M{"$first": "$overall_data"},
			"monthly_data": 1,
			// "ledger":       1,
		},
	}}

	ctx := context.TODO()
	var resp []schema.GetInfluencerDashboardResp
	cur, err := ii.DB.Collection(model.CommissionLedgerColl).Aggregate(ctx, mongo.Pipeline{matchStage, facetStage, projectStage})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query influencer dashboard data")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	if len(resp) == 0 {
		return nil, nil
	}
	resp[0].Ledger, err = ii.GetInfluencerLedger(&schema.GetInfluencerLedgerOpts{
		ID: opts.ID,
		// Page:      opts.Page,
		Type:      "credit",
		StartDate: opts.StartDate,
		EndDate:   opts.EndDate,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ledger")
	}
	// fmt.Println("response", resp)
	return &resp[0], nil
}

func (ii *InfluencerImpl) GetInfluencerLedger(opts *schema.GetInfluencerLedgerOpts) ([]schema.GetInfluencerLedgerResp, error) {
	ctx := context.TODO()

	var pipeline mongo.Pipeline
	if opts.Type == "credit" {
		pipeline = ii.commissionCreditPipeline(opts)
	} else {
		pipeline = ii.commissionDebitPipeline(opts)
	}
	var resp []schema.GetInfluencerLedgerResp
	fmt.Println(pipeline)
	cur, err := ii.DB.Collection(model.CommissionLedgerColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query influencer ledger data")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp, nil
}

func (ii *InfluencerImpl) GetInfluencerPayoutInfo(id primitive.ObjectID) (*schema.GetPayoutInfoResp, error) {

	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": id,
		},
	}}
	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         model.CommissionLedgerColl,
			"localField":   "_id",
			"foreignField": "influencer_id",
			"as":           "ledger",
		},
	}}
	projectStage := bson.D{{
		Key: "$project", Value: bson.M{
			"payout_information": 1,
			"balance": bson.M{
				"$last": "$ledger.balance",
			},
		},
	}}

	var resp []schema.GetPayoutInfoResp

	ctx := context.TODO()
	cur, err := ii.DB.Collection(model.InfluencerColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query influencer dashboard data")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}

	if len(resp) == 0 {
		return nil, nil
	}
	return &resp[0], nil
}

func (ii *InfluencerImpl) GetCommissionAndRevenue(opts *schema.GetCommissionAndRevenueOpts) (*schema.GetCommissionAndRevenueResp, error) {
	ctx := context.TODO()
	var matchStage bson.D
	if opts.StartDate.IsZero() && opts.EndDate.IsZero() {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"influencer_id": opts.ID,
			},
		}}
	} else {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"influencer_id": opts.ID,
				"created_at": bson.M{
					"$gte": opts.StartDate,
					"$lte": opts.EndDate,
				},
			},
		}}
	}
	fmt.Println(matchStage)
	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$influencer_id",
			"commission": bson.M{
				"$sum": "$commission_value",
			},
			"revenue": bson.M{
				"$sum": "$order_value.value",
			},
			"balance": bson.M{
				"$last": "$balance",
			},
		},
	}}

	var resp []schema.GetCommissionAndRevenueResp

	cur, err := ii.DB.Collection(model.CommissionLedgerColl).Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query influencer dashboard data")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}

	if len(resp) == 0 {
		return nil, nil
	}
	return &resp[0], nil
}

// EditInfluencerAppV2 updates existing influencer details
func (ii *InfluencerImpl) EditInfluencerAppV2(opts *schema.EditInfluencerAppV2Opts) (*schema.EditInfluencerResp, error) {
	var influencer model.Influencer
	var update bson.D
	if opts.Username != "" {
		err := ii.CheckInfluencerUsernameExists(opts.Username, nil)
		if err != nil {
			return nil, err
		}
		update = append(update, bson.E{Key: "username", Value: opts.Username})
	}
	if opts.PayoutInformation != nil {
		update = append(update, bson.E{Key: "payout_information", Value: model.PayoutInformation{
			UPIID:           opts.PayoutInformation.UPIID,
			BankInformation: opts.PayoutInformation.BankInformation,
			PanCard:         strings.ToUpper(opts.PayoutInformation.PanCard),
		}})
	}
	if opts.Name != "" {
		update = append(update, bson.E{Key: "name", Value: opts.Name})
	}
	if opts.Bio != "" {
		update = append(update, bson.E{Key: "bio", Value: opts.Bio})
	}
	if opts.ProfileImage != nil {
		img := model.IMG{
			SRC: opts.ProfileImage.SRC,
		}
		if err := img.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid profile image for influencer")
		}
		update = append(update, bson.E{Key: "profile_image", Value: img})
	}
	if opts.CoverImg != nil {
		img := model.IMG{
			SRC: opts.CoverImg.SRC,
		}
		if err := img.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid cover image for influencer")
		}
		update = append(update, bson.E{Key: "cover_img", Value: img})
	}
	if len(opts.ExternalLinks) != 0 {
		update = append(update, bson.E{Key: "external_links", Value: opts.ExternalLinks})
	}
	if opts.SocialAccount != nil {
		if opts.SocialAccount.Facebook != nil {
			facebook := &model.SocialMedia{
				URL:            opts.SocialAccount.Facebook.URL,
				FollowersCount: uint(opts.SocialAccount.Facebook.FollowersCount),
			}
			update = append(update, bson.E{Key: "social_account.facebook", Value: facebook})
		}
		if opts.SocialAccount.Instagram != nil {
			instagram := &model.SocialMedia{
				URL:            opts.SocialAccount.Instagram.URL,
				FollowersCount: uint(opts.SocialAccount.Instagram.FollowersCount),
			}
			update = append(update, bson.E{Key: "social_account.instagram", Value: instagram})
		}
		if opts.SocialAccount.Youtube != nil {
			youtube := &model.SocialMedia{
				URL:            opts.SocialAccount.Youtube.URL,
				FollowersCount: uint(opts.SocialAccount.Youtube.FollowersCount),
			}
			update = append(update, bson.E{Key: "social_account.youtube", Value: youtube})
		}
		if opts.SocialAccount.Twitter != nil {
			twitter := &model.SocialMedia{
				URL:            opts.SocialAccount.Twitter.URL,
				FollowersCount: uint(opts.SocialAccount.Twitter.FollowersCount),
			}
			update = append(update, bson.E{Key: "social_account.twitter", Value: twitter})
		}
	}
	if update == nil {
		return nil, errors.New("no fields found to update")
	}
	update = append(update, bson.E{Key: "updated_at", Value: time.Now().UTC()})

	filterQuery := bson.M{"_id": opts.ID}
	updateQuery := bson.M{"$set": update}

	queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := ii.DB.Collection(model.InfluencerColl).FindOneAndUpdate(context.TODO(), filterQuery, updateQuery, queryOpts).Decode(&influencer)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "influencer with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrapf(err, "failed to update influencer with id:%s", opts.ID.Hex())
	}

	resp := schema.EditInfluencerResp{
		ID:            influencer.ID,
		Name:          influencer.Name,
		Username:      influencer.Username,
		ExternalLinks: influencer.ExternalLinks,
		CoverImg:      influencer.CoverImg,
		ProfileImage:  influencer.ProfileImage,
		SocialAccount: influencer.SocialAccount,
		Bio:           influencer.Bio,
		CreatedAt:     influencer.CreatedAt,
		UpdatedAt:     influencer.UpdatedAt,
	}
	return &resp, nil
}

//InfluencerAccountRequestV2 creates a new influencer account request
func (ii *InfluencerImpl) InfluencerAccountRequestV2(opts *schema.InfluencerAccountRequestV2Opts) error {
	var request model.InfluencerAccountRequest
	ctx := context.TODO()
	filter := bson.M{
		"user_id":   opts.UserID,
		"is_active": true,
	}
	if err := ii.DB.Collection(model.InfluencerAccountRequestColl).FindOne(ctx, filter).Decode(&request); err != nil {
		if err != mongo.ErrNilDocument && err != mongo.ErrNoDocuments {
			return errors.Wrap(err, "failed to check for existing requests")
		}
	}
	if !request.ID.IsZero() {
		if request.Status == model.AcceptedStatus {
			return errors.Errorf("account already has influencer access")
		}
		return errors.Errorf("account upgrade request is already in active status")
	}
	if opts.Username == "" {
		opts.Username = GenerateUsernameInfluencer(opts.FullName)
	}
	err := ii.CheckInfluencerUsernameExists(opts.Username, nil)
	if err != nil {
		return err
	}
	r := model.InfluencerAccountRequest{
		UserID:     opts.UserID,
		CustomerID: opts.CustomerID,
		// InfluencerID: opts.InfluencerID,
		Name:     opts.FullName,
		Username: opts.Username,
		ProfileImage: &model.IMG{
			SRC: opts.ProfileImage.SRC,
		},
		CoverImage: &model.IMG{
			SRC: opts.CoverImage.SRC,
		},
		// Bio:       opts.Bio,
		// Website:   opts.Website,
		IsActive:        true,
		CreatedAt:       time.Now().UTC(),
		Status:          model.InReviewStatus,
		AreaOfExpertise: opts.AreaOfExpertise,
	}
	if err := r.ProfileImage.LoadFromURL(); err != nil {
		return errors.Wrap(err, "invalid profile image for influencer")
	}
	if err := r.CoverImage.LoadFromURL(); err != nil {
		return errors.Wrap(err, "invalid cover image for influencer")
	}

	if opts.SocialAccount != nil {
		r.SocialAccount = &model.SocialAccount{}
		if opts.SocialAccount.Facebook != nil {
			r.SocialAccount.Facebook = &model.SocialMedia{
				URL:            opts.SocialAccount.Facebook.URL,
				FollowersCount: uint(opts.SocialAccount.Facebook.FollowersCount),
			}
		}
		if opts.SocialAccount.Instagram != nil {
			r.SocialAccount.Instagram = &model.SocialMedia{
				URL:            opts.SocialAccount.Instagram.URL,
				FollowersCount: uint(opts.SocialAccount.Instagram.FollowersCount),
			}
		}
		if opts.SocialAccount.Youtube != nil {
			r.SocialAccount.Youtube = &model.SocialMedia{
				URL:            opts.SocialAccount.Youtube.URL,
				FollowersCount: uint(opts.SocialAccount.Youtube.FollowersCount),
			}
		}
		if opts.SocialAccount.Twitter != nil {
			r.SocialAccount.Twitter = &model.SocialMedia{
				URL:            opts.SocialAccount.Twitter.URL,
				FollowersCount: uint(opts.SocialAccount.Twitter.FollowersCount),
			}
		}
	}
	if _, err := ii.DB.Collection(model.InfluencerAccountRequestColl).InsertOne(ctx, r); err != nil {
		return errors.Wrap(err, "failed to create account upgrade request")
	}
	return nil
}

func (ii *InfluencerImpl) commissionCreditPipeline(opts *schema.GetInfluencerLedgerOpts) mongo.Pipeline {

	var matchStage bson.D
	if opts.StartDate.IsZero() && opts.EndDate.IsZero() {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"type":          opts.Type,
				"influencer_id": opts.ID,
			},
		}}
	} else {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"type":          opts.Type,
				"influencer_id": opts.ID,
				"order_date": bson.M{
					"$gte": opts.StartDate,
					"$lte": opts.EndDate,
				},
			},
		}}
	}

	sortStage1 := bson.D{{
		Key: "$sort", Value: bson.M{
			"_id": -1,
		},
	}}
	addFieldsStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"month": bson.M{
				"$let": bson.M{
					"vars": bson.M{
						"monthsInString": bson.A{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sept", "Oct", "Nov", "Dec"},
					},
					"in": bson.M{
						"$arrayElemAt": bson.A{"$$monthsInString", bson.M{"$month": "$order_date"}},
					},
				},
			},
		},
	}}

	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"date": bson.M{
				"$concat": bson.A{
					bson.M{"$toString": bson.M{"$dayOfMonth": "$order_date"}},
					" ",
					"$month",
					",",
					bson.M{"$toString": bson.M{"$year": "$order_date"}},
				},
			},
		},
	}}

	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$date",
			"ledger": bson.M{
				"$push": "$$ROOT",
			},
			"commission": bson.M{
				"$sum": "$commission_value",
			},
			"revenue": bson.M{
				"$sum": "$order_value.value",
			},
		},
	}}

	sortStage := bson.D{{
		Key: "$sort", Value: bson.M{
			"ledger.order_date": -1,
		},
	}}

	skipStage := bson.D{{
		Key: "$skip", Value: int64(opts.Page) * 10,
	}}

	limitStage := bson.D{{
		Key: "$limit", Value: 10,
	}}

	return mongo.Pipeline{matchStage, sortStage1, addFieldsStage, setStage, groupStage, sortStage, skipStage, limitStage}
}

func (ii *InfluencerImpl) commissionDebitPipeline(opts *schema.GetInfluencerLedgerOpts) mongo.Pipeline {

	var matchStage bson.D
	if opts.StartDate.IsZero() && opts.EndDate.IsZero() {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"type":          opts.Type,
				"influencer_id": opts.ID,
			},
		}}
	} else {
		matchStage = bson.D{{
			Key: "$match", Value: bson.M{
				"type":          opts.Type,
				"influencer_id": opts.ID,
				"created_at": bson.M{
					"$gte": opts.StartDate,
					"$lte": opts.EndDate,
				},
			},
		}}
	}

	sortStage1 := bson.D{{
		Key: "$sort", Value: bson.M{
			"_id": -1,
		},
	}}
	addFieldsStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"month": bson.M{
				"$let": bson.M{
					"vars": bson.M{
						"monthsInString": bson.A{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sept", "Oct", "Nov", "Dec"},
					},
					"in": bson.M{
						"$arrayElemAt": bson.A{"$$monthsInString", bson.M{"$month": "$created_at"}},
					},
				},
			},
		},
	}}

	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"date": bson.M{
				"$concat": bson.A{
					bson.M{"$toString": bson.M{"$dayOfMonth": "$created_at"}},
					" ",
					"$month",
					",",
					bson.M{"$toString": bson.M{"$year": "$created_at"}},
				},
			},
		},
	}}

	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$date",
			"ledger": bson.M{
				"$push": "$$ROOT",
			},
			"commission": bson.M{
				"$sum": "$commission_value",
			},
			"revenue": bson.M{
				"$sum": "$order_value.value",
			},
		},
	}}

	sortStage := bson.D{{
		Key: "$sort", Value: bson.M{
			"ledger.created_at": -1,
		},
	}}

	skipStage := bson.D{{
		Key: "$skip", Value: int64(opts.Page) * 10,
	}}

	limitStage := bson.D{{
		Key: "$limit", Value: 10,
	}}

	return mongo.Pipeline{matchStage, sortStage1, addFieldsStage, setStage, groupStage, sortStage, skipStage, limitStage}
}
