package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"sync"

	"go-app/server/auth"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Customer contains methods and functionality for customer related operations such as login, signup etc
type Customer interface {
	Login(*schema.EmailLoginCustomerOpts) (auth.Claim, error)
	SignUp(*schema.CreateUserOpts) (auth.Claim, error)
	UpdateCustomer(*schema.UpdateCustomerOpts) (auth.Claim, error)

	AddBrandFollowing(mongo.SessionContext, *schema.AddBrandFollowerOpts) error
	RemoveBrandFollowing(mongo.SessionContext, *schema.AddBrandFollowerOpts) error
	AddInfluencerFollowing(mongo.SessionContext, *schema.AddInfluencerFollowerOpts) error
	RemoveInfluencerFollowing(mongo.SessionContext, *schema.AddInfluencerFollowerOpts) error
	AddAddress(opts *schema.AddAddressOpts) (*schema.AddAddressResp, error)
	GetAddresses(primitive.ObjectID) ([]model.Address, error)
	GetAppCustomerInfo(id primitive.ObjectID) (*schema.GetCustomerProfileInfoResp, error)

	RemoveAddress(primitive.ObjectID, primitive.ObjectID) error
	EditAddress(opts *schema.EditAddressOpts) error
}

// CustomerImpl implements Customer interface methods
type CustomerImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// CustomerImplOpts contains args required to create a new instance of CustomerImpl
type CustomerImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitCustomer returns new instance of customer implementation
func InitCustomer(opts *CustomerImplOpts) Customer {
	ci := CustomerImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ci
}

// Login logs in customer use
func (ci *CustomerImpl) Login(opts *schema.EmailLoginCustomerOpts) (auth.Claim, error) {
	claim, err := ci.App.User.EmailLoginCustomerUser(opts)
	if err != nil {
		return nil, err
	}
	return claim, nil
}

// SignUp create a new customer
func (ci *CustomerImpl) SignUp(opts *schema.CreateUserOpts) (auth.Claim, error) {
	opts.Type = model.CustomerType
	user, err := ci.App.User.CreateUser(opts)
	if err != nil {
		return nil, err
	}
	customer := model.Customer{
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
	}
	res, err := ci.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)
	if err != nil {
		return nil, err
	}
	customer.ID = res.InsertedID.(primitive.ObjectID)

	claim := auth.UserClaim{
		ID:         user.ID.Hex(),
		CustomerID: customer.ID.Hex(),
		Type:       user.Type,
		Role:       model.UserRole,
		Email:      user.Email,
		PhoneNo:    user.PhoneNo,
	}
	return &claim, nil
}

// UpdateCustomer update existing customer fields
func (ci *CustomerImpl) UpdateCustomer(opts *schema.UpdateCustomerOpts) (auth.Claim, error) {
	var update bson.D
	var wg sync.WaitGroup
	if opts.Email != "" {
		ci.App.User.UpdateUserEmail(&schema.UpdateUserEmailOpts{ID: opts.UserID, Email: opts.Email})
	}
	if opts.PhoneNo != nil {
		ci.App.User.UpdateUserPhoneNo(&schema.UpdateUserPhoneNoOpts{ID: opts.UserID, PhoneNo: opts.PhoneNo})
	}

	if opts.FullName != "" {
		update = append(update, bson.E{Key: "full_name", Value: opts.FullName})
	}
	if !opts.DOB.IsZero() {
		update = append(update, bson.E{Key: "dob", Value: opts.DOB})
	}
	if opts.Gender != "" {
		gender := model.GetGender(opts.Gender)
		if gender == model.Invalid {
			return nil, errors.Errorf("%s is invalid gender value", opts.Gender)
		}
		update = append(update, bson.E{Key: "gender", Value: gender})
	}
	if opts.ProfileImage != nil {
		img := model.IMG{SRC: opts.ProfileImage.SRC}
		if err := img.LoadFromURL(); err != nil {
			return nil, err
		}
		update = append(update, bson.E{Key: "profile_image", Value: img})
	}
	var user *model.User
	var customer model.Customer
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, _ = ci.App.User.GetUserByID(opts.UserID)
	}()
	if update == nil {
		filter := bson.M{"_id": opts.ID}
		if err := ci.DB.Collection(model.CustomerColl).FindOne(context.TODO(), filter).Decode(&customer); err != nil {
			if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
				return nil, errors.Wrapf(err, "customer with id:%s not found", opts.ID.Hex())
			}
			return nil, errors.Wrapf(err, "failed to find customer with id:%s", opts.ID.Hex())
		}
	} else {
		filter := bson.M{"_id": opts.ID}
		updateQuery := bson.M{"$set": update}
		queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		if err := ci.DB.Collection(model.CustomerColl).FindOneAndUpdate(context.TODO(), filter, updateQuery, queryOpts).Decode(&customer); err != nil {
			if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
				return nil, errors.Wrapf(err, "customer with id:%s not found", opts.ID.Hex())
			}
			return nil, errors.Wrapf(err, "failed to find customer with id:%s", opts.ID.Hex())
		}
	}
	wg.Wait()
	claim := ci.App.User.GetUserClaim(user, &customer)
	return claim, nil
}

func (ci *CustomerImpl) AddBrandFollowing(sc mongo.SessionContext, opts *schema.AddBrandFollowerOpts) error {
	filter := bson.M{
		"_id": opts.CustomerID,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"brand_following": opts.BrandID,
		},
		"$inc": bson.M{
			"brand_follow_count": 1,
		},
	}

	res, err := ci.DB.Collection(model.CustomerColl).UpdateOne(sc, filter, update)
	if err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed add brand id into following field")
		return errors.Wrap(err, "failed to add brand following")
	}

	if res.MatchedCount == 0 {
		return errors.Errorf("customer with id:%s not found", opts.CustomerID.Hex())
	}

	return nil
}

func (ci *CustomerImpl) RemoveBrandFollowing(sc mongo.SessionContext, opts *schema.AddBrandFollowerOpts) error {
	filter := bson.M{
		"_id": opts.CustomerID,
	}

	update := bson.M{
		"$pull": bson.M{
			"brand_following": opts.BrandID,
		},
		"$inc": bson.M{
			"brand_follow_count": -1,
		},
	}

	res, err := ci.DB.Collection(model.CustomerColl).UpdateOne(sc, filter, update)
	if err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed remove brand id from following field")
		return errors.Wrap(err, "failed to remove brand following")
	}

	if res.MatchedCount == 0 {
		return errors.Errorf("customer with id:%s not found", opts.CustomerID.Hex())
	}

	return nil
}

func (ci *CustomerImpl) AddInfluencerFollowing(sc mongo.SessionContext, opts *schema.AddInfluencerFollowerOpts) error {
	filter := bson.M{
		"_id": opts.CustomerID,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"influencer_following": opts.InfluencerID,
		},
		"$inc": bson.M{
			"influencer_follow_count": 1,
		},
	}

	res, err := ci.DB.Collection(model.CustomerColl).UpdateOne(sc, filter, update)
	if err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed add influencer id into following field")
		return errors.Wrap(err, "failed to add influencer following")
	}

	if res.MatchedCount == 0 {
		return errors.Errorf("customer with user_id:%s not found", opts.CustomerID.Hex())
	}

	return nil
}

func (ci *CustomerImpl) RemoveInfluencerFollowing(sc mongo.SessionContext, opts *schema.AddInfluencerFollowerOpts) error {
	filter := bson.M{
		"_id": opts.CustomerID,
	}

	update := bson.M{
		"$pull": bson.M{
			"influencer_following": opts.InfluencerID,
		},
		"$inc": bson.M{
			"influencer_follow_count": -1,
		},
	}

	res, err := ci.DB.Collection(model.CustomerColl).UpdateOne(sc, filter, update)
	if err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed remove influencer id into following field")
		return errors.Wrap(err, "failed to remove influencer following")
	}

	if res.MatchedCount == 0 {
		return errors.Errorf("customer with id:%s not found", opts.CustomerID.Hex())
	}

	return nil
}

func (ci *CustomerImpl) AddAddress(opts *schema.AddAddressOpts) (*schema.AddAddressResp, error) {
	findQuery := bson.M{
		"user_id": opts.UserID,
	}

	//if line 2 exist addding comma (,) after it
	if opts.Line2 != "" {
		opts.Line2 += ", "

	}
	if opts.District != "" {
		opts.District += ", "
	}
	plainAddress := opts.Line1 + ", " + opts.Line2 + opts.District + opts.City + ", " + opts.State.Name + ", " + opts.Country.Name + " - " + opts.PostalCode

	address := model.Address{
		ID:                primitive.NewObjectID(),
		DisplayName:       opts.DisplayName,
		Line1:             opts.Line1,
		Line2:             opts.Line2,
		District:          opts.District,
		City:              opts.City,
		State:             opts.State,
		PostalCode:        opts.PostalCode,
		Country:           opts.Country,
		PlainAddress:      plainAddress,
		IsBillingAddress:  opts.IsBillingAddress,
		IsShippingAddress: opts.IsShippingAddress,
		IsDefaultAddress:  opts.IsDefaultAddress,
		ContactNumber:     opts.ContactNumber,
	}
	updateQuery := bson.M{
		"$push": bson.M{
			"addresses": address,
		},
	}
	res, err := ci.DB.Collection(model.CustomerColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to add address")
	}
	if res.MatchedCount == 0 {
		return nil, errors.Errorf("unable to find user with id: %s", opts.UserID.Hex())
	}
	addressRes := schema.AddAddressResp{
		ID:            address.ID,
		DisplayName:   address.DisplayName,
		ContactNumber: address.ContactNumber,
		Line1:         address.Line1,
		Line2:         address.Line2,
		District:      address.District,
		City:          address.City,
		State:         address.State,
		PostalCode:    address.PostalCode,
		PlainAddress:  address.PlainAddress,
	}
	return &addressRes, nil
}

func (ci *CustomerImpl) GetAddresses(id primitive.ObjectID) ([]model.Address, error) {
	findQuery := bson.M{
		"user_id": id,
	}
	var customer model.Customer
	err := ci.DB.Collection(model.CustomerColl).FindOne(context.TODO(), findQuery).Decode(&customer)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch addresses")
	}
	return customer.Addresses, nil
}

func (ci *CustomerImpl) GetAppCustomerInfo(id primitive.ObjectID) (*schema.GetCustomerProfileInfoResp, error) {
	var resp []schema.GetCustomerProfileInfoResp

	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"_id": id,
			},
		},
	}

	lookupStage := bson.D{
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
			},
		},
	}

	ctx := context.TODO()
	cur, err := ci.DB.Collection(model.CustomerColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, setStage})
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get customer profile info")
	}

	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to find customer info")
	}
	if len(resp) == 0 {
		return nil, nil
	}
	if !resp[0].UserInfo.EmailVerifiedAt.IsZero() {
		resp[0].UserInfo.EmailVerified = true
	}
	if !resp[0].UserInfo.PhoneVerifiedAt.IsZero() {
		resp[0].UserInfo.PhoneVerified = true
	}
	return &resp[0], nil
}

func (ci *CustomerImpl) RemoveAddress(userID, addressID primitive.ObjectID) error {

	filter := bson.M{
		"user_id":       userID,
		"addresses._id": addressID,
	}
	update := bson.M{
		"$pull": bson.M{
			"addresses": bson.M{
				"_id": addressID,
			},
		},
	}
	resp, err := ci.DB.Collection(model.CustomerColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrapf(err, "unable to remove address with id: %s, from user with id %s", addressID, userID)
	}
	if resp.MatchedCount == 0 {
		return errors.Errorf("unable to find from user with id %s", addressID, userID)
	}
	return nil
}

func (ci *CustomerImpl) EditAddress(opts *schema.EditAddressOpts) error {

	filter := bson.M{
		"user_id":       opts.UserID,
		"addresses._id": opts.AddressID,
	}

	//if line 2 or district exist addding comma (,) after it
	if opts.Line2 != "" {
		opts.Line2 += ", "

	}
	if opts.District != "" {
		opts.District += ", "
	}
	plainAddress := opts.Line1 + ", " + opts.Line2 + opts.District + opts.City + ", " + opts.State.Name + ", " + opts.Country.Name + " - " + opts.PostalCode

	address := model.Address{
		ID:                opts.AddressID,
		DisplayName:       opts.DisplayName,
		Line1:             opts.Line1,
		Line2:             opts.Line2,
		District:          opts.District,
		City:              opts.City,
		State:             opts.State,
		PostalCode:        opts.PostalCode,
		Country:           opts.Country,
		PlainAddress:      plainAddress,
		IsBillingAddress:  opts.IsBillingAddress,
		IsShippingAddress: opts.IsShippingAddress,
		IsDefaultAddress:  opts.IsDefaultAddress,
		ContactNumber:     opts.ContactNumber,
	}

	update := bson.M{
		"$set": bson.M{
			"addresses.$": address,
		},
	}

	resp, err := ci.DB.Collection(model.CustomerColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrapf(err, "unable to edit address with id: %s, from user with id %s", opts.AddressID, opts.UserID)
	}
	if resp.MatchedCount == 0 {
		return errors.Errorf("unable to find from user with id %s", opts.AddressID, opts.UserID)
	}
	return nil
}
