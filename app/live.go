package app

import (
	"context"
	"encoding/json"
	"go-app/model"
	"go-app/schema"
	"go-app/server/kafka"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ivs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	segKafka "github.com/segmentio/kafka-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Live contains methods to implement live functionality in the app
type Live interface {
	CreateLiveStream(*schema.CreateLiveStreamOpts) (*schema.CreateLiveStreamResp, error)
	StartLiveStream(primitive.ObjectID) (*schema.StartLiveStreamResp, error)
	DiscardLiveStream(primitive.ObjectID) error
	EndLiveStream(primitive.ObjectID) error
	JoinLiveStream(primitive.ObjectID) (*schema.JoinLiveStreamResp, error)

	PushComment(*schema.CreateLiveCommentOpts)
	PushOrder(opts *schema.PushNewOrderOpts)
	PushCatalog(opts *schema.PushCatalogOpts)
	PushJoin(opts *schema.PushJoinOpts)

	GetLiveStreamByID(primitive.ObjectID) (*schema.GetLiveStreamResp, error)
	GetLiveStreams(*schema.GetLiveStreamsFilter) ([]schema.GetLiveStreamResp, error)
	ConsumeComment(m kafka.Message)
	CreateLiveComment(*schema.CreateLiveCommentOpts)

	GetAppLiveStreams(*schema.GetAppLiveStreamsFilter) ([]schema.GetAppLiveStreamResp, error)
	GetAppLiveStreamByID(primitive.ObjectID) (*schema.GetAppLiveStreamResp, error)
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
	resp, err := li.IVS.CreateChannel(s.Slug)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate live stream")
	}
	ivs := model.IVS{
		Channel: &model.IVSChannel{
			ARN:                   *resp.Channel.Arn,
			Name:                  s.Slug,
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
func (li *LiveImpl) StartLiveStream(id primitive.ObjectID) (*schema.StartLiveStreamResp, error) {
	// Updating status in mongodb
	st := model.StreamStatus{
		Name:      model.ActiveStatus,
		CreatedAt: time.Now().UTC(),
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status": st,
		},
		"$push": bson.M{
			"status_history": st,
		},
	}
	var stream model.Live
	if err := li.DB.Collection(model.LiveColl).FindOneAndUpdate(context.TODO(), filter, update).Decode(&stream); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("failed to find live stream with id:%s", id.Hex())
		}
		return nil, errors.Errorf("failed to start live stream with id:%s", id.Hex())
	}
	return &schema.StartLiveStreamResp{
		StreamKey: stream.IVS.Ingestion.StreamKey,
		IngestURL: stream.IVS.Ingestion.IngestURL,
	}, nil
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
			"status": st,
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
			"status": st,
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

// JoinLiveStream returns a playback url to user can stream the live feed
func (li *LiveImpl) JoinLiveStream(id primitive.ObjectID) (*schema.JoinLiveStreamResp, error) {
	filter := bson.M{"_id": id}

	var stream model.Live
	opts := options.FindOne().SetProjection(bson.M{"ivs.playback": 1, "status": 1, "ivs.channel": 1})
	if err := li.DB.Collection(model.LiveColl).FindOne(context.TODO(), filter, opts).Decode(&stream); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Errorf("failed to find live stream with id:%s", id.Hex())
		}
		return nil, errors.Errorf("failed to end live stream with id:%s", id.Hex())
	}

	if stream.Status == nil {
		return nil, errors.New("stream is not active")
	}
	if stream.Status.Name != model.ActiveStatus {
		return nil, errors.New("stream is not active")
	}

	resp := schema.JoinLiveStreamResp{
		PlaybackURL: stream.IVS.Playback.PlaybackURL,
		ARN:         stream.IVS.Channel.ARN,
	}
	return &resp, nil
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
		filter = append(filter, bson.E{Key: "status.name", Value: bson.M{"$ne": model.EndStatus}})
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

// PushComment pushes the comment object into kafka topic and aws ivs meta data api
func (li *LiveImpl) PushComment(opts *schema.CreateLiveCommentOpts) {
	// Pushing comment to kafka topic
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		opts.Type = "comment"
		opts.CreatedAt = time.Now().UTC()
		bytes, err := json.Marshal(opts)
		if err == nil {
			li.App.LiveCommentProducer.Publish(segKafka.Message{
				Key:   nil,
				Value: bytes,
			})
			return
		}
		li.Logger.Err(err).Interface("opts", opts).Msg("failed to decode opts to bytes")
	}()

	// Pushing comment to IVS meta-data
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := schema.IVSMetaData{
			Type: "comment",
			Data: schema.CreateIVSCommentMetaData{
				ID:           opts.UserID,
				Name:         opts.Name,
				ProfileImage: opts.ProfileImage,
				Description:  opts.Description,
			},
		}
		bytes, err := json.Marshal(s)
		metaData := string(bytes)
		if err == nil {
			params := ivs.PutMetadataInput{
				ChannelArn: &opts.ARN,
				Metadata:   &metaData,
			}
			_, err := li.IVS.PutMetadata(&params)
			if err != nil {
				li.Logger.Err(err).RawJSON("metadata", bytes).Msg("failed to push comment in ivs metadata")
			}
			return
		}
		li.Logger.Err(err).Interface("metadata_struct", s).Msg("failed to convert struct to bytes")
	}()

	wg.Wait()
}

func (li *LiveImpl) PushJoin(opts *schema.PushJoinOpts) {
	s := schema.IVSMetaData{
		Type: "join",
		Data: schema.CreateIVSNewJoinMetaData{
			ID:   opts.ID,
			Name: opts.Name,
		},
	}
	bytes, err := json.Marshal(s)
	metaData := string(bytes)
	if err == nil {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			li.pushViewerCount(&schema.PushViewerCount{ARN: opts.ARN})
		}()
		params := ivs.PutMetadataInput{
			ChannelArn: &opts.ARN,
			Metadata:   &metaData,
		}
		_, err := li.IVS.PutMetadata(&params)
		if err != nil {
			li.Logger.Err(err).RawJSON("metadata", bytes).Msg("failed to push join in ivs metadata")
		}
		wg.Wait()
		return
	}
	li.Logger.Err(err).Interface("metadata_struct", s).Msg("failed to convert struct to bytes")
}

func (li *LiveImpl) pushViewerCount(opts *schema.PushViewerCount) {
	out, err := li.IVS.GetStream(opts.ARN)
	if err != nil {
		li.Logger.Err(err).Msgf("failed to get stream by arn: %s", opts.ARN)
		return
	}
	if bytes, err := json.Marshal(schema.IVSMetaData{Type: "live", Data: schema.ViewerCountMetadata{Count: uint(*out.Stream.ViewerCount)}}); err == nil {
		metadata := string(bytes)
		params := ivs.PutMetadataInput{
			ChannelArn: out.Stream.ChannelArn,
			Metadata:   &metadata,
		}
		li.IVS.PutMetadata(&params)
	}
}

func (li *LiveImpl) ConsumeComment(m kafka.Message) {
	message := m.(segKafka.Message).Value
	var opts schema.CreateLiveCommentOpts
	if err := json.Unmarshal(message, &opts); err != nil {
		li.Logger.Err(err).RawJSON("message", message).Msg("failed to read live comment data")
		return
	}
	doc := model.Comment{
		ResourceType: model.LiveColl,
		ResourceID:   opts.LiveID,
		Description:  opts.Description,
		UserID:       opts.UserID,
		CreatedAt:    opts.CreatedAt,
	}
	if _, err := li.DB.Collection(model.CommentColl).InsertOne(context.TODO(), &doc); err != nil {
		li.Logger.Err(err).Interface("opts", opts).Msg("failed to push comment")
		return
	}
	li.App.LiveComments.Commit(context.Background(), m)
}

func (li *LiveImpl) CreateLiveComment(opts *schema.CreateLiveCommentOpts) {
	s := schema.CreateCommentOpts{
		ResourceType: model.LiveType,
		ResourceID:   opts.LiveID,
		Description:  opts.Description,
		UserID:       opts.UserID,
		CreatedAt:    opts.CreatedAt,
	}

	if _, err := li.App.Content.CreateComment(&s); err != nil {
		li.Logger.Err(err).Interface("opts", opts).Msg("failed to create live comment")
	}
}

func (li *LiveImpl) PushCatalog(opts *schema.PushCatalogOpts) {
	s := schema.IVSMetaData{
		Type: "catalog",
		Data: schema.CreateIVSCatalogMetaData{
			ID: opts.ID,
		},
	}
	bytes, err := json.Marshal(s)
	metaData := string(bytes)
	if err == nil {
		params := ivs.PutMetadataInput{
			ChannelArn: &opts.ARN,
			Metadata:   &metaData,
		}
		_, err := li.IVS.PutMetadata(&params)
		if err != nil {
			li.Logger.Err(err).RawJSON("metadata", bytes).Msg("failed to push catalog in ivs metadata")
		}
		return
	}
	li.Logger.Err(err).Interface("metadata_struct", s).Msg("failed to convert struct to bytes")
}

func (li *LiveImpl) PushOrder(opts *schema.PushNewOrderOpts) {
	s := schema.IVSMetaData{
		Type: "order",
		Data: schema.CreateIVSOrderMetaData{
			Name:         opts.Name,
			ProfileImage: opts.ProfileImage,
		},
	}
	bytes, err := json.Marshal(s)
	metaData := string(bytes)
	if err == nil {
		params := ivs.PutMetadataInput{
			ChannelArn: &opts.ARN,
			Metadata:   &metaData,
		}
		_, err := li.IVS.PutMetadata(&params)
		if err != nil {
			li.Logger.Err(err).RawJSON("metadata", bytes).Msg("failed to push order in ivs metadata")
		}
		return
	}
	li.Logger.Err(err).Interface("metadata_struct", s).Msg("failed to convert struct to bytes")
}

func (li *LiveImpl) GetAppLiveStreamByID(id primitive.ObjectID) (*schema.GetAppLiveStreamResp, error) {
	var resp schema.GetAppLiveStreamResp
	filter := bson.M{"_id": id}
	if err := li.DB.Collection(model.LiveColl).FindOne(context.TODO(), filter).Decode(&resp); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, errors.Wrapf(err, "live stream by id:%s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "failed to find live stream by id:%s", id.Hex())
	}
	return &resp, nil
}

func (li *LiveImpl) GetAppLiveStreams(filterOpts *schema.GetAppLiveStreamsFilter) ([]schema.GetAppLiveStreamResp, error) {
	ctx := context.TODO()
	filter := bson.M{

		"$or": bson.A{
			bson.M{
				"status.name": model.ActiveStatus,
			},
			bson.M{
				"$and": bson.A{
					bson.M{
						"scheduled_at": bson.M{
							"$gte": time.Now().UTC(),
						},
					},
					bson.M{
						"status.name": bson.M{
							"$nin": bson.A{model.EndStatus, model.ActiveStatus},
						},
					},
				},
			},
		},
	}
	var resp []schema.GetAppLiveStreamResp
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
