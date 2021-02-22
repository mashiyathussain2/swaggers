//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_catalog.go -package=mock go-app/app KeeperCatalog

package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// KeeperCatalog service allows `Keeper` to execute admin operations.
type KeeperCatalog interface {
	CreateCatalog(*schema.CreateCatalogOpts) (*schema.CreateCatalogResp, error)
	EditCatalog(*schema.EditCatalogOpts) (*schema.EditCatalogResp, error)
	AddVariant(primitive.ObjectID, *schema.AddVariantOpts) (*schema.AddVariantResp, error)

	GetBasicCatalogInfo(*schema.GetBasicCatalogFilter) ([]schema.GetBasicCatalogResp, error)
	GetCatalogFilter() (*schema.GetCatalogFilterResp, error)
	// EditVariant(primitive.ObjectID, *schema.CreateVariantOpts)
	// DeleteVariant(primitive.ObjectID)
}

// UserCatalog service allows `app` or user api to perform operations on catalog.
type UserCatalog interface{}

// KeeperCatalogImpl implements keeper related operations
type KeeperCatalogImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// KeeperCatalogOpts contains arguments required to create a new instance of KeeperCatalog
type KeeperCatalogOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitKeeperCatalog returns KeeperCatalog instance
func InitKeeperCatalog(opts *KeeperCatalogOpts) KeeperCatalog {
	return &KeeperCatalogImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

// CreateCatalog inserts a new catalog document with specified data into the database
func (kc *KeeperCatalogImpl) CreateCatalog(opts *schema.CreateCatalogOpts) (*schema.CreateCatalogResp, error) {
	c := model.Catalog{
		Name:        opts.Name,
		LName:       strings.ToLower(opts.Name),
		Description: opts.Description,
		Keywords:    opts.Keywords,
		HSNCode:     opts.HSNCode,
		Slug:        UniqueSlug(opts.Name),
		BasePrice:   model.SetINRPrice(float32(opts.BasePrice)),
		RetailPrice: model.SetINRPrice(float32(opts.RetailPrice)),
		CreatedAt:   time.Now().UTC(),
	}

	// If variants are passed in the opts then setting variants in catalog model
	if opts.VariantType != "" {
		c.VariantType = opts.VariantType
		for _, variant := range opts.Variants {
			c.Variants = append(c.Variants, *kc.createVariant(&variant))
		}
	}

	// Checking if brands id exists otherwise setting up brandID
	exists, err := kc.App.Brand.CheckBrandIDExists(context.Background(), opts.BrandID)
	if err != nil || !exists {
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find brand with id: %s", opts.BrandID.Hex())
		}
		return nil, errors.Errorf("brand id %s does not exists", opts.BrandID.Hex())
	}
	c.BrandID = opts.BrandID

	// setting catalog specifications
	for _, specOpt := range opts.Specifications {
		c.Specifications = append(c.Specifications, model.Specification{Name: specOpt.Name, Value: specOpt.Value})
	}

	// setting up filter attributes
	for _, attr := range opts.FilterAttribute {
		c.FilterAttribute = append(c.FilterAttribute, model.Attribute{Name: attr.Name, Value: attr.Value})
	}

	// If eta is passed then setting up the eta
	if opts.ETA != nil {
		c.ETA = &model.ETA{
			Min:  int(opts.ETA.Min),
			Max:  int(opts.ETA.Max),
			Unit: opts.ETA.Unit,
		}
	}

	// Setting up category path
	for _, id := range opts.CategoryID {
		path, err := kc.App.Category.GetCategoryPath(id)
		if err != nil {
			return nil, err
		}
		c.Paths = append(c.Paths, path)
	}

	// Inserting the document in the DB
	res, err := kc.DB.Collection(model.CatalogColl).InsertOne(context.Background(), c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert catalog in db")
	}

	resp := &schema.CreateCatalogResp{
		ID:              res.InsertedID.(primitive.ObjectID),
		Name:            c.Name,
		Paths:           c.Paths,
		BrandID:         c.BrandID,
		Keywords:        c.Keywords,
		Description:     c.Description,
		FeaturedImage:   c.FeaturedImage,
		Specifications:  c.Specifications,
		FilterAttribute: c.FilterAttribute,
		VariantType:     c.VariantType,
		Variants:        c.Variants,
		BasePrice:       *c.BasePrice,
		RetailPrice:     *c.RetailPrice,
		HSNCode:         c.HSNCode,
		Status:          c.Status,
		ETA:             c.ETA,
		CreatedAt:       c.CreatedAt,
	}

	return resp, nil
}

// EditCatalog edits an existing catalog
func (kc *KeeperCatalogImpl) EditCatalog(opts *schema.EditCatalogOpts) (*schema.EditCatalogResp, error) {
	c := model.Catalog{}
	if opts.Name != "" {
		c.Name = opts.Name
	}
	if opts.Description != "" {
		c.Description = opts.Description
	}
	if len(opts.Keywords) != 0 {
		c.Keywords = opts.Keywords
	}
	if len(opts.CategoryID) != 0 {
		for _, id := range opts.CategoryID {
			path, err := kc.App.Category.GetCategoryPath(id)
			if err != nil {
				return nil, err
			}
			c.Paths = append(c.Paths, path)
		}
	}
	if opts.ETA != nil {
		c.ETA = &model.ETA{
			Min:  int(opts.ETA.Min),
			Max:  int(opts.ETA.Max),
			Unit: opts.ETA.Unit,
		}
	}
	if len(opts.Specifications) != 0 {
		for _, specOpt := range opts.Specifications {
			c.Specifications = append(c.Specifications, model.Specification{Name: specOpt.Name, Value: specOpt.Value})
		}
	}
	if len(opts.FilterAttribute) != 0 {
		for _, attr := range opts.FilterAttribute {
			c.FilterAttribute = append(c.FilterAttribute, model.Attribute{Name: attr.Name, Value: attr.Value})
		}
	}
	if opts.HSNCode != "" {
		c.HSNCode = opts.HSNCode
	}
	if opts.BasePrice != 0 {
		c.BasePrice = model.SetINRPrice(float32(opts.BasePrice))
	}
	if opts.RetailPrice != 0 {
		c.RetailPrice = model.SetINRPrice(float32(opts.RetailPrice))
	}

	if reflect.DeepEqual(model.Catalog{}, c) {
		return nil, errors.New("no fields found to update")
	}
	c.UpdatedAt = time.Now().UTC()
	filter := bson.M{
		"_id": opts.ID,
	}
	update := bson.M{"$set": c}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := kc.DB.Collection(model.CatalogColl).FindOneAndUpdate(context.TODO(), filter, update, qOpts).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("catalog with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrap(err, "failed to update catalog")
	}

	return &schema.EditCatalogResp{
		ID:              c.ID,
		Name:            c.Name,
		Description:     c.Description,
		Paths:           c.Paths,
		Keywords:        c.Keywords,
		Specifications:  c.Specifications,
		FilterAttribute: c.FilterAttribute,
		HSNCode:         c.HSNCode,
		BasePrice:       *c.BasePrice,
		RetailPrice:     *c.RetailPrice,
		ETA:             c.ETA,
		UpdatedAt:       c.UpdatedAt,
	}, nil
}

func (kc *KeeperCatalogImpl) createVariant(opts *schema.CreateVariantOpts) *model.Variant {
	return &model.Variant{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		SKU:       opts.SKU,
		Attribute: opts.Attribute,
	}
}

// AddVariant adds a new variant to an existing catalog
func (kc *KeeperCatalogImpl) AddVariant(catalogID primitive.ObjectID, opts *schema.AddVariantOpts) (*schema.AddVariantResp, error) {
	ctx := context.TODO()

	if err := kc.validateAddVariant(ctx, catalogID, opts); err != nil {
		return nil, err
	}

	v := model.Variant{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		SKU:       opts.SKU,
		Attribute: opts.Attribute,
	}
	// TODO: Add logic to create inventory

	filter := bson.M{"_id": catalogID}
	update := bson.M{
		"$push": bson.M{
			"variants": v,
		},
	}
	res, err := kc.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("catalog with id:%s not found", catalogID.Hex())
		}
		return nil, errors.Wrap(err, "failed to add variant in catalog")
	}

	if res.MatchedCount == 0 {
		return nil, errors.Errorf("catalog with id:%s not found", catalogID.Hex())
	}
	if res.ModifiedCount == 0 {
		return nil, errors.Wrap(err, "failed to add variant in catalog")
	}
	return &schema.CreateVariantResp{
		ID:        v.ID,
		SKU:       v.SKU,
		Attribute: v.Attribute,
	}, nil
}

func (kc *KeeperCatalogImpl) validateAddVariant(ctx context.Context, catalogID primitive.ObjectID, opts *schema.AddVariantOpts) error {
	var c model.Catalog
	if err := kc.DB.Collection(model.CatalogColl).FindOne(ctx, bson.M{"_id": catalogID}).Decode(&c); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return errors.Errorf("catalog with id:%s not found", catalogID.Hex())
		}
		return errors.Wrap(err, "failed to find catalog")
	}

	// Checking for same variant type if variant type already exists in DB
	if c.VariantType != "" {
		if c.VariantType != opts.VariantType {
			return errors.Errorf("cannot set variant type %s: catalog is set with variant type %s", opts.VariantType, c.VariantType)
		}
	}

	// checking for duplicate SKU and attribute
	for _, v := range c.Variants {
		if v.Attribute == opts.Attribute {
			return errors.Errorf("variant with attribute %s already exists", opts.Attribute)
		}
		if v.SKU == opts.SKU {
			return errors.Errorf("variant with sku %s already exists", opts.SKU)
		}
	}
	return nil
}

// GetBasicCatalogInfo returns list of catalog with basic detail such as name, thumbnail, category, retail price, status
/*
	Filters By brand_id, category_id
*/
func (kc *KeeperCatalogImpl) GetBasicCatalogInfo(filter *schema.GetBasicCatalogFilter) ([]schema.GetBasicCatalogResp, error) {
	ctx := context.TODO()
	queryFilter := bson.M{}
	if len(filter.BrandID) > 0 {
		queryFilter["brand_id"] = bson.M{
			"$in": filter.BrandID,
		}
	}
	if len(filter.CategoryID) > 0 {
		var regexID []string
		var regexStr string
		for _, id := range filter.CategoryID {
			regexID = append(regexID, id.Hex())
		}
		regexStr = strings.Join(regexID, "|")
		queryFilter["category_path"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: regexStr,
				Options: "i",
			},
		}
	}

	var res []schema.GetBasicCatalogResp
	cur, err := kc.DB.Collection(model.CatalogColl).Find(ctx, queryFilter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for catalog")
	}
	if err := cur.All(ctx, &res); err != nil {
		return nil, errors.Wrap(err, "failed to find catalog")
	}

	return res, nil
}

// GetCatalogFilter returns list of filter supported for filter and their respective values
func (kc *KeeperCatalogImpl) GetCatalogFilter() (*schema.GetCatalogFilterResp, error) {
	// ctx := context.TODO()
	c, err := kc.App.Category.GetCategoriesBasic()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get category filter")
	}

	return &schema.GetCatalogFilterResp{
		Category: c,
	}, nil
}
