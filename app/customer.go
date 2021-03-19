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
