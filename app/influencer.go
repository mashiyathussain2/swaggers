package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"regexp"
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
		"name": primitive.Regex{
			Pattern: name,
			Options: "i",
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
	if request.ID.IsZero() == false {
		if request.Status == model.AcceptedStatus {
			return errors.Errorf("account already has influencer access")
		}
		return errors.Errorf("account upgrade request is already in active status")
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
		if request.ID.IsZero() == true {
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

// func (ii *InfluencerImpl) GetInfluencerAccountRequest() ([]schema.InfluencerAccountRequestResp, error) {
// 	var resp []schema.InfluencerAccountRequestResp
// 	ctx := context.TODO()
// 	matchStage := bson.D{
// 		{
// 			Key: "$match",
// 			Value: bson.M{
// 				"is_active": true,
// 			},
// 		},
// 	}
// 	lookupStage := bson.D{
// 		{
// 			Key: "$lookup",
// 			Value: bson.M{
// 				"from":         model.InfluencerColl,
// 				"localField":   "influencer_id",
// 				"foreignField": "_id",
// 				"as":           "influencer_info",
// 			},
// 		},
// 	}
// 	lookupStage2 := bson.D{
// 		{
// 			Key: "$lookup",
// 			Value: bson.M{
// 				"from":         model.UserColl,
// 				"localField":   "user_id",
// 				"foreignField": "_id",
// 				"as":           "user_info",
// 			},
// 		},
// 	}
// 	setStage := bson.D{
// 		{
// 			Key: "$set",
// 			Value: bson.M{
// 				"influencer_info": bson.M{
// 					"$arrayElemAt": bson.A{
// 						"$influencer_info",
// 						0,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	setStage2 := bson.D{
// 		{
// 			Key: "$set",
// 			Value: bson.M{
// 				"user_info": bson.M{
// 					"$arrayElemAt": bson.A{
// 						"$user_info",
// 						0,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	cur, err := ii.DB.Collection(model.InfluencerAccountRequestColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, lookupStage2, setStage, setStage2})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to query influencer account requests")
// 	}
// 	if err := cur.All(ctx, &resp); err != nil {
// 		return nil, errors.Wrap(err, "failed to get results")
// 	}
// 	return resp, nil
// }

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
	// projectStage := bson.D{
	// 	{
	// 		Key: "$project",
	// 		Value: bson.M{
	// 			"name":           "$customer_info.full_name",
	// 			"profile_image":  1,
	// 			"cover_image":    1,
	// 			"bio":            1,
	// 			"website":        1,
	// 			"social_account": 1,
	// 			"email":          "$user_info.email",
	// 			"phone_no":       "$user_info.phone_no",
	// 			"gender":         "$customer_info.gender",
	// 			"dob":            "$customer_info.dob",
	// 			"created_at":     1,
	// 		},
	// 	},
	// }
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

	isAlpha := regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	if !isAlpha(username) {
		errors.Errorf("%s is not valid", username)
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
		return nil, errors.Wrap(err, "failed to add follower")
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
