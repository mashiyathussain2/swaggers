//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_catalog.go -package=mock go-app/app KeeperCatalog

package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// KeeperCatalog service allows `Keeper` to execute admin operations.
type KeeperCatalog interface {
	CreateCatalog(*schema.CreateCatalogOpts) (*schema.CreateCatalogResp, error)
	EditCatalog(*schema.EditCatalogOpts) (*schema.EditCatalogResp, error)
	// AddVariant(*schema.CreateVariantOpts)
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
	fmt.Printf("%+v\n", opts)
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
			return nil, errors.Wrapf(err, "failed to find brand with id: %s", opts.BrandID)
		}
		return nil, errors.Errorf("brand id %s does not exists", opts.BrandID.Hex())
	}
	c.BrandID = opts.BrandID

	// setting catalog specifications
	for _, specOpt := range opts.Specifications {
		c.Specifications = append(c.Specifications, model.Specification{Name: specOpt.Name, Value: specOpt.Value})
	}

	// If eta is passed then setting up the eta
	if opts.ETA != nil {
		c.ETA = &model.ETA{
			Min:  int(opts.ETA.Min),
			Max:  int(opts.ETA.Max),
			Unit: opts.ETA.Unit,
		}
	}

	// TODO: add logic to set catalog category path

	// Inserting the document in the DB
	res, err := kc.DB.Collection(model.CatalogColl).InsertOne(context.Background(), c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert catalog in db")
	}

	resp := &schema.CreateCatalogResp{
		ID:              res.InsertedID.(primitive.ObjectID),
		Name:            c.Name,
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
	// Add edit catalog logic here

	return nil, nil
}

func (kc *KeeperCatalogImpl) createVariant(opts *schema.CreateVariantOpts) *model.Variant {
	return &model.Variant{
		ID:  primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		SKU: opts.SKU,
	}
}
