package app

import (
	"go-app/model"
	"go-app/schema"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// KeeperCatalog service allows `Keeper` to execute admin operations.
type KeeperCatalog interface {
	CreateCatalog(*schema.CreateCatalogOpts)
	EditCatalog(*schema.CreateCatalogOpts)
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
func (kc *KeeperCatalogImpl) CreateCatalog(opts *schema.CreateCatalogOpts) {
	c := model.Catalog{
		Name:        opts.Name,
		Description: opts.Description,
		HSNCode:     opts.HSNCode,
	}
	// If variants are passed in the opts then setting variants in catalog model
	if opts.VariantType != "" {
		c.VariantType = opts.VariantType
		for _, variant := range opts.Variants {
			c.Variants = append(c.Variants, kc.createVariant(&variant))
		}
	}
}

// EditCatalog edits an existing catalog
func (kc *KeeperCatalogImpl) EditCatalog(opts *schema.CreateCatalogOpts) {}

/* handler funcs here */

func (kc *KeeperCatalogImpl) createVariant(opts *schema.CreateVariantOpts) model.Variant {
	return model.Variant{
		ID:          primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		SKU:         opts.SKU,
		BasePrice:   model.SetINRPrice(float32(opts.BasePrice)),
		RetailPrice: model.SetINRPrice(float32(opts.RetailPrice)),
	}
}

func (kc *KeeperCatalogImpl) uploadCatalogFeaturedImage() {}

func (kc *KeeperCatalogImpl) prepareCatalogFeaturedImage() {}
