package app

import (
	"context"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestContentImpl_generateS3UploadToken(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}

	type args struct {
		videoID   string
		videoType string
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       string
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockS3)
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
			args: args{
				videoID:   primitive.NewObjectIDFromTimestamp(time.Now()).Hex(),
				videoType: ".mp4",
			},
			buildStubs: func(tt *TC, mc *mock.MockS3) {
				resp := "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd"
				mc.EXPECT().GetPutObjectRequestURL(gomock.Any()).Times(1).Return(resp, nil)
				tt.want = resp
			},
		},
		{
			name: "[Error] Failed to generate token",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				videoID:   primitive.NewObjectIDFromTimestamp(time.Now()).Hex(),
				videoType: ".mp4",
			},
			wantErr: true,
			buildStubs: func(tt *TC, mc *mock.MockS3) {
				mc.EXPECT().GetPutObjectRequestURL(gomock.Any()).Times(1).Return("", errors.New("failed to generate token"))
				tt.err = errors.Wrap(errors.New("failed to generate token"), "failed to generate upload token")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &MediaImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			mockS3 := mock.NewMockS3(ctrl)
			ci.App.S3 = mockS3

			tt.fields.App.Media = ci
			tt.buildStubs(&tt, mockS3)

			got, err := ci.generateS3UploadToken(tt.args.videoID, tt.args.videoType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.generateS3UploadToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ContentImpl.generateS3UploadToken() = %v, want %v", got, tt.want)
			}
			if tt.wantErr {
				assert.Equal(t, "", got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestContentImpl_GenerateVideoUploadToken(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.GenerateVideoUploadTokenOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.GenerateVideoUploadTokenResp
		wantErr    bool
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockS3)
		validate   func(*testing.T, *TC, *schema.GenerateVideoUploadTokenResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GenerateVideoUploadTokenOpts{
					FileName: "test.mov",
				},
			},
			buildStubs: func(tt *TC, mc *mock.MockS3) {
				resp := "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd"
				mc.EXPECT().GetPutObjectRequestURL(gomock.Any()).Times(1).Return(resp, nil)
				tt.want = &schema.GenerateVideoUploadTokenResp{
					Token: resp,
				}
			},
			prepare: func(tt *TC) {},
			validate: func(t *testing.T, tt *TC, resp *schema.GenerateVideoUploadTokenResp) {
				assert.Equal(t, tt.want.Token, resp.Token)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &MediaImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			mockS3 := mock.NewMockS3(ctrl)
			ci.App.S3 = mockS3
			tt.fields.App.Media = ci
			tt.buildStubs(&tt, mockS3)

			got, err := ci.GenerateVideoUploadToken(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.GenerateVideoUploadToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestMediaImpl_DeleteMedia(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		id primitive.ObjectID
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			want: true,
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateVideoOpts()
				resp, _ := tt.fields.App.Media.CreateVideoMedia(opts)
				tt.args.id = resp.ID
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				assert.Equal(t, tt.want, resp)
				count, err := tt.fields.DB.Collection(model.MediaColl).CountDocuments(context.TODO(), bson.M{"_id": tt.args.id})
				assert.Nil(t, err)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "[Error] When id does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: true,
			args:    args{},
			prepare: func(tt *TC) {
				tt.args.id = primitive.NewObjectIDFromTimestamp(time.Now())
				tt.err = errors.Errorf("media with id:%s not found", tt.args.id.Hex())
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				assert.Equal(t, tt.want, resp)
				count, err := tt.fields.DB.Collection(model.MediaColl).CountDocuments(context.TODO(), bson.M{"_id": tt.args.id})
				assert.Nil(t, err)
				assert.Equal(t, int64(0), count)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := &MediaImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Media = mi
			tt.prepare(&tt)
			got, err := mi.DeleteMedia(tt.args.id)
			if tt.wantErr {
				assert.False(t, got, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestMediaImpl_CreateVideoMedia(t *testing.T) {

	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateVideoOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.CreateVideoResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.CreateVideoResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateVideoOpts()
				tt.args.opts = opts
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateVideoResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.args.opts.FileName, resp.FileName)
				assert.Equal(t, tt.args.opts.GUID, resp.GUID)
				assert.Equal(t, tt.args.opts.SRCBucket, resp.SRCBucket)
				assert.Equal(t, tt.args.opts.SRCHeight, resp.Dimensions.Height)
				assert.Equal(t, tt.args.opts.SRCWidth, resp.Dimensions.Width)
				assert.Equal(t, tt.args.opts.IsPortrait, resp.IsPortrait)
				assert.Equal(t, tt.args.opts.CloudFrontURL, resp.CloudfrontURL)
				assert.Equal(t, tt.args.opts.Duration, resp.Duration)
				assert.Equal(t, tt.args.opts.Framerate, resp.Framerate)
				assert.Equal(t, tt.args.opts.PlaybackBucket, resp.PlaybackBucket)
				assert.Equal(t, tt.args.opts.PlaybackURL, resp.PlaybackURL)
				assert.Equal(t, tt.args.opts.ThumbnailBuckets, resp.ThumbnailBuckets)
				assert.Equal(t, tt.args.opts.ThumbnailURLS, resp.ThumbnailURLS)
				assert.Equal(t, tt.args.opts.ProcessedAt, resp.ProcessedAt)
				assert.WithinDuration(t, time.Now().UTC(), resp.CreatedAt, time.Millisecond*100)

				var doc model.Video

				err := tt.fields.DB.Collection(model.MediaColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, doc.ID, resp.ID)
				assert.Equal(t, doc.Type, model.VideoType)
				assert.Equal(t, doc.FileName, resp.FileName)
				assert.Equal(t, doc.GUID, resp.GUID)
				assert.Equal(t, doc.SRCBucket, resp.SRCBucket)
				assert.Equal(t, doc.Dimensions, resp.Dimensions)
				assert.Equal(t, doc.IsPortrait, resp.IsPortrait)
				assert.Equal(t, doc.CloudfrontURL, resp.CloudfrontURL)
				assert.Equal(t, doc.Duration, resp.Duration)
				assert.Equal(t, doc.Framerate, resp.Framerate)
				assert.Equal(t, doc.PlaybackBucket, resp.PlaybackBucket)
				assert.Equal(t, doc.PlaybackURL, resp.PlaybackURL)
				assert.Equal(t, doc.ThumbnailBuckets, resp.ThumbnailBuckets)
				assert.Equal(t, doc.ThumbnailURLS, resp.ThumbnailURLS)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := &MediaImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Media = mi
			tt.prepare(&tt)
			got, err := mi.CreateVideoMedia(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("MediaImpl.CreateVideoMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestMediaImpl_GetVideoMediaByID(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreateVideoOpts
		createResp *schema.CreateVideoResp
		id         primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.GetMediaResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockS3)
		validate   func(*testing.T, *TC, *schema.GetMediaResp)
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
				opts := schema.GetRandomCreateVideoOpts()
				tt.args.createOpts = opts
				resp, _ := tt.fields.App.Media.CreateVideoMedia(opts)
				tt.args.createResp = resp
				tt.args.id = resp.ID
				tt.want = &schema.GetMediaResp{
					ID:            resp.ID,
					FileName:      resp.FileName,
					CloudfrontURL: resp.CloudfrontURL,
					SRCBucket:     resp.SRCBucket,
					Dimensions:    resp.Dimensions,
					Duration:      resp.Duration,
					Framerate:     resp.Framerate,
					IsPortrait:    resp.IsPortrait,
					PlaybackURL:   resp.PlaybackURL,
					ThumbnailURLS: resp.ThumbnailURLS,
					CreatedAt:     resp.CreatedAt.UTC(),
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.GetMediaResp) {
				assert.Equal(t, tt.want.ID, resp.ID)
				assert.Equal(t, tt.want.FileName, resp.FileName)
				assert.Equal(t, tt.want.CloudfrontURL, resp.CloudfrontURL)
				assert.Equal(t, tt.want.SRCBucket, resp.SRCBucket)
				assert.Equal(t, tt.want.Dimensions, resp.Dimensions)
				assert.Equal(t, tt.want.Duration, resp.Duration)
				assert.Equal(t, tt.want.Framerate, resp.Framerate)
				assert.Equal(t, tt.want.IsPortrait, resp.IsPortrait)
				assert.Equal(t, tt.want.PlaybackURL, resp.PlaybackURL)
				assert.Equal(t, tt.want.ThumbnailURLS, resp.ThumbnailURLS)
				assert.WithinDuration(t, tt.want.CreatedAt, resp.CreatedAt, time.Microsecond*1000)
			},
		},
		{
			name: "[Error] Media ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.MediaConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: true,
			args:    args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateVideoOpts()
				tt.args.createOpts = opts
				resp, _ := tt.fields.App.Media.CreateVideoMedia(opts)
				tt.args.createResp = resp
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("media with id:%s not found", tt.args.id.Hex())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := &MediaImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Media = mi
			tt.prepare(&tt)
			got, err := mi.GetVideoMediaByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MediaImpl.GetVideoMediaByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestMediaImpl_CreateImageMedia(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateImageMediaOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.CreateImageMediaResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockS3)
		validate   func(*testing.T, *TC, *schema.CreateImageMediaResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateImageMediaOpts()
				tt.args.opts = opts
			},
			buildStubs: func(tt *TC, mc *mock.MockS3) {
				mc.EXPECT().PutObject(gomock.Any()).Times(1).Return(nil, nil)
			},
			validate: func(t *testing.T, tt *TC, got *schema.CreateImageMediaResp) {
				assert.False(t, got.ID.IsZero())
				assert.Contains(t, got.FileName, tt.args.opts.FileName)
				assert.Equal(t, uint(225), got.Dimensions.Height)
				assert.Equal(t, uint(225), got.Dimensions.Width)
				assert.NotEmpty(t, got.URL)
			},
		},
		{
			name: "[Error] Invalid Base64",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			buildStubs: func(tt *TC, mc *mock.MockS3) {
				mc.EXPECT().PutObject(gomock.Any()).Times(0).Return(nil, nil)
			},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateImageMediaOpts()
				opts.Base64SRC = opts.Base64SRC + "somerandomdata"
				tt.args.opts = opts
				tt.err = errors.New("failed to un-base image string: illegal base64 data at input byte 1892")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockS3 := mock.NewMockS3(ctrl)

			mi := &MediaImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Media = mi
			tt.fields.App.S3 = mockS3
			tt.buildStubs(&tt, mockS3)
			tt.prepare(&tt)
			got, err := mi.CreateImageMedia(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("MediaImpl.CreateImageMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
