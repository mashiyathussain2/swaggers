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

func TestBrandImpl_CreateBrand(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateBrandOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.CreateBrandResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.CreateBrandResp)
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
				opts := schema.GetRandomCreateBrandOpts()
				tt.args.opts = opts
			},
			validate: func(t *testing.T, tt *TC, got *schema.CreateBrandResp) {
				assert.Equal(t, tt.args.opts.Name, got.Name)
				assert.Equal(t, tt.args.opts.RegisteredName, got.RegisteredName)
				assert.Equal(t, tt.args.opts.Domain, got.Domain)
				assert.Equal(t, tt.args.opts.Website, got.Website)
				assert.Equal(t, tt.args.opts.FulfillmentCCEmail, got.FulfillmentCCEmail)
				assert.Equal(t, tt.args.opts.FulfillmentEmail, got.FulfillmentEmail)
				assert.Equal(t, tt.args.opts.Logo.SRC, got.Logo.SRC)
				assert.Equal(t, 200, got.Logo.Width)
				assert.Equal(t, 200, got.Logo.Height)
				assert.Equal(t, tt.args.opts.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.Equal(t, tt.args.opts.Bio, got.Bio)
				if tt.args.opts.SocialAccount != nil {
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Facebook.FollowersCount), got.SocialAccount.Facebook.FollowersCount)
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Instagram.FollowersCount), got.SocialAccount.Instagram.FollowersCount)
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Youtube.FollowersCount), got.SocialAccount.Youtube.FollowersCount)
					assert.Equal(t, uint(tt.args.opts.SocialAccount.Twitter.FollowersCount), got.SocialAccount.Twitter.FollowersCount)
				}
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 3*time.Second)

				var doc model.Brand
				err := tt.fields.DB.Collection(model.BrandColl).FindOne(context.TODO(), bson.M{"_id": got.ID}).Decode(&doc)
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
				opts := schema.GetRandomCreateBrandOpts()
				opts.SocialAccount = nil
				tt.args.opts = opts
			},
			validate: func(t *testing.T, tt *TC, got *schema.CreateBrandResp) {
				assert.Equal(t, tt.args.opts.Name, got.Name)
				assert.Equal(t, tt.args.opts.RegisteredName, got.RegisteredName)
				assert.Equal(t, tt.args.opts.Domain, got.Domain)
				assert.Equal(t, tt.args.opts.Website, got.Website)
				assert.Equal(t, tt.args.opts.FulfillmentCCEmail, got.FulfillmentCCEmail)
				assert.Equal(t, tt.args.opts.FulfillmentEmail, got.FulfillmentEmail)
				assert.Equal(t, tt.args.opts.Logo.SRC, got.Logo.SRC)
				assert.Equal(t, 200, got.Logo.Width)
				assert.Equal(t, 200, got.Logo.Height)
				assert.Equal(t, tt.args.opts.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.Equal(t, tt.args.opts.Bio, got.Bio)
				assert.Nil(t, got.SocialAccount)
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 3*time.Second)

				var doc model.Brand
				err := tt.fields.DB.Collection(model.BrandColl).FindOne(context.TODO(), bson.M{"_id": got.ID}).Decode(&doc)
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = bi
			tt.prepare(&tt)
			got, err := bi.CreateBrand(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrandImpl.CreateBrand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestBrandImpl_EditBrand(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreateBrandOpts
		createResp *schema.CreateBrandResp
		opts       *schema.EditBrandOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.EditBrandResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.EditBrandResp)
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
				opts := schema.GetRandomCreateBrandOpts()
				resp, _ := tt.fields.App.Brand.CreateBrand(opts)
				tt.args.createOpts = opts
				tt.args.createResp = resp
				tt.args.opts = &schema.EditBrandOpts{
					ID:               resp.ID,
					Name:             faker.Company().Name(),
					FulfillmentEmail: faker.Internet().SafeEmail(),
					Logo: &schema.Img{
						SRC: faker.Avatar().Url("png", 100, 100),
					},
					Bio: faker.Lorem().Sentence(1),
				}
			},
			validate: func(t *testing.T, tt *TC, got *schema.EditBrandResp) {
				assert.Equal(t, tt.args.opts.Name, got.Name)
				assert.Equal(t, tt.args.createResp.RegisteredName, got.RegisteredName)
				assert.Equal(t, tt.args.createResp.Domain, got.Domain)
				assert.Equal(t, tt.args.createResp.Website, got.Website)
				assert.Equal(t, tt.args.createResp.FulfillmentCCEmail, got.FulfillmentCCEmail)
				assert.Equal(t, tt.args.opts.FulfillmentEmail, got.FulfillmentEmail)
				assert.Equal(t, tt.args.opts.Logo.SRC, got.Logo.SRC)
				assert.Equal(t, 100, got.Logo.Width)
				assert.Equal(t, 100, got.Logo.Height)
				assert.Equal(t, tt.args.createResp.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.Equal(t, tt.args.opts.Bio, got.Bio)
				if tt.args.createResp.SocialAccount != nil {
					assert.Equal(t, uint(tt.args.createResp.SocialAccount.Facebook.FollowersCount), got.SocialAccount.Facebook.FollowersCount)
					assert.Equal(t, uint(tt.args.createResp.SocialAccount.Instagram.FollowersCount), got.SocialAccount.Instagram.FollowersCount)
					assert.Equal(t, uint(tt.args.createResp.SocialAccount.Youtube.FollowersCount), got.SocialAccount.Youtube.FollowersCount)
					assert.Equal(t, uint(tt.args.createResp.SocialAccount.Twitter.FollowersCount), got.SocialAccount.Twitter.FollowersCount)
				}
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 4*time.Second)
				assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, 4*time.Second)
			},
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
				opts := schema.GetRandomCreateBrandOpts()
				opts.SocialAccount = nil
				resp, _ := tt.fields.App.Brand.CreateBrand(opts)
				tt.args.createOpts = opts
				tt.args.createResp = resp
				tt.args.opts = &schema.EditBrandOpts{
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
			validate: func(t *testing.T, tt *TC, got *schema.EditBrandResp) {
				assert.Equal(t, tt.args.createResp.Name, got.Name)
				assert.Equal(t, tt.args.createResp.RegisteredName, got.RegisteredName)
				assert.Equal(t, tt.args.createResp.Domain, got.Domain)
				assert.Equal(t, tt.args.createResp.Website, got.Website)
				assert.Equal(t, tt.args.createResp.FulfillmentCCEmail, got.FulfillmentCCEmail)
				assert.Equal(t, tt.args.createResp.FulfillmentEmail, got.FulfillmentEmail)
				assert.Equal(t, tt.args.createResp.Logo.SRC, got.Logo.SRC)
				assert.Equal(t, 200, got.Logo.Width)
				assert.Equal(t, 200, got.Logo.Height)
				assert.Equal(t, tt.args.createResp.CoverImg.SRC, got.CoverImg.SRC)
				assert.Equal(t, 200, got.CoverImg.Width)
				assert.Equal(t, 400, got.CoverImg.Height)
				assert.NotNil(t, got.SocialAccount)
				assert.Equal(t, uint(tt.args.opts.SocialAccount.Facebook.FollowersCount), got.SocialAccount.Facebook.FollowersCount)
				assert.Equal(t, uint(tt.args.opts.SocialAccount.Youtube.FollowersCount), got.SocialAccount.Youtube.FollowersCount)
				assert.Nil(t, got.SocialAccount.Instagram)
				assert.Nil(t, got.SocialAccount.Twitter)
				assert.Equal(t, tt.args.createResp.Bio, got.Bio)
				assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, 3*time.Second)
				assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, 4*time.Second)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = bi
			tt.prepare(&tt)
			got, err := bi.EditBrand(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrandImpl.EditBrand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestBrandImpl_GetBrandByID(t *testing.T) {
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
		want     *schema.GetBrandResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.GetBrandResp)
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
				opts := schema.GetRandomCreateBrandOpts()
				resp, _ := tt.fields.App.Brand.CreateBrand(opts)

				var want model.Brand
				tt.fields.DB.Collection(model.BrandColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&want)

				tt.args.id = resp.ID
				tt.want = &schema.GetBrandResp{
					ID:                 want.ID,
					Name:               want.Name,
					LName:              want.LName,
					RegisteredName:     want.RegisteredName,
					FulfillmentEmail:   want.FulfillmentEmail,
					FulfillmentCCEmail: want.FulfillmentCCEmail,
					Domain:             want.Domain,
					Website:            want.Website,
					Logo:               want.Logo,
					CoverImg:           want.CoverImg,
					SocialAccount:      want.SocialAccount,
				}
			},
			validate: func(t *testing.T, tt *TC, got *schema.GetBrandResp) {
				assert.Equal(t, tt.want, got)
			},
		},
		{
			name: "[Error] brand id does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateBrandOpts()
				resp, _ := tt.fields.App.Brand.CreateBrand(opts)

				var want model.Brand
				tt.fields.DB.Collection(model.BrandColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&want)

				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("brand with id:%s not found: mongo: no documents in result", tt.args.id.Hex())
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = bi
			tt.prepare(&tt)
			got, err := bi.GetBrandByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrandImpl.GetBrandByID() error = %v, wantErr %v", err, tt.wantErr)
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

func TestBrandImpl_BrandByIDCheck(t *testing.T) {
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
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		err     error
		prepare func(*TC)
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
				opts := schema.GetRandomCreateBrandOpts()
				resp, _ := tt.fields.App.Brand.CreateBrand(opts)
				tt.args.id = resp.ID
				tt.want = true
			},
		},
		{
			name: "[Ok] brand id does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GetRandomCreateBrandOpts()
				_, _ = tt.fields.App.Brand.CreateBrand(opts)

				tt.args.id = primitive.NewObjectID()
				tt.want = false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = bi
			tt.prepare(&tt)
			got, err := bi.CheckBrandByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrandImpl.BrandByIDCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
			if tt.wantErr {
				assert.False(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestBrandImpl_GetBrandsByID(t *testing.T) {
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
		want     []schema.GetBrandResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, []schema.GetBrandResp)
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
				opts := schema.GetRandomCreateBrandOpts()
				resp, _ := tt.fields.App.Brand.CreateBrand(opts)

				opts1 := schema.GetRandomCreateBrandOpts()
				_, _ = tt.fields.App.Brand.CreateBrand(opts1)

				opts2 := schema.GetRandomCreateBrandOpts()
				resp2, _ := tt.fields.App.Brand.CreateBrand(opts2)

				var want []model.Brand
				cur, _ := tt.fields.DB.Collection(model.BrandColl).Find(context.TODO(), bson.M{"_id": bson.M{"$in": bson.A{resp.ID, resp2.ID}}})
				cur.All(context.TODO(), &want)

				tt.args.ids = []primitive.ObjectID{resp.ID, resp2.ID}

				tt.want = []schema.GetBrandResp{
					{
						ID:                 want[0].ID,
						Name:               want[0].Name,
						LName:              want[0].LName,
						RegisteredName:     want[0].RegisteredName,
						FulfillmentEmail:   want[0].FulfillmentEmail,
						FulfillmentCCEmail: want[0].FulfillmentCCEmail,
						Domain:             want[0].Domain,
						Website:            want[0].Website,
						Logo:               want[0].Logo,
						CoverImg:           want[0].CoverImg,
						Bio:                want[0].Bio,
						SocialAccount:      want[0].SocialAccount,
					},
					{
						ID:                 want[1].ID,
						Name:               want[1].Name,
						LName:              want[1].LName,
						RegisteredName:     want[1].RegisteredName,
						FulfillmentEmail:   want[1].FulfillmentEmail,
						FulfillmentCCEmail: want[1].FulfillmentCCEmail,
						Domain:             want[1].Domain,
						Website:            want[1].Website,
						Logo:               want[1].Logo,
						Bio:                want[1].Bio,
						CoverImg:           want[1].CoverImg,
						SocialAccount:      want[1].SocialAccount,
					},
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetBrandResp) {
				assert.Equal(t, tt.want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = bi
			tt.prepare(&tt)
			got, err := bi.GetBrandsByID(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrandImpl.GetBrandsByID() error = %v, wantErr %v", err, tt.wantErr)
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

func TestBrandImpl_GetBrands(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     []schema.GetBrandResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, []schema.GetBrandResp)
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
				opts := schema.GetRandomCreateBrandOpts()
				_, _ = tt.fields.App.Brand.CreateBrand(opts)

				opts1 := schema.GetRandomCreateBrandOpts()
				_, _ = tt.fields.App.Brand.CreateBrand(opts1)

				opts2 := schema.GetRandomCreateBrandOpts()
				_, _ = tt.fields.App.Brand.CreateBrand(opts2)

				var want []model.Brand
				cur, _ := tt.fields.DB.Collection(model.BrandColl).Find(context.TODO(), bson.M{})
				cur.All(context.TODO(), &want)

				tt.want = []schema.GetBrandResp{
					{
						ID:                 want[0].ID,
						Name:               want[0].Name,
						LName:              want[0].LName,
						RegisteredName:     want[0].RegisteredName,
						FulfillmentEmail:   want[0].FulfillmentEmail,
						FulfillmentCCEmail: want[0].FulfillmentCCEmail,
						Domain:             want[0].Domain,
						Website:            want[0].Website,
						Logo:               want[0].Logo,
						CoverImg:           want[0].CoverImg,
						Bio:                want[0].Bio,
						SocialAccount:      want[0].SocialAccount,
					},
					{
						ID:                 want[1].ID,
						Name:               want[1].Name,
						LName:              want[1].LName,
						RegisteredName:     want[1].RegisteredName,
						FulfillmentEmail:   want[1].FulfillmentEmail,
						FulfillmentCCEmail: want[1].FulfillmentCCEmail,
						Domain:             want[1].Domain,
						Website:            want[1].Website,
						Logo:               want[1].Logo,
						Bio:                want[1].Bio,
						CoverImg:           want[1].CoverImg,
						SocialAccount:      want[1].SocialAccount,
					},
					{
						ID:                 want[2].ID,
						Name:               want[2].Name,
						LName:              want[2].LName,
						RegisteredName:     want[2].RegisteredName,
						FulfillmentEmail:   want[2].FulfillmentEmail,
						FulfillmentCCEmail: want[2].FulfillmentCCEmail,
						Domain:             want[2].Domain,
						Website:            want[2].Website,
						Logo:               want[2].Logo,
						Bio:                want[2].Bio,
						CoverImg:           want[2].CoverImg,
						SocialAccount:      want[2].SocialAccount,
					},
				}
			},
			validate: func(t *testing.T, tt *TC, got []schema.GetBrandResp) {
				assert.Equal(t, tt.want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = bi
			tt.prepare(&tt)
			got, err := bi.GetBrands()
			if (err != nil) != tt.wantErr {
				t.Errorf("BrandImpl.GetBrandsByID() error = %v, wantErr %v", err, tt.wantErr)
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
