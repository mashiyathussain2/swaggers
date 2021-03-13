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

// Live contains methods to implement live functionality in the app
type Live interface {
	CreateLiveStream(*schema.CreateLiveStreamOpts) (*schema.CreateLiveStreamResp, error)
	StartLiveStream(primitive.ObjectID) (string, error)
	DiscardLiveStream(primitive.ObjectID) error
	EndLiveStream(primitive.ObjectID) error
}

// LiveImpl implemethods Live interface methods
type LiveImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
	IVS    IVS
}

// LiveOpts contains args required to create a new instance of LiveImpl
type LiveOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
	IVS    IVS
}

// InitLive returns new instance of live implementation
func InitLive(opts *LiveOpts) Live {
	li := LiveImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
		IVS:    opts.IVS,
	}
	return &li
}

func (li *LiveImpl) validateCreateLiveStream(s *model.Live) error {

	// Verifying featured_image
	if err := s.FeaturedImage.LoadFromURL(); err != nil {
		return errors.Wrap(err, "failed to load featured_image")
	}
	// Verifying stream_end_imag
	if err := s.StreamEndImage.LoadFromURL(); err != nil {
		return errors.Wrap(err, "failed to load stream_end_image")
	}

	// Validating Stream Start Time
	if s.ScheduledAt.Before(time.Now().UTC().Add(5 * time.Minute)) {
		return errors.Errorf("invalid stream_scheduled_at: should be after %s", time.Now().UTC().Add(5*time.Minute))
	}
	return nil
}

// CreateLiveStream creates a new upcoming live stream in DB
func (li *LiveImpl) CreateLiveStream(opts *schema.CreateLiveStreamOpts) (*schema.CreateLiveStreamResp, error) {
	s := model.Live{
		Name:           opts.Name,
		Slug:           UniqueSlug(opts.Name),
		FeaturedImage:  &model.IMG{SRC: opts.FeaturedImage.SRC},
		StreamEndImage: &model.IMG{SRC: opts.StreamEndImage.SRC},
		ScheduledAt:    opts.ScheduledAt,
		InfluencerIDs:  opts.InfluencerIDs,
		CatalogIDs:     opts.CatalogIDs,
		CreatedAt:      time.Now().UTC(),
	}
	if err := li.validateCreateLiveStream(&s); err != nil {
		return nil, err
	}
	resp, err := li.IVS.CreateChannel(opts.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate live stream")
	}
	ivs := model.IVS{
		Channel: &model.IVSChannel{
			ARN:                   *resp.Channel.Arn,
			Name:                  *resp.Channel.Name,
			Type:                  *resp.Channel.Type,
			LatencyMode:           *resp.Channel.LatencyMode,
			PlaybackAuthorization: *resp.Channel.Authorized,
		},
		Playback: &model.IVSPlayback{
			PlaybackURL: *resp.Channel.PlaybackUrl,
		},
		Ingestion: &model.IVSIngest{
			IngestURL: *resp.Channel.IngestEndpoint,
			StreamKey: *resp.StreamKey.Value,
		},
	}
	s.IVS = &ivs

	res, err := li.DB.Collection(model.LiveColl).InsertOne(context.TODO(), s)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create live stream")
	}
	s.ID = res.InsertedID.(primitive.ObjectID)
	return &schema.CreateLiveStreamResp{
		ID:             s.ID,
		Name:           s.Name,
		Slug:           s.Slug,
		InfluencerIDs:  s.InfluencerIDs,
		CatalogIDs:     s.CatalogIDs,
		ScheduledAt:    s.ScheduledAt,
		FeaturedImage:  s.FeaturedImage,
		StreamEndImage: s.StreamEndImage,
		CreatedAt:      s.CreatedAt,
	}, nil
}

// StartLiveStream starts a live stream
func (li *LiveImpl) StartLiveStream(id primitive.ObjectID) (string, error) {
	// Updating status in mongodb
	st := model.StreamStatus{
		Name:      model.ActiveStatus,
		CreatedAt: time.Now().UTC(),
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"stream_status": st,
		},
		"$push": bson.M{
			"status_history": st,
		},
	}
	var stream model.Live
	if err := li.DB.Collection(model.LiveColl).FindOneAndUpdate(context.TODO(), filter, update).Decode(&stream); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return "", errors.Errorf("failed to find live stream with id:%s", id.Hex())
		}
		return "", errors.Errorf("failed to start live stream with id:%s", id.Hex())
	}
	return stream.IVS.Ingestion.StreamKey, nil
}

// DiscardLiveStream marks an upcoming livestream discarded
func (li *LiveImpl) DiscardLiveStream(id primitive.ObjectID) error {
	st := model.StreamStatus{
		Name:      model.DiscardStatus,
		CreatedAt: time.Now().UTC(),
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"stream_status": st,
		},
		"$push": bson.M{
			"status_history": st,
		},
	}
	res, err := li.DB.Collection(model.LiveColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrapf(err, "failed to discard live stream id:%s", id.Hex())
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("failed to find live stream with id:%s", id.Hex())
	}
	if res.ModifiedCount == 0 {
		return errors.Errorf("failed to update live stream with id:%s", id.Hex())
	}
	return nil
}

// EndLiveStream marks an active livestream ended
func (li *LiveImpl) EndLiveStream(id primitive.ObjectID) error {
	// Marking live stream status as ended
	st := model.StreamStatus{
		Name:      model.EndStatus,
		CreatedAt: time.Now().UTC(),
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"stream_status": st,
		},
		"$push": bson.M{
			"status_history": st,
		},
	}

	var stream model.Live
	if err := li.DB.Collection(model.LiveColl).FindOneAndUpdate(context.TODO(), filter, update).Decode(&stream); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return errors.Errorf("failed to find live stream with id:%s", id.Hex())
		}
		return errors.Errorf("failed to end live stream with id:%s", id.Hex())
	}

	// Stopping aws stream channel
	_, err := li.IVS.StopStream(stream.IVS.Channel.ARN)
	if err != nil {
		li.Logger.Err(err).Msgf("failed to stop aws channel stream with id:%s", id.Hex())
		return errors.Wrapf(err, "failed to stop aws channel stream with id:%s", id.Hex())
	}

	return nil
}

// GetLiveStreamByID returns specific live stream info matched by id
func (li *LiveImpl) GetLiveStreamByID(id primitive.ObjectID) (*schema.GetLiveStreamResp, error) {
	var resp schema.GetLiveStreamResp
	filter := bson.M{"_id": id}
	if err := li.DB.Collection(model.LiveColl).FindOne(context.TODO(), filter).Decode(&resp); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Wrapf(err, "live stream by id:%s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "failed to find live stream by id:%s", id.Hex())
	}
	return &resp, nil
}

// GetLiveStreams returns live streams
func (li *LiveImpl) GetLiveStreams(filterOpts *schema.GetLiveStreamsFilter) ([]schema.GetLiveStreamResp, error) {
	var filter bson.D
	if len(filterOpts.Status) > 0 {
		filter = append(filter, bson.E{Key: "status.name", Value: bson.M{"$in": filterOpts.Status}})
	}
	if !filterOpts.CreatedAtFrom.IsZero() {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$gte": filterOpts.CreatedAtFrom}})
	}
	if !filterOpts.CreatedAtTo.IsZero() {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$lte": filterOpts.CreatedAtTo}})
	}
	if !filterOpts.ScheduledAtFrom.IsZero() {
		filter = append(filter, bson.E{Key: "scheduled_at", Value: bson.M{"$gte": filterOpts.ScheduledAtFrom}})
	}
	if !filterOpts.ScheduledAtTo.IsZero() {
		filter = append(filter, bson.E{Key: "scheduled_at", Value: bson.M{"$lte": filterOpts.ScheduledAtTo}})
	}

	if filter == nil {
		filter = bson.D{}
	}
	var resp []schema.GetLiveStreamResp
	ctx := context.TODO()
	queryOpts := options.Find().SetSkip(int64(filterOpts.Page * 10)).SetLimit(10)
	cur, err := li.DB.Collection(model.LiveColl).Find(ctx, filter, queryOpts)
	if err != nil {
		return nil, errors.Wrap(err, "query failed to get live streams")
	}
	if err := cur.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to get live streams")
	}
	return resp, nil
}
