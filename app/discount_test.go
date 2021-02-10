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

func TestDiscountImpl_validateCreateDiscount(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		ctx  context.Context
		opts *schema.CreateDiscountOpts
	}
	type TC struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
		prepare func(*TC)
	}
	tests := []TC{
		{
			name: "[Error] Invalid saleID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID:  primitive.NewObjectIDFromTimestamp(time.Now()),
					SaleID:     primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:       model.FlatOffType,
					Value:      100,
				}
				tt.err = errors.Errorf("sale with id:%s not found", tt.args.opts.SaleID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Error] Discount exists for a catalog from x1 datetime to x2 datetime",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				createDiscount := model.Discount{
					CatalogID: primitive.NewObjectIDFromTimestamp(time.Now()),
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:        model.FlatOffType,
					Value:       100,
					IsActive:    true,
					ValidAfter:  validAfter.UTC(),
					ValidBefore: validBefore.UTC(),
					CreatedAt:   time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.DiscountColl).InsertOne(context.TODO(), createDiscount)

				validBeforeOpts, _ := time.Parse(time.RFC3339, "2021-02-14T00:00:00+00:00")
				validAfterOpts, _ := time.Parse(time.RFC3339, "2021-02-12T00:00:00+00:00")

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID: createDiscount.CatalogID,
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  createDiscount.VariantsID,
					Type:        model.FlatOffType,
					Value:       200,
					ValidAfter:  validAfterOpts,
					ValidBefore: validBeforeOpts,
				}

				tt.err = errors.Errorf("discount from %s to %s already exists for catalog id: %s", validAfterOpts.String(), validBeforeOpts.String(), tt.args.opts.CatalogID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Error] Discount exists for a catalog from x1 datetime to x2 datetime",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				createDiscount := model.Discount{
					CatalogID: primitive.NewObjectIDFromTimestamp(time.Now()),
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:        model.FlatOffType,
					Value:       100,
					IsActive:    true,
					ValidAfter:  validAfter.UTC(),
					ValidBefore: validBefore.UTC(),
					CreatedAt:   time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.DiscountColl).InsertOne(context.TODO(), createDiscount)

				validBeforeOpts, _ := time.Parse(time.RFC3339, "2021-02-19T00:00:00+00:00")
				validAfterOpts, _ := time.Parse(time.RFC3339, "2021-02-12T00:00:00+00:00")

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID: createDiscount.CatalogID,
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  createDiscount.VariantsID,
					Type:        model.FlatOffType,
					Value:       200,
					ValidAfter:  validAfterOpts,
					ValidBefore: validBeforeOpts,
				}

				tt.err = errors.Errorf("discount from %s to %s already exists for catalog id: %s", validAfterOpts.String(), validBeforeOpts.String(), tt.args.opts.CatalogID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Error] Discount exists for a catalog from x1 datetime to x2 datetime",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				createDiscount := model.Discount{
					CatalogID: primitive.NewObjectIDFromTimestamp(time.Now()),
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:        model.FlatOffType,
					Value:       100,
					IsActive:    true,
					ValidAfter:  validAfter.UTC(),
					ValidBefore: validBefore.UTC(),
					CreatedAt:   time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.DiscountColl).InsertOne(context.TODO(), createDiscount)

				validBeforeOpts, _ := time.Parse(time.RFC3339, "2021-02-19T00:00:00+00:00")
				validAfterOpts, _ := time.Parse(time.RFC3339, "2021-02-06T00:00:00+00:00")

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID: createDiscount.CatalogID,
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  createDiscount.VariantsID,
					Type:        model.FlatOffType,
					Value:       200,
					ValidAfter:  validAfterOpts,
					ValidBefore: validBeforeOpts,
				}

				tt.err = errors.Errorf("discount from %s to %s already exists for catalog id: %s", validAfterOpts.String(), validBeforeOpts.String(), tt.args.opts.CatalogID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Error] Discount exists for a catalog from x1 datetime to x2 datetime",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				createDiscount := model.Discount{
					CatalogID: primitive.NewObjectIDFromTimestamp(time.Now()),
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:        model.FlatOffType,
					Value:       100,
					IsActive:    true,
					ValidAfter:  validAfter.UTC(),
					ValidBefore: validBefore.UTC(),
					CreatedAt:   time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.DiscountColl).InsertOne(context.TODO(), createDiscount)

				validBeforeOpts, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfterOpts, _ := time.Parse(time.RFC3339, "2021-02-06T00:00:00+00:00")

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID: createDiscount.CatalogID,
					// SaleID:      primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  createDiscount.VariantsID,
					Type:        model.FlatOffType,
					Value:       200,
					ValidAfter:  validAfterOpts,
					ValidBefore: validBeforeOpts,
				}

				tt.err = errors.Errorf("discount from %s to %s already exists for catalog id: %s", validAfterOpts.String(), validBeforeOpts.String(), tt.args.opts.CatalogID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Ok] 2nd-Discount 1 Sec After 1st-Discount Expiration",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				createDiscount := model.Discount{
					CatalogID:   primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID:  []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:        model.FlatOffType,
					Value:       100,
					IsActive:    true,
					ValidAfter:  validAfter.UTC(),
					ValidBefore: validBefore.UTC(),
					CreatedAt:   time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.DiscountColl).InsertOne(context.TODO(), createDiscount)

				validBeforeOpts, _ := time.Parse(time.RFC3339, "2021-02-20T00:00:00+00:00")
				validAfterOpts, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:01+00:00")

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID:   createDiscount.CatalogID,
					VariantsID:  createDiscount.VariantsID,
					Type:        model.FlatOffType,
					Value:       200,
					ValidAfter:  validAfterOpts,
					ValidBefore: validBeforeOpts,
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := &DiscountImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Discount = di
			tt.prepare(&tt)

			err := di.validateCreateDiscount(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscountImpl.validateCreateDiscount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestDiscountImpl_CreateDiscount(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateDiscountOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.CreateDiscountResp)
	}

	tests := []TC{
		{
			name: "[Ok] FixedType | Without SaleID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID:  primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:       model.FlatOffType,
					Value:      100,
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateDiscountResp) {
				assert.False(t, resp.ID.IsZero())
				assert.WithinDuration(t, time.Now().UTC(), resp.CreatedAt, 100*time.Millisecond)
				assert.True(t, resp.IsActive)
				assert.True(t, resp.SaleID.IsZero())
				assert.Equal(t, tt.args.opts.CatalogID, resp.CatalogID)
				assert.Equal(t, tt.args.opts.VariantsID, resp.VariantsID)
				assert.Equal(t, tt.args.opts.ValidAfter, resp.ValidAfter)
				assert.Equal(t, tt.args.opts.ValidBefore, resp.ValidBefore)
				assert.Equal(t, tt.args.opts.Value, resp.Value)
				assert.Equal(t, tt.args.opts.Type, resp.Type)
				assert.Equal(t, tt.args.opts.MaxValue, resp.MaxValue)
			},
		},
		{
			name: "[Ok] FixedType | SaleID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")

				createSale := model.Sale{
					Name:        "test",
					ValidAfter:  validAfter,
					ValidBefore: validBefore,
					Slug:        "slug",
					Banner: model.IMG{
						SRC:    faker.Avatar().Url("png", 400, 400),
						Width:  400,
						Height: 400,
					},
					CreatedAt: time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.SaleColl).InsertOne(context.TODO(), createSale)

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID:  primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:       model.FlatOffType,
					SaleID:     createSale.ID,
					Value:      100,
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateDiscountResp) {
				assert.False(t, resp.ID.IsZero())
				assert.WithinDuration(t, time.Now().UTC(), resp.CreatedAt, 100*time.Millisecond)
				assert.True(t, resp.IsActive)
				assert.True(t, resp.SaleID.IsZero())
				assert.Equal(t, tt.args.opts.CatalogID, resp.CatalogID)
				assert.Equal(t, tt.args.opts.VariantsID, resp.VariantsID)
				assert.Equal(t, tt.args.opts.ValidAfter, resp.ValidAfter)
				assert.Equal(t, tt.args.opts.ValidBefore, resp.ValidBefore)
				assert.Equal(t, tt.args.opts.Value, resp.Value)
				assert.Equal(t, tt.args.opts.Type, resp.Type)
				assert.Equal(t, tt.args.opts.MaxValue, resp.MaxValue)
			},
		},
		{
			name: "[Ok] PercentType | SaleID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")

				createSale := model.Sale{
					Name:        "test",
					ValidAfter:  validAfter,
					ValidBefore: validBefore,
					Slug:        "slug",
					Banner: model.IMG{
						SRC:    faker.Avatar().Url("png", 400, 400),
						Width:  400,
						Height: 400,
					},
					CreatedAt: time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.SaleColl).InsertOne(context.TODO(), createSale)

				tt.args.opts = &schema.CreateDiscountOpts{
					CatalogID:  primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:       model.FlatOffType,
					SaleID:     createSale.ID,
					Value:      100,
					MaxValue:   500,
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateDiscountResp) {
				assert.False(t, resp.ID.IsZero())
				assert.WithinDuration(t, time.Now().UTC(), resp.CreatedAt, 100*time.Millisecond)
				assert.True(t, resp.IsActive)
				assert.True(t, resp.SaleID.IsZero())
				assert.Equal(t, tt.args.opts.CatalogID, resp.CatalogID)
				assert.Equal(t, tt.args.opts.VariantsID, resp.VariantsID)
				assert.Equal(t, tt.args.opts.ValidAfter, resp.ValidAfter)
				assert.Equal(t, tt.args.opts.ValidBefore, resp.ValidBefore)
				assert.Equal(t, tt.args.opts.Value, resp.Value)
				assert.Equal(t, tt.args.opts.Type, resp.Type)
				assert.Equal(t, tt.args.opts.MaxValue, resp.MaxValue)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := &DiscountImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Discount = di
			tt.prepare(&tt)
			got, err := di.CreateDiscount(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscountImpl.CreateDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestDiscountImpl_DeactivateDiscount(t *testing.T) {
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
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				opts := &schema.CreateDiscountOpts{
					CatalogID:  primitive.NewObjectIDFromTimestamp(time.Now()),
					VariantsID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
					Type:       model.FlatOffType,
					Value:      100,
				}
				resp, _ := tt.fields.App.Discount.CreateDiscount(opts)
				tt.args.id = resp.ID
			},
			validate: func(t *testing.T, tt *TC) {
				var d model.Discount
				err := tt.fields.DB.Collection(model.DiscountColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&d)
				assert.Nil(t, err)
				assert.False(t, d.IsActive)
			},
		},
		{
			name: "[Error] Invalid Discount ID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("discount id: %s not found", tt.args.id.Hex())
			},
			validate: func(t *testing.T, tt *TC) {
				var d model.Discount
				err := tt.fields.DB.Collection(model.DiscountColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&d)
				assert.Nil(t, err)
				assert.False(t, d.IsActive)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := &DiscountImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Discount = di
			tt.prepare(&tt)
			err := di.DeactivateDiscount(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscountImpl.DeactivateDiscount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				tt.validate(t, &tt)
			}
		})
	}
}

func TestDiscountImpl_CreateSale(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateSaleOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.CreateSaleResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				opts := schema.CreateSaleOpts{
					Name: faker.Commerce().Department(),
					Banner: schema.Img{
						SRC: faker.Avatar().Url(faker.RandomChoice([]string{"png", "jpg", "jpeg"}), 200, 200),
					},
					ValidAfter:  validAfter,
					ValidBefore: validBefore,
				}
				tt.args.opts = &opts
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateSaleResp) {
				assert.False(t, resp.ID.IsZero())
				assert.WithinDuration(t, time.Now().UTC(), resp.CreatedAt, 2*time.Second)
				assert.Equal(t, tt.args.opts.Name, resp.Name)
				assert.NotEmpty(t, resp.Slug)
				assert.Equal(t, tt.args.opts.Banner.SRC, resp.Banner.SRC)
				assert.Equal(t, 200, resp.Banner.Width)
				assert.Equal(t, 200, resp.Banner.Height)
			},
		},
		{
			name: "[Error] Invalid Banner Image",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.DiscountConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				validBefore, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				validAfter, _ := time.Parse(time.RFC3339, "2021-02-07T00:00:00+00:00")
				opts := schema.CreateSaleOpts{
					Name: faker.Commerce().Department(),
					Banner: schema.Img{
						SRC: "png@gmail.com",
					},
					ValidAfter:  validAfter,
					ValidBefore: validBefore,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("failed to load banner image: Get \"%s\": unsupported protocol scheme \"\"", opts.Banner.SRC)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := &DiscountImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Discount = di
			tt.prepare(&tt)
			got, err := di.CreateSale(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscountImpl.CreateSale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
