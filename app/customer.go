package app

import (
	"context"
	"go-app/model"
	"go-app/schema"

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
	UpdateCustomer(*schema.UpdateCustomerOpts) (*schema.GetCustomerInfoResp, error)

	AddBrandFollowing(mongo.SessionContext, *schema.AddBrandFollowerOpts) error
	AddInfluencerFollowing(mongo.SessionContext, *schema.AddInfluencerFollowerOpts) error
	AddAddress(opts *schema.AddAddressOpts) error
	GetAddresses(primitive.ObjectID) ([]model.Address, error)
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
func (ci *CustomerImpl) UpdateCustomer(opts *schema.UpdateCustomerOpts) (*schema.GetCustomerInfoResp, error) {
	var update bson.D
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

	if update == nil {
		return nil, errors.New("no field update found for customer")
	}
	filter := bson.M{"_id": opts.ID}
	updateQuery := bson.M{"$set": update}
	queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var resp schema.GetCustomerInfoResp
	if err := ci.DB.Collection(model.CustomerColl).FindOneAndUpdate(context.TODO(), filter, updateQuery, queryOpts).Decode(&resp); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "customer with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrapf(err, "failed to find customer with id:%s", opts.ID.Hex())
	}
	return &resp, nil
}

func (ci *CustomerImpl) AddBrandFollowing(sc mongo.SessionContext, opts *schema.AddBrandFollowerOpts) error {
	filter := bson.M{
		"user_id": opts.UserID,
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
		return errors.Errorf("customer with user_id:%s not found", opts.UserID.Hex())
	}

	return nil
}

func (ci *CustomerImpl) AddInfluencerFollowing(sc mongo.SessionContext, opts *schema.AddInfluencerFollowerOpts) error {
	filter := bson.M{
		"user_id": opts.UserID,
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
		return errors.Errorf("customer with user_id:%s not found", opts.UserID.Hex())
	}

	return nil
}

func (ci *CustomerImpl) AddAddress(opts *schema.AddAddressOpts) error {
	findQuery := bson.M{
		"user_id": opts.UserID,
	}
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
		PlainAddress:      opts.PlainAddress,
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
		return errors.Wrapf(err, "unable to add address")
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find user with id: %s", opts.UserID.Hex())
	}
	return nil
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
