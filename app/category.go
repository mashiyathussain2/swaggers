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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Category defines all the methods for category service to implement
type Category interface {
	GetAncestorsByID(id primitive.ObjectID) ([]primitive.ObjectID, error)
	CreateCategory(*schema.CreateCategoryOpts) (*schema.CreateCategoryResp, error)
}

// CategoryOpts contains args required to create a new instance of category service
type CategoryOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// CategoryImpl implements Category service methods
type CategoryImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitCategory returns a new instance of CategoryImpl
func InitCategory(opts *CategoryOpts) Category {
	return &CategoryImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

// CreateCategory creates a new category in the db category collection
func (c *CategoryImpl) CreateCategory(opts *schema.CreateCategoryOpts) (*schema.CreateCategoryResp, error) {
	category := model.Category{
		Name: opts.Name,
		Slug: UniqueSlug(opts.Name),
		FeaturedImage: &model.IMG{
			SRC: opts.FeaturedImage.SRC,
		},
		Thumbnail: &model.IMG{
			SRC: opts.FeaturedImage.SRC,
		},
		IsMain: opts.IsMain,
	}

	// Setting featured image and thumbnail from image source
	if err := category.FeaturedImage.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "invalid featured image url")
	}
	if err := category.Thumbnail.LoadFromURL(); err != nil {
		return nil, errors.Wrap(err, "invalid thumbnail url")
	}

	// Verify and setting parentID and ancestors if parentID is passed
	if !opts.ParentID.IsZero() {
		category.ParentID = opts.ParentID
		parentAncestors, err := c.GetAncestorsByID(category.ParentID)
		if err != nil {
			return nil, err
		}
		// Setting up new category ancestors by appending parent id into parent category's ancestor
		category.AncestorID = append(parentAncestors, category.ParentID)
	}

	res, err := c.DB.Collection(model.CategoryColl).InsertOne(context.Background(), category)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert new category")
	}

	return &schema.CreateCategoryResp{
		ID:            res.InsertedID.(primitive.ObjectID),
		Name:          category.Name,
		Slug:          category.Slug,
		ParentID:      category.ParentID,
		AncestorID:    category.AncestorID,
		Thumbnail:     category.Thumbnail,
		FeaturedImage: category.FeaturedImage,
		IsMain:        category.IsMain,
	}, nil
}

// GetAncestorsByID returns ancestors list (list of objectIDs) of matching id
func (c *CategoryImpl) GetAncestorsByID(id primitive.ObjectID) ([]primitive.ObjectID, error) {
	var category model.Category

	fields := bson.M{"ancestors_id": 1}
	filter := bson.M{"_id": id}

	opts := options.FindOne().SetProjection(fields)
	err := c.DB.Collection(model.CategoryColl).FindOne(context.Background(), filter, opts).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("category with id:%s not found", id)
		}
		return nil, errors.Wrapf(err, "failed to find category: %s", id)
	}
	return category.AncestorID, nil
}
