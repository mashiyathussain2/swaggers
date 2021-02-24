//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_content.go -package=mock go-app/app Content

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
		S3     S3
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
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
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
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
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
				S3:     tt.fields.S3,
			}
			mockS3 := mock.NewMockS3(ctrl)
			ci.S3 = mockS3

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
		S3     S3
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
				DB:     app.MongoDB.Client.Database(app.Config.ContentConfig.DBName),
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
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.want.Token, resp.Token)
				var res model.Video
				err := tt.fields.DB.Collection(model.VideoContentColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&res)
				assert.Nil(t, err)
				assert.Empty(t, res.CloudfrontURL)
				assert.Empty(t, res.DestBucket)
				assert.Empty(t, res.Dimensions)
				assert.Empty(t, res.FileName)
				assert.Empty(t, res.Framerate)
				assert.Empty(t, res.SRCBucket)
				assert.Empty(t, res.Duration)
				assert.Empty(t, res.GUID)
				assert.False(t, res.IsPortrait)
				assert.False(t, res.IsProcessed)
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
				S3:     tt.fields.S3,
			}
			mockS3 := mock.NewMockS3(ctrl)
			ci.S3 = mockS3
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
