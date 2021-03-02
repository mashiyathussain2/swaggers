//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_content.go -package=mock go-app/app Content

package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Content contains methods to implement and operation pebble(video-only) content
type Content interface {
	ProcessVideoContent(*schema.ProcessVideoContentOpts) (bool, error)
	GetContentByID(primitive.ObjectID) (*schema.GetContentResp, error)

	CreatePebble(*schema.CreatePebbleOpts) (*schema.CreatePebbleResp, error)
	EditPebble(*schema.EditPebbleOpts) (*schema.EditPebbleResp, error)
	DeletePebble(primitive.ObjectID) (bool, error)
}

// ContentImpl implements `Pebble` functionality
type ContentImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// ContentOpts contains args required to create a new instance of `PebbleImpl`
type ContentOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitContent returns a new instance of `Pebble` Implementation
func InitContent(opts *ContentOpts) Content {
	p := ContentImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &p
}

// CreatePebble creates new create a new pebble document in the db, and generates and returns a token to upload video
func (ci *ContentImpl) CreatePebble(opts *schema.CreatePebbleOpts) (*schema.CreatePebbleResp, error) {
	pebble := model.Content{
		Type:          model.PebbleType,
		MediaType:     model.VideoType,
		Caption:       opts.Caption,
		Hashtags:      ParseHashtag(opts.Caption),
		InfluencerIDs: opts.InfluencerIDs,
		BrandIDs:      opts.BrandIDs,
		CatalogIDs:    opts.CatalogIDs,
		Label: &model.Label{
			Interests: opts.Label.Interests,
			AgeGroups: opts.Label.AgeGroup,
			Genders:   opts.Label.Gender,
		},
		CreatedAt: time.Now().UTC(),
	}

	res1, err1 := ci.DB.Collection(model.ContentColl).InsertOne(context.TODO(), pebble)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed to create pebble document")
	}

	// Getting s3 upload token with provided args
	// This token is then used by frontend to directly upload media to s3
	res0, err0 := ci.App.Media.GenerateVideoUploadToken(
		&schema.GenerateVideoUploadTokenOpts{
			FileName: res1.InsertedID.(primitive.ObjectID).Hex(),
		},
	)
	if err0 != nil {
		return nil, err0
	}
	return &schema.CreatePebbleResp{ID: res1.InsertedID.(primitive.ObjectID), Token: res0.Token}, nil
}

// EditPebble updates the pebble document Fields
/*
	Fields available to update
		InfluencerIDs
		BrandIDs
		Caption
		Hashtags
		CatalogIDs
		IsActive
		Label
*/
func (ci *ContentImpl) EditPebble(opts *schema.EditPebbleOpts) (*schema.EditPebbleResp, error) {
	var update bson.D

	if opts.Caption != "" {
		update = append(update, bson.E{Key: "caption", Value: opts.Caption})
		update = append(update, bson.E{Key: "hashtags", Value: ParseHashtag(opts.Caption)})
	}
	if len(opts.BrandIDs) > 0 {
		update = append(update, bson.E{Key: "brand_ids", Value: opts.BrandIDs})
	}
	if len(opts.CatalogIDs) > 0 {
		update = append(update, bson.E{Key: "catalog_ids", Value: opts.CatalogIDs})
	}
	if len(opts.InfluencerIDs) > 0 {
		update = append(update, bson.E{Key: "influencer_ids", Value: opts.InfluencerIDs})
	}
	if opts.Label != nil {
		if len(opts.Label.AgeGroup) > 0 {
			update = append(update, bson.E{Key: "label.age_groups", Value: opts.Label.AgeGroup})
		}
		if len(opts.Label.Gender) > 0 {
			update = append(update, bson.E{Key: "label.genders", Value: opts.Label.Gender})
		}
		if len(opts.Label.Interests) > 0 {
			update = append(update, bson.E{Key: "label.interests", Value: opts.Label.Interests})
		}
	}
	if opts.IsActive != nil {
		update = append(update, bson.E{Key: "is_active", Value: opts.IsActive})
	}

	filter := bson.M{"_id": opts.ID}
	updateQuery := bson.M{"$set": update}

	res, err := ci.DB.Collection(model.ContentColl).UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to update content")
	}
	if res.MatchedCount == 0 {
		return nil, errors.Wrap(err, "failed to find content")
	}
	if res.ModifiedCount == 0 {
		return nil, errors.Wrap(err, "failed to update content")
	}
	return &schema.EditPebbleResp{
		ID:            opts.ID,
		Caption:       opts.Caption,
		InfluencerIDs: opts.InfluencerIDs,
		BrandIDs:      opts.BrandIDs,
		CatalogIDs:    opts.CatalogIDs,
		Label:         opts.Label,
		IsActive:      opts.IsActive,
	}, nil
}

// DeletePebble removes the pebble instance from DB
func (ci *ContentImpl) DeletePebble(id primitive.ObjectID) (bool, error) {
	ctx := context.TODO()
	var c model.Content

	findOpts := options.FindOne().SetProjection(bson.M{"_id": 1, "media_type": 1, "media_id": 1})
	findFilter := bson.M{"_id": id}
	if err := ci.DB.Collection(model.ContentColl).FindOne(ctx, findFilter, findOpts).Decode(&c); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, errors.Wrapf(err, "content with id: %s not found", id.Hex())
		}
		return false, errors.Wrapf(err, "failed to find content with id: %s", id.Hex())
	}

	// Deleting media document reference from media collection
	if _, err := ci.App.Media.DeleteMedia(c.MediaID); err != nil {
		return false, err
	}

	// Deleting content document from cotent collection
	if _, err := ci.DB.Collection(model.ContentColl).DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		ci.Logger.Err(err).Msgf("failed to delete content with id: %s but media with id: %s is deleted", id.Hex(), c.MediaID.Hex())
		return false, errors.Wrapf(err, "failed to delete content with id: %s", id.Hex())
	}

	return true, nil
}

// ProcessVideoContent mark video content as processed
func (ci *ContentImpl) ProcessVideoContent(opts *schema.ProcessVideoContentOpts) (bool, error) {
	// Extracting content id from filename EG: 283782738273823.mp4
	cID, err := primitive.ObjectIDFromHex(strings.Split(opts.FileName, ".")[0])
	if err != nil {
		ci.Logger.Err(err).Interface("opts_info", opts).Msgf("failed to parse id from filename:%s while processing video content", opts.FileName)
		return false, errors.Wrapf(err, "failed to parse id from filename:%s while processing video content", opts.FileName)
	}

	// Creating media object from data received
	res, err := ci.App.Media.CreateVideoMedia(opts)
	if err != nil {
		ci.Logger.Err(err).Msg("failed to create video media")
		return false, errors.Wrap(err, "failed to create video media")
	}

	// Updating content as IsProcessed true and linking content with media received from above
	var c model.Content
	filter := bson.M{"_id": cID}
	update := bson.M{
		"$set": bson.M{
			"is_processed": true,
			"processed_at": time.Now().UTC(),
			"media_type":   model.VideoType,
			"media_id":     res.ID,
		},
	}
	if err := ci.DB.Collection(model.ContentColl).FindOneAndUpdate(context.TODO(), filter, update).Decode(&c); err != nil {
		ci.Logger.Err(err).Interface("media_info", res).Msgf("failed to mark content:%s as processed", cID.Hex())
		return false, errors.Wrapf(err, "failed to mark content:%s as processed", cID.Hex())
	}
	return true, nil
}

// GetContentByID returns the content document matching with the id
func (ci *ContentImpl) GetContentByID(id primitive.ObjectID) (*schema.GetContentResp, error) {
	var c schema.GetContentResp
	if err := ci.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&c); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "content with id: %s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "query failed to find content with id: %s", id.Hex())
	}

	// Joning content with image/video info
	if !c.MediaID.IsZero() {
		switch c.MediaType {
		case model.VideoType:
			vi, err := ci.App.Media.GetVideoMediaByID(c.MediaID)
			if err != nil {
				return nil, err
			}
			c.MediaInfo = vi
		}
	}
	return &c, nil
}

// GetContent returns multiple content object based on applied filter
func (ci *ContentImpl) GetContent(filterOpts *schema.GetContentFilter) ([]schema.GetContentResp, error) {
	var filter bson.D

	// Setting up filters
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
	var res []schema.GetContentResp
	cur, err := ci.DB.Collection(model.ContentColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get content")
	}
	if err := cur.All(ctx, &res); err != nil {
		return nil, errors.Wrap(err, "failed to get content")
	}

	return res, nil
}
