package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InfluencerCollection service allows `InfluencerCollection` to execute admin operations.
type InfluencerCollection interface {
	CreateInfluencerCollection(opts *schema.CreateInfluencerCollectionOpts) (*schema.InfluencerCollectionResp, error)
	KeeperGetInfluencerCollections(opts *schema.GetInfluencerCollectionsOpts) ([]schema.InfluencerCollectionResp, error)
	EditInfluencerCollection(opts *schema.EditInfluencerCollectionOpts) (*schema.InfluencerCollectionResp, error)
	EditInfluencerCollectionApp(opts *schema.EditInfluencerCollectionAppOpts) (*schema.InfluencerCollectionResp, error)
	// GetActiveInfluencerCollections()
}

// InfluencerCollectionImpl implements Influencercollection related operations
type InfluencerCollectionImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InfluencerCollectionOpts contains arguments required to create a new instance of InfluencerCollection
type InfluencerCollectionOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

//InitInfluencerCollection returns InfluencerCollection instance
func InitInfluencerCollection(opts *InfluencerCollectionOpts) InfluencerCollection {
	return &InfluencerCollectionImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

//CreateInfluencerCollection inserts a new catalog document with specified data into the database
func (ii *InfluencerCollectionImpl) CreateInfluencerCollection(opts *schema.CreateInfluencerCollectionOpts) (*schema.InfluencerCollectionResp, error) {
	t := time.Now().UTC()
	ctx := context.TODO()
	collection := model.InfluncerCollection{
		InfluencerID: opts.InfluencerID,
		Name:         opts.Name,
		Slug:         UniqueSlug(opts.Name),
		CatalogIDs:   opts.CatalogIDs,
		CreatedAt:    t,
		Status:       model.Draft,
		Order:        -1,
	}
	image := model.IMG{
		SRC: opts.Image.SRC,
	}
	if err := image.LoadFromURL(); err != nil {
		return nil, errors.Wrapf(err, "unable to process image for collection %s", opts.Name)
	}
	collection.Image = &image
	res, err := ii.DB.Collection(model.InfluencerCollectionColl).InsertOne(ctx, collection)
	if err != nil {
		return nil, err
	}
	collectionResp := schema.InfluencerCollectionResp{
		ID:           res.InsertedID.(primitive.ObjectID),
		InfluencerID: collection.InfluencerID,
		Name:         collection.Name,
		Slug:         collection.Slug,
		Image:        collection.Image,
		CatalogIDs:   collection.CatalogIDs,
		Status:       collection.Status,
		Order:        collection.Order,
		CreatedAt:    collection.CreatedAt,
		UpdatedAt:    collection.UpdatedAt,
	}
	return &collectionResp, nil
}

//GetInfluencerCollections inserts a new catalog document with specified data into the database
func (ii *InfluencerCollectionImpl) KeeperGetInfluencerCollections(opts *schema.GetInfluencerCollectionsOpts) ([]schema.InfluencerCollectionResp, error) {
	ctx := context.TODO()
	var filter bson.D
	fmt.Println(1)

	if opts.InfluencerID != "" {
		iid, err := primitive.ObjectIDFromHex(opts.InfluencerID)
		if err != nil {
			return nil, errors.Wrapf(err, "influencer id is incorrect")
		}
		if iid != primitive.NilObjectID {
			filter = append(filter, bson.E{Key: "influencer_id", Value: iid})
		}
	}
	if opts.Status != "" {
		filter = append(filter, bson.E{Key: "status", Value: opts.Status})
	}
	if filter == nil {
		filter = bson.D{}
	}
	fmt.Println(1)
	queryOpts := options.Find().SetSkip(int64(ii.App.Config.PageSize * opts.Page)).SetLimit(int64(ii.App.Config.PageSize)).SetSort(bson.M{"_id": -1})
	fmt.Println(1)
	cur, err := ii.DB.Collection(model.InfluencerCollectionColl).Find(ctx, filter, queryOpts)
	if err != nil {
		return nil, err
	}
	fmt.Println(1)
	var collectionResp []schema.InfluencerCollectionResp
	if err := cur.All(ctx, &collectionResp); err != nil {
		return nil, err
	}
	return collectionResp, nil
}

func (ii *InfluencerCollectionImpl) EditInfluencerCollection(opts *schema.EditInfluencerCollectionOpts) (*schema.InfluencerCollectionResp, error) {
	ctx := context.TODO()
	var collection model.InfluncerCollection
	t := time.Now()
	if opts.Image != nil {
		image := model.IMG{
			SRC: opts.Image.SRC,
		}
		if err := image.LoadFromURL(); err != nil {
			return nil, errors.Wrapf(err, "unable to process image for collection %s", opts.Name)
		}
		collection.Image = &image
	}
	if opts.Name != "" {
		collection.Name = opts.Name
		collection.Slug = UniqueSlug(opts.Name)
	}
	if opts.CatalogIDs != nil {
		collection.CatalogIDs = opts.CatalogIDs
	}
	if opts.Order != -1 {
		collection.Order = opts.Order
	}
	if opts.Status != "" {
		collection.Status = opts.Status
	}
	collection.UpdatedAt = t
	filter := bson.M{
		"_id":    opts.ID,
		"status": bson.M{"$ne": model.Archive},
	}
	update := bson.M{
		"$set": collection,
	}
	var collectionResp schema.InfluencerCollectionResp
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := ii.DB.Collection(model.InfluencerCollectionColl).FindOneAndUpdate(ctx, filter, update, qOpts).Decode(&collectionResp)
	if err != nil {
		return nil, err
	}
	return &collectionResp, nil
}

func (ii *InfluencerCollectionImpl) EditInfluencerCollectionApp(opts *schema.EditInfluencerCollectionAppOpts) (*schema.InfluencerCollectionResp, error) {
	ctx := context.TODO()
	var collection model.InfluncerCollection
	t := time.Now()
	if opts.Image != nil {
		image := model.IMG{
			SRC: opts.Image.SRC,
		}
		if err := image.LoadFromURL(); err != nil {
			return nil, errors.Wrapf(err, "unable to process image for collection %s", opts.Name)
		}
		collection.Image = &image
	}
	if opts.Name != "" {
		collection.Name = opts.Name
		collection.Slug = UniqueSlug(opts.Name)
	}
	if opts.CatalogIDs != nil {
		collection.CatalogIDs = opts.CatalogIDs
	}
	if opts.Order != -1 {
		collection.Order = opts.Order
	}
	if opts.Status != "" {
		collection.Status = opts.Status
	}
	collection.UpdatedAt = t
	filter := bson.M{
		"_id":           opts.ID,
		"status":        bson.M{"$ne": model.Archive},
		"influencer_id": opts.InfluencerID,
	}
	update := bson.M{
		"$set": collection,
	}
	var collectionResp schema.InfluencerCollectionResp
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := ii.DB.Collection(model.InfluencerCollectionColl).FindOneAndUpdate(ctx, filter, update, qOpts).Decode(&collectionResp)
	if err != nil {
		return nil, err
	}
	return &collectionResp, nil
}

// func (ii *CollectionImpl) AddCatalogInfoToInfluencerCollection(id primitive.ObjectID) {
// 	var collection model.InfluncerCollection
// 	ctx := context.TODO()
// 	filter := bson.M{
// 		"_id": id,
// 	}
// 	if err := ii.DB.Collection(model.InfluencerCollectionColl).FindOne(ctx, filter).Decode(&collection); err != nil {
// 		ii.Logger.Err(err).Msgf("failed to collection with id: %s", id.Hex())
// 		return
// 	}
// 	catalogInfo, err := ii.App.KeeperCatalog.GetCollectionCatalogInfo(collection.CatalogIDs)
// 	if err != nil {
// 		ii.Logger.Err(err).Msgf("failed to find catalog for influencer collection with id: %s", collection.ID.Hex())
// 		return
// 	}
// 	if catalogInfo == nil {
// 		ii.Logger.Err(err).Msgf("empty catalog info for influencer collection with id: %s", collection.ID.Hex())
// 		return
// 	}
// 	_, err := ii.DB.Collection(model.InfluencerCollectionColl).InsertOne(context.TODO(), operations, &bulkOption)
// 	if err != nil {
// 		ii.Logger.Err(err).Msgf("failed to add catalog info inside collection with id:%s", id.Hex())
// 	}
// }
