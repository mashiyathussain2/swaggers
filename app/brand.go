//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_brand.go -package=mock go-app/app Brand

package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Brand service contains all the CRUD operations related to brands
type Brand interface {
	CreateBrand(*schema.CreateBrandOpts) (*schema.CreateBrandResp, error)
	EditBrand(*schema.EditBrandOpts) (*schema.EditBrandResp, error)
	CheckBrandIDExists(ctx context.Context, id primitive.ObjectID) (bool, error)
}

// BrandImpl implements brand service methods
type BrandImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// BrandOpts contains args required to create a new instance of brand service
type BrandOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitBrand returns brand service implementation instance
func InitBrand(opts *BrandOpts) Brand {
	return &BrandImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

// CreateBrand create a new brand document in the database
func (b *BrandImpl) CreateBrand(opts *schema.CreateBrandOpts) (*schema.CreateBrandResp, error) {

	brand := model.Brand{
		Name:           opts.Name,
		RegisteredName: strings.ToLower(opts.RegisteredName),
		Slug:           UniqueSlug(opts.Name),
		Description:    opts.Description,
		Fulfillment: &model.Fulfillment{
			Email: opts.FulfillmentEmail,
		},
		WebsiteLink: opts.WebsiteLink,
	}

	res, err := b.DB.Collection(model.BrandColl).InsertOne(context.Background(), brand)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert brand in db")
	}

	return &schema.CreateBrandResp{
		ID:          res.InsertedID.(primitive.ObjectID),
		Name:        brand.Name,
		Description: brand.Description,
		Slug:        brand.Slug,
		WebsiteLink: brand.WebsiteLink,
		Fulfillment: brand.Fulfillment,
	}, nil
}

// EditBrand edits the brand document and saves in the database
func (b *BrandImpl) EditBrand(opts *schema.EditBrandOpts) (*schema.EditBrandResp, error) {
	var updateField model.Brand = model.Brand{}
	if opts.Name != "" {
		updateField.Name = opts.Name
	}
	if opts.Description != "" {
		updateField.Description = opts.Description
	}
	if opts.FulfillmentEmail != "" {
		updateField.Fulfillment = &model.Fulfillment{
			Email: opts.FulfillmentEmail,
		}
	}
	if opts.WebsiteLink != "" {
		updateField.WebsiteLink = opts.WebsiteLink
	}
	if updateField == (model.Brand{}) {
		return nil, errors.New("no fields found to update")
	}

	var res schema.EditBrandResp
	filter := bson.M{"_id": opts.ID}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := b.DB.Collection(model.BrandColl).FindOneAndUpdate(context.Background(), filter, bson.M{"$set": updateField}, qOpts).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("brand with id:%s not found", opts.ID)
		}
		return nil, errors.Wrap(err, "failed to update brand")
	}

	return &res, nil
}

// CheckBrandIDExists return true/false based on if passed id exists in brand collection
func (b *BrandImpl) CheckBrandIDExists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"_id": id,
	}
	count, err := b.DB.Collection(model.BrandColl).CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}
