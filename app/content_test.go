package app

import (
	"context"
	"fmt"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"syreclabs.com/go/faker"
)

func TestContentImpl_CreatePebble(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreatePebbleOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.CreatePebbleResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, *schema.CreatePebbleResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			want: &schema.CreatePebbleResp{},
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreatePebbleOpts()
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
				tt.want.Token = resp.Token
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreatePebbleResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.want.Token, resp.Token)

				var res model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&res)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.opts.Caption, res.Caption)
				assert.Equal(t, tt.args.opts.BrandIDs, res.BrandIDs)
				assert.Equal(t, tt.args.opts.InfluencerIDs, res.InfluencerIDs)
				assert.Equal(t, tt.args.opts.CatalogIDs, res.CatalogIDs)
				assert.Equal(t, primitive.NilObjectID, res.CustomerID)
				assert.Equal(t, tt.args.opts.Label.Gender, res.Label.Genders)
				assert.Equal(t, tt.args.opts.Label.AgeGroup, res.Label.AgeGroups)
				assert.Equal(t, tt.args.opts.Label.Interests, res.Label.Interests)
				assert.Equal(t, model.PebbleType, res.Type)
				assert.WithinDuration(t, time.Now().UTC(), res.CreatedAt, time.Millisecond*100)
				assert.True(t, res.UpdatedAt.IsZero())
				assert.Nil(t, res.Hashtags)
				assert.False(t, res.IsProcessed)
				assert.False(t, res.IsActive)
			},
		},
		{
			name: "[Ok] With hashtags",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			want: &schema.CreatePebbleResp{},
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreatePebbleOpts()
				tt.args.opts.Caption += "This string #also has #hashtags"

			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
				tt.want.Token = resp.Token
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreatePebbleResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.want.Token, resp.Token)

				var res model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&res)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.opts.Caption, res.Caption)
				assert.Equal(t, tt.args.opts.BrandIDs, res.BrandIDs)
				assert.Equal(t, tt.args.opts.InfluencerIDs, res.InfluencerIDs)
				assert.Equal(t, tt.args.opts.CatalogIDs, res.CatalogIDs)
				assert.Equal(t, primitive.NilObjectID, res.CustomerID)
				assert.Equal(t, tt.args.opts.Label.Gender, res.Label.Genders)
				assert.Equal(t, tt.args.opts.Label.AgeGroup, res.Label.AgeGroups)
				assert.Equal(t, tt.args.opts.Label.Interests, res.Label.Interests)
				assert.Equal(t, model.PebbleType, res.Type)
				assert.WithinDuration(t, time.Now().UTC(), res.CreatedAt, time.Millisecond*100)
				assert.True(t, res.UpdatedAt.IsZero())
				assert.Equal(t, []string{"#also", "#hashtags"}, res.Hashtags)
				assert.False(t, res.IsProcessed)
				assert.False(t, res.IsActive)
			},
		},
		{
			name: "[Error] error on GenerateVideoUploadToken",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreatePebbleOpts()
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				err := errors.Errorf("cannot generate upload token")
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(nil, err)
				tt.err = err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pi := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = pi
			mockMedia := mock.NewMockMedia(ctrl)
			pi.App.Media = mockMedia
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockMedia)
			got, err := pi.CreatePebble(tt.args.opts)
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestContentImpl_EditPebble(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreatePebbleOpts
		createResp *schema.CreatePebbleResp
		opts       *schema.EditPebbleOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.EditPebbleResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, *schema.EditPebbleResp)
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
				createOpts := schema.GetRandomCreatePebbleOpts()
				createResp, _ := tt.fields.App.Content.CreatePebble(createOpts)
				tt.args.opts = &schema.EditPebbleOpts{
					ID:      createResp.ID,
					Caption: createOpts.Caption + " Edited",
				}
				tt.args.createOpts = createOpts
				tt.want = &schema.EditPebbleResp{
					ID:      createResp.ID,
					Caption: createOpts.Caption + " Edited",
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp *schema.EditPebbleResp) {
				assert.Equal(t, tt.want, resp)
				var doc model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.opts.Caption, doc.Caption)
				assert.Equal(t, tt.args.createOpts.BrandIDs, doc.BrandIDs)
				assert.Equal(t, tt.args.createOpts.InfluencerIDs, doc.InfluencerIDs)
				assert.Equal(t, tt.args.createOpts.CatalogIDs, doc.CatalogIDs)
				assert.Equal(t, false, doc.IsActive)
				assert.Equal(t, false, doc.IsProcessed)
				assert.Equal(t, tt.args.createOpts.Label.AgeGroup, doc.Label.AgeGroups)
				assert.Equal(t, tt.args.createOpts.Label.Gender, doc.Label.Genders)
				assert.Equal(t, tt.args.createOpts.Label.Interests, doc.Label.Interests)
			},
		},
		{
			name: "[Ok] Updating caption with hashtag",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreatePebbleOpts()
				createResp, _ := tt.fields.App.Content.CreatePebble(createOpts)
				tt.args.opts = &schema.EditPebbleOpts{
					ID:      createResp.ID,
					Caption: "Edited #Edited",
				}
				tt.args.createOpts = createOpts
				tt.args.createResp = createResp
				tt.want = &schema.EditPebbleResp{
					ID:      createResp.ID,
					Caption: "Edited #Edited",
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp *schema.EditPebbleResp) {
				assert.Equal(t, tt.want, resp)
				var doc model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.opts.Caption, doc.Caption)
				assert.Equal(t, tt.args.createOpts.BrandIDs, doc.BrandIDs)
				assert.Equal(t, tt.args.createOpts.InfluencerIDs, doc.InfluencerIDs)
				assert.Equal(t, tt.args.createOpts.CatalogIDs, doc.CatalogIDs)
				assert.Equal(t, false, doc.IsActive)
				assert.Equal(t, false, doc.IsProcessed)
				assert.Equal(t, tt.args.createOpts.Label.AgeGroup, doc.Label.AgeGroups)
				assert.Equal(t, tt.args.createOpts.Label.Gender, doc.Label.Genders)
				assert.Equal(t, tt.args.createOpts.Label.Interests, doc.Label.Interests)
				assert.Equal(t, []string{"#Edited"}, doc.Hashtags)
			},
		},
		{
			name: "[Ok] Updating Label",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreatePebbleOpts()
				createResp, _ := tt.fields.App.Content.CreatePebble(createOpts)
				tt.args.opts = &schema.EditPebbleOpts{
					ID: createResp.ID,
					Label: &schema.EditLabelOpts{
						Interests: []string{"footwear", "women"},
						AgeGroup:  []string{"16-20"},
						Gender:    []string{"F"},
					},
				}
				tt.args.createOpts = createOpts
				tt.args.createResp = createResp
				tt.want = &schema.EditPebbleResp{
					ID: createResp.ID,
					Label: &schema.EditLabelOpts{
						Interests: []string{"footwear", "women"},
						AgeGroup:  []string{"16-20"},
						Gender:    []string{"F"},
					},
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp *schema.EditPebbleResp) {
				assert.Equal(t, tt.want, resp)
				var doc model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.createOpts.Caption, doc.Caption)
				assert.Equal(t, tt.args.createOpts.BrandIDs, doc.BrandIDs)
				assert.Equal(t, tt.args.createOpts.InfluencerIDs, doc.InfluencerIDs)
				assert.Equal(t, tt.args.createOpts.CatalogIDs, doc.CatalogIDs)
				assert.Equal(t, false, doc.IsActive)
				assert.Equal(t, false, doc.IsProcessed)
				assert.Equal(t, tt.args.opts.Label.AgeGroup, doc.Label.AgeGroups)
				assert.Equal(t, tt.args.opts.Label.Gender, doc.Label.Genders)
				assert.Equal(t, tt.args.opts.Label.Interests, doc.Label.Interests)
			},
		},
		{
			name: "[Ok] Updating Label Interests",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreatePebbleOpts()
				createResp, _ := tt.fields.App.Content.CreatePebble(createOpts)
				tt.args.opts = &schema.EditPebbleOpts{
					ID: createResp.ID,
					Label: &schema.EditLabelOpts{
						Interests: []string{"footwear", "women"},
					},
				}
				tt.args.createOpts = createOpts
				tt.args.createResp = createResp
				tt.want = &schema.EditPebbleResp{
					ID: createResp.ID,
					Label: &schema.EditLabelOpts{
						Interests: []string{"footwear", "women"},
					},
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp *schema.EditPebbleResp) {
				assert.Equal(t, tt.want, resp)
				var doc model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.createOpts.Caption, doc.Caption)
				assert.Equal(t, tt.args.createOpts.BrandIDs, doc.BrandIDs)
				assert.Equal(t, tt.args.createOpts.InfluencerIDs, doc.InfluencerIDs)
				assert.Equal(t, tt.args.createOpts.CatalogIDs, doc.CatalogIDs)
				assert.Equal(t, false, doc.IsActive)
				assert.Equal(t, false, doc.IsProcessed)
				assert.Equal(t, tt.args.createOpts.Label.AgeGroup, doc.Label.AgeGroups)
				assert.Equal(t, tt.args.createOpts.Label.Gender, doc.Label.Genders)
				assert.Equal(t, tt.args.opts.Label.Interests, doc.Label.Interests)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = ci
			mockMedia := mock.NewMockMedia(ctrl)
			tt.fields.App.Media = mockMedia
			tt.buildStubs(&tt, mockMedia)
			tt.prepare(&tt)
			got, err := ci.EditPebble(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.EditPebble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestContentImpl_DeletePebble(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreatePebbleOpts
		createResp *schema.CreatePebbleResp
		id         primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       bool
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			want: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreatePebbleOpts()
				createResp, _ := tt.fields.App.Content.CreatePebble(createOpts)
				tt.args.id = createResp.ID
				tt.args.createOpts = createOpts
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			got, err := ci.DeletePebble(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.DeletePebble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ContentImpl.DeletePebble() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentImpl_ProcessVideoContent(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts    *schema.CreatePebbleOpts
		createResp    *schema.CreatePebbleResp
		mockMediaResp *schema.CreateVideoResp
		opts          *schema.ProcessVideoContentOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       bool
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Error] when content does not match filename",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateVideoOpts()
				tt.args.opts = opts
				id, _ := primitive.ObjectIDFromHex(strings.Split(tt.args.opts.FileName, ".")[0])
				tt.err = errors.Wrapf(errors.New("mongo: no documents in result"), "failed to mark content:%s as processed", id.Hex())
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := &schema.CreateVideoResp{
					ID:               primitive.NewObjectIDFromTimestamp(time.Now()),
					GUID:             tt.args.opts.GUID,
					FileName:         tt.args.opts.FileName,
					SRCBucket:        tt.args.opts.SRCBucket,
					DestBucket:       tt.args.opts.DestBucket,
					CloudfrontURL:    tt.args.opts.CloudFrontURL,
					IsPortrait:       tt.args.opts.IsPortrait,
					Duration:         tt.args.opts.Duration,
					Framerate:        tt.args.opts.Framerate,
					PlaybackBucket:   tt.args.opts.PlaybackBucket,
					PlaybackURL:      tt.args.opts.PlaybackURL,
					ThumbnailBuckets: tt.args.opts.ThumbnailBuckets,
					ThumbnailURLS:    tt.args.opts.ThumbnailURLS,
					Dimensions:       &model.Dimensions{Width: tt.args.opts.SRCWidth, Height: tt.args.opts.SRCHeight},
					CreatedAt:        time.Now().UTC(),
					ProcessedAt:      time.Now().UTC(),
				}
				tt.args.mockMediaResp = resp
				mc.EXPECT().CreateVideoMedia(tt.args.opts).Times(1).Return(resp, nil)
			},
		},
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreatePebbleOpts()
				tt.args.createOpts = createOpts
				opts := schema.GetRandomCreateVideoOpts()
				tt.args.opts = opts

				tt.want = true
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {

				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}, nil)
				resp0, _ := tt.fields.App.Content.CreatePebble(tt.args.createOpts)
				tt.args.createResp = resp0
				tt.args.opts.FileName = fmt.Sprintf("%s.%s", resp0.ID.Hex(), strings.Split(tt.args.createOpts.FileName, ".")[1])

				resp1 := &schema.CreateVideoResp{
					ID:               primitive.NewObjectIDFromTimestamp(time.Now()),
					GUID:             tt.args.opts.GUID,
					FileName:         tt.args.opts.FileName,
					SRCBucket:        tt.args.opts.SRCBucket,
					DestBucket:       tt.args.opts.DestBucket,
					CloudfrontURL:    tt.args.opts.CloudFrontURL,
					IsPortrait:       tt.args.opts.IsPortrait,
					Duration:         tt.args.opts.Duration,
					Framerate:        tt.args.opts.Framerate,
					PlaybackBucket:   tt.args.opts.PlaybackBucket,
					PlaybackURL:      tt.args.opts.PlaybackURL,
					ThumbnailBuckets: tt.args.opts.ThumbnailBuckets,
					ThumbnailURLS:    tt.args.opts.ThumbnailURLS,
					Dimensions:       &model.Dimensions{Width: tt.args.opts.SRCWidth, Height: tt.args.opts.SRCHeight},
					CreatedAt:        time.Now().UTC(),
					ProcessedAt:      time.Now().UTC(),
				}
				tt.args.mockMediaResp = resp1
				mc.EXPECT().CreateVideoMedia(tt.args.opts).Times(1).Return(resp1, nil)
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				assert.Equal(t, tt.want, resp)
				var doc model.Content
				id, _ := primitive.ObjectIDFromHex(strings.Split(tt.args.opts.FileName, ".")[0])
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, true, doc.IsProcessed)
				assert.Equal(t, model.VideoType, doc.MediaType)
				assert.False(t, doc.MediaID.IsZero())
				assert.WithinDuration(t, time.Now().UTC(), doc.ProcessedAt, 100*time.Millisecond)
				assert.Equal(t, model.VideoType, doc.MediaType)
				assert.Equal(t, tt.args.mockMediaResp.ID, doc.MediaID)
			},
		},
		{
			name: "[Error] while creating media",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreatePebbleOpts()
				tt.args.createOpts = createOpts
				opts := schema.GetRandomCreateVideoOpts()
				tt.args.opts = opts

				tt.want = true
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {

				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}, nil)
				resp0, _ := tt.fields.App.Content.CreatePebble(tt.args.createOpts)
				tt.args.createResp = resp0
				tt.args.opts.FileName = fmt.Sprintf("%s.%s", resp0.ID.Hex(), strings.Split(tt.args.createOpts.FileName, ".")[1])
				err := errors.New("DB query failed")
				mc.EXPECT().CreateVideoMedia(tt.args.opts).Times(1).Return(nil, err)
				tt.err = errors.Wrap(err, "failed to create video media")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = ci
			mockMedia := mock.NewMockMedia(ctrl)
			tt.fields.App.Media = mockMedia
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockMedia)
			got, err := ci.ProcessVideoContent(tt.args.opts)
			if tt.wantErr {
				assert.False(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestContentImpl_GetContentByID(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createContentOpts *model.Content
		id                primitive.ObjectID
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.GetContentResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, *schema.GetContentResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				var InsertedContent []model.Content
				for i := 0; i < faker.RandomInt(1, 10); i++ {
					contentOpts := schema.GetRandomCreatePebbleOpts()
					c := model.Content{
						ID:            primitive.NewObjectIDFromTimestamp(time.Now()),
						Type:          model.PebbleType,
						MediaType:     model.VideoType,
						MediaID:       primitive.NewObjectIDFromTimestamp(time.Now()),
						InfluencerIDs: contentOpts.InfluencerIDs,
						BrandIDs:      contentOpts.BrandIDs,
						Label: &model.Label{
							Interests: contentOpts.Label.Interests,
							AgeGroups: contentOpts.Label.AgeGroup,
							Genders:   contentOpts.Label.Gender,
						},
						IsProcessed: true,
						Caption:     contentOpts.Caption,
						Hashtags:    []string{"#test", "#unitest"},
						CatalogIDs:  contentOpts.CatalogIDs,
						CreatedAt:   time.Now().UTC(),
						ProcessedAt: time.Now().UTC(),
					}
					opts := options.FindOneAndUpdate().SetUpsert(true)
					tt.fields.DB.Collection(model.ContentColl).FindOneAndUpdate(context.TODO(), bson.M{"_id": c.ID}, bson.M{"$set": c}, opts).Decode(&c)
					InsertedContent = append(InsertedContent, c)
				}

				tt.args.createContentOpts = &InsertedContent[faker.RandomInt(0, len(InsertedContent)-1)]
				tt.args.id = tt.args.createContentOpts.ID
				tt.want = &schema.GetContentResp{
					ID:            tt.args.createContentOpts.ID,
					Type:          tt.args.createContentOpts.Type,
					MediaType:     tt.args.createContentOpts.MediaType,
					MediaID:       tt.args.createContentOpts.MediaID,
					InfluencerIDs: tt.args.createContentOpts.InfluencerIDs,
					BrandIDs:      tt.args.createContentOpts.BrandIDs,
					CustomerID:    tt.args.createContentOpts.CustomerID,
					CatalogIDs:    tt.args.createContentOpts.CatalogIDs,
					Label:         tt.args.createContentOpts.Label,
					IsActive:      tt.args.createContentOpts.IsActive,
					Caption:       tt.args.createContentOpts.Caption,
					Hashtags:      tt.args.createContentOpts.Hashtags,
					CreatedAt:     tt.args.createContentOpts.CreatedAt,
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				mediaOpts := schema.GetRandomCreateVideoOpts()
				resp := schema.GetMediaResp{
					ID:            tt.args.createContentOpts.MediaID,
					CloudfrontURL: mediaOpts.CloudFrontURL,
					SRCBucket:     mediaOpts.SRCBucket,
					FileName:      mediaOpts.FileName,
					IsPortrait:    mediaOpts.IsPortrait,
					Dimensions: &model.Dimensions{
						Height: mediaOpts.SRCHeight,
						Width:  mediaOpts.SRCWidth,
					},
					Duration:      mediaOpts.Duration,
					Framerate:     mediaOpts.Framerate,
					PlaybackURL:   mediaOpts.PlaybackURL,
					ThumbnailURLS: mediaOpts.ThumbnailURLS,
					CreatedAt:     time.Now().UTC(),
				}
				mc.EXPECT().GetVideoMediaByID(gomock.Any()).Times(1).Return(&resp, nil)
				tt.want.MediaInfo = &resp
			},
			validate: func(t *testing.T, tt *TC, resp *schema.GetContentResp) {
				assert.WithinDuration(t, tt.want.CreatedAt, resp.CreatedAt, 1000*time.Microsecond)
				tt.want.CreatedAt = time.Time{}
				resp.CreatedAt = time.Time{}
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Error] when invalid media_id is linked with content",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				var InsertedContent []model.Content
				for i := 0; i < faker.RandomInt(1, 10); i++ {
					contentOpts := schema.GetRandomCreatePebbleOpts()
					c := model.Content{
						ID:            primitive.NewObjectIDFromTimestamp(time.Now()),
						Type:          model.PebbleType,
						MediaType:     model.VideoType,
						MediaID:       primitive.NewObjectIDFromTimestamp(time.Now()),
						InfluencerIDs: contentOpts.InfluencerIDs,
						BrandIDs:      contentOpts.BrandIDs,
						Label: &model.Label{
							Interests: contentOpts.Label.Interests,
							AgeGroups: contentOpts.Label.AgeGroup,
							Genders:   contentOpts.Label.Gender,
						},
						IsProcessed: true,
						Caption:     contentOpts.Caption,
						Hashtags:    []string{"#test", "#unitest"},
						CatalogIDs:  contentOpts.CatalogIDs,
						CreatedAt:   time.Now().UTC(),
						ProcessedAt: time.Now().UTC(),
					}
					opts := options.FindOneAndUpdate().SetUpsert(true)
					tt.fields.DB.Collection(model.ContentColl).FindOneAndUpdate(context.TODO(), bson.M{"_id": c.ID}, bson.M{"$set": c}, opts).Decode(&c)
					InsertedContent = append(InsertedContent, c)
				}

				tt.args.createContentOpts = &InsertedContent[0]
				tt.args.id = InsertedContent[0].ID
				tt.err = errors.Errorf("media with id %s not found", InsertedContent[0].MediaID.Hex())
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				mc.EXPECT().GetVideoMediaByID(gomock.Any()).Times(1).Return(nil, errors.Errorf("media with id %s not found", tt.args.createContentOpts.MediaID.Hex()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = ci
			mockMedia := mock.NewMockMedia(ctrl)
			tt.fields.App.Media = mockMedia
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockMedia)
			got, err := ci.GetContentByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.GetContent() error = %v, wantErr %v", err, tt.wantErr)
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

func TestContentImpl_GetContent(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	// Adding random content and media into DB
	var contents []model.Content
	var videoMedias []model.Video
	ctx := context.TODO()

	// Adding Processed
	for i := 0; i < 5; i++ {
		c := model.GetRandomContent()
		m := model.GetRandomVideoMedia()
		c.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		c.MediaType = m.Type
		m.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		c.MediaID = m.ID
		m.FileName = fmt.Sprintf("%s%s", c.ID.Hex(), faker.RandomChoice([]string{".mp4", ".mov"}))

		opts := options.FindOneAndUpdate().SetUpsert(true)
		f1 := bson.M{"_id": c.ID}
		f2 := bson.M{"_id": m.ID}
		u1 := bson.M{"$set": c}
		u2 := bson.M{"$set": m}
		_ = app.MongoDB.Client.Database(app.Config.ContentConfig.DBName).Collection(model.ContentColl).FindOneAndUpdate(ctx, f1, u1, opts).Decode(&c)
		_ = app.MongoDB.Client.Database(app.Config.MediaConfig.DBName).Collection(model.MediaColl).FindOneAndUpdate(ctx, f2, u2, opts).Decode(&m)
		contents = append(contents, *c)
		videoMedias = append(videoMedias, *m)
	}

	// Adding Processed And Active
	for i := 0; i < 5; i++ {
		c := model.GetRandomContent()
		m := model.GetRandomVideoMedia()
		c.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		c.MediaType = m.Type
		m.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		c.MediaID = m.ID
		c.IsActive = true
		m.FileName = fmt.Sprintf("%s%s", c.ID.Hex(), faker.RandomChoice([]string{".mp4", ".mov"}))

		opts := options.FindOneAndUpdate().SetUpsert(true)
		f1 := bson.M{"_id": c.ID}
		f2 := bson.M{"_id": m.ID}
		u1 := bson.M{"$set": c}
		u2 := bson.M{"$set": m}
		_ = app.MongoDB.Client.Database(app.Config.ContentConfig.DBName).Collection(model.ContentColl).FindOneAndUpdate(ctx, f1, u1, opts).Decode(&c)
		_ = app.MongoDB.Client.Database(app.Config.MediaConfig.DBName).Collection(model.MediaColl).FindOneAndUpdate(ctx, f2, u2, opts).Decode(&m)
		contents = append(contents, *c)
		videoMedias = append(videoMedias, *m)
	}

	// Adding UnProcessed
	for i := 0; i < 5; i++ {
		c := model.GetRandomContent()
		// m := model.GetRandomVideoMedia()
		c.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		// c.MediaType = m.Type
		// m.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		// c.MediaID = m.ID
		c.IsProcessed = false
		// m.FileName = fmt.Sprintf("%s%s", c.ID.Hex(), faker.RandomChoice([]string{".mp4", ".mov"}))

		opts := options.FindOneAndUpdate().SetUpsert(true)
		f1 := bson.M{"_id": c.ID}
		// f2 := bson.M{"_id": m.ID}
		u1 := bson.M{"$set": c}
		// u2 := bson.M{"$set": m}
		_ = app.MongoDB.Client.Database(app.Config.ContentConfig.DBName).Collection(model.ContentColl).FindOneAndUpdate(ctx, f1, u1, opts).Decode(&c)
		// _ = app.MongoDB.Client.Database(app.Config.ContentConfig.DBName).Collection(model.MediaColl).FindOneAndUpdate(ctx, f2, u2, opts).Decode(&m)
		contents = append(contents, *c)
		// videoMedias = append(videoMedias, *m)
	}

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		filterOpts *schema.GetContentFilter
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     []schema.GetContentResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, []schema.GetContentResp)
	}

	tests := []TC{
		{
			name: "[Ok] With IsActive Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				t := true
				opts := schema.GetContentFilter{
					IsActive: &t,
				}
				tt.args.filterOpts = &opts
				for i, c := range contents[5:10] {
					tt.want = append(tt.want, schema.GetContentResp{
						ID:            c.ID,
						Type:          c.Type,
						MediaType:     c.MediaType,
						MediaID:       c.MediaID,
						InfluencerIDs: c.InfluencerIDs,
						CustomerID:    c.CustomerID,
						BrandIDs:      c.BrandIDs,
						CatalogIDs:    c.CatalogIDs,
						Label:         c.Label,
						Caption:       c.Caption,
						Hashtags:      c.Hashtags,
						IsActive:      c.IsActive,
						CreatedAt:     c.CreatedAt,
						MediaInfo: &schema.GetMediaResp{
							ID:            videoMedias[5+i].ID,
							CloudfrontURL: videoMedias[5+i].CloudfrontURL,
							SRCBucket:     videoMedias[5+i].SRCBucket,
							Dimensions:    videoMedias[5+i].Dimensions,
							Duration:      videoMedias[5+i].Duration,
							FileName:      videoMedias[5+i].FileName,
							IsPortrait:    videoMedias[5+i].IsPortrait,
							Framerate:     videoMedias[5+i].Framerate,
							PlaybackURL:   videoMedias[5+i].PlaybackURL,
							ThumbnailURLS: videoMedias[5+i].ThumbnailURLS,
							CreatedAt:     videoMedias[5+i].CreatedAt,
						},
					})
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetContentResp) {
				assert.Equal(t, len(tt.want), len(got))
				for _, expected := range tt.want {
					for _, resp := range got {
						if expected.ID != resp.ID {
							continue
						}
						assert.WithinDuration(t, expected.MediaInfo.CreatedAt, resp.MediaInfo.CreatedAt, 100*time.Millisecond)
						assert.WithinDuration(t, expected.CreatedAt, resp.CreatedAt, 100*time.Millisecond)
						x := time.Time{}
						expected.CreatedAt = x
						expected.MediaInfo.CreatedAt = x
						resp.CreatedAt = x
						resp.MediaInfo.CreatedAt = x
						assert.Equal(t, expected, resp)
					}
				}
			},
		},
		{
			name: "[Ok] With IsActive & IsProcessed Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				t := true
				opts := schema.GetContentFilter{
					IsActive: &t,
				}
				tt.args.filterOpts = &opts
				for i, c := range contents[5:10] {
					tt.want = append(tt.want, schema.GetContentResp{
						ID:            c.ID,
						Type:          c.Type,
						MediaType:     c.MediaType,
						MediaID:       c.MediaID,
						InfluencerIDs: c.InfluencerIDs,
						CustomerID:    c.CustomerID,
						BrandIDs:      c.BrandIDs,
						CatalogIDs:    c.CatalogIDs,
						Label:         c.Label,
						Caption:       c.Caption,
						Hashtags:      c.Hashtags,
						IsActive:      c.IsActive,
						CreatedAt:     c.CreatedAt,
						MediaInfo: &schema.GetMediaResp{
							ID:            videoMedias[5+i].ID,
							CloudfrontURL: videoMedias[5+i].CloudfrontURL,
							SRCBucket:     videoMedias[5+i].SRCBucket,
							Dimensions:    videoMedias[5+i].Dimensions,
							Duration:      videoMedias[5+i].Duration,
							FileName:      videoMedias[5+i].FileName,
							IsPortrait:    videoMedias[5+i].IsPortrait,
							Framerate:     videoMedias[5+i].Framerate,
							PlaybackURL:   videoMedias[5+i].PlaybackURL,
							ThumbnailURLS: videoMedias[5+i].ThumbnailURLS,
							CreatedAt:     videoMedias[5+i].CreatedAt,
						},
					})
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetContentResp) {
				for i, expected := range tt.want {
					assert.WithinDuration(t, expected.MediaInfo.CreatedAt, got[i].MediaInfo.CreatedAt, 100*time.Millisecond)
					assert.WithinDuration(t, expected.CreatedAt, got[i].CreatedAt, 100*time.Millisecond)
					x := time.Time{}
					expected.CreatedAt = x
					expected.MediaInfo.CreatedAt = x
					got[i].CreatedAt = x
					got[i].MediaInfo.CreatedAt = x
					assert.Equal(t, expected, got[i])
				}
			},
		},
		{
			name: "[Ok] With UnProcessed Content Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				// t := true
				f := false
				opts := schema.GetContentFilter{
					IsProcessed: &f,
				}
				tt.args.filterOpts = &opts
				for _, c := range contents[10:15] {
					tt.want = append(tt.want, schema.GetContentResp{
						ID:            c.ID,
						Type:          c.Type,
						MediaType:     c.MediaType,
						MediaID:       c.MediaID,
						InfluencerIDs: c.InfluencerIDs,
						CustomerID:    c.CustomerID,
						BrandIDs:      c.BrandIDs,
						CatalogIDs:    c.CatalogIDs,
						Label:         c.Label,
						Caption:       c.Caption,
						Hashtags:      c.Hashtags,
						IsActive:      c.IsActive,
						CreatedAt:     c.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetContentResp) {
				for i, expected := range tt.want {
					assert.WithinDuration(t, expected.CreatedAt, got[i].CreatedAt, 100*time.Millisecond)
					x := time.Time{}
					expected.CreatedAt = x
					got[i].CreatedAt = x
					assert.Equal(t, expected, got[i])
				}
			},
		},
		{
			name: "[Ok] With BrandIds & CatalogIDs",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				// t := true
			ALLOCATE_RANDOM_INT:
				randInt0 := faker.RandomInt(0, len(contents))
				randInt1 := faker.RandomInt(0, len(contents))
				if randInt0 == randInt1 {
					goto ALLOCATE_RANDOM_INT
				}

				var brandIDs []primitive.ObjectID
				brandIDs = append(brandIDs, contents[randInt0].BrandIDs...)
				brandIDs = append(brandIDs, contents[randInt1].BrandIDs...)
				opts := schema.GetContentFilter{
					BrandIDs: brandIDs,
				}
				tt.args.filterOpts = &opts

				tt.want = []schema.GetContentResp{
					{
						ID:            contents[randInt0].ID,
						Type:          contents[randInt0].Type,
						MediaType:     contents[randInt0].MediaType,
						MediaID:       contents[randInt0].MediaID,
						InfluencerIDs: contents[randInt0].InfluencerIDs,
						CustomerID:    contents[randInt0].CustomerID,
						BrandIDs:      contents[randInt0].BrandIDs,
						CatalogIDs:    contents[randInt0].CatalogIDs,
						Label:         contents[randInt0].Label,
						Caption:       contents[randInt0].Caption,
						Hashtags:      contents[randInt0].Hashtags,
						IsActive:      contents[randInt0].IsActive,
						CreatedAt:     contents[randInt0].CreatedAt,
					},
					{
						ID:            contents[randInt1].ID,
						Type:          contents[randInt1].Type,
						MediaType:     contents[randInt1].MediaType,
						MediaID:       contents[randInt1].MediaID,
						InfluencerIDs: contents[randInt1].InfluencerIDs,
						CustomerID:    contents[randInt1].CustomerID,
						BrandIDs:      contents[randInt1].BrandIDs,
						CatalogIDs:    contents[randInt1].CatalogIDs,
						Label:         contents[randInt1].Label,
						Caption:       contents[randInt1].Caption,
						Hashtags:      contents[randInt1].Hashtags,
						IsActive:      contents[randInt1].IsActive,
						CreatedAt:     contents[randInt1].CreatedAt,
					},
				}

			},
			validate: func(t *testing.T, tt *TC, got []schema.GetContentResp) {
				assert.Len(t, got, len(tt.want))
				for i, expected := range tt.want {
					got[i].MediaInfo = nil
					assert.WithinDuration(t, expected.CreatedAt, got[i].CreatedAt, 100*time.Millisecond)
					x := time.Time{}
					expected.CreatedAt = x
					got[i].CreatedAt = x
					assert.Equal(t, expected, got[i])
				}
			},
		},
		{
			name: "[Ok] No Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetContentFilter{}
				tt.args.filterOpts = &opts
				for i, c := range contents {
					if i == 10 {
						break
					}
					tt.want = append(tt.want, schema.GetContentResp{
						ID:            c.ID,
						Type:          c.Type,
						MediaType:     c.MediaType,
						MediaID:       c.MediaID,
						InfluencerIDs: c.InfluencerIDs,
						CustomerID:    c.CustomerID,
						BrandIDs:      c.BrandIDs,
						CatalogIDs:    c.CatalogIDs,
						Label:         c.Label,
						Caption:       c.Caption,
						Hashtags:      c.Hashtags,
						IsActive:      c.IsActive,
						CreatedAt:     c.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetContentResp) {
				assert.Len(t, got, 10)
				for _, expected := range tt.want {
					for _, resp := range got {
						if expected.ID != resp.ID {
							continue
						}
						resp.MediaInfo = nil
						assert.WithinDuration(t, expected.CreatedAt, resp.CreatedAt, 100*time.Millisecond)
						x := time.Time{}
						expected.CreatedAt = x
						resp.CreatedAt = x
						assert.Equal(t, expected, resp)
					}
				}
			},
		},
		{
			name: "[Ok] Pagination",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetContentFilter{
					Page: 1,
				}
				tt.args.filterOpts = &opts
				for _, c := range contents[10:] {
					tt.want = append(tt.want, schema.GetContentResp{
						ID:            c.ID,
						Type:          c.Type,
						MediaType:     c.MediaType,
						MediaID:       c.MediaID,
						InfluencerIDs: c.InfluencerIDs,
						CustomerID:    c.CustomerID,
						BrandIDs:      c.BrandIDs,
						CatalogIDs:    c.CatalogIDs,
						Label:         c.Label,
						Caption:       c.Caption,
						Hashtags:      c.Hashtags,
						IsActive:      c.IsActive,
						CreatedAt:     c.CreatedAt,
					})
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetContentResp) {
				assert.Len(t, got, 5)
				for i, expected := range tt.want {
					got[i].MediaInfo = nil
					assert.WithinDuration(t, expected.CreatedAt, got[i].CreatedAt, 100*time.Millisecond)
					x := time.Time{}
					expected.CreatedAt = x
					got[i].CreatedAt = x
					assert.Equal(t, expected, got[i])
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = ci
			tt.prepare(&tt)
			got, err := ci.GetContent(tt.args.filterOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestContentImpl_CreateCatalogVideoContent(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateVideoCatalogContentOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.CreateVideoCatalogContentResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, *schema.CreateVideoCatalogContentResp)
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
				tt.args.opts = schema.GetRandomCreateCatalogContentOpts()
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
				tt.want = &schema.CreateVideoCatalogContentResp{Token: resp.Token}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreatePebbleResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.want.Token, resp.Token)

				var res model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&res)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.opts.BrandID, res.BrandIDs[0])
				assert.Nil(t, res.InfluencerIDs)
				assert.Equal(t, tt.args.opts.CatalogID, res.CatalogIDs[0])
				assert.Equal(t, primitive.NilObjectID, res.CustomerID)
				assert.Equal(t, tt.args.opts.Label.Gender, res.Label.Genders)
				assert.Equal(t, tt.args.opts.Label.AgeGroup, res.Label.AgeGroups)
				assert.Equal(t, tt.args.opts.Label.Interests, res.Label.Interests)
				assert.Equal(t, model.CatalogContentType, res.Type)
				assert.WithinDuration(t, time.Now().UTC(), res.CreatedAt, time.Millisecond*200)
				assert.True(t, res.UpdatedAt.IsZero())
				assert.Nil(t, res.Hashtags)
				assert.False(t, res.IsProcessed)
				assert.False(t, res.IsActive)
			},
		},
		{
			name: "[Error] error on GenerateVideoUploadToken",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreateCatalogContentOpts()
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				err := errors.Errorf("cannot generate upload token")
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(nil, err)
				tt.err = err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = ci
			mockMedia := mock.NewMockMedia(ctrl)
			ci.App.Media = mockMedia
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockMedia)

			got, err := ci.CreateCatalogVideoContent(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.CreateCatalogVideoContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestContentImpl_EditCatalogContent(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreateVideoCatalogContentOpts
		createResp *schema.CreateVideoCatalogContentResp
		opts       *schema.EditCatalogContentOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.EditCatalogContentResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockMedia)
		validate   func(*testing.T, *TC, *schema.EditCatalogContentResp)
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
				createOpts := schema.GetRandomCreateCatalogContentOpts()
				createResp, _ := tt.fields.App.Content.CreateCatalogVideoContent(createOpts)
				trueBool := true
				tt.args.opts = &schema.EditCatalogContentOpts{
					ID:       createResp.ID,
					IsActive: &trueBool,
				}
				tt.args.createOpts = createOpts
				tt.want = &schema.EditCatalogContentResp{
					ID:       createResp.ID,
					IsActive: &trueBool,
				}
			},
			buildStubs: func(tt *TC, mc *mock.MockMedia) {
				resp := schema.GenerateVideoUploadTokenResp{
					Token: "https://hypd-vod-source-16jim3me9cmrc.s3.ap-south-1.amazonaws.com/5fbb7f1f7f10f60aaffa2598.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA5LCMVADVOIOHO66X%2F20201123%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20201123T092135Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4ed77bbf055c9dbdfa0f45c6e352859f0ed0cf3dad2175c19469427e0f7c82dd",
				}
				mc.EXPECT().GenerateVideoUploadToken(gomock.Any()).Times(1).Return(&resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp *schema.EditCatalogContentResp) {
				assert.Equal(t, tt.want, resp)
				var doc model.Content
				err := tt.fields.DB.Collection(model.ContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.createOpts.BrandID, doc.BrandIDs[0])
				assert.Equal(t, primitive.NilObjectID, doc.CustomerID)
				assert.Nil(t, doc.InfluencerIDs)
				assert.Equal(t, tt.args.createOpts.CatalogID, doc.CatalogIDs[0])
				assert.Equal(t, true, doc.IsActive)
				assert.Equal(t, false, doc.IsProcessed)
				assert.Equal(t, tt.args.createOpts.Label.AgeGroup, doc.Label.AgeGroups)
				assert.Equal(t, tt.args.createOpts.Label.Gender, doc.Label.Genders)
				assert.Equal(t, tt.args.createOpts.Label.Interests, doc.Label.Interests)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &ContentImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Content = ci
			mockMedia := mock.NewMockMedia(ctrl)
			tt.fields.App.Media = mockMedia
			tt.buildStubs(&tt, mockMedia)
			tt.prepare(&tt)

			got, err := ci.EditCatalogContent(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentImpl.EditCatalogContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}
