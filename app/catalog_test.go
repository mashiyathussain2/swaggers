package app

import (
	"context"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

func validateCreateCatalogResp(t *testing.T, opts *schema.CreateCatalogOpts, resp *schema.CreateCatalogResp) {
	assert.False(t, resp.ID.IsZero())
	assert.Equal(t, opts.Name, resp.Name)
	assert.Equal(t, opts.BrandID, resp.BrandID)
	assert.Equal(t, len(opts.CategoryID), len(resp.Paths))
	assert.Equal(t, opts.Description, resp.Description)
	assert.Equal(t, opts.Keywords, resp.Keywords)
	assert.Equal(t, opts.HSNCode, resp.HSNCode)
	if opts.ETA != nil {
		// assert.Equal(t, opts.ETA, resp.ETA)
	} else {
		assert.Nil(t, resp.ETA)
	}
	assert.Equal(t, len(opts.Specifications), len(resp.Specifications))
	for i, spec := range opts.Specifications {
		assert.Equal(t, spec.Name, resp.Specifications[i].Name)
		assert.Equal(t, spec.Value, resp.Specifications[i].Value)
	}

	for i, attr := range opts.FilterAttribute {
		assert.Equal(t, attr.Name, resp.FilterAttribute[i].Name)
		assert.Equal(t, attr.Value, resp.FilterAttribute[i].Value)
	}

	assert.Equal(t, opts.BasePrice, uint32(resp.BasePrice.Value))
	assert.Equal(t, opts.RetailPrice, uint32(resp.RetailPrice.Value))
	assert.Equal(t, model.INR, resp.RetailPrice.CurrencyISO)
	assert.Equal(t, model.INR, resp.BasePrice.CurrencyISO)

	if opts.Variants != nil {
		validateCreateVariantResp(t, opts.Variants, resp)
	}
}

func validateCreateVariantResp(t *testing.T, opts []schema.CreateVariantOpts, resp *schema.CreateCatalogResp) {
	assert.Equal(t, len(opts), len(resp.Variants))
	for i, variantOpt := range opts {
		assert.Equal(t, variantOpt.SKU, resp.Variants[i].SKU)
	}
}

func TestKeeperCatalogImpl_CreateCatalog(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	opts := schema.GetRandomCreateCatalogOpts()
	optsWithVariants := schema.GetRandomCreateCatalogOpts()
	optsWithVariants.VariantType = model.SizeType
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		optsWithVariants.Variants = append(optsWithVariants.Variants, *schema.GetRandomCreateVariantOpts())
	}

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateCatalogOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockBrand, *mock.MockCategory)
		validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			validator: func(t *testing.T, s1 *schema.CreateCatalogOpts, s2 *schema.CreateCatalogResp) {
				validateCreateCatalogResp(t, s1, s2)
			},
			buildStubs: func(tt *TC, b *mock.MockBrand, c *mock.MockCategory) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), opts.BrandID).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range tt.args.opts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := c.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
			},
		},
		{
			name: "[Ok] With Variants",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: optsWithVariants,
			},
			validator: func(t *testing.T, s1 *schema.CreateCatalogOpts, s2 *schema.CreateCatalogResp) {
				validateCreateCatalogResp(t, s1, s2)
			},
			buildStubs: func(tt *TC, b *mock.MockBrand, c *mock.MockCategory) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), optsWithVariants.BrandID).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range tt.args.opts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := c.EXPECT().GetCategoryPath(id).Times(1).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
			},
		},
		{
			name: "[Error] When brandID does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: optsWithVariants,
			},
			wantErr: true,
			err:     errors.Errorf("brand id %s does not exists", optsWithVariants.BrandID.Hex()),
			buildStubs: func(tt *TC, b *mock.MockBrand, c *mock.MockCategory) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), optsWithVariants.BrandID).Times(1).Return(false, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KeeperCatalogImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBrand := mock.NewMockBrand(ctrl)
			mockCategory := mock.NewMockCategory(ctrl)
			kc.App.Brand = mockBrand
			kc.App.Category = mockCategory
			tt.buildStubs(&tt, mockBrand, mockCategory)

			resp, err := kc.CreateCatalog(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				tt.validator(t, tt.args.opts, resp)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestKeeperCatalogImpl_EditCatalog(t *testing.T) {

	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreateCatalogOpts
		createResp *schema.CreateCatalogResp
		opts       *schema.EditCatalogOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.EditCatalogResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)
		validate   func(*testing.T, TC, *schema.EditCatalogResp)
	}

	tests := []TC{
		{
			name: "[Ok] Editing Name",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)

				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp
			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.EditCatalogOpts{
					ID:   resp.ID,
					Name: resp.Name + " Edited",
				}
				tt.args.opts = &opts

				want := schema.EditCatalogResp{
					ID:              resp.ID,
					Name:            opts.Name,
					Paths:           resp.Paths,
					Description:     resp.Description,
					Keywords:        resp.Keywords,
					Specifications:  resp.Specifications,
					FilterAttribute: resp.FilterAttribute,
					HSNCode:         resp.HSNCode,
					BasePrice:       resp.BasePrice,
					RetailPrice:     resp.RetailPrice,
					ETA:             resp.ETA,
				}
				tt.want = &want
			},
			validate: func(t *testing.T, tt TC, resp *schema.EditCatalogResp) {
				assert.Equal(t, tt.want.ID, resp.ID)
				assert.NotEqual(t, tt.args.createOpts.Name, resp.Name)
				assert.Equal(t, tt.want.Name, resp.Name)
				assert.Equal(t, tt.want.Paths, resp.Paths)
				assert.Equal(t, tt.want.Description, resp.Description)
				assert.Equal(t, tt.want.Keywords, resp.Keywords)
				assert.Equal(t, tt.want.Specifications, resp.Specifications)
				assert.Equal(t, tt.want.FilterAttribute, resp.FilterAttribute)
				assert.Equal(t, tt.want.HSNCode, resp.HSNCode)
				assert.Equal(t, tt.want.BasePrice, resp.BasePrice)
				assert.Equal(t, tt.want.RetailPrice, resp.RetailPrice)
				assert.Equal(t, tt.want.ETA, resp.ETA)
				assert.WithinDuration(t, time.Now().UTC(), resp.UpdatedAt, time.Millisecond*100)
			},
		},
		{
			name: "[Ok] Updating Category",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.EditCatalogOpts{
					CategoryID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
				},
			},
			want: &schema.EditCatalogResp{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}

				for _, id := range tt.args.opts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
					tt.want.Paths = append(tt.want.Paths, path)
				}

				gomock.InOrder(categoryCalls...)

				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp
			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				tt.args.opts.ID = resp.ID

				tt.want.ID = resp.ID
				tt.want.Name = resp.Name
				tt.want.Description = resp.Description
				tt.want.Keywords = resp.Keywords
				tt.want.Specifications = resp.Specifications
				tt.want.FilterAttribute = resp.FilterAttribute
				tt.want.HSNCode = resp.HSNCode
				tt.want.BasePrice = resp.BasePrice
				tt.want.RetailPrice = resp.RetailPrice
				tt.want.ETA = resp.ETA

			},
			validate: func(t *testing.T, tt TC, resp *schema.EditCatalogResp) {
				assert.Equal(t, tt.want.ID, resp.ID)
				assert.Equal(t, tt.want.Name, resp.Name)
				assert.Equal(t, tt.want.Paths, resp.Paths)
				assert.NotEqual(t, tt.args.createResp.Paths, resp.Paths)
				assert.Equal(t, tt.want.Description, resp.Description)
				assert.Equal(t, tt.want.Keywords, resp.Keywords)
				assert.Equal(t, tt.want.Specifications, resp.Specifications)
				assert.Equal(t, tt.want.FilterAttribute, resp.FilterAttribute)
				assert.Equal(t, tt.want.HSNCode, resp.HSNCode)
				assert.Equal(t, tt.want.BasePrice, resp.BasePrice)
				assert.Equal(t, tt.want.RetailPrice, resp.RetailPrice)
				assert.Equal(t, tt.want.ETA, resp.ETA)
				assert.WithinDuration(t, time.Now().UTC(), resp.UpdatedAt, time.Millisecond*100)
			},
		},
		{
			name: "[Ok] Updating all fields except category",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)

				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp
			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp

				tt.args.opts = schema.GetRandomEditCatalogOpts(resp.ID)

				want := schema.EditCatalogResp{
					ID:          resp.ID,
					Name:        tt.args.opts.Name,
					Paths:       resp.Paths,
					Description: tt.args.opts.Description,
					Keywords:    tt.args.opts.Keywords,
					HSNCode:     tt.args.opts.HSNCode,
				}
				for _, spec := range tt.args.opts.Specifications {
					want.Specifications = append(want.Specifications, model.Specification{Name: spec.Name, Value: spec.Value})
				}
				for _, attr := range tt.args.opts.FilterAttribute {
					want.FilterAttribute = append(want.FilterAttribute, model.Attribute{Name: attr.Name, Value: attr.Value})
				}
				want.ETA = &model.ETA{
					Min:  int(tt.args.opts.ETA.Min),
					Max:  int(tt.args.opts.ETA.Max),
					Unit: tt.args.opts.ETA.Unit,
				}
				want.BasePrice = *model.SetINRPrice(float32(tt.args.opts.BasePrice))
				want.RetailPrice = *model.SetINRPrice(float32(tt.args.opts.RetailPrice))
				tt.want = &want
			},
			validate: func(t *testing.T, tt TC, resp *schema.EditCatalogResp) {
				assert.Equal(t, tt.want.ID, resp.ID)
				assert.NotEqual(t, tt.args.createOpts.Name, resp.Name)
				assert.Equal(t, tt.want.Name, resp.Name)
				assert.Equal(t, tt.want.Paths, resp.Paths)
				assert.Equal(t, tt.want.Description, resp.Description)
				assert.Equal(t, tt.want.Keywords, resp.Keywords)
				assert.Equal(t, tt.want.Specifications, resp.Specifications)
				assert.Equal(t, tt.want.FilterAttribute, resp.FilterAttribute)
				assert.Equal(t, tt.want.HSNCode, resp.HSNCode)
				assert.Equal(t, tt.want.BasePrice, resp.BasePrice)
				assert.Equal(t, tt.want.RetailPrice, resp.RetailPrice)
				assert.Equal(t, tt.want.ETA, resp.ETA)
				assert.WithinDuration(t, time.Now().UTC(), resp.UpdatedAt, time.Millisecond*100)
			},
		},
		{
			name: "[Error] With Invalid ID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)

				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp
			},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomEditCatalogOpts(primitive.NewObjectIDFromTimestamp(time.Now()))
				tt.err = errors.Errorf("catalog with id:%s not found", tt.args.opts.ID.Hex())
			},
		},
		{
			name: "[Error] Updating Category With Invalid CategoryID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.EditCatalogOpts{
					CategoryID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
				},
			},
			wantErr: true,
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}

				for _, id := range tt.args.opts.CategoryID {
					call := ct.EXPECT().GetCategoryPath(id).Return("", errors.Errorf("category with id:%s not found", id.Hex()))
					categoryCalls = append(categoryCalls, call)
					break
				}

				gomock.InOrder(categoryCalls...)

				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp
			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				tt.args.opts.ID = resp.ID
				tt.err = errors.Errorf("category with id:%s not found", tt.args.opts.CategoryID[0].Hex())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			kc := &KeeperCatalogImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.KeeperCatalog = kc
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			got, err := kc.EditCatalog(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.EditCatalog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, tt, got)
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestKeeperCatalogImpl_AddVariant(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		catalogID primitive.ObjectID
		opts      *schema.AddVariantOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.AddVariantResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)
		validate   func(*testing.T, *TC, *schema.AddVariantResp)
	}

	tests := []TC{
		{
			name: "[Ok] Adding 1 Variant When No Variant Exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.catalogID = createCatalogResp.ID

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					SKU:         faker.Lorem().Word(),
					Attribute:   faker.Commerce().Color(),
					VariantType: model.SizeType,
				}

				tt.want = &schema.AddVariantResp{
					SKU:       tt.args.opts.SKU,
					Attribute: tt.args.opts.Attribute,
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.AddVariantResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.want.Attribute, resp.Attribute)
				assert.Equal(t, tt.want.SKU, resp.SKU)

				var dbResp model.Catalog
				err := tt.fields.DB.Collection(model.CatalogColl).FindOne(context.Background(), bson.M{"variant._id": resp.ID}).Decode(&dbResp)
				assert.Nil(t, err)
				assert.Len(t, dbResp.Variants, 1)
				assert.Equal(t, resp.ID, dbResp.Variants[0].ID)
				assert.Equal(t, model.SizeType, dbResp.VariantType)
			},
		},
		{
			name: "[Ok] Adding Variant When Variants Exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.catalogID = createCatalogResp.ID

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					SKU:         faker.Lorem().Word(),
					Attribute:   faker.Commerce().Color(),
					VariantType: model.SizeType,
				}

				tt.want = &schema.AddVariantResp{
					SKU:       tt.args.opts.SKU,
					Attribute: tt.args.opts.Attribute,
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.AddVariantResp) {
				assert.False(t, resp.ID.IsZero())
				assert.Equal(t, tt.want.Attribute, resp.Attribute)
				assert.Equal(t, tt.want.SKU, resp.SKU)

				var dbResp model.Catalog
				err := tt.fields.DB.Collection(model.CatalogColl).FindOne(context.Background(), bson.M{"variant._id": resp.ID}).Decode(&dbResp)
				assert.Nil(t, err)
				assert.Len(t, dbResp.Variants, 2)
				assert.Equal(t, resp.ID, dbResp.Variants[1].ID)
				assert.Equal(t, model.SizeType, dbResp.VariantType)
			},
		},
		{
			name: "[Error] Adding 1 Variant When No Variant Exists With Invalid Type",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.catalogID = createCatalogResp.ID

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					SKU:         faker.Lorem().Word(),
					Attribute:   faker.Commerce().Color(),
					VariantType: model.ColorType,
				}
				tt.err = errors.Errorf("cannot set variant type %s: catalog is set with variant type %s", tt.args.opts.VariantType, model.SizeType)
			},
			wantErr: true,
		},
		{
			name: "[Error] Invalid Catalog ID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				// createCatalogOpts.VariantType = model.SizeType
				// createVariantOpts := schema.GetRandomCreateVariantOpts()
				// createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				// b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				// var categoryCalls []*gomock.Call
				// for _, id := range createCatalogOpts.CategoryID {
				// 	path := schema.GetRandomGetCategoryPath(id)
				// 	call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
				// 	categoryCalls = append(categoryCalls, call)
				// }
				// gomock.InOrder(categoryCalls...)
				// createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.catalogID = primitive.NewObjectIDFromTimestamp(time.Now())

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					SKU:         faker.Lorem().Word(),
					Attribute:   faker.Commerce().Color(),
					VariantType: model.ColorType,
				}
				tt.err = errors.Errorf("catalog with id:%s not found", tt.args.catalogID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Error] Duplicate SKU",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createVariantOpts.SKU = "1"
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.catalogID = createCatalogResp.ID

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.AddVariantOpts{
					SKU:         "1",
					Attribute:   faker.Commerce().Color(),
					VariantType: model.SizeType,
				}
				tt.err = errors.Errorf("variant with sku %s already exists", "1")
			},
			wantErr: true,
		},
		{
			name: "[Error] Duplicate Attribute",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createVariantOpts.Attribute = "red"
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)
				tt.args.catalogID = createCatalogResp.ID
			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					SKU:         faker.Lorem().Word(),
					Attribute:   "red",
					VariantType: model.SizeType,
				}
				tt.err = errors.Errorf("variant with attribute %s already exists", "red")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			kc := &KeeperCatalogImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			mockBrand := mock.NewMockBrand(ctrl)
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Brand = mockBrand
			tt.fields.App.Category = mockCategory
			tt.fields.App.KeeperCatalog = kc
			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			got, err := kc.AddVariant(tt.args.catalogID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.AddVariant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
