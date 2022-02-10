//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_content.go -package=mock go-app/app Content

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
	GetContent(*schema.GetContentFilter) ([]schema.GetContentResp, error)
	ChangeContentStatus(*schema.ChangeContentStatusOpts) (bool, error)
	DeleteContent(primitive.ObjectID) (bool, error)

	CreatePebble(*schema.CreatePebbleOpts) (*schema.CreatePebbleResp, error)
	EditPebble(*schema.EditPebbleOpts) (*schema.EditPebbleResp, error)
	GetPebbles(opts *schema.GetPebblesKeeperFilter) ([]schema.GetContentResp, error)

	CreateCatalogVideoContent(*schema.CreateVideoCatalogContentOpts) (*schema.CreateVideoCatalogContentResp, error)
	CreateCatalogImageContent(*schema.CreateImageCatalogContentOpts) (*schema.CreateImageCatalogContentResp, error)
	EditCatalogContent(*schema.EditCatalogContentOpts) (*schema.EditCatalogContentResp, error)

	CreateVideoReviewContent(opts *schema.CreateVideoReviewContentOpts) (*schema.CreateVideoReviewContentResp, error)

	CreateComment(*schema.CreateCommentOpts) (*schema.CreateCommentResp, error)
	CreateView(*schema.CreateViewOpts) error
	CreateLike(*schema.CreateLikeOpts) error

	UpdateContentBrandInfo(*schema.UpdateContentBrandInfoOpts)
	UpdateContentInfluencerInfo(*schema.UpdateContentInfluencerInfoOpts)
	UpdateContentCatalogInfo(*schema.UpdateContentCatalogInfoOpts)
	AddContentComment(opts *schema.ProcessCommentOpts)
	DeleteContentLike(opts *schema.ProcessLikeOpts)
	AddContentLike(opts *schema.ProcessLikeOpts)
	AddContentView(opts *schema.ProcessViewOpts)

	// External API functions
	GetBrandInfo([]string) ([]model.BrandInfo, error)
	GetInfluencerInfo([]string) ([]model.InfluencerInfo, error)
	GetCatalogInfo([]string) ([]model.CatalogInfo, error)

	SearchPebbleByCaption(*schema.SearchPebbleByCaption) ([]schema.GetPebbleSearchCaptionResp, error)

	//APIs for Creators on the App
	CreatePebbleApp(opts *schema.CreatePebbleAppOpts) (*schema.CreatePebbleResp, error)
	EditPebbleApp(opts *schema.EditPebbleAppOpts) (*schema.EditPebbleAppResp, error)
	GetPebblesForCreator(opts *schema.GetPebblesCreatorFilter) ([]schema.CreatorGetContentResp, error)
	ContentProcessFail(opts *schema.ContentProcessFail)
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
	// Setting up category path
	for _, id := range opts.CategoryID {
		path, err := ci.App.Category.GetCategoryPath(id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create pebble document, error fetching category path")
		}
		pebble.Paths = append(pebble.Paths, path)
	}

	res1, err1 := ci.DB.Collection(model.ContentColl).InsertOne(context.TODO(), pebble)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed to create pebble document")
	}

	fType, err := FileTypeFromFileName(opts.FileName)
	if err != nil {
		return nil, errors.Wrap(err, "invalid file type: missing file extension")
	}
	// Getting s3 upload token with provided args
	// This token is then used by frontend to directly upload media to s3
	res0, err0 := ci.App.Media.GenerateVideoUploadToken(
		&schema.GenerateVideoUploadTokenOpts{
			FileName: fmt.Sprintf("%s.%s", res1.InsertedID.(primitive.ObjectID).Hex(), fType),
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
	if len(opts.CategoryID) > 0 {
		var paths []model.Path
		// Setting up category path
		for _, id := range opts.CategoryID {
			path, err := ci.App.Category.GetCategoryPath(id)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update pebble document, error fetching category path")
			}
			paths = append(paths, path)
		}
		update = append(update, bson.E{Key: "category_path", Value: paths})
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

// DeleteContent removes the content instance from DB
func (ci *ContentImpl) DeleteContent(id primitive.ObjectID) (bool, error) {
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
	if !c.MediaID.IsZero() {
		ci.App.Media.DeleteMedia(c.MediaID)
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
	ctx := context.TODO()
	cID, err := primitive.ObjectIDFromHex(strings.Split(opts.FileName, ".")[0])
	if err != nil {
		ci.Logger.Err(err).Interface("opts_info", opts).Msgf("failed to parse id from filename:%s while processing video content", opts.FileName)
		return false, errors.Wrapf(err, "failed to parse id from filename:%s while processing video content", opts.FileName)
	}

	var content model.Content
	if err := ci.DB.Collection(model.ContentColl).FindOne(ctx, bson.M{"_id": cID}).Decode(&content); err != nil {
		ci.Logger.Err(err).Interface("opts_info", opts).Msgf("failed to find content from id:%s while processing video content", cID.Hex())
		return false, errors.Wrapf(err, "failed to find content with id:%s", cID.Hex())
	}

	if &content == nil {
		return false, errors.Wrapf(err, "failed to find content with id:%s", cID.Hex())
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
	var update bson.M
	if content.Type != model.CatalogContentType && content.Type != model.ReviewStoryType {
		if content.CreatorID != primitive.NilObjectID {
			update = bson.M{
				"$set": bson.M{
					"is_processed": true,
					"is_active":    true,
					"processed_at": time.Now().UTC(),
					"media_type":   model.VideoType,
					"media_id":     res.ID,
				},
			}
		} else {
			update = bson.M{
				"$set": bson.M{
					"is_processed": true,
					"processed_at": time.Now().UTC(),
					"media_type":   model.VideoType,
					"media_id":     res.ID,
				},
			}
		}

	} else {
		update = bson.M{
			"$set": bson.M{
				"is_processed": true,
				"is_active":    true,
				"processed_at": time.Now().UTC(),
				"media_type":   model.VideoType,
				"media_id":     res.ID,
			},
		}
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

// CreateCatalogVideoContent creates video content for catalog
func (ci *ContentImpl) CreateCatalogVideoContent(opts *schema.CreateVideoCatalogContentOpts) (*schema.CreateVideoCatalogContentResp, error) {
	cc := model.Content{
		Type:       model.CatalogContentType,
		MediaType:  model.VideoType,
		BrandIDs:   []primitive.ObjectID{opts.BrandID},
		CatalogIDs: []primitive.ObjectID{opts.CatalogID},
		CreatedAt:  time.Now().UTC(),
	}

	res, err := ci.DB.Collection(model.ContentColl).InsertOne(context.TODO(), cc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create catalog content")
	}
	cc.ID = res.InsertedID.(primitive.ObjectID)
	fType, err := FileTypeFromFileName(opts.FileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get video file extension")
	}

	// Getting s3 upload token with provided args
	// This token is then used by frontend to directly upload media to s3
	res1, err1 := ci.App.Media.GenerateVideoUploadToken(
		&schema.GenerateVideoUploadTokenOpts{
			FileName: fmt.Sprintf("%s.%s", cc.ID.Hex(), fType),
		},
	)
	if err1 != nil {
		return nil, err1
	}
	return &schema.CreateVideoCatalogContentResp{ID: cc.ID, Token: res1.Token}, nil
}

// CreateCatalogVideoContent creates video content for catalog
func (ci *ContentImpl) CreateVideoReviewContent(opts *schema.CreateVideoReviewContentOpts) (*schema.CreateVideoReviewContentResp, error) {
	cc := model.Content{
		Type:       model.ReviewStoryType,
		MediaType:  model.VideoType,
		BrandIDs:   []primitive.ObjectID{opts.BrandID},
		CatalogIDs: []primitive.ObjectID{opts.CatalogID},
		UserID:     opts.UserID,
		CreatedAt:  time.Now().UTC(),
	}

	res, err := ci.DB.Collection(model.ContentColl).InsertOne(context.TODO(), cc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create review content")
	}
	cc.ID = res.InsertedID.(primitive.ObjectID)
	fType, err := FileTypeFromFileName(opts.FileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get video file extension")
	}

	// Getting s3 upload token with provided args
	// This token is then used by frontend to directly upload media to s3
	res1, err1 := ci.App.Media.GenerateVideoUploadToken(
		&schema.GenerateVideoUploadTokenOpts{
			FileName: fmt.Sprintf("%s.%s", cc.ID.Hex(), fType),
		},
	)
	if err1 != nil {
		return nil, err1
	}
	return &schema.CreateVideoCatalogContentResp{ID: cc.ID, Token: res1.Token}, nil
}

// EditCatalogContent updates the catalog content allowed editable fields
/*
	Allowed Fields:
		IsActive
*/
func (ci *ContentImpl) EditCatalogContent(opts *schema.EditCatalogContentOpts) (*schema.EditCatalogContentResp, error) {
	var update bson.D
	if opts.IsActive != nil {
		update = append(update, bson.E{Key: "is_active", Value: opts.IsActive})
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

	return &schema.EditCatalogContentResp{
		ID:       opts.ID,
		Label:    opts.Label,
		IsActive: opts.IsActive,
	}, nil
}

func (ci *ContentImpl) CreateCatalogImageContent(opts *schema.CreateImageCatalogContentOpts) (*schema.CreateImageCatalogContentResp, error) {
	cc := model.Content{
		Type:        model.CatalogContentType,
		MediaType:   model.ImageType,
		MediaID:     opts.MediaID,
		BrandIDs:    []primitive.ObjectID{opts.BrandID},
		CatalogIDs:  []primitive.ObjectID{opts.CatalogID},
		IsProcessed: true,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}

	res, err := ci.DB.Collection(model.ContentColl).InsertOne(context.TODO(), cc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create catalog content")
	}
	cc.ID = res.InsertedID.(primitive.ObjectID)

	return &schema.CreateImageCatalogContentResp{
		ID: cc.ID,
	}, nil
}

func (ci *ContentImpl) CreateComment(opts *schema.CreateCommentOpts) (*schema.CreateCommentResp, error) {
	c := model.Comment{
		ResourceType: opts.ResourceType,
		ResourceID:   opts.ResourceID,
		Description:  opts.Description,
		UserID:       opts.UserID,
		CreatedAt:    opts.CreatedAt,
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	res, err := ci.DB.Collection(model.CommentColl).InsertOne(context.TODO(), c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert comment")
	}

	return &schema.CreateCommentResp{
		ID:           res.InsertedID.(primitive.ObjectID),
		ResourceType: c.ResourceType,
		ResourceID:   c.ResourceID,
		Description:  c.Description,
		UserID:       c.UserID,
		CreatedAt:    c.CreatedAt,
	}, nil
}

func (ci *ContentImpl) CreateView(opts *schema.CreateViewOpts) error {
	v := model.View{
		ResourceType: opts.ResourceType,
		ResourceID:   opts.ResourceID,
		Duration:     opts.Duration,
		UserID:       opts.UserID,
		CreatedAt:    opts.CreatedAt,
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now().UTC()
	}

	_, err := ci.DB.Collection(model.ViewColl).InsertOne(context.TODO(), v)
	if err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to create view")
		return errors.Wrap(err, "failed to create view")
	}
	return nil
}

// CreateLike register a new like if like does not exists for that specific user else remove the like
func (ci *ContentImpl) CreateLike(opts *schema.CreateLikeOpts) error {
	ctx := context.TODO()
	filter := bson.M{"resource_type": opts.ResourceType, "resource_id": opts.ResourceID, "user_id": opts.UserID}
	exists, err := ci.DB.Collection(model.LikeColl).CountDocuments(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "failed to check if like exists")
	}

	// like does not exists for the user for the specific resource thus creating a new like and creating multiple likes in case of live
	if exists == 0 || opts.ResourceType == model.LiveType {
		v := model.Like{
			ResourceType: opts.ResourceType,
			ResourceID:   opts.ResourceID,
			CreatedAt:    time.Now().UTC(),
			UserID:       opts.UserID,
		}
		if _, err := ci.DB.Collection(model.LikeColl).InsertOne(ctx, v); err != nil {
			ci.Logger.Err(err).Interface("opts", opts).Msg("failed to create like")
			return errors.Wrap(err, "failed to create like")
		}
		return nil
	}

	// like exists thus removing the like if content is pebble
	if opts.ResourceType == model.PebbleType {
		if _, err = ci.DB.Collection(model.LikeColl).DeleteOne(ctx, filter); err != nil {
			return errors.Wrap(err, "failed to unlike")
		}
		return nil
	}
	return nil
}

func (ci *ContentImpl) UpdateContentBrandInfo(opts *schema.UpdateContentBrandInfoOpts) {
	filter := bson.M{
		"brand_ids": opts.ID,
		"is_active": true,
		"type":      model.PebbleType,
	}
	update := bson.M{
		"$set": bson.M{
			"last_sync": time.Now().UTC(),
			"sync_type": "brand",
		},
	}
	if _, err := ci.DB.Collection(model.ContentColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to update brand last sync")
		return
	}
}

func (ci *ContentImpl) UpdateContentInfluencerInfo(opts *schema.UpdateContentInfluencerInfoOpts) {
	filter := bson.M{
		"influencer_ids": opts.ID,
		"is_active":      true,
		"type":           model.PebbleType,
	}
	update := bson.M{
		"$set": bson.M{
			"last_sync": time.Now().UTC(),
			"sync_type": "influencer",
		},
	}
	if _, err := ci.DB.Collection(model.ContentColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to update influencer last sync")
	}
}

func (ci *ContentImpl) UpdateContentCatalogInfo(opts *schema.UpdateContentCatalogInfoOpts) {
	filter := bson.M{
		"catalog_ids": opts.ID,
		"is_active":   true,
		"type":        model.PebbleType,
	}
	update := bson.M{
		"$set": bson.M{
			"last_sync": time.Now().UTC(),
			"sync_type": "catalog",
		},
	}
	if _, err := ci.DB.Collection(model.ContentColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to update influencer last sync")
	}
}

func (ci *ContentImpl) AddContentLike(opts *schema.ProcessLikeOpts) {
	var resourceColl string
	switch opts.ResourceType {
	case model.PebbleType:
		resourceColl = model.ContentColl
	case model.LiveType:
		resourceColl = model.LiveColl
	default:
		ci.Logger.Err(errors.New("invalid resource type")).Interface("opts", opts).Msg("failed to add like")
		return
	}
	filter := bson.M{
		"_id": opts.ResourceID,
	}
	update := bson.M{
		"$push": bson.M{
			"like_ids": opts.ID,
		},
		"$addToSet": bson.M{
			"liked_by": opts.UserID,
		},
		"$inc": bson.M{
			"like_count": 1,
		},
	}
	if _, err := ci.DB.Collection(resourceColl).UpdateOne(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to add like")
	}
}

func (ci *ContentImpl) AddContentView(opts *schema.ProcessViewOpts) {
	var resourceColl string
	switch opts.ResourceType {
	case model.PebbleType:
		resourceColl = model.ContentColl
	case model.LiveType:
		resourceColl = model.LiveColl
	}
	filter := bson.M{
		"_id": opts.ResourceID,
	}
	update := bson.M{
		"$inc": bson.M{
			"view_count": 1,
		},
	}
	if _, err := ci.DB.Collection(resourceColl).UpdateOne(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to add view")
	}
}

func (ci *ContentImpl) DeleteContentLike(opts *schema.ProcessLikeOpts) {
	filter := bson.M{
		"like_ids": opts.ID,
	}
	update := bson.M{
		"$pull": bson.M{
			"like_ids": opts.ID,
		},
		"$inc": bson.M{
			"like_count": -1,
		},
	}
	if _, err := ci.DB.Collection(model.ContentColl).UpdateOne(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to delete like")
		return
	}
}

func (ci *ContentImpl) AddContentComment(opts *schema.ProcessCommentOpts) {
	var resourceColl string
	switch opts.ResourceType {
	case model.PebbleType:
		resourceColl = model.ContentColl
	case model.LiveType:
		resourceColl = model.LiveColl
	}
	filter := bson.M{
		"_id": opts.ResourceID,
	}
	update := bson.M{
		"$inc": bson.M{
			"comment_count": 1,
		},
	}
	if _, err := ci.DB.Collection(resourceColl).UpdateOne(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to update like count")
	}
}

func (ci *ContentImpl) GetBrandInfo(ids []string) ([]model.BrandInfo, error) {
	var s schema.GetBrandInfoResp
	url := ci.App.Config.HypdAPIConfig.EntityAPI + "/api/keeper/brand/get"
	postBody, _ := json.Marshal(map[string][]string{
		"id": ids,
	})
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to get brand info")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ci.App.Config.HypdAPIConfig.Token)
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		ci.Logger.Err(err).Str("responseBody", string(postBody)).Msgf("failed to send request to api %s", url)
		return nil, errors.Wrap(err, "failed to get brandinfo")
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).Str("responseBody", string(postBody)).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrap(err, "failed to get brandinfo")
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		ci.Logger.Err(errors.New("success false from entity")).Str("body", string(body)).Msg("got success false response from entity")
		return nil, errors.New("got success false response from entity")
	}
	return s.Payload, nil
}

func (ci *ContentImpl) GetInfluencerInfo(ids []string) ([]model.InfluencerInfo, error) {
	var s schema.GetInfluencerInfoResp
	url := ci.App.Config.HypdAPIConfig.EntityAPI + "/api/keeper/influencer/get"
	postBody, _ := json.Marshal(map[string][]string{
		"id": ids,
	})
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to get influencer info")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ci.App.Config.HypdAPIConfig.Token)
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		ci.Logger.Err(err).Str("responseBody", string(postBody)).Msgf("failed to send request to api %s", url)
		return nil, errors.Wrapf(err, "failed to send request to api %s", url)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).Str("responseBody", string(postBody)).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrapf(err, "failed to read response from api %s", url)
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		ci.Logger.Err(errors.New("success false from entity")).Str("body", string(body)).Msg("got success false response from entity")
		return nil, errors.New("got success false response from entity")
	}
	return s.Payload, nil
}

func (ci *ContentImpl) GetCatalogInfo(ids []string) ([]model.CatalogInfo, error) {
	var s schema.GetCatalogInfoResp
	url := ci.App.Config.HypdAPIConfig.CatalogAPI + "/api/keeper/catalog/get/ids"
	postBody, _ := json.Marshal(map[string][]string{
		"id": ids,
	})
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to get catalog info")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ci.App.Config.HypdAPIConfig.Token)
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		ci.Logger.Err(err).Str("responseBody", string(postBody)).Msgf("failed to send request to api %s", url)
		return nil, errors.Wrapf(err, "failed to send request to api %s", url)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).Str("responseBody", string(postBody)).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrapf(err, "failed to read response from api %s", url)
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		ci.Logger.Err(errors.New("success false from catalog")).Str("body", string(body)).Msg("got success false response from catalog")
		return nil, errors.New("got success false response from catalog")
	}
	return s.Payload, nil
}

func (ci *ContentImpl) GetPebbles(opts *schema.GetPebblesKeeperFilter) ([]schema.GetContentResp, error) {
	var resp []schema.GetContentResp
	var matchFilter bson.D
	matchFilter = append(matchFilter, bson.E{Key: "type", Value: opts.Type})
	if len(opts.InfluencerIDs) > 0 {
		matchFilter = append(matchFilter, bson.E{Key: "influencer_ids", Value: bson.M{"$in": opts.InfluencerIDs}})
	}
	if len(opts.CatalogIDs) > 0 {
		matchFilter = append(matchFilter, bson.E{Key: "catalog_ids", Value: bson.M{"$in": opts.CatalogIDs}})
	}
	if opts.Caption != "" {
		matchFilter = append(matchFilter, bson.E{Key: "caption", Value: bson.M{"$regex": primitive.Regex{Pattern: opts.Caption, Options: "i"}}})
	}
	var matchStage bson.D
	matchStage = append(matchStage, bson.E{Key: "$match", Value: matchFilter})
	sortStage := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"_id": -1,
		},
	}}
	skipStage := bson.D{
		{
			Key:   "$skip",
			Value: opts.Page * 20,
		},
	}

	limitStage := bson.D{
		{
			Key:   "$limit",
			Value: 20,
		},
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

	ctx := context.TODO()
	cur, err := ci.DB.Collection(model.ContentColl).Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage, lookupStage, setStage})
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get pebbles")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, " failed to get pebbles")
	}
	return resp, nil
}

func (ci *ContentImpl) ChangeContentStatus(opts *schema.ChangeContentStatusOpts) (bool, error) {
	filter := bson.M{
		"_id": opts.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"is_active": opts.IsActive,
		},
	}

	if _, err := ci.DB.Collection(model.ContentColl).UpdateOne(context.TODO(), filter, update); err != nil {
		return false, errors.Wrap(err, "failed to update content status")
	}

	return true, nil
}

func (ci *ContentImpl) SearchPebbleByCaption(opts *schema.SearchPebbleByCaption) ([]schema.GetPebbleSearchCaptionResp, error) {

	ctx := context.TODO()

	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"is_active": true,
			"type":      "pebble",
			"caption": bson.M{
				"$regex": primitive.Regex{Pattern: opts.Caption, Options: "i"},
			},
		},
	}}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "media",
			"localField":   "media_id",
			"foreignField": "_id",
			"as":           "media_info",
		},
	}}

	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"media_info": bson.M{
				"$first": "$media_info",
			},
		},
	}}

	skipStage := bson.D{
		{
			Key:   "$skip",
			Value: 20 * opts.Page,
		},
	}

	limitStage := bson.D{
		{
			Key:   "$limit",
			Value: 20,
		},
	}

	var resp []schema.GetPebbleSearchCaptionResp

	cur, err := ci.DB.Collection(model.ContentColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, setStage, skipStage, limitStage})
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get content")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode pebbles")
	}
	return resp, nil
}

// CreatePebbleApp creates new create a new pebble document in the db, and generates and returns a token to upload video
func (ci *ContentImpl) CreatePebbleApp(opts *schema.CreatePebbleAppOpts) (*schema.CreatePebbleResp, error) {
	pebble := model.Content{
		Type:          model.PebbleType,
		MediaType:     model.VideoType,
		Caption:       opts.Caption,
		Hashtags:      ParseHashtag(opts.Caption),
		CreatorID:     opts.CreatorID,
		InfluencerIDs: opts.InfluencerIDs,
		BrandIDs:      opts.BrandIDs,
		CatalogIDs:    opts.CatalogIDs,
		// Label: &model.Label{
		// 	AgeGroups: opts.Label.AgeGroup,
		// 	Genders:   opts.Label.Gender,
		// },
		CreatedAt: time.Now().UTC(),
	}
	// Setting up category path
	for _, id := range opts.CategoryID {
		path, err := ci.App.Category.GetCategoryPath(id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create pebble document, error fetching category path")
		}
		pebble.Paths = append(pebble.Paths, path)
	}

	res1, err1 := ci.DB.Collection(model.ContentColl).InsertOne(context.TODO(), pebble)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed to create pebble document")
	}

	fType, err := FileTypeFromFileName(opts.FileName)
	if err != nil {
		return nil, errors.Wrap(err, "invalid file type: missing file extension")
	}
	// Getting s3 upload token with provided args
	// This token is then used by frontend to directly upload media to s3
	res0, err0 := ci.App.Media.GenerateVideoUploadToken(
		&schema.GenerateVideoUploadTokenOpts{
			FileName: fmt.Sprintf("%s.%s", res1.InsertedID.(primitive.ObjectID).Hex(), fType),
		},
	)
	if err0 != nil {
		return nil, err0
	}
	return &schema.CreatePebbleResp{ID: res1.InsertedID.(primitive.ObjectID), Token: res0.Token}, nil
}

// EditPebbleApp updates the pebble document Fields
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
func (ci *ContentImpl) EditPebbleApp(opts *schema.EditPebbleAppOpts) (*schema.EditPebbleAppResp, error) {
	var update bson.D

	if opts.Caption != "" {
		update = append(update, bson.E{Key: "caption", Value: opts.Caption})
		update = append(update, bson.E{Key: "hashtags", Value: ParseHashtag(opts.Caption)})
	}
	if opts.BrandIDs != nil {
		update = append(update, bson.E{Key: "brand_ids", Value: opts.BrandIDs})
	}
	if opts.CatalogIDs != nil {
		update = append(update, bson.E{Key: "catalog_ids", Value: opts.CatalogIDs})
	}
	if len(opts.InfluencerIDs) > 0 {
		update = append(update, bson.E{Key: "influencer_ids", Value: opts.InfluencerIDs})
	}
	// if opts.Label != nil {
	// 	if len(opts.Label.AgeGroup) > 0 {
	// 		update = append(update, bson.E{Key: "label.age_groups", Value: opts.Label.AgeGroup})
	// 	}
	// 	if len(opts.Label.Gender) > 0 {
	// 		update = append(update, bson.E{Key: "label.genders", Value: opts.Label.Gender})
	// 	}
	// }
	if opts.IsActive != nil {
		update = append(update, bson.E{Key: "is_active", Value: opts.IsActive})
	}
	if len(opts.CategoryID) > 0 {
		var paths []model.Path
		// Setting up category path
		for _, id := range opts.CategoryID {
			path, err := ci.App.Category.GetCategoryPath(id)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update pebble document, error fetching category path")
			}
			paths = append(paths, path)
		}
		update = append(update, bson.E{Key: "category_path", Value: paths})
	}
	filter := bson.M{"_id": opts.ID, "creator_id": opts.CreatorID}
	updateQuery := bson.M{"$set": update}

	res, err := ci.DB.Collection(model.ContentColl).UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update pebble, please try again in few minutes")
	}
	if res.MatchedCount == 0 {
		return nil, errors.Wrap(err, "failed to find pebble, please check if you are the creator")
	}
	if res.ModifiedCount == 0 {
		return nil, errors.Wrap(err, "failed to update content")
	}
	return &schema.EditPebbleAppResp{
		ID:            opts.ID,
		Caption:       opts.Caption,
		InfluencerIDs: opts.InfluencerIDs,
		BrandIDs:      opts.BrandIDs,
		CatalogIDs:    opts.CatalogIDs,
		// Label:         opts.Label,
		IsActive: opts.IsActive,
	}, nil
}

func (ci *ContentImpl) GetPebblesForCreator(opts *schema.GetPebblesCreatorFilter) ([]schema.CreatorGetContentResp, error) {
	var resp []schema.CreatorGetContentResp
	var matchFilter bson.D
	matchFilter = append(matchFilter, bson.E{Key: "type", Value: "pebble"})
	if len(opts.InfluencerIDs) > 0 {
		matchFilter = append(matchFilter, bson.E{Key: "influencer_ids", Value: bson.M{"$in": opts.InfluencerIDs}})
	}
	if len(opts.CatalogIDs) > 0 {
		matchFilter = append(matchFilter, bson.E{Key: "catalog_ids", Value: bson.M{"$in": opts.CatalogIDs}})
	}
	//showing only processed video
	matchFilter = append(matchFilter, bson.E{Key: "is_processed", Value: true})

	var matchStage bson.D
	matchStage = append(matchStage, bson.E{Key: "$match", Value: matchFilter})
	sortStage := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"_id": -1,
		},
	}}
	skipStage := bson.D{
		{
			Key:   "$skip",
			Value: opts.Page * 20,
		},
	}

	limitStage := bson.D{
		{
			Key:   "$limit",
			Value: 20,
		},
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

	ctx := context.TODO()
	cur, err := ci.DB.Collection(model.ContentColl).Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage, lookupStage, setStage})
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get pebbles")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, " failed to get pebbles")
	}
	return resp, nil
}

func (ci *ContentImpl) ContentProcessFail(opts *schema.ContentProcessFail) {

	ctx := context.TODO()
	sid := opts.Event["srcVideo"]
	fmt.Println("sid", sid)
	sid = strings.Split(sid, ".")[0]
	fmt.Println(sid)
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		ci.Logger.Err(err).Msg("error getting media id from request")
		return
	}
	filter := bson.M{
		"_id": id,
	}
	update := bson.M{
		"$set": bson.M{
			"is_active":    false,
			"is_processed": true,
			"is_failed":    true,
			"error_msg":    opts.ErrorMessage,
		},
	}
	_, err = ci.DB.Collection(model.ContentColl).UpdateOne(ctx, filter, update)
	if err != nil {
		ci.Logger.Err(err).Msg("error updating process fail data to DB")
		return
	}
	ci.Logger.Log().Msg("Process failed data saved successfully")
}
