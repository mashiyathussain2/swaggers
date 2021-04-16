//go:generate $GOBIN/mockgen -destination=./../mock/mock_catalog.go -package=mock go-app/app KeeperCatalog

package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
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
	GetBasicCatalogInfo(*schema.GetBasicCatalogFilter) ([]schema.GetBasicCatalogResp, error)
	GetCatalogFilter() (*schema.GetCatalogFilterResp, error)
	AddVariant(*schema.AddVariantOpts) (*schema.AddVariantResp, error)
	KeeperSearchCatalog(*schema.KeeperSearchCatalogOpts) ([]schema.KeeperSearchCatalogResp, error)
	DeleteVariant(*schema.DeleteVariantOpts) error
	UpdateCatalogStatus(*schema.UpdateCatalogStatusOpts) ([]schema.UpdateCatalogStatusResp, error)
	CheckCatalogIDsExists(context.Context, []primitive.ObjectID) (int64, error)
	GetCatalogByIDs(context.Context, []primitive.ObjectID) ([]schema.GetCatalogResp, error)
	AddCatalogContent(*schema.AddCatalogContentOpts) (*schema.PayloadVideo, []error)
	AddCatalogContentImage(*schema.AddCatalogContentImageOpts) []error
	GetKeeperCatalogContent(primitive.ObjectID) ([]schema.CatalogContentInfoResp, error)
	GetCatalogContent(id primitive.ObjectID) ([]schema.CatalogContentInfoResp, error)
	GetCatalogsByFilter(*schema.GetCatalogsByFilterOpts) ([]schema.GetCatalogResp, error)
	GetCatalogBySlug(string) (*schema.GetCatalogResp, error)
	GetAllCatalogInfo(primitive.ObjectID) (*schema.GetAllCatalogInfoResp, error)
	GetCollectionCatalogInfo(ids []primitive.ObjectID) ([]schema.GetAllCatalogInfoResp, error)
	GetPebbleCatalogInfo(ids []primitive.ObjectID) ([]schema.GetAllCatalogInfoResp, error)
	SyncCatalog(primitive.ObjectID)
	SyncCatalogs([]primitive.ObjectID)
	SyncCatalogContent(id primitive.ObjectID)
	GetCatalogVariant(primitive.ObjectID, primitive.ObjectID) (*schema.GetCatalogVariantResp, error)
	RemoveContent(*schema.RemoveContentOpts) error
	// EditVariant(primitive.ObjectID, *schema.CreateVariantOpts)
	// DeleteVariant(primitive.ObjectID)
}

// UserCatalog service allows `app` or user api to perform operations on catalog.
type UserCatalog interface {
}

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

	currentTime := time.Now().UTC()
	c := model.Catalog{
		ID:            primitive.NewObjectID(),
		Name:          opts.Name,
		LName:         strings.ToLower(opts.Name),
		Description:   opts.Description,
		Keywords:      opts.Keywords,
		HSNCode:       opts.HSNCode,
		Slug:          UniqueSlug(opts.Name),
		BasePrice:     model.SetINRPrice(float32(opts.BasePrice)),
		RetailPrice:   model.SetINRPrice(float32(opts.RetailPrice)),
		TransferPrice: model.SetINRPrice(float32(opts.TransferPrice)),
		CreatedAt:     currentTime,
	}

	tax := &model.Tax{
		Type: opts.Tax.Type,
	}
	if opts.Tax.Type == model.SingleTax {
		if opts.Tax.Rate == 0 {
			return nil, errors.Errorf("tax rate cannot be 0")
		}
		tax.Rate = opts.Tax.Rate
	} else {
		if len(opts.Tax.TaxRanges) == 0 {
			return nil, errors.Errorf("tax range cannot be empty")
		}
		tax.TaxRanges = opts.Tax.TaxRanges
	}
	c.Tax = tax

	c.FeaturedImage = &model.IMG{
		SRC: opts.FeaturedImage.SRC,
	}

	if err := c.FeaturedImage.LoadFromURL(); err != nil {
		return nil, errors.Wrapf(err, "unable to process featured image for catalog")
	}

	// If variants are passed in the opts then setting variants in catalog model
	if opts.VariantType != "" {
		c.VariantType = opts.VariantType
		for _, variant := range opts.Variants {
			v, err := kc.createVariant(c.ID, &variant)
			if err != nil {
				return nil, err
			}
			c.Variants = append(c.Variants, *v)
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
	c.Status = &model.Status{
		Name:      "Draft",
		Value:     "draft",
		CreatedAt: currentTime,
	}

	// Setting up category path
	for _, id := range opts.CategoryID {
		path, err := kc.App.Category.GetCategoryPath(id)
		if err != nil {
			return nil, err
		}
		c.Paths = append(c.Paths, path)
	}
	c.StatusHistory = []model.Status{
		{
			Name:      "Draft",
			Value:     model.Draft,
			CreatedAt: currentTime,
		},
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
		TransferPrice:   *c.TransferPrice,
		Tax:             c.Tax,
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
		c.LName = strings.ToLower(opts.Name)
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
	if opts.FeaturedImage != nil {
		img := model.IMG{SRC: opts.FeaturedImage.SRC}
		if err := img.LoadFromURL(); err != nil {
			return nil, errors.Wrap(err, "failed to load featured image")
		}
		c.FeaturedImage = &img
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
	if opts.TransferPrice != 0 {
		c.TransferPrice = model.SetINRPrice(float32(opts.TransferPrice))
	}
	if opts.Tax != nil {
		c.Tax = &model.Tax{
			Type: opts.Tax.Type,
		}
		if opts.Tax.Type == model.SingleTax {
			if opts.Tax.Rate == 0 {
				return nil, errors.Errorf("rate cannot be 0")
			}
			c.Tax.Rate = opts.Tax.Rate
			c.Tax.TaxRanges = []model.TaxRange{}
		}
		if opts.Tax.Type == model.MultipleTax {
			if len(opts.Tax.TaxRanges) == 0 {
				return nil, errors.Errorf("tax ranges cannot be empty")
			}
			c.Tax.TaxRanges = opts.Tax.TaxRanges
			c.Tax.Rate = 0
		}

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
		FeaturedImage:   c.FeaturedImage,
		Specifications:  c.Specifications,
		FilterAttribute: c.FilterAttribute,
		HSNCode:         c.HSNCode,
		BasePrice:       *c.BasePrice,
		RetailPrice:     *c.RetailPrice,
		TransferPrice:   *c.TransferPrice,
		ETA:             c.ETA,
		UpdatedAt:       c.UpdatedAt,
		Tax:             *c.Tax,
	}, nil
}

func (kc *KeeperCatalogImpl) createVariant(id primitive.ObjectID, opts *schema.CreateVariantOpts) (*model.Variant, error) {

	cOpts := schema.CreateInventoryOpts{
		CatalogID: id,
		VariantID: primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		SKU:       opts.SKU,
		Unit:      opts.Unit,
	}
	inv, err := kc.App.Inventory.CreateInventory(&cOpts)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create inventory")
	}
	return &model.Variant{
		ID:          cOpts.VariantID,
		SKU:         opts.SKU,
		Attribute:   opts.Attribute,
		InventoryID: inv,
	}, nil
}

// AddVariant adds a new variant to an existing catalog
func (kc *KeeperCatalogImpl) AddVariant(opts *schema.AddVariantOpts) (*schema.AddVariantResp, error) {
	ctx := context.TODO()

	// var catalog model.Catalog
	// err := kc.DB.Collection(model.CatalogColl).FindOne(ctx, bson.M{"_id": opts.ID}).Decode(&catalog)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return nil, errors.Errorf("catalog with id:%s not found", opts.ID.Hex())
	// 	}
	// 	return nil, errors.Wrapf(err, "unable to query for catalog")
	// }
	// if catalog.VariantType != opts.VariantType {
	// 	return nil, errors.Errorf("variant type do not match")
	// }
	if err := kc.validateAddVariant(ctx, opts); err != nil {
		return nil, err
	}

	v := model.Variant{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		SKU:       opts.SKU,
		Attribute: opts.Attribute,
	}

	cOpts := schema.CreateInventoryOpts{
		CatalogID: opts.ID,
		VariantID: v.ID,
		SKU:       v.SKU,
		Unit:      opts.Unit,
	}
	inv, err := kc.App.Inventory.CreateInventory(&cOpts)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create inventory")
	}
	v.InventoryID = inv

	filter := bson.M{"_id": opts.ID}
	update := bson.M{
		"$push": bson.M{
			"variants": v,
		},
		"$set": bson.M{
			"variant_type": opts.VariantType,
		},
	}
	res, err := kc.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("catalog with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrap(err, "failed to add variant in catalog")
	}

	if res.MatchedCount == 0 {
		return nil, errors.Errorf("catalog with id:%s not found", opts.ID.Hex())
	}
	if res.ModifiedCount == 0 {
		return nil, errors.Wrap(err, "failed to add variant in catalog")
	}

	return &schema.CreateVariantResp{
		ID:        v.ID,
		SKU:       v.SKU,
		Attribute: v.Attribute,
		Unit:      opts.Unit,
	}, nil
}

func (kc *KeeperCatalogImpl) validateAddVariant(ctx context.Context, opts *schema.AddVariantOpts) error {
	var c model.Catalog
	if err := kc.DB.Collection(model.CatalogColl).FindOne(ctx, bson.M{"_id": opts.ID}).Decode(&c); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return errors.Errorf("catalog with id:%s not found", opts.ID.Hex())
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

// KeeperSearchCatalog searches catalog based on name and paginates results
func (kc *KeeperCatalogImpl) KeeperSearchCatalog(keeperSearchCatalogOpts *schema.KeeperSearchCatalogOpts) ([]schema.KeeperSearchCatalogResp, error) {
	var keeperSearchCatalogResp []schema.KeeperSearchCatalogResp

	ctx := context.TODO()
	//search using Regex, searches for part
	filter := bson.M{"lname": bson.M{"$regex": strings.ToLower(keeperSearchCatalogOpts.Name)}}

	// filter := bson.M{"$text": bson.M{"$search": keeperSearchCatalogOpts.Name}}

	opts := options.Find().SetProjection(bson.M{
		"catalog_id":     1,
		"name":           1,
		"category_path":  1,
		"base_price":     1,
		"retail_price":   1,
		"status":         1,
		"variants":       1,
		"variant_type":   1,
		"transfer_price": 1,
	}).SetSkip(int64(kc.App.Config.PageSize) * keeperSearchCatalogOpts.Page).SetLimit(int64(kc.App.Config.PageSize))

	cursor, err := kc.DB.Collection(model.CatalogColl).Find(ctx, filter, opts)
	if err != nil {
		kc.Logger.Err(err)
		return nil, err
	}
	if err = cursor.All(ctx, &keeperSearchCatalogResp); err != nil {
		kc.Logger.Err(err)
		return nil, err
	}
	return keeperSearchCatalogResp, nil

}

//DeleteVariant deletes variant from the catalog
func (kc *KeeperCatalogImpl) DeleteVariant(opts *schema.DeleteVariantOpts) error {
	ctx := context.TODO()
	filter := bson.M{"_id": opts.CatalogID, "variants._id": opts.VariantID}
	deleteQuery := bson.M{"$set": bson.M{
		"variants.$.is_deleted": true,
	},
	}
	resp, err := kc.DB.Collection(model.CatalogColl).UpdateOne(ctx, filter, deleteQuery)
	if err != nil {
		return errors.Wrap(err, "failed to delete variant in catalog")
	}

	if resp.ModifiedCount == 0 {
		return errors.Errorf("Failed to delete Variant with id %s", opts.VariantID.Hex())
	}
	return nil
}

// UpdateCatalogStatus updates status of the Catalog
func (kc *KeeperCatalogImpl) UpdateCatalogStatus(opts *schema.UpdateCatalogStatusOpts) ([]schema.UpdateCatalogStatusResp, error) {

	var catalog model.Catalog
	updateStatusValue := strings.ToLower(opts.Status)

	ctx := context.TODO()
	filter := bson.M{
		"_id": opts.CatalogID,
	}
	err := kc.DB.Collection(model.CatalogColl).FindOne(ctx, filter).Decode(&catalog)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Find Catalog")
	}
	currentStatusValue := catalog.Status.Value

	//Checking if status change is allowed
	if currentStatusValue == model.Draft && updateStatusValue == model.Unlist {
		return nil, errors.Errorf("Status change not allowed from %s to %s", currentStatusValue, updateStatusValue)
	}
	if currentStatusValue == model.Publish && updateStatusValue == model.Draft {
		return nil, errors.Errorf("Status change not allowed from %s to %s", currentStatusValue, updateStatusValue)
	}
	if currentStatusValue == model.Unlist && updateStatusValue == model.Draft {
		return nil, errors.Errorf("Status change not allowed from %s to %s", currentStatusValue, updateStatusValue)
	}
	if currentStatusValue == model.Archive {
		return nil, errors.Errorf("Status change not allowed from %s to %s", currentStatusValue, updateStatusValue)
	}

	//Draft to Publish
	//Checking All Data is Available to Publish Catalog
	var resp []schema.UpdateCatalogStatusResp
	isRequiredString := " is a required field"
	if currentStatusValue == model.Draft && updateStatusValue == model.Publish {
		if catalog.Name == "" {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Name" + isRequiredString,
				Field:   "Name",
			})
		}
		if catalog.Description == "" {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Description" + isRequiredString,
				Field:   "Description",
			})
		}
		if catalog.Paths == nil {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Category" + isRequiredString,
				Field:   "Category",
			})
		}
		if len(catalog.Keywords) == 0 {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Keywords" + isRequiredString,
				Field:   "Keywords",
			})
		}
		if catalog.FeaturedImage == nil {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Featured Image" + isRequiredString,
				Field:   "Featured Image",
			})
		}
		// if catalog.FilterAttribute == nil {
		// 	resp = append(resp, schema.UpdateCatalogStatusResp{
		// 		Type:    "Field Missing",
		// 		Message: "Filter Attribute" + isRequiredString,
		// 		Field:   "Filter Attribute",
		// 	})
		// }
		if len(catalog.Variants) == 0 {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "At least one Variant is Required",
				Field:   "Variants",
			})
		}
		if catalog.VariantType == "" {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Variant Type" + isRequiredString,
				Field:   "Variant Type",
			})
		}
		if catalog.ETA == nil {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "ETA" + isRequiredString,
				Field:   "ETA",
			})
		}
		if catalog.HSNCode == "" {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "HSN Code" + isRequiredString,
				Field:   "HSN Code",
			})
		}
		if catalog.BasePrice == nil {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Base Price" + isRequiredString,
				Field:   "Base Price",
			})
		}
		if catalog.RetailPrice == nil {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "Retail Price" + isRequiredString,
				Field:   "Retail Price",
			})
		}
		if len(catalog.CatalogContent) == 0 {
			resp = append(resp, schema.UpdateCatalogStatusResp{
				Type:    "Field Missing",
				Message: "At least one Catalog Content is required",
				Field:   "Catalog Content",
			})
		}

	}
	if len(resp) > 0 {
		return resp, errors.Errorf("Catalog Data not Complete")
	}

	updateStatus := model.Status{
		Name:      strings.Title(updateStatusValue),
		Value:     updateStatusValue,
		CreatedAt: time.Now(),
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"status": updateStatus,
		},
		"$push": bson.M{
			"status_history": updateStatus,
		},
	}
	updateResp, err := kc.DB.Collection(model.CatalogColl).UpdateOne(ctx, filter, updateQuery)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to update Status")
	}
	if updateResp.ModifiedCount == 0 {
		return nil, errors.Errorf("Unable to update Status")
	}

	return nil, nil
}

// CheckCatalogIDsExists return count based on if passed id exists in catalog collection
func (kc *KeeperCatalogImpl) CheckCatalogIDsExists(ctx context.Context, ids []primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	count, err := kc.DB.Collection(model.CatalogColl).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

//GetCatalogByIDs searches Catalogs by ID
func (kc *KeeperCatalogImpl) GetCatalogByIDs(ctx context.Context, ids []primitive.ObjectID) ([]schema.GetCatalogResp, error) {

	filterQuery := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := kc.DB.Collection(model.CatalogColl).Find(ctx, filterQuery)
	if err != nil {
		return nil, err
	}
	var catalogs []schema.GetCatalogResp
	if err = cursor.All(ctx, &catalogs); err != nil {
		return nil, err
	}
	return catalogs, nil
}

//AddCatalogContent takes catalog and content details, and returns token to keeper to upload content
func (kc *KeeperCatalogImpl) AddCatalogContent(opts *schema.AddCatalogContentOpts) (*schema.PayloadVideo, []error) {

	ctx := context.TODO()

	catalogs, err := kc.GetCatalogByIDs(ctx, []primitive.ObjectID{opts.CatalogID})
	if err != nil {
		return nil, []error{err}
	}
	if len(catalogs) == 0 {
		return nil, []error{errors.Errorf("unable to find the catalog with id: %s", opts.CatalogID.Hex())}
	}
	opts.BrandID = catalogs[0].BrandID
	requestByte, _ := json.Marshal(opts)
	url := kc.App.Config.HypdApiConfig.CmsApi + "/api/keeper/content/catalog/video"
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestByte))
	req.Header.Add("Authorization", kc.App.Config.HypdApiConfig.Token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to generate request to create catalog video")}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, []error{err}
	}
	var res schema.AddCatalogContentResp
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, []error{err}
	}
	if !res.Success {
		var errs []error
		for i := 0; i < len(res.Error); i++ {
			errs = append(errs, errors.Errorf(res.Error[i].Message))
		}
		return nil, errs
	}
	filter := bson.M{
		"_id": opts.CatalogID,
	}
	updateQuery := bson.M{
		"$push": bson.M{
			"catalog_content": res.Payload.ID,
		},
	}
	upRes, err := kc.DB.Collection(model.CatalogColl).UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return nil, []error{err}
	}
	if upRes.ModifiedCount == 0 {
		return nil, []error{errors.Errorf("error adding content to the catalog")}
	}
	return &res.Payload, nil
}

//AddCatalogContentImage takes catalog and content details, and returns token to keeper to upload content
func (kc *KeeperCatalogImpl) AddCatalogContentImage(opts *schema.AddCatalogContentImageOpts) []error {

	ctx := context.TODO()

	catalogs, err := kc.GetCatalogByIDs(ctx, []primitive.ObjectID{opts.CatalogID})
	if err != nil {
		return []error{err}
	}
	if len(catalogs) == 0 {
		return []error{errors.Errorf("unable to find the catalog with id: %s", opts.CatalogID.Hex())}
	}
	requestData := map[string]interface{}{
		"media_id":   opts.MediaID.Hex(),
		"brand_id":   catalogs[0].BrandID.Hex(),
		"catalog_id": opts.CatalogID.Hex(),
	}

	requestByte, _ := json.Marshal(requestData)
	url := kc.App.Config.HypdApiConfig.CmsApi + "/api/keeper/content/catalog/image"
	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestByte))
	if err != nil {
		return []error{errors.Wrap(err, "failed to generate request to create catalog image")}
	}
	req.Header.Add("Authorization", kc.App.Config.HypdApiConfig.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []error{err}
	}
	var res schema.AddCatalogContentImageResp
	// var test interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return []error{err}
	}
	if !res.Success {
		var errs []error
		for i := 0; i < len(res.Error); i++ {
			errs = append(errs, errors.Errorf(res.Error[i].Message))
		}
		return errs
	}

	filter := bson.M{
		"_id": opts.CatalogID,
	}
	updateQuery := bson.M{
		"$push": bson.M{
			"catalog_content": res.Payload.ID,
		},
	}
	upRes, err := kc.DB.Collection(model.CatalogColl).UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return []error{err}
	}
	if upRes.ModifiedCount == 0 {
		return []error{errors.Errorf("error adding content to the catalog")}
	}
	return nil
}

//GetCatalogsByFilter returns catalogs based on the filters entered
// func (kc *KeeperCatalogImpl) GetCatalogsByFilter(opts *schema.GetCatalogsByFilterOpts) ([]schema.GetCatalogResp, error) {

// 	var cur *mongo.Cursor
// 	var err error
// 	ctx := context.TODO()
// 	var filterQuery bson.D
// 	if len(opts.BrandIDs) > 0 {
// 		bQuery := bson.E{
// 			Key: "brand_id", Value: bson.M{
// 				"$in": opts.BrandIDs,
// 			},
// 		}
// 		filterQuery = append(filterQuery, bQuery)
// 	}
// 	if len(opts.Status) > 0 {
// 		sQuery := bson.E{
// 			Key: "status.value", Value: bson.M{
// 				"$in": opts.Status,
// 			},
// 		}
// 		filterQuery = append(filterQuery, sQuery)
// 	}
// 	if opts.Name != "" {
// 		nQuery := bson.E{
// 			Key: "lname", Value: bson.M{
// 				"$regex": strings.ToLower(opts.Name),
// 			},
// 		}
// 		filterQuery = append(filterQuery, nQuery)
// 	}
// 	// filter := bson.M{"lname": bson.M{"$regex": strings.ToLower(keeperSearchCatalogOpts.Name)}}

// 	var catalogs []schema.GetCatalogResp

// 	pageSize := kc.App.Config.PageSize
// 	skip := int64(pageSize * opts.Page)
// 	limit := int64(pageSize)
// 	findOpts := options.Find().SetSkip(skip).SetLimit(limit)
// 	fmt.Println(len(filterQuery))
// 	if len(filterQuery) == 0 {
// 		cur, err = kc.DB.Collection(model.CatalogColl).Find(ctx, bson.M{}, findOpts)
// 	} else {
// 		cur, err = kc.DB.Collection(model.CatalogColl).Find(ctx, filterQuery, findOpts)
// 	}
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error finding catalogs")
// 	}
// 	if err := cur.All(ctx, &catalogs); err != nil {
// 		return nil, err
// 	}

// 	return catalogs, nil
// }

func (kc *KeeperCatalogImpl) GetCatalogsByFilter(opts *schema.GetCatalogsByFilterOpts) ([]schema.GetCatalogResp, error) {
	var err error
	ctx := context.TODO()

	pipeline := mongo.Pipeline{}
	if len(opts.BrandIDs) > 0 {
		bMatchStage := bson.D{{
			Key: "$match", Value: bson.M{
				"brand_id": bson.M{
					"$in": opts.BrandIDs,
				},
			},
		}}
		pipeline = append(pipeline, bMatchStage)
	}

	if len(opts.Status) > 0 {
		sMatchStage := bson.D{{
			Key: "$match", Value: bson.M{
				"status.value": bson.M{
					"$in": opts.Status,
				},
			},
		}}
		pipeline = append(pipeline, sMatchStage)
	}
	if opts.Name != "" {
		nMatchStage := bson.D{{
			Key: "$match", Value: bson.M{
				"lname": bson.M{
					"$regex": strings.ToLower(opts.Name),
				},
			},
		}}
		pipeline = append(pipeline, nMatchStage)
	}
	limitStage := bson.D{
		{Key: "$limit", Value: kc.App.Config.PageSize},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: kc.App.Config.PageSize * opts.Page},
	}

	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path":                       "$variants",
			"preserveNullAndEmptyArrays": true,
		},
	}}
	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "inventory",
			"localField":   "variants.inventory_id",
			"foreignField": "_id",
			"as":           "inventory_info",
		},
	}}
	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"variants.inventory_info": bson.M{
				"$first": "$inventory_info",
			},
		},
	}}
	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$_id",
			"catalogs": bson.M{
				"$push": "$$ROOT",
			},
			"variants": bson.M{
				"$push": "$variants",
			},
		},
	}}

	addFieldsStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"catalog": bson.M{
				"$arrayElemAt": bson.A{
					"$catalogs",
					0,
				},
			},
		},
	}}

	setStage2 := bson.D{{
		Key: "$set", Value: bson.M{
			"catalog.variants": "$variants",
		},
	}}

	replaceRootStage := bson.D{{
		Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$catalog",
		},
	}}
	sortStage := bson.D{{
		Key: "$sort", Value: bson.M{
			"updated_at": -1,
		},
	}}

	pipeline = append(pipeline, mongo.Pipeline{
		skipStage,
		limitStage,
		unwindStage,
		lookupStage,
		setStage,
		groupStage,
		addFieldsStage,
		setStage2,
		replaceRootStage, sortStage}...)
	cur, err := kc.DB.Collection(model.CatalogColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var catalogResp []schema.GetCatalogResp
	if err := cur.All(ctx, &catalogResp); err != nil {
		return nil, err
	}

	return catalogResp, nil
}

//GetCatalogBySlug finds and return the catalog with given slug
func (kc *KeeperCatalogImpl) GetCatalogBySlug(slug string) (*schema.GetCatalogResp, error) {
	var catalog *schema.GetCatalogResp
	err := kc.DB.Collection(model.CatalogColl).FindOne(context.TODO(), bson.M{"slug": slug}).Decode(&catalog)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Errorf("unable to find the catalog with slug: %s", slug)
		}
		return nil, err
	}
	return catalog, nil
}

func (kc *KeeperCatalogImpl) GetAllCatalogInfo(id primitive.ObjectID) (*schema.GetAllCatalogInfoResp, error) {
	var wg sync.WaitGroup
	var contentInfo []schema.CatalogContentInfoResp
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": id,
		},
	}}
	lookupGroupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         model.GroupColl,
			"localField":   "_id",
			"foreignField": "catalog_ids",
			"as":           "group_info",
		},
	}}
	lookupDiscountStage := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": "discount",
			"let":  bson.M{"catalog_id": "$_id"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$catalog_id", "$$catalog_id"}},
								bson.M{"$eq": bson.A{"$is_active", true}},
							},
						},
					},
				},
			},
			"as": "discount_info",
		},
	}}
	setStage0 := bson.D{{
		Key: "$set", Value: bson.M{
			"discount_info": bson.M{
				"$first": "$discount_info",
			},
		},
	}}
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path":                       "$variants",
			"preserveNullAndEmptyArrays": true,
		},
	}}
	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "inventory",
			"localField":   "variants.inventory_id",
			"foreignField": "_id",
			"as":           "inventory_info",
		},
	}}
	setStage1 := bson.D{{
		Key: "$set", Value: bson.M{
			"variants.inventory_info": bson.M{
				"$first": "$inventory_info",
			},
		},
	}}
	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$_id",
			"catalogs": bson.M{
				"$push": "$$ROOT",
			},
			"variants": bson.M{
				"$push": "$variants",
			},
		},
	}}

	addFieldsStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"catalog": bson.M{
				"$arrayElemAt": bson.A{
					"$catalogs",
					0,
				},
			},
		},
	}}

	setStage2 := bson.D{{
		Key: "$set", Value: bson.M{
			"catalog.variants": "$variants",
		},
	}}

	replaceRootStage := bson.D{{
		Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$catalog",
		},
	}}

	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := kc.GetCatalogContent(id)
		if err != nil {
			kc.App.Logger.Err(err).Msgf("failed to get catalog content for id: %s", id.Hex())
			return
		}
		contentInfo = info
	}()

	catalogsCursor, err := kc.DB.Collection(model.CatalogColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupDiscountStage,
		setStage0,
		lookupGroupStage,
		unwindStage,
		lookupStage,
		setStage1,
		groupStage,
		addFieldsStage,
		setStage2,
		replaceRootStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query for catalog with id:%s", id.Hex())
	}

	var catalog []schema.GetAllCatalogInfoResp
	if err := catalogsCursor.All(ctx, &catalog); err != nil {
		return nil, errors.Wrap(err, "error decoding Catalogs")
	}

	wg.Wait()

	if len(catalog) != 0 {
		var brandInfo *schema.BrandInfoResp
		brandInfo, err = kc.App.Brand.GetBrandInfo([]string{catalog[0].BrandID.Hex()})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get brand-info")
		}

		catalog[0].ContentInfo = contentInfo
		catalog[0].BrandInfo = brandInfo
		return &catalog[0], nil
	}

	return nil, errors.Errorf("unable to find info for catalog with id: %s", id.Hex())
}

func (kc *KeeperCatalogImpl) GetKeeperCatalogContent(id primitive.ObjectID) ([]schema.CatalogContentInfoResp, error) {
	url := kc.App.Config.HypdApiConfig.CmsApi + "/api/keeper/content"
	data, err := json.Marshal(map[string]interface{}{
		"type":        "catalog_content",
		"catalog_ids": []string{id.Hex()},
		"page":        999,
	})
	if err != nil {
		kc.Logger.Err(err).Msg("failed to prepare request to get catalog content")
		return nil, errors.Wrap(err, "failed to prepare request to get catalog content")
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to catalog content")
	}
	req.Header.Add("Authorization", kc.App.Config.HypdApiConfig.Token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		kc.Logger.Err(err).Str("responseBody", string(data)).Msgf("failed to send request to api %s", url)
		return nil, errors.Wrapf(err, "failed to send request to api %s", url)
	}

	defer resp.Body.Close()
	//Read the response body

	var s schema.GetCatalogContentInfoResp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kc.Logger.Err(err).Str("responseBody", string(data)).Msgf("failed to read response from api %s", url)
	}
	if err := json.Unmarshal(body, &s); err != nil {
		kc.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		kc.Logger.Err(errors.New("success false from cms")).Str("body", string(body)).Msg("got success false response from cms")
		return nil, errors.New("got success false response from cms")
	}
	return s.Payload, nil
}

func (kc *KeeperCatalogImpl) GetCatalogContent(id primitive.ObjectID) ([]schema.CatalogContentInfoResp, error) {
	url := kc.App.Config.HypdApiConfig.CmsApi + "/api/keeper/content"
	data, err := json.Marshal(map[string]interface{}{
		"is_active":   true,
		"type":        "catalog_content",
		"catalog_ids": []string{id.Hex()},
		"page":        999,
	})
	if err != nil {
		kc.Logger.Err(err).Msg("failed to prepare request to get catalog content")
		return nil, errors.Wrap(err, "failed to prepare request to get catalog content")
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to get catalog content")
	}
	req.Header.Add("Authorization", kc.App.Config.HypdApiConfig.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		kc.Logger.Err(err).Str("responseBody", string(data)).Msgf("failed to send request to api %s", url)
		return nil, errors.Wrapf(err, "failed to send request to api %s", url)
	}

	defer resp.Body.Close()
	//Read the response body

	var s schema.GetCatalogContentInfoResp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kc.Logger.Err(err).Str("responseBody", string(data)).Msgf("failed to read response from api %s", url)
	}
	if err := json.Unmarshal(body, &s); err != nil {
		kc.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		kc.Logger.Err(errors.New("success false from cms")).Str("body", string(body)).Msg("got success false response from cms")
		return nil, errors.New("got success false response from cms")
	}
	return s.Payload, nil
}

func (kc *KeeperCatalogImpl) SyncCatalog(id primitive.ObjectID) {
	filter := bson.M{
		"_id":          id,
		"status.value": model.Publish,
	}
	update := bson.M{
		"$set": bson.M{
			"last_sync": time.Now().UTC(),
		},
	}
	if _, err := kc.DB.Collection(model.CatalogColl).UpdateMany(context.TODO(), filter, update); err != nil {
		kc.Logger.Err(err).Interface("opts", id).Msg("failed to sync catalog")
	}
}

func (kc *KeeperCatalogImpl) SyncCatalogs(ids []primitive.ObjectID) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"status.value": model.Publish,
	}
	update := bson.M{
		"$set": bson.M{
			"last_sync": time.Now().UTC(),
		},
	}
	if _, err := kc.DB.Collection(model.CatalogColl).UpdateMany(context.TODO(), filter, update); err != nil {
		kc.Logger.Err(err).Interface("opts", ids).Msg("failed to sync catalogs")
	}
}

func (kc *KeeperCatalogImpl) SyncCatalogContent(id primitive.ObjectID) {
	filter := bson.M{
		"catalog_content": id,
		"status.value":    model.Publish,
	}
	update := bson.M{
		"$set": bson.M{
			"last_sync": time.Now().UTC(),
		},
	}
	if _, err := kc.DB.Collection(model.CatalogColl).UpdateMany(context.TODO(), filter, update); err != nil {
		kc.Logger.Err(err).Interface("opts", id).Msg("failed to sync catalog content")
	}
}

func (kc *KeeperCatalogImpl) GetPebbleCatalogInfo(ids []primitive.ObjectID) ([]schema.GetAllCatalogInfoResp, error) {
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": bson.M{
				"$in": ids,
			},
			"status.value": model.Publish,
		},
	}}
	lookupDiscountStage := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": "discount",
			"let":  bson.M{"catalog_id": "$_id"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$catalog_id", "$$catalog_id"}},
								bson.M{"$eq": bson.A{"$is_active", true}},
							},
						},
					},
				},
			},
			"as": "discount_info",
		},
	}}
	setStage0 := bson.D{{
		Key: "$set", Value: bson.M{
			"discount_info": bson.M{
				"$first": "$discount_info",
			},
		},
	}}

	catalogsCursor, err := kc.DB.Collection(model.CatalogColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupDiscountStage,
		setStage0,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query for catalog with id")
	}

	var catalogs []schema.GetAllCatalogInfoResp
	if err := catalogsCursor.All(ctx, &catalogs); err != nil {
		return nil, errors.Wrap(err, "error decoding Catalogs")
	}

	// if len(catalogs) == 0 {
	// 	return nil, errors.Errorf("unable to find info for catalog for collection")
	// }
	for i, catalog := range catalogs {
		bi, err := kc.App.Brand.GetBrandInfo([]string{catalog.BrandID.Hex()})
		if err != nil {
			kc.Logger.Err(err).Msgf("failed to get brand info for catalog with brand-id: %s", catalog.BrandID.Hex())
			continue
		}
		catalogs[i].BrandInfo = bi
	}
	return catalogs, nil
}

func (kc *KeeperCatalogImpl) GetCollectionCatalogInfo(ids []primitive.ObjectID) ([]schema.GetAllCatalogInfoResp, error) {
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		},
	}}
	lookupDiscountStage := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": "discount",
			"let":  bson.M{"catalog_id": "$_id"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$catalog_id", "$$catalog_id"}},
								bson.M{"$eq": bson.A{"$is_active", true}},
							},
						},
					},
				},
			},
			"as": "discount_info",
		},
	}}
	setStage0 := bson.D{{
		Key: "$set", Value: bson.M{
			"discount_info": bson.M{
				"$first": "$discount_info",
			},
		},
	}}

	catalogsCursor, err := kc.DB.Collection(model.CatalogColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupDiscountStage,
		setStage0,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query for catalog with id")
	}

	var catalogs []schema.GetAllCatalogInfoResp
	if err := catalogsCursor.All(ctx, &catalogs); err != nil {
		return nil, errors.Wrap(err, "error decoding Catalogs")
	}
	if len(catalogs) == 0 {
		return nil, errors.Errorf("unable to find info for catalog for collection")
	}
	for i, catalog := range catalogs {
		bi, err := kc.App.Brand.GetBrandInfo([]string{catalog.BrandID.Hex()})
		if err != nil {
			kc.Logger.Err(err).Msgf("failed to get brand info for catalog with brand-id: %s", catalog.BrandID.Hex())
			continue
		}
		catalogs[i].BrandInfo = bi
	}
	return catalogs, nil
}

func (kc *KeeperCatalogImpl) GetCatalogVariant(cat_id, var_id primitive.ObjectID) (*schema.GetCatalogVariantResp, error) {

	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id":          cat_id,
			"variants._id": var_id,
		},
	}}
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path": "$variants",
		},
	}}
	matchStage2 := bson.D{{
		Key: "$match", Value: bson.M{
			"variants._id": var_id,
		},
	}}
	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from": model.DiscountColl,
			"let": bson.M{
				"variant_id": "$variants._id",
			},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr":     bson.M{"$in": bson.A{"$$variant_id", "$variants_id"}},
						"is_active": true,
					}},
			},
			"as": "discount_info",
		},
	}}
	unwindStage2 := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path":                       "$discount_info",
			"preserveNullAndEmptyArrays": true,
		},
	}}
	inventoryLookUpStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "inventory",
			"localField":   "variants.inventory_id",
			"foreignField": "_id",
			"as":           "inventory_info",
		},
	}}
	projectStage :=
		bson.D{{
			Key: "$project", Value: bson.M{
				"_id":                     1,
				"name":                    1,
				"base_price":              1,
				"retail_price":            1,
				"transfer_price":          1,
				"discount_info._id":       1,
				"discount_info.value":     1,
				"discount_info.type":      1,
				"discount_info.max_value": 1,
				"variant_type":            1,
				"variant":                 "$variants",
				"featured_image":          1,
				"inventory_info":          bson.M{"$arrayElemAt": bson.A{"$inventory_info", 0}},
			},
		}}

	ctx := context.TODO()

	catalogsCursor, err := kc.DB.Collection(model.CatalogColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		unwindStage,
		matchStage2,
		lookupStage,
		unwindStage2,
		inventoryLookUpStage,
		projectStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query for catalog with id:%s", cat_id.Hex())
	}
	var catalog []schema.GetCatalogVariantResp
	if err := catalogsCursor.All(ctx, &catalog); err != nil {
		return nil, errors.Wrap(err, "error decoding Catalogs")
	}
	if len(catalog) > 0 {
		return &catalog[0], nil
	}
	return nil, nil
}

func (kc *KeeperCatalogImpl) RemoveContent(opts *schema.RemoveContentOpts) error {
	filter := bson.M{"_id": opts.CatalogID}
	updateQuery := bson.M{
		"$pull": bson.M{
			"catalog_content": opts.ContentID,
		},
	}
	res, err := kc.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "error removing content with id: %s", opts.ContentID.Hex())
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("error finding catalog with id: %s", opts.CatalogID.Hex())
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/keeper/content/%s", kc.App.Config.HypdApiConfig.CmsApi, opts.ContentID.Hex()), nil)
	if err != nil {
		return errors.Wrap(err, "error sending deleting content request to cms")
	}
	req.Header.Add("Authorization", kc.App.Config.HypdApiConfig.Token)
	if _, err := client.Do(req); err != nil {
		return errors.Wrap(err, "error deleting content from cms")
	}
	return nil
}
