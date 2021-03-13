package app

import (
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ivs"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"syreclabs.com/go/faker"
)

func TestLiveImpl_validateCreateLiveStream(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	li := &LiveImpl{
		App:    app,
		DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
		Logger: app.Logger,
	}

	type args struct {
		s *model.Live
	}
	tests := []struct {
		name    string
		li      *LiveImpl
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "[Ok] No error",
			li:   li,
			args: args{
				s: &model.Live{
					FeaturedImage:  &model.IMG{SRC: "https://i.ytimg.com/vi/ndfhyMX-ja4/maxresdefault.jpg"},
					StreamEndImage: &model.IMG{SRC: "https://i2.wp.com/www.movieslantern.com/wp-content/uploads/2019/10/maxresdefault-170.jpg?fit=768%2C432&ssl=1"},
					ScheduledAt:    time.Now().UTC().Add(10 * time.Minute),
				},
			},
		},
		{
			name: "[Ok] Start time is less than now+5Min",
			li:   li,
			args: args{
				s: &model.Live{
					FeaturedImage:  &model.IMG{SRC: "https://i.ytimg.com/vi/ndfhyMX-ja4/maxresdefault.jpg"},
					StreamEndImage: &model.IMG{SRC: "https://i2.wp.com/www.movieslantern.com/wp-content/uploads/2019/10/maxresdefault-170.jpg?fit=768%2C432&ssl=1"},
					ScheduledAt:    time.Now().UTC().Add(2 * time.Minute),
				},
			},
			wantErr: true,
		},
		{
			name: "[Ok] Invalid Featured Image",
			li:   li,
			args: args{
				s: &model.Live{
					FeaturedImage:  &model.IMG{SRC: "https://i.ytim.com/vi/ndfhyMX-ja4/maxresdefault.jpg"},
					StreamEndImage: &model.IMG{SRC: "https://i2.wp.com/www.movieslantern.com/wp-content/uploads/2019/10/maxresdefault-170.jpg?fit=768%2C432&ssl=1"},
					ScheduledAt:    time.Now().UTC().Add(2 * time.Minute),
				},
			},
			wantErr: true,
		},
		{
			name: "[Ok] Invalid Stream End Image",
			li:   li,
			args: args{
				s: &model.Live{
					FeaturedImage:  &model.IMG{SRC: "https://i.ytimg.com/vi/ndfhyMX-ja4/maxresdefault.jpg"},
					StreamEndImage: &model.IMG{SRC: "https://i2.com/www.movieslantern.com/wp-content/uploads/2019/10/maxresdefault-170.jpg?fit=768%2C432&ssl=1"},
					ScheduledAt:    time.Now().UTC().Add(2 * time.Minute),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &LiveImpl{}
			err := li.validateCreateLiveStream(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.validateCreateLiveStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLiveImpl_CreateLiveStream(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
		IVS    IVS
	}
	type args struct {
		ivsResp *ivs.CreateChannelOutput
		opts    *schema.CreateLiveStreamOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.CreateLiveStreamResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockIVS)
		validate   func(*testing.T, *TC, *schema.CreateLiveStreamResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateLiveStreamOpts()
				tt.args.opts = opts
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				name := tt.args.opts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
				tt.args.ivsResp = &resp
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateLiveStreamResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.args.opts.Name, resp.Name)
				assert.Equal(t, tt.args.opts.InfluencerIDs, resp.InfluencerIDs)
				assert.Equal(t, tt.args.opts.CatalogIDs, resp.CatalogIDs)
				assert.Equal(t, tt.args.opts.ScheduledAt, resp.ScheduledAt)
				assert.Equal(t, tt.args.opts.FeaturedImage.SRC, resp.FeaturedImage.SRC)
				assert.Equal(t, 704, resp.FeaturedImage.Height)
				assert.Equal(t, 1200, resp.FeaturedImage.Width)
				assert.Equal(t, tt.args.opts.StreamEndImage.SRC, resp.StreamEndImage.SRC)
				assert.Equal(t, 432, resp.StreamEndImage.Height)
				assert.Equal(t, 768, resp.StreamEndImage.Width)
				assert.NotEmpty(t, resp.Slug)
				assert.WithinDuration(t, time.Now(), resp.CreatedAt, 100*time.Second)

				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, *tt.args.ivsResp.Channel.Arn, doc.IVS.Channel.ARN)
				assert.Equal(t, *tt.args.ivsResp.Channel.LatencyMode, doc.IVS.Channel.LatencyMode)
				assert.Equal(t, *tt.args.ivsResp.Channel.Name, doc.IVS.Channel.Name)
				assert.Equal(t, *tt.args.ivsResp.Channel.Type, doc.IVS.Channel.Type)
				assert.Equal(t, *tt.args.ivsResp.Channel.Authorized, doc.IVS.Channel.PlaybackAuthorization)
				assert.Equal(t, *tt.args.ivsResp.Channel.IngestEndpoint, doc.IVS.Ingestion.IngestURL)
				assert.Equal(t, *tt.args.ivsResp.StreamKey.Value, doc.IVS.Ingestion.StreamKey)
				assert.Equal(t, *tt.args.ivsResp.Channel.PlaybackUrl, doc.IVS.Playback.PlaybackURL)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIVS := mock.NewMockIVS(ctrl)
			li := &LiveImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
				IVS:    mockIVS,
			}
			tt.fields.App.Live = li
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockIVS)

			got, err := li.CreateLiveStream(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.CreateLiveStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestLiveImpl_DiscardLiveStream(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
		IVS    IVS
	}
	type args struct {
		createOpts *schema.CreateLiveStreamOpts
		id         primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockIVS)
		validate   func(*testing.T, *TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				resp, _ := tt.fields.App.Live.CreateLiveStream(tt.args.createOpts)
				tt.args.id = resp.ID
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				tt.args.createOpts = schema.GetRandomCreateLiveStreamOpts()
				name := tt.args.createOpts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
			},
			validate: func(t *testing.T, tt *TC) {
				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&doc)
				assert.Nil(t, err)
				assert.Len(t, doc.StatusHistory, 1)
				assert.Equal(t, &doc.StatusHistory[0], doc.Status)
				assert.Equal(t, doc.Status.Name, model.DiscardStatus)
				assert.WithinDuration(t, time.Now().UTC(), doc.Status.CreatedAt, 200*time.Millisecond)
			},
		},
		{
			name: "[Error] Id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				_, _ = tt.fields.App.Live.CreateLiveStream(tt.args.createOpts)
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("failed to find live stream with id:%s", tt.args.id.Hex())
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				tt.args.createOpts = schema.GetRandomCreateLiveStreamOpts()
				name := tt.args.createOpts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
			},
			wantErr: true,
			validate: func(t *testing.T, tt *TC) {
				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&doc)
				assert.NotNil(t, err)
				err = tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{}).Decode(&doc)
				assert.Nil(t, err)
				assert.Len(t, doc.StatusHistory, 1)
				assert.Equal(t, &doc.StatusHistory[0], doc.Status)
				assert.Equal(t, doc.Status.Name, model.DiscardStatus)
				assert.WithinDuration(t, time.Now().UTC(), doc.Status.CreatedAt, 300*time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIVS := mock.NewMockIVS(ctrl)
			li := &LiveImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
				IVS:    mockIVS,
			}

			tt.fields.App.Live = li
			tt.buildStubs(&tt, mockIVS)
			tt.prepare(&tt)
			err := li.DiscardLiveStream(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.DiscardLiveStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.validate(t, &tt)
		})
	}
}

func TestLiveImpl_EndLiveStream(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
		IVS    IVS
	}
	type args struct {
		createOpts *schema.CreateLiveStreamOpts
		id         primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockIVS)
		validate   func(*testing.T, *TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				resp, _ := tt.fields.App.Live.CreateLiveStream(tt.args.createOpts)
				tt.args.id = resp.ID
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				tt.args.createOpts = schema.GetRandomCreateLiveStreamOpts()
				name := tt.args.createOpts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
				resp1 := ivs.StopStreamOutput{}
				mc.EXPECT().StopStream(arn).Times(1).Return(&resp1, nil)
			},
			validate: func(t *testing.T, tt *TC) {
				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&doc)
				assert.Nil(t, err)
				assert.Len(t, doc.StatusHistory, 1)
				assert.Equal(t, &doc.StatusHistory[0], doc.Status)
				assert.Equal(t, doc.Status.Name, model.EndStatus)
				assert.WithinDuration(t, time.Now().UTC(), doc.Status.CreatedAt, 300*time.Millisecond)
			},
		},
		{
			name: "[Error] Id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				_, _ = tt.fields.App.Live.CreateLiveStream(tt.args.createOpts)
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("failed to find live stream with id:%s", tt.args.id.Hex())
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				tt.args.createOpts = schema.GetRandomCreateLiveStreamOpts()
				name := tt.args.createOpts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
			},
			wantErr: true,
			validate: func(t *testing.T, tt *TC) {
				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&doc)
				assert.NotNil(t, err)
				err = tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{}).Decode(&doc)
				assert.Nil(t, err)
				assert.Len(t, doc.StatusHistory, 1)
				assert.Equal(t, &doc.StatusHistory[0], doc.Status)
				assert.Equal(t, doc.Status.Name, model.EndStatus)
				assert.WithinDuration(t, time.Now().UTC(), doc.Status.CreatedAt, 300*time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIVS := mock.NewMockIVS(ctrl)
			li := &LiveImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
				IVS:    mockIVS,
			}
			tt.fields.App.Live = li
			tt.buildStubs(&tt, mockIVS)
			tt.prepare(&tt)
			err := li.EndLiveStream(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.DiscardLiveStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.validate(t, &tt)
		})
	}
}

func TestLiveImpl_StartLiveStream(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
		IVS    IVS
	}
	type args struct {
		createOpts *schema.CreateLiveStreamOpts
		ivsResp    *ivs.CreateChannelOutput
		id         primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       string
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockIVS)
		validate   func(*testing.T, *TC, string)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				resp, _ := tt.fields.App.Live.CreateLiveStream(tt.args.createOpts)
				tt.args.id = resp.ID
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				tt.args.createOpts = schema.GetRandomCreateLiveStreamOpts()
				name := tt.args.createOpts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
				tt.args.ivsResp = &resp
			},
			validate: func(t *testing.T, tt *TC, resp string) {
				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&doc)
				assert.Nil(t, err)
				assert.Len(t, doc.StatusHistory, 1)
				assert.Equal(t, &doc.StatusHistory[0], doc.Status)
				assert.Equal(t, doc.Status.Name, model.ActiveStatus)
				assert.WithinDuration(t, time.Now().UTC(), doc.Status.CreatedAt, 300*time.Millisecond)
				assert.Equal(t, *tt.args.ivsResp.StreamKey.Value, doc.IVS.Ingestion.StreamKey)
				assert.Equal(t, doc.IVS.Ingestion.StreamKey, resp)
			},
		},
		{
			name: "[Error] Id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("failed to find live stream with id:%s", tt.args.id.Hex())
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {

			},
			wantErr: true,
			validate: func(t *testing.T, tt *TC, resp string) {
				assert.Equal(t, "", resp)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIVS := mock.NewMockIVS(ctrl)
			li := &LiveImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
				IVS:    mockIVS,
			}

			tt.fields.App.Live = li
			tt.buildStubs(&tt, mockIVS)
			tt.prepare(&tt)
			got, err := li.StartLiveStream(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.DiscardLiveStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Empty(t, got)
			}
			tt.validate(t, &tt, got)
		})
	}
}

func TestLiveImpl_GetLiveStreamByID(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
		IVS    IVS
	}
	type args struct {
		createOpts *schema.CreateLiveStreamOpts
		ivsResp    *ivs.CreateChannelOutput
		id         primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.GetLiveStreamResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockIVS)
		validate   func(*testing.T, *TC, *schema.GetLiveStreamResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				resp, _ := tt.fields.App.Live.CreateLiveStream(tt.args.createOpts)
				tt.args.id = resp.ID
				tt.want = &schema.GetLiveStreamResp{
					ID:             resp.ID,
					Name:           resp.Name,
					Slug:           resp.Slug,
					FeaturedImage:  resp.FeaturedImage,
					StreamEndImage: resp.StreamEndImage,
					CatalogIDs:     resp.CatalogIDs,
					InfluencerIDs:  resp.InfluencerIDs,
					IVS: &model.IVS{
						Channel: &model.IVSChannel{
							ARN:                   *tt.args.ivsResp.Channel.Arn,
							Name:                  *tt.args.ivsResp.Channel.Name,
							Type:                  *tt.args.ivsResp.Channel.Type,
							LatencyMode:           *tt.args.ivsResp.Channel.LatencyMode,
							PlaybackAuthorization: *tt.args.ivsResp.Channel.Authorized,
						},
						Ingestion: &model.IVSIngest{
							IngestURL: *tt.args.ivsResp.Channel.IngestEndpoint,
							StreamKey: *tt.args.ivsResp.StreamKey.Value,
						},
						Playback: &model.IVSPlayback{
							PlaybackURL: *tt.args.ivsResp.Channel.PlaybackUrl,
						},
					},
					ScheduledAt: resp.ScheduledAt,
					CreatedAt:   resp.CreatedAt,
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {
				tt.args.createOpts = schema.GetRandomCreateLiveStreamOpts()
				name := tt.args.createOpts.Name
				arn := faker.Letterify("???-????-?????")
				authorized := false
				channelType := faker.RandomChoice([]string{"STANDARD"})
				ingestEndpoint := faker.Letterify("rtmp://??.???.??.??:8000")
				latencyMode := faker.RandomChoice([]string{"LOW"})
				playbackURL := faker.Letterify("https://?????.??????.com/????.m3u8")
				streamKey := faker.RandomString(20)

				resp := ivs.CreateChannelOutput{
					Channel: &ivs.Channel{
						Arn:            &arn,
						Authorized:     &authorized,
						IngestEndpoint: &ingestEndpoint,
						LatencyMode:    &latencyMode,
						Name:           &name,
						Type:           &channelType,
						PlaybackUrl:    &playbackURL,
					},
					StreamKey: &ivs.StreamKey{
						Value: &streamKey,
					},
				}
				mc.EXPECT().CreateChannel(name).Times(1).Return(&resp, nil)
				tt.args.ivsResp = &resp
			},
			validate: func(t *testing.T, tt *TC, resp *schema.GetLiveStreamResp) {
				assert.WithinDuration(t, tt.want.CreatedAt, resp.CreatedAt, 100*time.Millisecond)
				assert.WithinDuration(t, tt.want.ScheduledAt, resp.ScheduledAt, 100*time.Millisecond)
				tt.want.CreatedAt = time.Time{}
				tt.want.ScheduledAt = time.Time{}
				resp.CreatedAt = time.Time{}
				resp.ScheduledAt = time.Time{}
				assert.Equal(t, tt.want, resp)
				var doc model.Live
				err := tt.fields.DB.Collection(model.LiveColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&doc)
				assert.Nil(t, err)
			},
		},
		{
			name: "[Error] Id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("failed to find live stream with id:%s", tt.args.id.Hex())
			},
			buildStubs: func(tt *TC, mc *mock.MockIVS) {

			},
			wantErr: true,
			validate: func(t *testing.T, tt *TC, resp *schema.GetLiveStreamResp) {
				assert.Nil(t, resp)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIVS := mock.NewMockIVS(ctrl)
			li := &LiveImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
				IVS:    mockIVS,
			}

			tt.fields.App.Live = li
			tt.buildStubs(&tt, mockIVS)
			tt.prepare(&tt)
			got, err := li.GetLiveStreamByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.GetLiveStreamByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Empty(t, got)
			}
			tt.validate(t, &tt, got)
		})
	}
}

func TestLiveImpl_GetLiveStreams(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	ctx := context.Background()
	var liveDoc []model.Live
	DB := app.MongoDB.Client.Database(app.Config.LiveConfig.DBName)
	// Document with active status and scheduled for next 24hour
	for i := 0; i < 10; i++ {
		s := model.GetRandomLive()
		s.ScheduledAt = time.Now().Add(24 * time.Hour)
		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		_ = DB.Collection(model.LiveColl).FindOneAndUpdate(ctx, bson.M{"_id": s.ID}, bson.M{"$set": s}, opts).Decode(&s)
		liveDoc = append(liveDoc, *s)
	}

	// Document with discard status and created 48hours ago
	for i := 0; i < 10; i++ {
		s := model.GetRandomLive()
		s.CreatedAt = time.Now().Add(-48 * time.Hour)
		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		_ = DB.Collection(model.LiveColl).FindOneAndUpdate(ctx, bson.M{"_id": s.ID}, bson.M{"$set": s}, opts).Decode(&s)
		liveDoc = append(liveDoc, *s)
	}

	// Document with end status and created 48hours ago
	for i := 0; i < 10; i++ {
		s := model.GetRandomLive()
		s.Status = &model.StreamStatus{Name: model.EndStatus, CreatedAt: time.Now().Add(-23 * time.Hour)}
		s.StatusHistory = []model.StreamStatus{*s.Status}
		s.CreatedAt = time.Now().Add(-48 * time.Hour)
		s.ScheduledAt = time.Now().Add(-24 * time.Hour)
		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		_ = DB.Collection(model.LiveColl).FindOneAndUpdate(ctx, bson.M{"_id": s.ID}, bson.M{"$set": s}, opts).Decode(&s)
		liveDoc = append(liveDoc, *s)
	}

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
		IVS    IVS
	}
	type args struct {
		filterOpts *schema.GetLiveStreamsFilter
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       []schema.GetLiveStreamResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockIVS)
		validate   func(*testing.T, *TC, []schema.GetLiveStreamResp)
	}

	tests := []TC{
		{
			name: "[Ok] No Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[:10] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] Page 1",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Page: 1,
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[10:20] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] Page 3",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Page: 3,
				}
				tt.args.filterOpts = &opts
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] Active Status Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Status: []string{model.ActiveStatus},
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[:10] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] Active & End Status Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Status: []string{model.ActiveStatus, model.EndStatus},
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[:10] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] Active & End Status Filter Page 2",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Page:   1,
					Status: []string{model.ActiveStatus, model.EndStatus},
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[10:20] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] CreatedAtFrom Past Date",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Page:          0,
					CreatedAtFrom: time.Now().Add(-49 * time.Hour),
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[0:10] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] CreatedAtFrom Past Date Page 1",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Page:          1,
					CreatedAtFrom: time.Now().Add(-49 * time.Hour),
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[10:20] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Ok] ScheduledAt Next 24 hours",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.LiveConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetLiveStreamsFilter{
					Page:          0,
					ScheduledAtTo: time.Now(),
				}
				tt.args.filterOpts = &opts
				for _, doc := range liveDoc[20:30] {
					tt.want = append(tt.want, schema.GetLiveStreamResp{
						ID:             doc.ID,
						Name:           doc.Name,
						Slug:           doc.Slug,
						InfluencerIDs:  doc.InfluencerIDs,
						CatalogIDs:     doc.CatalogIDs,
						ScheduledAt:    doc.ScheduledAt,
						FeaturedImage:  doc.FeaturedImage,
						StreamEndImage: doc.StreamEndImage,
						IVS:            doc.IVS,
						CreatedAt:      doc.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, resp []schema.GetLiveStreamResp) {
				assert.Len(t, resp, len(tt.want))
				assert.Equal(t, tt.want, resp)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			li := &LiveImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Live = li
			tt.prepare(&tt)
			got, err := li.GetLiveStreams(tt.args.filterOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiveImpl.GetLiveStreams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}
