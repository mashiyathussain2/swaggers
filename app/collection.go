package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection service allows `Collection` to execute admin operations.
type Collection interface {
	CreateCollection(*schema.CreateCollectionOpts) (*schema.CreateCollectionResp, []error)
	DeleteCollection(primitive.ObjectID) error
	AddSubCollection(*schema.AddSubCollectionOpts) (*schema.CreateCollectionResp, []error)
	DeleteSubCollection(primitive.ObjectID, primitive.ObjectID) error
	EditCollection(*schema.EditCollectionOpts) (*schema.CreateCollectionResp, error)
	UpdateSubCollectionImage(opts *schema.UpdateSubCollectionImageOpts) error
	AddCatalogsToSubCollection(*schema.UpdateCatalogsInSubCollectionOpts) []error
	RemoveCatalogsFromSubCollection(*schema.UpdateCatalogsInSubCollectionOpts) []error
	GetCollections(int) ([]schema.CreateCollectionResp, error)

	AddCatalogInfoToCollection(id primitive.ObjectID)
	UpdateCollectionCatalogInfo(id primitive.ObjectID)

	// GetActiveCollections()
}

// CollectionImpl implements collection related operations
type CollectionImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// CollectionOpts contains arguments required to create a new instance of Collection
type CollectionOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

//InitCollection returns Collection instance
func InitCollection(opts *CollectionOpts) Collection {
	return &CollectionImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

//CreateCollection inserts a new catalog document with specified data into the database
func (ci *CollectionImpl) CreateCollection(opts *schema.CreateCollectionOpts) (*schema.CreateCollectionResp, []error) {
	t := time.Now().UTC()
	ctx := context.TODO()
	var subCollections []model.SubCollection

	for _, sc := range opts.SubCollection {
		err := ci.checkCatalogs(sc.CatalogIDs)
		if err != nil && len(err) > 0 {
			return nil, err
		}
		subCollection := model.SubCollection{
			ID:         primitive.NewObjectIDFromTimestamp(time.Now()),
			Name:       sc.Name,
			CatalogIDs: sc.CatalogIDs,
			CreatedAt:  t,
		}
		if sc.Image != nil {
			image := model.IMG{
				SRC: sc.Image.SRC,
			}
			if err := image.LoadFromURL(); err != nil {
				return nil, []error{errors.Wrapf(err, "unable to process image for sub collection %s", sc.Name)}
			}
			subCollection.Image = &image
		}
		subCollections = append(subCollections, subCollection)

	}
	collection := model.Collection{
		Name:           UniqueSlug(opts.Title),
		Type:           opts.Type,
		Genders:        opts.Genders,
		Title:          opts.Title,
		SubCollections: subCollections,
		CreatedAt:      t,
		Status:         model.Draft,
		Order:          -1,
	}
	res, err := ci.DB.Collection(model.CollectionColl).InsertOne(ctx, collection)

	if err != nil {
		return nil, []error{err}
	}

	collectionResp := schema.CreateCollectionResp{
		ID:             res.InsertedID.(primitive.ObjectID),
		Name:           collection.Name,
		Type:           collection.Type,
		Genders:        collection.Genders,
		Title:          collection.Title,
		SubCollections: collection.SubCollections,
		Order:          -1,
	}

	return &collectionResp, nil
}

//DeleteCollection deletes the collection from the database with given collectionID
func (ci *CollectionImpl) DeleteCollection(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	deleteQuery := bson.M{
		"$set": bson.M{
			"status": model.Disable,
		},
	}
	res, err := ci.DB.Collection(model.CollectionColl).UpdateOne(context.TODO(), filter, deleteQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to query for collection with id: %s", id.Hex())
	}
	if res.ModifiedCount == 0 {
		return errors.Errorf("unable to delete collection with id: %s", id.Hex())
	}
	return nil
}

//AddSubCollection adds a sub collection to the collection with given id
func (ci *CollectionImpl) AddSubCollection(opts *schema.AddSubCollectionOpts) (*schema.CreateCollectionResp, []error) {

	err := ci.checkCatalogs(opts.SubCollection.CatalogIDs)
	if err != nil && len(err) > 0 {
		return nil, err
	}

	subCollection := model.SubCollection{
		ID:         primitive.NewObjectIDFromTimestamp(time.Now()),
		Name:       opts.SubCollection.Name,
		CatalogIDs: opts.SubCollection.CatalogIDs,
		CreatedAt:  time.Now(),
	}

	if opts.SubCollection.Image != nil {
		image := model.IMG{
			SRC: opts.SubCollection.Image.SRC,
		}
		if err := image.LoadFromURL(); err != nil {
			return nil, []error{errors.Wrapf(err, "unable to process image for sub collection %s", opts.SubCollection.Name)}
		}
		subCollection.Image = &image
	}

	var collectionModel model.Collection
	findQuery := bson.M{"_id": opts.ID}
	updateQuery := bson.M{
		"$push": bson.M{
			"sub_collections": subCollection,
		},
	}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	errResp := ci.DB.Collection(model.CollectionColl).FindOneAndUpdate(context.TODO(), findQuery, updateQuery, qOpts).Decode(&collectionModel)
	if errResp != nil {
		if errResp == mongo.ErrNoDocuments || errResp == mongo.ErrNilDocument {
			return nil, []error{errors.Errorf("collection with id:%s not found", opts.ID.Hex())}
		}
		return nil, []error{errors.Wrap(errResp, "failed to update catalog")}
	}

	collection := schema.CreateCollectionResp{
		ID:             collectionModel.ID,
		Type:           collectionModel.Type,
		Name:           collectionModel.Name,
		Genders:        collectionModel.Genders,
		Title:          collectionModel.Title,
		SubCollections: collectionModel.SubCollections,
	}
	return &collection, nil
}

//DeleteSubCollection deletes the sub collection from the given collection
func (ci *CollectionImpl) DeleteSubCollection(collID primitive.ObjectID, subID primitive.ObjectID) error {
	filter := bson.M{"_id": collID}
	query := bson.M{
		"$pull": bson.M{
			"sub_collections": bson.M{"_id": subID},
		},
	}
	res, err := ci.DB.Collection(model.CollectionColl).UpdateOne(context.TODO(), filter, query)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find collection with id - %s ", collID)
	}
	if res.ModifiedCount == 0 {
		return errors.Errorf("unable to delete sub collection")
	}
	return nil
}

//EditCollection edits the collection details such as title, name, genders
func (ci *CollectionImpl) EditCollection(opts *schema.EditCollectionOpts) (*schema.CreateCollectionResp, error) {
	collection := model.Collection{}

	if opts.Title != "" {
		collection.Title = opts.Title
	}
	if opts.Genders != nil {
		collection.Genders = opts.Genders
	}
	if opts.Order != 0 {
		collection.Order = opts.Order
	}

	if reflect.DeepEqual(model.Collection{}, collection) {
		return nil, errors.New("no fields found to update")
	}
	collection.UpdatedAt = time.Now().UTC()
	filter := bson.M{
		"_id": opts.ID,
	}
	update := bson.M{"$set": collection}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := ci.DB.Collection(model.CollectionColl).FindOneAndUpdate(context.TODO(), filter, update, qOpts).Decode(&collection)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("catalog with id:%s not found", opts.ID.Hex())
		}
		return nil, errors.Wrap(err, "failed to update catalog")
	}

	collectionResp := &schema.CreateCollectionResp{
		ID:             collection.ID,
		Title:          collection.Title,
		Type:           collection.Type,
		Name:           collection.Name,
		Genders:        collection.Genders,
		SubCollections: collection.SubCollections,
		Order:          collection.Order,
	}

	return collectionResp, nil
}

//UpdateSubCollectionImage updates the sub collection image
func (ci *CollectionImpl) UpdateSubCollectionImage(opts *schema.UpdateSubCollectionImageOpts) error {

	findQuery := bson.M{"_id": opts.ColID, "sub_collection._id": opts.SubID}
	img := model.IMG{
		SRC: opts.Image,
	}
	err := img.LoadFromURL()
	if err != nil {
		return errors.Wrapf(err, "unable to load image")
	}
	updateQuery := bson.M{"$set": bson.M{"sub_collection.$.image": img}}
	res, err := ci.DB.Collection(model.CollectionColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find the sub collection with id - %s", opts.SubID.Hex())
	}
	if res.ModifiedCount == 0 {
		return errors.Errorf("unable to the update the sub collection with id - %s", opts.SubID.Hex())
	}
	return nil
}

//AddCatalogsToSubCollection adds catalogs to the sub collectionUpdateCatalogsToSubCollection
func (ci *CollectionImpl) AddCatalogsToSubCollection(opts *schema.UpdateCatalogsInSubCollectionOpts) []error {

	findQuery := bson.M{"_id": opts.ColID, "sub_collection._id": opts.SubID}

	err := ci.checkCatalogs(opts.CatalogIDs)
	if err != nil && len(err) > 0 {
		return err
	}

	updateQuery := bson.M{"$addToSet": bson.M{"sub_collection.$.catalog_ids": bson.M{
		"$each": opts.CatalogIDs,
	}}}

	res, errResp := ci.DB.Collection(model.CollectionColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if errResp != nil {
		return []error{errResp}
	}
	if res.MatchedCount == 0 {
		return []error{errors.Errorf("unable to find the sub collection with id - %s", opts.SubID.Hex())}
	}
	if res.ModifiedCount == 0 {
		return []error{errors.Errorf("unable to the update the sub collection with id - %s", opts.SubID.Hex())}
	}

	return nil
}

//RemoveCatalogsFromSubCollection adds catalogs to the sub collection
func (ci *CollectionImpl) RemoveCatalogsFromSubCollection(opts *schema.UpdateCatalogsInSubCollectionOpts) []error {

	findQuery := bson.M{"_id": opts.ColID, "sub_collection._id": opts.SubID}

	err := ci.checkCatalogs(opts.CatalogIDs)
	if err != nil && len(err) > 0 {
		return err
	}

	updateQuery := bson.M{"$pull": bson.M{"sub_collection.$.catalog_ids": bson.M{
		"$in": opts.CatalogIDs,
	}}}

	res, errResp := ci.DB.Collection(model.CollectionColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if errResp != nil {
		return []error{errResp}
	}
	if res.MatchedCount == 0 {
		return []error{errors.Errorf("unable to find the sub collection with id - %s", opts.SubID.Hex())}
	}
	if res.ModifiedCount == 0 {
		return []error{errors.Errorf("unable to the update the sub collection with id - %s", opts.SubID.Hex())}
	}

	return nil
}

func (ci *CollectionImpl) checkCatalogs(opts []primitive.ObjectID) []error {
	var errorRes []error

	catalogs, err := ci.App.KeeperCatalog.GetCatalogByIDs(context.TODO(), opts)
	if err != nil {
		return []error{errors.Wrap(err, "Unable to query for Catalogs")}
	}
	catalogMap := make(map[primitive.ObjectID]schema.GetCatalogResp)
	for i := 0; i < len(catalogs); i++ {
		catalogMap[catalogs[i].ID] = catalogs[i]
	}

	if len(catalogs) != len(opts) {
		for i := 0; i < len(opts); i++ {
			_, ok := catalogMap[opts[i]]
			if !ok {
				errorRes = append(errorRes, errors.Errorf("catalog with id: %s not found", opts[i].Hex()))
			}
		}
	}
	return errorRes
}

func (ci *CollectionImpl) GetCollections(page int) ([]schema.CreateCollectionResp, error) {

	ctx := context.TODO()
	opts := options.Find().SetSkip(int64(ci.App.Config.PageSize * page)).SetLimit(int64(ci.App.Config.PageSize)).SetSort(bson.D{{Key: "order", Value: 1}})
	cur, err := ci.DB.Collection(model.CollectionColl).Find(ctx, bson.M{}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("no collections found")
		}
		return nil, errors.Wrapf(err, "error querying the database")
	}
	var collectionResp []schema.CreateCollectionResp
	if err := cur.All(ctx, &collectionResp); err != nil {
		return nil, err
	}

	return collectionResp, nil
}

func (ci *CollectionImpl) AddCatalogInfoToCollection(id primitive.ObjectID) {
	var collection model.Collection
	ctx := context.TODO()
	filter := bson.M{
		"_id": id,
	}

	if err := ci.DB.Collection(model.CollectionColl).FindOne(ctx, filter).Decode(&collection); err != nil {
		ci.Logger.Err(err).Msgf("failed to collection with id: %s", id.Hex())
		return
	}

	if collection.Type != model.ProductCollection {
		ci.Logger.Info().Msg("sub collection is not product collection, thus skipping catalog_info linking")
		return
	}

	var operations []mongo.WriteModel
	for _, subColl := range collection.SubCollections {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{"_id": id, "sub_collections._id": subColl.ID})
		catalogInfo, err := ci.App.KeeperCatalog.GetCollectionCatalogInfo(subColl.CatalogIDs)
		if err != nil {
			ci.Logger.Err(err).Msgf("failed to find catalog for subcollection with id: %s", subColl.ID.Hex())
			return
		}
		operation.SetUpdate(bson.M{
			"$set": bson.M{
				"sub_collections.$.catalog_info": catalogInfo,
			},
		})
		operations = append(operations, operation)
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	_, err := ci.DB.Collection(model.CollectionColl).BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to add catalog info inside collection with id:%s", id.Hex())
	}
}

func (ci *CollectionImpl) UpdateCollectionCatalogInfo(id primitive.ObjectID) {
	filter := bson.M{
		"sub_collections.catalog_ids": id,
	}

	catalogInfo, err := ci.App.KeeperCatalog.GetCollectionCatalogInfo([]primitive.ObjectID{id})
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to find catalog with id: %s", id.Hex())
		return
	}

	update := bson.M{
		"$set": bson.M{
			"sub_collections.catalog_info.$": catalogInfo[0],
		},
	}

	if _, err = ci.DB.Collection(model.CollectionColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Msgf("failed to update catalog info for catalog id:%s", id.Hex())
		return
	}
}
