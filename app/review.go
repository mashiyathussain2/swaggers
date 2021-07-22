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
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Review interface {
	CreateReviewStory(*schema.CreateReviewStoryOpts) (*schema.CreateReviewStoryResp, error)
	ProcessReviewStory(storyID primitive.ObjectID)
	GetReviewInfo(id primitive.ObjectID) (*schema.ReviewStoryFullMessage, error)
}

type ReviewImpl struct {
	App    *App
	Logger *zerolog.Logger
	DB     *mongo.Database
}

type ReviewOpts struct {
	App    *App
	Logger *zerolog.Logger
	DB     *mongo.Database
}

func InitReview(opts *ReviewOpts) Review {
	return &ReviewImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

func (ri *ReviewImpl) getReviewStoryUploadURL(opts *schema.CreateReviewStoryOpts) (*schema.GetReviewStoryUploadURLResp, error) {
	url := ri.App.Config.HypdApiConfig.CmsApi + "/api/keeper/content/review/video"
	reqOpts := schema.CreateVideoReviewContentOpts{
		FileName:  opts.FileName,
		CatalogID: opts.CatalogID,
		UserID:    opts.UserID,
		BrandID:   opts.BrandID,
	}
	reqData, _ := json.Marshal(reqOpts)
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqData))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to create review video")
	}
	req.Header.Add("Authorization", ri.App.Config.HypdApiConfig.Token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to create review video")
	}
	var res schema.GetReviewStoryUploadURLBodyResp
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.Wrap(err, "failed to decode response for create video review")
	}
	if !res.Success {
		return nil, errors.New(res.Error[0].Message)
	}
	return res.Payload, nil
}

func (ri *ReviewImpl) validateCreateReview(opts *schema.CreateReviewStoryOpts) error {
	// checking if review already exists
	filter := bson.M{
		"user_id":    opts.UserID,
		"catalog_id": opts.CatalogID,
	}
	count, err := ri.DB.Collection(model.ReviewColl).CountDocuments(context.TODO(), filter)
	if err != nil {
		return errors.Errorf("failed to validate review request")
	}
	if count != 0 {
		return errors.Errorf("review for this catalog already exists")
	}

	var catalog model.Catalog
	queryOpts := options.FindOne().SetProjection(bson.M{"_id": 1, "brand_id": 1})
	if err := ri.DB.Collection(model.CatalogColl).FindOne(context.TODO(), bson.M{"_id": opts.CatalogID}, queryOpts).Decode(&catalog); err != nil {
		return errors.Wrap(err, "failed to find catalog")
	}
	opts.BrandID = catalog.BrandID
	return nil
}

func (ri *ReviewImpl) CreateReviewStory(opts *schema.CreateReviewStoryOpts) (*schema.CreateReviewStoryResp, error) {
	var storyMediaResp *schema.GetReviewStoryUploadURLResp
	var err error
	if err = ri.validateCreateReview(opts); err != nil {
		return nil, err
	}

	// getting story upload url from cms service
	if storyMediaResp, err = ri.getReviewStoryUploadURL(opts); err != nil {
		return nil, errors.Wrap(err, "failed to create review upload url")
	}

	review := model.ReviewStory{
		CatalogID: opts.CatalogID,
		UserID:    opts.UserID,
		Rating:    opts.Rating,
		BrandID:   opts.BrandID,
		StoryID:   storyMediaResp.MediaID,
		CreatedAt: time.Now().UTC(),
	}

	// saving review in db
	res, err := ri.DB.Collection(model.ReviewColl).InsertOne(context.TODO(), review)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user review")
	}
	resp := &schema.CreateReviewStoryResp{
		ID:        res.InsertedID.(primitive.ObjectID),
		UploadURL: storyMediaResp.UploadURL,
	}
	return resp, nil
}

func (ri *ReviewImpl) ProcessReviewStory(storyID primitive.ObjectID) {
	ctx := context.TODO()
	filter := bson.M{
		"story_id": storyID,
	}
	update := bson.M{
		"$set": bson.M{
			"is_processed": true,
		},
	}

	var review model.ReviewStory
	if err := ri.DB.Collection(model.ReviewColl).FindOneAndUpdate(ctx, filter, update).Decode(&review); err != nil {
		ri.Logger.Err(err).Msgf("failed to update review with story_id: %s", storyID.Hex())
		return
	}

	// getting total review and avg review count
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"catalog_id":   review.CatalogID,
				"is_processed": true,
			},
		},
	}
	groupStage := bson.D{
		{
			Key: "$group",
			Value: bson.M{
				"_id": "$catalog_id",
				"avg_rating": bson.M{
					"$avg": "$rating",
				},
				"total_rating_count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	var res []interface{}
	cur, err := ri.DB.Collection(model.ReviewColl).Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		ri.Logger.Err(err).Interface("storyID", storyID.Hex()).Msgf("aggregation query failed for review with catalog_id: %s", review.CatalogID.Hex())
		return
	}

	if err := cur.All(ctx, &res); err != nil {
		ri.Logger.Err(err).Interface("storyID", storyID.Hex()).Msgf("failed to aggregation for review with catalog_id: %s", review.CatalogID.Hex())
		return
	}

	// Updating catalog avg_rating and total user review
	if len(res) == 1 {
		catalogRating := res[0].(primitive.D).Map()
		filter = bson.M{
			"_id": review.CatalogID,
		}
		update = bson.M{
			"$set": bson.M{
				"avg_rating":         Round(catalogRating["avg_rating"].(float64), 0.5, 1),
				"total_rating_count": catalogRating["total_rating_count"].(int32),
			},
		}
		res, err := ri.DB.Collection(model.CatalogColl).UpdateOne(ctx, filter, update)
		fmt.Println(res, err)
	}
}

func (ri *ReviewImpl) GetReviewInfo(id primitive.ObjectID) (*schema.ReviewStoryFullMessage, error) {
	ctx := context.TODO()
	filter := bson.M{
		"_id":          id,
		"is_processed": true,
	}

	var review schema.ReviewStoryFullMessage
	if err := ri.DB.Collection(model.ReviewColl).FindOne(ctx, filter).Decode(&review); err != nil {
		return nil, errors.Wrapf(err, "failed to get review with id: %s", id.Hex())
	}
	if &review == nil {
		return nil, errors.Errorf("review with id: %s not found", id.Hex())
	}
	storyInfo, err := ri.App.KeeperCatalog.GetReviewStoryByID(review.StoryID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get story info by id: %s", review.ID.Hex())
	}
	if len(storyInfo) == 0 {
		return nil, errors.Wrapf(err, "got 0 story info by id: %s", review.ID.Hex())
	}
	review.StoryInfo = &storyInfo[0]

	userInfo, err := ri.getUserInfo(review.UserID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user info by userid: %s", review.UserID.Hex())
	}
	review.UserInfo = userInfo

	return &review, nil
}

func (ri *ReviewImpl) getUserInfo(id primitive.ObjectID) (*schema.ReviewUserInfo, error) {
	var s schema.ReviewUserInfoResp
	url := ri.App.Config.HypdApiConfig.EntityApi + "/api/keeper/user/get"
	postBody, _ := json.Marshal(map[string]string{
		"id": id.Hex(),
	})
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request to get user")
	}
	req.Header.Add("Authorization", ri.App.Config.HypdApiConfig.Token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		ri.Logger.Err(err).Msgf("failed to send request to brand api %s", url)
		return nil, errors.Wrapf(err, "failed to send request to brand api %s", url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ri.Logger.Err(err).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrapf(err, "failed to read response from api %s", url)
	}

	if err := json.Unmarshal(body, &s); err != nil {
		ri.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}

	if !s.Success {
		ri.Logger.Err(errors.New("success false from entity")).Str("body", string(body)).Msg("got success false response from entity")
		return nil, errors.New("got success false response from entity")
	}
	return s.Payload, nil
}
