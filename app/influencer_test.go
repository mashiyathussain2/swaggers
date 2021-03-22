package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

func TestInfluencerImpl_CreateInfluencer(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateInfluencerOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.CreateInfluencerResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.CreateInfluencerResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.InfluencerConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateInfluencerOpts()
				tt.args.opts = opts
			},
			validate: func(t *testing.T, tt *TC, got *schema.CreateInfluencerResp) {
				assert.Equal(t, tt.args.opts.Name, got.Name)
				assert.Equal(t, tt.args.opts.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.Equal(t, tt.args.opts.Bio, got.Bio)
				assert.Equal(t, tt.args.opts.ExternalLinks, got.ExternalLinks)
				assert.Equal(t, tt.args.opts.ProfileImage.SRC, got.ProfileImage.SRC)
				assert.Equal(t, 200, got.ProfileImage.Width)
				assert.Equal(t, 400, got.ProfileImage.Height)
				if tt.args.opts.SocialAccount != nil {
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Facebook.FollowersCount), got.SocialAccount.Facebook.FollowersCount)
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Instagram.FollowersCount), got.SocialAccount.Instagram.FollowersCount)
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Youtube.FollowersCount), got.SocialAccount.Youtube.FollowersCount)
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Twitter.FollowersCount), got.SocialAccount.Twitter.FollowersCount)
				}
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 3*time.Second)

				var doc model.Influencer
				err := tt.fields.DB.Collection(model.InfluencerColl).FindOne(context.TODO(), bson.M{"_id": got.ID}).Decode(&doc)
				assert.Nil(t, err)
			},
		},
		{
			name: "[Ok] Without Social Account",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateInfluencerOpts()
				opts.SocialAccount = nil
				tt.args.opts = opts
			},
			validate: func(t *testing.T, tt *TC, got *schema.CreateInfluencerResp) {
				assert.Equal(t, tt.args.opts.Name, got.Name)
				assert.Equal(t, tt.args.opts.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.Equal(t, tt.args.opts.Bio, got.Bio)
				assert.Equal(t, tt.args.opts.ExternalLinks, got.ExternalLinks)
				assert.Nil(t, got.SocialAccount)
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 3*time.Second)

				var doc model.Influencer
				err := tt.fields.DB.Collection(model.InfluencerColl).FindOne(context.TODO(), bson.M{"_id": got.ID}).Decode(&doc)
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ii := &InfluencerImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Influencer = ii
			tt.prepare(&tt)
			got, err := ii.CreateInfluencer(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("InfluencerImpl.CreateInfluencer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestInfluencerImpl_EditInfluencer(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreateInfluencerOpts
		createResp *schema.CreateInfluencerResp
		opts       *schema.EditInfluencerOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.EditInfluencerResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.EditInfluencerResp)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.InfluencerConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateInfluencerOpts()
				resp, _ := tt.fields.App.Influencer.CreateInfluencer(opts)

				var want model.Influencer
				tt.fields.DB.Collection(model.InfluencerColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&want)

				tt.args.opts = &schema.EditInfluencerOpts{
					ID:   want.ID,
					Name: faker.Name().Name(),
					Bio:  faker.Lorem().Sentence(2),
					CoverImg: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
					ExternalLinks: []string{faker.Internet().Url()},
				}
				tt.want = &schema.EditInfluencerResp{
					ID:   want.ID,
					Name: tt.args.opts.Name,
					CoverImg: &model.IMG{
						SRC:    tt.args.opts.CoverImg.SRC,
						Width:  50,
						Height: 50,
					},
					ProfileImage: &model.IMG{
						SRC:    tt.args.opts.ProfileImage.SRC,
						Width:  50,
						Height: 50,
					},
					SocialAccount: resp.SocialAccount,
					Bio:           tt.args.opts.Bio,
					ExternalLinks: tt.args.opts.ExternalLinks,
					CreatedAt:     resp.CreatedAt,
					UpdatedAt:     time.Now().UTC(),
				}
			},
			validate: func(t *testing.T, tt *TC, got *schema.EditInfluencerResp) {
				assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, 100*time.Millisecond)
				assert.WithinDuration(t, tt.want.CreatedAt, got.CreatedAt, 100*time.Millisecond)
				got.UpdatedAt = time.Time{}
				tt.want.UpdatedAt = time.Time{}
				got.CreatedAt = time.Time{}
				tt.want.CreatedAt = time.Time{}
				assert.Equal(t, tt.want, got)
			},
		},
		{
			name: "[Error] No fields",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.InfluencerConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateInfluencerOpts()
				resp, _ := tt.fields.App.Influencer.CreateInfluencer(opts)

				var want model.Influencer
				tt.fields.DB.Collection(model.InfluencerColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&want)

				tt.args.opts = &schema.EditInfluencerOpts{
					ID: want.ID,
				}
				tt.err = errors.New("no fields found to update")
			},
			wantErr: true,
		},
		{
			name: "[Ok] With social account when social account does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateInfluencerOpts()
				opts.SocialAccount = nil
				resp, _ := tt.fields.App.Influencer.CreateInfluencer(opts)
				tt.args.createOpts = opts
				tt.args.createResp = resp
				tt.args.opts = &schema.EditInfluencerOpts{
					ID: resp.ID,
					SocialAccount: &schema.SocialAccountOpts{
						Facebook: &schema.SocialMediaOpts{
							FollowersCount: faker.RandomInt(0, 10000),
						},
						Youtube: &schema.SocialMediaOpts{
							FollowersCount: faker.RandomInt(0, 10000),
						},
					},
				}
			},
			validate: func(t *testing.T, tt *TC, got *schema.EditInfluencerResp) {
				assert.Equal(t, tt.args.createResp.Name, got.Name)
				assert.Equal(t, tt.args.createResp.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.NotNil(t, got.SocialAccount)
				assert.Equal(t, uint(tt.args.opts.SocialAccount.Facebook.FollowersCount), got.SocialAccount.Facebook.FollowersCount)
				assert.Equal(t, uint(tt.args.opts.SocialAccount.Youtube.FollowersCount), got.SocialAccount.Youtube.FollowersCount)
				assert.Nil(t, got.SocialAccount.Instagram)
				assert.Nil(t, got.SocialAccount.Twitter)
				assert.Equal(t, tt.args.createResp.Bio, got.Bio)
				assert.Equal(t, tt.args.createResp.ExternalLinks, got.ExternalLinks)
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 3*time.Second)
				assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, 4*time.Second)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ii := &InfluencerImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Influencer = ii
			tt.prepare(&tt)
			got, err := ii.EditInfluencer(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("InfluencerImpl.EditInfluencer() error = %v, wantErr %v", err, tt.wantErr)
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

func TestInfluencerImpl_GetInfluencersByID(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		ids []primitive.ObjectID
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     []schema.GetInfluencerResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, []schema.GetInfluencerResp)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateInfluencerOpts()
				resp, _ := tt.fields.App.Influencer.CreateInfluencer(opts)

				opts1 := schema.GetRandomCreateInfluencerOpts()
				_, _ = tt.fields.App.Influencer.CreateInfluencer(opts1)

				opts2 := schema.GetRandomCreateInfluencerOpts()
				resp2, _ := tt.fields.App.Influencer.CreateInfluencer(opts2)

				var want []model.Influencer
				cur, _ := tt.fields.DB.Collection(model.InfluencerColl).Find(context.TODO(), bson.M{"_id": bson.M{"$in": bson.A{resp.ID, resp2.ID}}})
				cur.All(context.TODO(), &want)

				tt.args.ids = []primitive.ObjectID{resp.ID, resp2.ID}

				tt.want = []schema.GetInfluencerResp{
					{
						ID:            want[0].ID,
						Name:          want[0].Name,
						CoverImg:      want[0].CoverImg,
						ProfileImage:  want[0].ProfileImage,
						SocialAccount: want[0].SocialAccount,
						Bio:           want[0].Bio,
						ExternalLinks: want[0].ExternalLinks,
					},
					{
						ID:            want[1].ID,
						Name:          want[1].Name,
						CoverImg:      want[1].CoverImg,
						ProfileImage:  want[1].ProfileImage,
						SocialAccount: want[1].SocialAccount,
						Bio:           want[1].Bio,
						ExternalLinks: want[1].ExternalLinks,
					},
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetInfluencerResp) {
				assert.Equal(t, tt.want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ii := &InfluencerImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Influencer = ii
			tt.prepare(&tt)
			got, err := ii.GetInfluencersByID(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("InfluencerImpl.GetInfluencersByID() error = %v, wantErr %v", err, tt.wantErr)
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
