package app

import (
	"context"
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

// Collection contains methods to implement and operation pebble(video-only) Collection
type Collection interface {
	CreateCollection(opts *schema.CreateCollectionOpts) error
	UpdateCollection(opts *schema.UpdateCollectionOpts) error
	KeeperGetCollections(opts *schema.GetCollectionsKeeperFilter) ([]schema.CollectionResp, error)
}

// CollectionImpl implements `Pebble` functionality
type CollectionImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// CollectionOpts contains args required to create a new instance of `PebbleImpl`
type CollectionOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitCollection returns a new instance of `Pebble` Implementation
func InitCollection(opts *CollectionOpts) Collection {
	p := CollectionImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &p
}

func (ci *CollectionImpl) CreateCollection(opts *schema.CreateCollectionOpts) error {
	coll := model.Collection{
		Name:      opts.Name,
		Type:      opts.Type,
		Genders:   opts.Genders,
		Status:    model.Draft,
		CreatedAt: time.Now(),
	}
	switch opts.Type {
	case model.HashTagCollection:
		coll.Hashtags = opts.Hashtags

	case model.BrandCollection:
		coll.BrandIDs = opts.BrandIDs

	case model.InfluencerCollection:
		coll.InfluencerIDs = opts.InfluencerIDs

	case model.SeriesCollection:
		for _, f := range opts.SeriesSubCollection {
			sc := model.SeriesSubCollection{
				ID:        primitive.NewObjectID(),
				Thumbnail: model.IMG{SRC: f.Thumbnail.SRC},
				SeriesIDs: f.SeriesIDs,
			}
			if err := sc.Thumbnail.LoadFromURL(); err != nil {
				return errors.Wrap(err, "failed to load thumbnail image")
			}
			coll.SeriesSubCollection = append(coll.SeriesSubCollection, sc)
		}
	}
	_, err := ci.DB.Collection(model.CollectionColl).InsertOne(context.TODO(), coll)
	if err != nil {
		return errors.Wrapf(err, "error creating a new collection ")
	}
	return nil
}

func (ci *CollectionImpl) UpdateCollection(opts *schema.UpdateCollectionOpts) error {

	c := model.Collection{}
	ctx := context.TODO()
	filter := bson.M{
		"_id": opts.ID,
	}

	err := ci.DB.Collection(model.CollectionColl).FindOne(ctx, filter).Decode(&c)
	c.UpdatedAt = time.Now()
	if opts.Name != "" {
		c.Name = opts.Name
	}
	if len(opts.Genders) != 0 {
		c.Genders = opts.Genders
	}
	if c.Type == model.HashTagCollection {
		if len(opts.Hashtags) != 0 {
			c.Hashtags = opts.Hashtags
		}
	} else if c.Type == model.BrandCollection {
		if len(opts.BrandIDs) != 0 {
			c.BrandIDs = opts.BrandIDs
		}
	} else if c.Type == model.InfluencerCollection {
		if len(opts.InfluencerIDs) != 0 {
			c.InfluencerIDs = opts.InfluencerIDs
		}
	} else if c.Type == model.SeriesCollection {
		if len(opts.SeriesSubCollection) > 0 {
			var s []model.SeriesSubCollection
			for _, f := range opts.SeriesSubCollection {
				sc := model.SeriesSubCollection{
					ID:        primitive.NewObjectID(),
					Thumbnail: model.IMG{SRC: f.Thumbnail.SRC},
					SeriesIDs: f.SeriesIDs,
				}
				if err := sc.Thumbnail.LoadFromURL(); err != nil {
					return errors.Wrap(err, "failed to load thumbnail image")
				}
				s = append(s, sc)
			}
			c.SeriesSubCollection = s
		}
	}

	if opts.Status != "" {
		if c.Status != model.Archive && (opts.Status == model.Publish || opts.Status == model.Unlist || opts.Status == model.Archive) {
			c.Status = opts.Status
		}
	}

	update := bson.M{"$set": c}

	err = ci.DB.Collection(model.CollectionColl).FindOneAndUpdate(context.TODO(), filter, update).Decode(&c)
	if err != nil {
		return errors.Wrapf(err, "error updating the collection ")
	}
	return nil
}

func (ci *CollectionImpl) KeeperGetCollections(opts *schema.GetCollectionsKeeperFilter) ([]schema.CollectionResp, error) {
	var queryFilter bson.M
	if len(opts.Status) > 0 {
		queryFilter = bson.M{
			"status": bson.M{
				"$in": opts.Status,
			},
		}
	}
	ctx := context.TODO()
	queryOpts := options.Find().SetSkip(int64(10 * opts.Page)).SetLimit(int64(10)).SetSort(bson.D{{Key: "_id", Value: -1}})
	cur, err := ci.DB.Collection(model.CollectionColl).Find(ctx, queryFilter, queryOpts)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("no collections found")
		}
		return nil, errors.Wrapf(err, "error querying the database")
	}

	var collectionResp []schema.CollectionResp
	if err := cur.All(ctx, &collectionResp); err != nil {
		return nil, err
	}

	for i, cr := range collectionResp {
		if cr.Type == model.InfluencerCollection {
			influencer_info, err := ci.App.Content.GetInfluencerInfo(cr.InfluencerIDs)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to get influencer info")
			}
			collectionResp[i].InfluenncerInfo = influencer_info
		}
	}

	return collectionResp, nil
}
