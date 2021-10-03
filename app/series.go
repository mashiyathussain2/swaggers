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

type Series interface {

	//Pebble Series
	CreateSeries(opts *schema.CreateSeriesOpts) error
	UpdateSeries(opts *schema.UpdateSeriesOpts) error
	GetContentForPebbleSeries(filterOpts *schema.GetContentFilter) ([]schema.ContentForSeries, error)
	KeeperGetSeries(opts *schema.GetSeriesKeeperFilter) ([]schema.PebbleSeriesResp, error)
	KeeperGetSeriesByID(opts *schema.KeeperGetSeriesByID) (*schema.KeeperPebbleSeriesResp, error)
	// SearchSeriesByName(opts *schema.SearchSeriesByName) ([]schema.PebbleSeriesResp, error)
	KeeperGetSeriesBasic(opts *schema.KeeperGetSeriesBasic) ([]schema.KeeperGetSeriesBasicResp, error)
	UpdateSeriesLastSync(id primitive.ObjectID) error
}

// SeriesImpl implements `Pebble Series` functionality
type SeriesImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// SeriesOpts contains args required to create a new instance of `SeriesImpl`
type SeriesOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitSeries returns a new instance of `Series` Implementation
func InitSeries(opts *SeriesOpts) Series {
	p := SeriesImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &p
}

func (si *SeriesImpl) CreateSeries(opts *schema.CreateSeriesOpts) error {
	ctx := context.TODO()
	s := model.PebbleSeries{
		Name:      opts.Name,
		PebbleIds: opts.PebbleIds,
		Thumbnail: &model.IMG{SRC: opts.Thumbnail.SRC},
		IsActive:  false,
		CreatedAt: time.Now(),
		Label: &model.SeriesLabel{
			Genders: opts.Label.Gender,
		},
	}
	if err := s.Thumbnail.LoadFromURL(); err != nil {
		return errors.Wrap(err, "failed to load thumbnail image")
	}
	res, err := si.DB.Collection(model.PebbleSeriesColl).InsertOne(ctx, s)
	if err != nil {
		return errors.Wrapf(err, "unable to create new series")
	}

	pebbleFilter := bson.M{
		"_id": bson.M{
			"$in": opts.PebbleIds,
		},
	}
	pebbleUpdate := bson.M{
		"$addToSet": bson.M{
			"series_ids": res.InsertedID,
		},
	}
	_, err = si.DB.Collection(model.ContentColl).UpdateMany(ctx, pebbleFilter, pebbleUpdate)
	if err != nil {
		return errors.Wrapf(err, "unable to add series ids to pebble")
	}
	return nil
}

func (si *SeriesImpl) UpdateSeries(opts *schema.UpdateSeriesOpts) error {
	ctx := context.TODO()
	s := model.PebbleSeries{
		UpdatedAt: time.Now(),
	}
	if opts.Name != "" {
		s.Name = opts.Name
	}
	if opts.Thumbnail != nil && opts.Thumbnail.SRC != "" {
		s.Thumbnail = &model.IMG{SRC: opts.Thumbnail.SRC}
		if err := s.Thumbnail.LoadFromURL(); err != nil {
			return errors.Wrap(err, "failed to load thumbnail image")
		}
	}
	if len(opts.PebbleIds) != 0 {
		s.PebbleIds = opts.PebbleIds
	}
	if opts.IsActive != nil {
		s.IsActive = *opts.IsActive
	}
	if opts.Label != nil {
		s.Label = &model.SeriesLabel{
			Genders: opts.Label.Gender,
		}
	}
	filter := bson.M{
		"_id": opts.ID,
	}
	update := bson.M{
		"$set": s,
	}
	err := si.DB.Collection(model.PebbleSeriesColl).FindOneAndUpdate(ctx, filter, update).Decode(&s)
	if err != nil {
		return errors.Wrapf(err, "unable to create update series")
	}
	if len(opts.PebbleIds) > 0 {
		pebbleFilter := bson.M{
			"_id": bson.M{
				"$in": opts.PebbleIds,
			},
		}
		pebbleUpdate := bson.M{
			"$addToSet": bson.M{
				"series_ids": opts.ID,
			},
		}
		_, err = si.DB.Collection(model.ContentColl).UpdateMany(ctx, pebbleFilter, pebbleUpdate)
		if err != nil {
			return errors.Wrapf(err, "unable to add series ids to pebble")
		}
	}
	return nil
}

// GetContent returns multiple content object based on applied filter
func (si *SeriesImpl) GetContentForPebbleSeries(filterOpts *schema.GetContentFilter) ([]schema.ContentForSeries, error) {
	var filter bson.D

	// Setting up filters
	if len(filterOpts.IDs) > 0 {
		filter = append(filter, bson.E{Key: "_id", Value: bson.M{"$in": filterOpts.IDs}})
	}
	if len(filterOpts.BrandIDs) > 0 {
		filter = append(filter, bson.E{Key: "brand_ids", Value: bson.M{"$in": filterOpts.BrandIDs}})
	}
	if len(filterOpts.CatalogIDs) > 0 {
		filter = append(filter, bson.E{Key: "catalog_ids", Value: bson.M{"$in": filterOpts.CatalogIDs}})
	}
	if filterOpts.IsActive != nil {
		filter = append(filter, bson.E{Key: "is_active", Value: *filterOpts.IsActive})
	}
	if filterOpts.IsProcessed != nil {
		filter = append(filter, bson.E{Key: "is_processed", Value: *filterOpts.IsProcessed})
	}
	if filterOpts.MediaType != "" {
		filter = append(filter, bson.E{Key: "media_type", Value: filterOpts.MediaType})
	}
	if filterOpts.Type != "" {
		filter = append(filter, bson.E{Key: "type", Value: filterOpts.Type})
	}
	if len(filterOpts.Hashtags) > 0 {
		filter = append(filter, bson.E{Key: "hashtags", Value: bson.M{"$in": filterOpts.Hashtags}})
	}
	if !filterOpts.From.IsZero() {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$gte": filterOpts.From}})
	}
	if !filterOpts.To.IsZero() {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$lte": filterOpts.From}})
	}

	var pipeline mongo.Pipeline

	if filter != nil {
		matchStage := bson.D{
			{
				Key:   "$match",
				Value: filter,
			},
		}
		pipeline = append(pipeline, matchStage)
	}

	// when page is set == 999 will return all the matching documents and skip pagination
	if filterOpts.Page != 999 {
		skipStage := bson.D{
			{
				Key:   "$skip",
				Value: 10 * filterOpts.Page,
			},
		}
		pipeline = append(pipeline, skipStage)

		limitStage := bson.D{
			{
				Key:   "$limit",
				Value: 10,
			},
		}
		pipeline = append(pipeline, limitStage)
	}

	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         model.MediaColl,
				"localField":   "media_id",
				"foreignField": "_id",
				"as":           "media_info",
			},
		},
	}
	pipeline = append(pipeline, lookupStage)

	setStage := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"media_info": bson.M{
					"$arrayElemAt": bson.A{
						"$media_info",
						0,
					},
				},
			},
		},
	}
	pipeline = append(pipeline, setStage)

	ctx := context.TODO()
	var res []schema.ContentForSeries
	cur, err := si.DB.Collection(model.ContentColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get content")
	}
	if err := cur.All(ctx, &res); err != nil {
		return nil, errors.Wrap(err, "failed to get content")
	}
	for i, r := range res {
		var ids []string
		for _, id := range r.InfluencerIDs {
			ids = append(ids, id.Hex())
		}
		influencerInfo, err := si.App.Content.GetInfluencerInfo(ids)
		if err != nil {
			si.Logger.Err(err).Interface("data", r).Msg("failed to get content influencer info")
			return nil, err
		}
		res[i].InfluencerInfo = influencerInfo

		var cat_ids []string
		for _, id := range r.CatalogIDs {
			cat_ids = append(cat_ids, id.Hex())
		}
		catalogInfo, err := si.App.Content.GetCatalogInfo(cat_ids)
		if err != nil {
			si.Logger.Err(err).Interface("data", r).Msg("failed to get content catalog info")
		}
		res[i].CatalogInfo = catalogInfo
	}
	return res, nil
}

func (si *SeriesImpl) KeeperGetSeries(opts *schema.GetSeriesKeeperFilter) ([]schema.PebbleSeriesResp, error) {
	var queryFilter bson.D

	if opts.Name != "" {
		queryFilter = append(queryFilter, bson.E{Key: "name", Value: bson.M{"$regex": primitive.Regex{Pattern: opts.Name, Options: "i"}}})
	}

	if opts.IsActive {
		queryFilter = append(queryFilter, bson.E{Key: "is_active", Value: opts.IsActive})
	}
	if queryFilter == nil {
		queryFilter = bson.D{{}}
	}
	ctx := context.TODO()
	queryOpts := options.Find().SetSkip(int64(10 * opts.Page)).SetLimit(int64(10)).SetSort(bson.D{{Key: "_id", Value: -1}})
	cur, err := si.DB.Collection(model.PebbleSeriesColl).Find(ctx, queryFilter, queryOpts)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("no series found")
		}
		return nil, errors.Wrapf(err, "error querying the database")
	}
	var seriesResp []schema.PebbleSeriesResp
	if err := cur.All(ctx, &seriesResp); err != nil {
		return nil, err
	}
	return seriesResp, nil
}

func (si *SeriesImpl) KeeperGetSeriesByID(opts *schema.KeeperGetSeriesByID) (*schema.KeeperPebbleSeriesResp, error) {

	ctx := context.TODO()
	var resp []schema.KeeperPebbleSeriesResp
	id, err := primitive.ObjectIDFromHex(opts.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "id is incorrect")

	}
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": id,
		},
	}}

	lookupStage1 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "content",
			"localField":   "pebble_ids",
			"foreignField": "_id",
			"as":           "pebble_info",
		},
	}}

	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path": "$pebble_info",
		},
	}}

	lookupStage2 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "media",
			"localField":   "pebble_info.media_id",
			"foreignField": "_id",
			"as":           "pebble_info.media_info",
		},
	}}

	addFieldStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"pebble_info.media_info": bson.M{
				"$first": "$pebble_info.media_info",
			},
		},
	}}
	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "_$id",
			"pebble_info": bson.M{
				"$push": "$pebble_info",
			},
			"root": bson.M{
				"$push": "$$ROOT",
			},
		},
	}}

	setStage1 := bson.D{{
		Key: "$set", Value: bson.M{
			"root": bson.M{"$first": "$root"},
		},
	}}
	setStage2 := bson.D{{
		Key: "$set", Value: bson.M{
			"root.pebble_info": "$pebble_info",
		},
	}}

	replaceRootStage := bson.D{{
		Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$root",
		},
	}}

	pipeline := mongo.Pipeline{matchStage, lookupStage1, unwindStage, lookupStage2, addFieldStage, groupStage, setStage1, setStage2, replaceRootStage}

	cursor, err := si.DB.Collection(model.PebbleSeriesColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get series data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding series info")
	}
	if len(resp) > 0 {
		return &resp[0], nil
	}
	return nil, nil
}

// func (si *SeriesImpl) SearchSeriesByName(opts *schema.SearchSeriesByName) ([]schema.PebbleSeriesResp, error) {

// 	ctx := context.TODO()
// 	filter := bson.M{
// 		"name": bson.M{
// 			"$regex": primitive.Regex{Pattern: opts.Name, Options: "i"},
// 		},
// 	}
// 	filterOpts := options.Find().SetSkip(int64(opts.Page * 10)).SetLimit(10)
// 	cur, err := si.DB.Collection(model.PebbleSeriesColl).Find(ctx, filter, filterOpts)
// 	if err != nil {
// 		si.Logger.Err(err).Msg("failed to get series")
// 		return nil, errors.Wrapf(err, "failed to get series")
// 	}
// 	var resp []schema.PebbleSeriesResp
// 	if err := cur.All(ctx, &resp); err != nil {
// 		return nil, errors.Wrap(err, "failed to decode series")
// 	}
// 	return resp, nil
// }

func (si *SeriesImpl) KeeperGetSeriesBasic(opts *schema.KeeperGetSeriesBasic) ([]schema.KeeperGetSeriesBasicResp, error) {

	ctx := context.TODO()
	filter := bson.M{
		"_id": bson.M{"$in": opts.IDs},
	}

	cur, err := si.DB.Collection(model.PebbleSeriesColl).Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("no series found")
		}
		return nil, errors.Wrapf(err, "error querying the database")
	}
	var seriesResp []schema.KeeperGetSeriesBasicResp
	if err := cur.All(ctx, &seriesResp); err != nil {
		return nil, err
	}
	return seriesResp, nil
}

func (si *SeriesImpl) UpdateSeriesLastSync(id primitive.ObjectID) error {
	ctx := context.TODO()
	s := bson.M{
		"last_sync": time.Now(),
	}

	filter := bson.M{
		"pebble_ids": id,
	}
	update := bson.M{
		"$set": s,
	}
	err := si.DB.Collection(model.PebbleSeriesColl).FindOneAndUpdate(ctx, filter, update).Decode(&s)
	if err != nil {
		return errors.Wrapf(err, "unable to create update series")
	}

	return nil
}
