//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_category.go -package=mock go-app/app Category

package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"reflect"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Category defines all the methods for category service to implement
type Category interface {
	CreateCategory(*schema.CreateCategoryOpts) (*schema.CreateCategoryResp, error)
	EditCategory(*schema.EditCategoryOpts) (*schema.EditCategoryResp, error)

	GetAncestorsByID(id primitive.ObjectID) ([]primitive.ObjectID, error)
	GetMainCategoriesMap() (map[string]schema.GetMainCategoriesMapResp, error)
	GetCategoryPath(primitive.ObjectID) (string, error)

	GetCategories() ([]schema.GetCategoriesResp, error)
	GetCategoriesBasic() ([]schema.GetCategoriesBasicResp, error)
	GetMainParentCategories() ([]schema.GetParentCategoriesResp, error)
	GetMainCategoriesByParentID(primitive.ObjectID) ([]schema.GetMainCategoriesByParentIDResp, error)
	GetSubCategoriesByParentID(primitive.ObjectID) ([]schema.GetSubCategoriesByParentIDResp, error)
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
func (ci *CategoryImpl) CreateCategory(opts *schema.CreateCategoryOpts) (*schema.CreateCategoryResp, error) {
	category := model.Category{
		Name:   opts.Name,
		Slug:   UniqueSlug(opts.Name),
		IsMain: &opts.IsMain,
	}

	if opts.FeaturedImage != nil {
		category.FeaturedImage = &model.IMG{
			SRC: opts.FeaturedImage.SRC,
		}
		// Setting featured image and thumbnail from image source
		if err := category.FeaturedImage.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid featured image url")
		}
	}
	if opts.Thumbnail != nil {
		category.Thumbnail = &model.IMG{
			SRC: opts.Thumbnail.SRC,
		}
		if err := category.Thumbnail.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid thumbnail url")
		}
	}

	// Verify and setting parentID and ancestors if parentID is passed
	if !opts.ParentID.IsZero() {
		category.ParentID = opts.ParentID
		parentAncestors, err := ci.GetAncestorsByID(category.ParentID)
		if err != nil {
			return nil, err
		}
		// Setting up new category ancestors by appending parent id into parent category's ancestor
		category.AncestorsID = append(parentAncestors, category.ParentID)
	}

	res, err := ci.DB.Collection(model.CategoryColl).InsertOne(context.Background(), category)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert new category")
	}

	return &schema.CreateCategoryResp{
		ID:            res.InsertedID.(primitive.ObjectID),
		Name:          category.Name,
		Slug:          category.Slug,
		ParentID:      category.ParentID,
		AncestorID:    category.AncestorsID,
		Thumbnail:     category.Thumbnail,
		FeaturedImage: category.FeaturedImage,
		IsMain:        *category.IsMain,
	}, nil
}

// GetAncestorsByID returns ancestors list (list of objectIDs) of matching id
func (ci *CategoryImpl) GetAncestorsByID(id primitive.ObjectID) ([]primitive.ObjectID, error) {
	var category model.Category

	fields := bson.M{"ancestors_id": 1}
	filter := bson.M{"_id": id}

	opts := options.FindOne().SetProjection(fields)
	err := ci.DB.Collection(model.CategoryColl).FindOne(context.Background(), filter, opts).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("category with id:%s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "failed to find category: %s", id.Hex())
	}
	return category.AncestorsID, nil
}

// EditCategory updates an existing category document
/*
	Changing the parentID or ancestorsID is a very expensive operations thus not supported.
*/
func (ci *CategoryImpl) EditCategory(opts *schema.EditCategoryOpts) (*schema.EditCategoryResp, error) {
	var updateField model.Category = model.Category{}

	if opts.Name != "" {
		updateField.Name = opts.Name
	}

	if opts.FeaturedImage != nil {
		updateField.FeaturedImage = &model.IMG{SRC: opts.FeaturedImage.SRC}
		if err := updateField.FeaturedImage.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid featured image url")
		}
	}
	if opts.Thumbnail != nil {
		updateField.Thumbnail = &model.IMG{SRC: opts.Thumbnail.SRC}
		if err := updateField.Thumbnail.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "invalid thumbnail url")
		}
	}

	if opts.IsMain != nil {
		updateField.IsMain = opts.IsMain
	}

	if reflect.DeepEqual(updateField, model.Category{}) {
		return nil, errors.New("no fields found to update")
	}

	filter := bson.M{"_id": opts.ID}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := ci.DB.Collection(model.CategoryColl).FindOneAndUpdate(context.Background(), filter, bson.M{"$set": updateField}, qOpts).Decode(&updateField); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("category with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrap(err, "failed to update category")
	}

	return &schema.EditCategoryResp{
		ID:            updateField.ID,
		Name:          updateField.Name,
		Slug:          updateField.Slug,
		FeaturedImage: updateField.FeaturedImage,
		Thumbnail:     updateField.Thumbnail,
		IsMain:        *updateField.IsMain,
		ParentID:      updateField.ParentID,
		AncestorID:    updateField.AncestorsID,
	}, nil
}

// GetCategoryPath takes category id as input and returns its path from topmost parent to itself as string
func (ci *CategoryImpl) GetCategoryPath(id primitive.ObjectID) (string, error) {
	var path string = ""

	ancestorsID, err := ci.GetAncestorsByID(id)
	if err != nil {
		return "", err
	}

	for _, id := range ancestorsID {
		path += fmt.Sprintf("/%s", id.Hex())
	}

	path += fmt.Sprintf("/%s", id.Hex())

	return path, nil
}

// GetMainCategoriesMap takes multiple category id as input,
// query the DB for the id and returns the map with id as key and DB response as value
func (ci *CategoryImpl) GetMainCategoriesMap() (map[string]schema.GetMainCategoriesMapResp, error) {
	ctx := context.Background()
	res := make(map[string]schema.GetMainCategoriesMapResp)

	filter := bson.M{"is_main": true}
	cur, err := ci.DB.Collection(model.CategoryColl).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for categories")
	}
	for cur.Next(ctx) {
		var result model.Category
		err := cur.Decode(&result)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode category")
		}
		res[result.ID.Hex()] = schema.GetMainCategoriesMapResp{
			ID:            result.ID,
			Name:          result.Name,
			ParentID:      result.ParentID,
			AncestorID:    result.AncestorsID,
			Thumbnail:     result.Thumbnail,
			FeaturedImage: result.FeaturedImage,
		}
	}
	return res, nil
}

// GetMainParentCategories returns all the parent categories which are is_main true
func (ci *CategoryImpl) GetMainParentCategories() ([]schema.GetParentCategoriesResp, error) {
	var categories []schema.GetParentCategoriesResp
	ctx := context.Background()

	filter := bson.M{"is_main": true, "parent_id": nil}
	cur, err := ci.DB.Collection(model.CategoryColl).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to find parent categories")
	}
	if err := cur.All(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "failed to decode categories")
	}
	return categories, nil
}

// GetMainCategoriesByParentID returns all the categories which are direct parent of category passed in argument
func (ci *CategoryImpl) GetMainCategoriesByParentID(id primitive.ObjectID) ([]schema.GetMainCategoriesByParentIDResp, error) {
	var categories []schema.GetMainCategoriesByParentIDResp
	ctx := context.Background()
	filter := bson.M{"is_main": true, "parent_id": id}
	cur, err := ci.DB.Collection(model.CategoryColl).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrapf(err, "query failed to find categories with parent:%s", id)
	}
	if err := cur.All(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "failed to decode categories")
	}
	return categories, nil
}

// GetSubCategoriesByParentID returns only the name of all the children categories matching with id passed as parent_id
func (ci *CategoryImpl) GetSubCategoriesByParentID(id primitive.ObjectID) ([]schema.GetSubCategoriesByParentIDResp, error) {
	var categories []schema.GetSubCategoriesByParentIDResp
	ctx := context.Background()
	filter := bson.M{"is_main": true, "parent_id": id}
	cur, err := ci.DB.Collection(model.CategoryColl).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrapf(err, "query failed to find categories with parent:%s", id)
	}
	if err := cur.All(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "failed to decode categories")
	}
	return categories, nil
}

// GetCategories returns all the categories document
func (ci *CategoryImpl) GetCategories() ([]schema.GetCategoriesResp, error) {
	var categories []schema.GetCategoriesResp
	ctx := context.Background()
	filter := bson.M{}
	cur, err := ci.DB.Collection(model.CategoryColl).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to find categories")
	}
	if err := cur.All(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "failed to decode categories")
	}
	return categories, nil
}

// GetCategoriesBasic returns all the categories document but only id,name and is_main field
func (ci *CategoryImpl) GetCategoriesBasic() ([]schema.GetCategoriesBasicResp, error) {
	var categories []schema.GetCategoriesBasicResp
	ctx := context.Background()
	opts := options.Find().SetProjection(bson.M{"_id": 1, "name": 1, "is_main": 1})
	filter := bson.M{}
	cur, err := ci.DB.Collection(model.CategoryColl).Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to find categories")
	}
	if err := cur.All(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "failed to decode categories")
	}
	return categories, nil
}
