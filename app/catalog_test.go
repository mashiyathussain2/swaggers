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

func validateDeleteVariant(t *testing.T, opts schema.DeleteVariantOpts, resp model.Catalog) {
	variants := resp.Variants
	for i := 0; i < len(variants); i++ {
		if opts.VariantID.Hex() == variants[i].ID.Hex() {
			assert.True(t, variants[i].IsDeleted)
		}
		// assert.Equal(t, opts.VariantID.Hex(), variants[i].ID.Hex())
	}
	assert.Equal(t, len(variants), 3)
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
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_1"),
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
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_2"),
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
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_3"),
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
		createResp *schema.CreateCatalogResp
		opts       *schema.AddVariantOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.AddVariantResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand, *mock.MockInventory)
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
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, inv *mock.MockInventory) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				inv.EXPECT().CreateInventory(gomock.Any()).Times(1).Return(primitive.NewObjectID(), nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)
				tt.args.createResp = createCatalogResp

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					ID:          tt.args.createResp.ID,
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
				err := tt.fields.DB.Collection(model.CatalogColl).FindOne(context.Background(), bson.M{"_id": tt.args.createResp.ID}).Decode(&dbResp)
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
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, inv *mock.MockInventory) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)

				inv.EXPECT().CreateInventory(gomock.Any()).Times(1).Return(primitive.NewObjectID(), nil)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.createResp = createCatalogResp

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					ID:          tt.args.createResp.ID,
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
				err := tt.fields.DB.Collection(model.CatalogColl).FindOne(context.Background(), bson.M{"_id": tt.args.createResp.ID}).Decode(&dbResp)
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
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, inv *mock.MockInventory) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				// inv.EXPECT().CreateInventory(gomock.Any()).Times(1).Return(primitive.NewObjectID(), nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.createResp = createCatalogResp

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					ID:          tt.args.createResp.ID,
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
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, inv *mock.MockInventory) {
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

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					ID:          primitive.NewObjectID(),
					SKU:         faker.Lorem().Word(),
					Attribute:   faker.Commerce().Color(),
					VariantType: model.ColorType,
				}
				tt.err = errors.Errorf("catalog with id:%s not found", tt.args.opts.ID.Hex())
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
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, inv *mock.MockInventory) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createVariantOpts.SKU = "1"
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				// inv.EXPECT().CreateInventory(gomock.Any()).Times(1).Return(primitive.NewObjectID(), nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)

				tt.args.createResp = createCatalogResp

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.AddVariantOpts{
					ID:          tt.args.createResp.ID,
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
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, inv *mock.MockInventory) {
				createCatalogOpts := schema.GetRandomCreateCatalogOpts()
				createCatalogOpts.VariantType = model.SizeType
				createVariantOpts := schema.GetRandomCreateVariantOpts()
				createVariantOpts.Attribute = "red"
				createCatalogOpts.Variants = append(createCatalogOpts.Variants, *createVariantOpts)
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				// inv.EXPECT().CreateInventory(gomock.Any()).Times(1).Return(primitive.NewObjectID(), nil)
				var categoryCalls []*gomock.Call
				for _, id := range createCatalogOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				createCatalogResp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createCatalogOpts)
				tt.args.createResp = createCatalogResp
			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.AddVariantOpts{
					ID:          tt.args.createResp.ID,
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
			mockInventory := mock.NewMockInventory(ctrl)
			tt.fields.App.Inventory = mockInventory
			tt.fields.App.Brand = mockBrand
			tt.fields.App.Category = mockCategory
			tt.fields.App.KeeperCatalog = kc
			tt.buildStubs(&tt, mockCategory, mockBrand, mockInventory)
			tt.prepare(&tt)

			got, err := kc.AddVariant(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.AddVariant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
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
func TestKeeperCatalogImpl_KeeperSearchCatalog(t *testing.T) {
	t.Parallel()

	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts       *schema.KeeperSearchCatalogOpts
		createOpts *schema.CreateCatalogOpts
		createResp *schema.CreateCatalogResp
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       []schema.KeeperSearchCatalogResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)

		// validator func(*testing.T, *schema.KeeperSearchCatalogOpts, []schema.KeeperSearchCatalogResp)
	}
	tests := []TC{

		// {
		// 	name: "[Ok] Random",
		// 	fields: fields{
		// 		App:    app,
		// 		DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
		// 		Logger: app.Logger,
		// 	},
		// 	args: args{},
		// 	buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
		// 		b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
		// 		createOpts := schema.GetRandomCreateCatalogOpts()
		// 		var categoryCalls []*gomock.Call
		// 		for _, id := range createOpts.CategoryID {
		// 			path := schema.GetRandomGetCategoryPath(id)
		// 			call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
		// 			categoryCalls = append(categoryCalls, call)
		// 		}
		// 		gomock.InOrder(categoryCalls...)
		// 		resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
		// 		tt.args.createOpts = createOpts
		// 		tt.args.createResp = resp
		// 	},
		// 	prepare: func(tt *TC) {
		// 		tt.args.opts = opts
		// 	},
		// 	wantErr: false,
		// },

		{
			name: "[OK] Mock With 1 Catalog",
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
				opts := schema.KeeperSearchCatalogOpts{
					Name: resp.Name,
					Page: 0,
				}
				tt.args.opts = &opts
				want := []schema.KeeperSearchCatalogResp{
					{
						ID:          resp.ID,
						Name:        opts.Name,
						Path:        resp.Paths,
						BasePrice:   resp.BasePrice,
						RetailPrice: resp.RetailPrice,
						Status:      resp.Status,
						Variants:    resp.Variants,
						VariantType: resp.VariantType,
					},
				}
				tt.want = want
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

			resp, err := kc.KeeperSearchCatalog(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)

			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestKeeperCatalogImpl_GetBasicCatalogInfoWithBothFilter(t *testing.T) {

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts   []schema.CreateCatalogOpts
		resp   []schema.CreateCatalogResp
		filter *schema.GetBasicCatalogFilter
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       []schema.GetBasicCatalogResp
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockBrand, *mock.MockCategory)
		validator  func(*testing.T, *TC, []schema.GetBasicCatalogResp)
	}
	tests := []TC{
		{
			name: "[Ok] With Category and BrandID Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_1"),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				randIndex0 := faker.RandomInt(0, len(tt.args.resp)-1)

				strCategoryID := strings.Split(tt.args.resp[randIndex0].Paths[faker.RandomInt(0, len(tt.args.resp[randIndex0].Paths)-1)], "/")
				randCategoryID0, _ := primitive.ObjectIDFromHex(faker.RandomChoice(strCategoryID))

				tt.args.filter = &schema.GetBasicCatalogFilter{
					BrandID:    []primitive.ObjectID{tt.args.resp[randIndex0].BrandID},
					CategoryID: []primitive.ObjectID{randCategoryID0},
				}

				tt.want = []schema.GetBasicCatalogResp{
					{
						ID:          tt.args.resp[randIndex0].ID,
						Paths:       tt.args.resp[randIndex0].Paths,
						Name:        tt.args.resp[randIndex0].Name,
						Description: tt.args.resp[randIndex0].Description,
						RetailPrice: tt.args.resp[randIndex0].RetailPrice,
					},
				}
			},
			buildStubs: func(tt *TC, b *mock.MockBrand, c *mock.MockCategory) {
				var opts []schema.CreateCatalogOpts
				for i := 0; i < faker.RandomInt(4, 20); i++ {
					opt := schema.GetRandomCreateCatalogOpts()
					opts = append(opts, *opt)
				}
				tt.args.opts = opts

				for _, opt := range tt.args.opts {
					call := b.EXPECT().CheckBrandIDExists(gomock.Any(), opt.BrandID).Return(true, nil)
					gomock.InOrder(call)
					var categoryCalls []*gomock.Call
					for _, id := range opt.CategoryID {
						path := schema.GetRandomGetCategoryPath(id)
						call := c.EXPECT().GetCategoryPath(id).Return(path, nil)
						categoryCalls = append(categoryCalls, call)
					}
					gomock.InOrder(categoryCalls...)
				}

				var resp []schema.CreateCatalogResp
				for _, opt := range tt.args.opts {
					res, _ := tt.fields.App.KeeperCatalog.CreateCatalog(&opt)
					resp = append(resp, *res)
				}
				tt.args.resp = resp

			},
			validator: func(t *testing.T, tt *TC, res []schema.GetBasicCatalogResp) {
				assert.Len(t, res, len(tt.want))
				assert.Equal(t, tt.want, res)
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
			tt.fields.App.Brand = mockBrand
			tt.fields.App.Category = mockCategory
			tt.fields.App.KeeperCatalog = kc

			tt.buildStubs(&tt, mockBrand, mockCategory)
			tt.prepare(&tt)

			got, err := kc.GetBasicCatalogInfo(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.GetBasicCatalogInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validator(t, &tt, got)
			}
		})
	}
}

func TestKeeperCatalogImpl_KeeperSearchMultipleCatalog(t *testing.T) {
	t.Parallel()

	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts       *schema.KeeperSearchCatalogOpts
		createOpts []schema.CreateCatalogOpts
		createResp []schema.CreateCatalogResp
		fakeStart  string
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       []schema.KeeperSearchCatalogResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)

		// validator func(*testing.T, *schema.KeeperSearchCatalogOpts, []schema.KeeperSearchCatalogResp)
	}
	tests := []TC{
		{
			name: "[OK] Mock With Multiple Catalog",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {

				for i := 0; i < faker.RandomInt(12, 15); i++ {
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
					tt.args.createOpts = append(tt.args.createOpts, *createOpts)
					tt.args.createResp = append(tt.args.createResp, *resp)
				}

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.KeeperSearchCatalogOpts{
					Name: resp[11].Name,
					Page: 0,
				}
				tt.args.opts = &opts
				want := []schema.KeeperSearchCatalogResp{
					{
						ID:          resp[11].ID,
						Name:        opts.Name,
						Path:        resp[11].Paths,
						BasePrice:   resp[11].BasePrice,
						RetailPrice: resp[11].RetailPrice,
						Status:      resp[11].Status,
						Variants:    resp[11].Variants,
						VariantType: resp[11].VariantType,
					},
				}
				tt.want = want
			},
			wantErr: false,
		},
		{
			name: "[OK] Multiple Catalog with Same Start Name",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				tt.args.fakeStart = faker.Team().Name()
				for i := 0; i < faker.RandomInt(12, 15); i++ {
					b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
					createOpts := schema.GetRandomCreateCatalogOpts()
					createOpts.Name = tt.args.fakeStart + " " + createOpts.Name
					var categoryCalls []*gomock.Call
					for _, id := range createOpts.CategoryID {
						path := schema.GetRandomGetCategoryPath(id)
						call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
						categoryCalls = append(categoryCalls, call)
					}
					gomock.InOrder(categoryCalls...)
					resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)
					tt.args.createOpts = append(tt.args.createOpts, *createOpts)
					tt.args.createResp = append(tt.args.createResp, *resp)
				}

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.KeeperSearchCatalogOpts{
					Name: tt.args.fakeStart,
					Page: 1,
				}
				tt.args.opts = &opts
				var want []schema.KeeperSearchCatalogResp
				for i := 10; i < len(resp); i++ {
					want = append(want, schema.KeeperSearchCatalogResp{
						ID:          resp[i].ID,
						Name:        resp[i].Name,
						Path:        resp[i].Paths,
						BasePrice:   resp[i].BasePrice,
						RetailPrice: resp[i].RetailPrice,
						Status:      resp[i].Status,
						Variants:    resp[i].Variants,
						VariantType: resp[i].VariantType,
					})
				}

				tt.want = want
			},
			wantErr: false,
		},
		{
			name: "[Error] Catalog Name do not Exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {

				for i := 0; i < faker.RandomInt(12, 15); i++ {
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
					tt.args.createOpts = append(tt.args.createOpts, *createOpts)
					tt.args.createResp = append(tt.args.createResp, *resp)
				}

			},
			prepare: func(tt *TC) {
				// resp := tt.args.createResp
				opts := schema.KeeperSearchCatalogOpts{
					Name: "resp[11].Name",
					Page: 0,
				}
				tt.args.opts = &opts
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
			tt.fields.App.KeeperCatalog = kc

			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			resp, err := kc.KeeperSearchCatalog(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)
			}

			if tt.wantErr {
				// assert.NotNil(t, err)
				assert.Nil(t, resp)
				// assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestKeeperCatalogImpl_GetBasicCatalogInfoWithCategoryFilter(t *testing.T) {

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts   []schema.CreateCatalogOpts
		resp   []schema.CreateCatalogResp
		filter *schema.GetBasicCatalogFilter
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       []schema.GetBasicCatalogResp
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockBrand, *mock.MockCategory)
		validator  func(*testing.T, *TC, []schema.GetBasicCatalogResp)
	}
	tests := []TC{
		{
			name: "[Ok] With Category Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_2"),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				randIndex0 := faker.RandomInt(0, len(tt.args.resp)-1)
				randIndex1 := faker.RandomInt(0, len(tt.args.resp)-1)

				strCategoryID0 := strings.Split(tt.args.resp[randIndex0].Paths[faker.RandomInt(0, len(tt.args.resp[randIndex0].Paths)-1)], "/")
				randCategoryID0, _ := primitive.ObjectIDFromHex(faker.RandomChoice(strCategoryID0))

				strCategoryID1 := strings.Split(tt.args.resp[randIndex1].Paths[faker.RandomInt(0, len(tt.args.resp[randIndex1].Paths)-1)], "/")
				randCategoryID1, _ := primitive.ObjectIDFromHex(faker.RandomChoice(strCategoryID1))

				tt.args.filter = &schema.GetBasicCatalogFilter{
					CategoryID: []primitive.ObjectID{randCategoryID0, randCategoryID1},
				}
				t.Log(tt.args.filter)

				tt.want = []schema.GetBasicCatalogResp{
					{
						ID:          tt.args.resp[randIndex0].ID,
						Paths:       tt.args.resp[randIndex0].Paths,
						Name:        tt.args.resp[randIndex0].Name,
						Description: tt.args.resp[randIndex0].Description,
						RetailPrice: tt.args.resp[randIndex0].RetailPrice,
					},
					{
						ID:          tt.args.resp[randIndex1].ID,
						Paths:       tt.args.resp[randIndex1].Paths,
						Name:        tt.args.resp[randIndex1].Name,
						Description: tt.args.resp[randIndex1].Description,
						RetailPrice: tt.args.resp[randIndex1].RetailPrice,
					},
				}
			},
			buildStubs: func(tt *TC, b *mock.MockBrand, c *mock.MockCategory) {
				var opts []schema.CreateCatalogOpts
				for i := 0; i < faker.RandomInt(10, 20); i++ {
					opt := schema.GetRandomCreateCatalogOpts()
					opts = append(opts, *opt)
				}
				tt.args.opts = opts
				var resp []schema.CreateCatalogResp
				for _, opt := range tt.args.opts {
					b.EXPECT().CheckBrandIDExists(gomock.Any(), opt.BrandID).Times(1).Return(true, nil)
					var categoryCalls []*gomock.Call
					for _, id := range opt.CategoryID {
						path := schema.GetRandomGetCategoryPath(id)
						call := c.EXPECT().GetCategoryPath(id).Return(path, nil)
						categoryCalls = append(categoryCalls, call)
					}
					gomock.InOrder(categoryCalls...)
					res, _ := tt.fields.App.KeeperCatalog.CreateCatalog(&opt)
					resp = append(resp, *res)
				}
				tt.args.resp = resp
			},
			validator: func(t *testing.T, tt *TC, res []schema.GetBasicCatalogResp) {
				assert.Len(t, res, len(tt.want))
				assert.Equal(t, tt.want, res)
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
			tt.fields.App.Brand = mockBrand
			tt.fields.App.Category = mockCategory
			tt.fields.App.KeeperCatalog = kc

			tt.buildStubs(&tt, mockBrand, mockCategory)
			tt.prepare(&tt)

			got, err := kc.GetBasicCatalogInfo(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.GetBasicCatalogInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validator(t, &tt, got)
			}
		})
	}
}

func TestKeeperCatalogImpl_DeleteVariant(t *testing.T) {
	t.Parallel()
	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts       *schema.DeleteVariantOpts
		createOpts *schema.CreateCatalogOpts
		createResp *schema.CreateCatalogResp
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       error
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)

		// validator func(*testing.T, *schema.KeeperSearchCatalogOpts, []schema.KeeperSearchCatalogResp)
	}
	tests := []TC{
		{
			name: "[OK] Single Catalog with Multiple Variants",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.SizeType
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
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
				opts := schema.DeleteVariantOpts{
					CatalogID: resp.ID,
					VariantID: resp.Variants[1].ID,
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[Error] Catalog With NO Variant",
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
				fakeVariantID := primitive.NewObjectIDFromTimestamp(time.Now())
				resp := tt.args.createResp
				opts := schema.DeleteVariantOpts{
					CatalogID: resp.ID,
					VariantID: fakeVariantID,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Failed to delete Variant with id %s", fakeVariantID.Hex())
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
			tt.fields.App.KeeperCatalog = kc

			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			err := kc.DeleteVariant(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				var catalogResp model.Catalog
				fmt.Println(tt.args.opts.CatalogID)
				err := tt.fields.DB.Collection(model.CatalogColl).FindOne(context.TODO(), bson.M{"_id": tt.args.opts.CatalogID}).Decode(&catalogResp)
				assert.Nil(t, err)
				validateDeleteVariant(t, *tt.args.opts, catalogResp)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestKeeperCatalogImpl_GetBasicCatalogInfoWithBrandFilter(t *testing.T) {

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts   []schema.CreateCatalogOpts
		resp   []schema.CreateCatalogResp
		filter *schema.GetBasicCatalogFilter
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       []schema.GetBasicCatalogResp
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockBrand, *mock.MockCategory)
		validator  func(*testing.T, *TC, []schema.GetBasicCatalogResp)
	}
	tests := []TC{
		{
			name: "[Ok] With Brand Filter",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_3"),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				randIndex0 := faker.RandomInt(0, len(tt.args.resp)-1)
				randIndex1 := faker.RandomInt(0, len(tt.args.resp)-1)

				tt.args.filter = &schema.GetBasicCatalogFilter{
					BrandID: []primitive.ObjectID{tt.args.resp[randIndex0].BrandID, tt.args.resp[randIndex1].BrandID},
				}

				tt.want = []schema.GetBasicCatalogResp{
					{
						ID:          tt.args.resp[randIndex0].ID,
						Paths:       tt.args.resp[randIndex0].Paths,
						Name:        tt.args.resp[randIndex0].Name,
						Description: tt.args.resp[randIndex0].Description,
						RetailPrice: tt.args.resp[randIndex0].RetailPrice,
					},
					{
						ID:          tt.args.resp[randIndex1].ID,
						Paths:       tt.args.resp[randIndex1].Paths,
						Name:        tt.args.resp[randIndex1].Name,
						Description: tt.args.resp[randIndex1].Description,
						RetailPrice: tt.args.resp[randIndex1].RetailPrice,
					},
				}
			},
			buildStubs: func(tt *TC, b *mock.MockBrand, c *mock.MockCategory) {
				var opts []schema.CreateCatalogOpts
				for i := 0; i < faker.RandomInt(4, 20); i++ {
					opt := schema.GetRandomCreateCatalogOpts()
					opts = append(opts, *opt)
				}
				tt.args.opts = opts

				for _, opt := range tt.args.opts {
					call := b.EXPECT().CheckBrandIDExists(gomock.Any(), opt.BrandID).Return(true, nil)
					gomock.InOrder(call)
					var categoryCalls []*gomock.Call
					for _, id := range opt.CategoryID {
						path := schema.GetRandomGetCategoryPath(id)
						call := c.EXPECT().GetCategoryPath(id).Return(path, nil)
						categoryCalls = append(categoryCalls, call)
					}
					gomock.InOrder(categoryCalls...)
				}

				var resp []schema.CreateCatalogResp
				for _, opt := range tt.args.opts {
					res, _ := tt.fields.App.KeeperCatalog.CreateCatalog(&opt)
					resp = append(resp, *res)
				}
				tt.args.resp = resp

			},
			validator: func(t *testing.T, tt *TC, res []schema.GetBasicCatalogResp) {
				assert.Len(t, res, len(tt.want))
				assert.Equal(t, tt.want, res)
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
			tt.fields.App.Brand = mockBrand
			tt.fields.App.Category = mockCategory
			tt.fields.App.KeeperCatalog = kc

			tt.buildStubs(&tt, mockBrand, mockCategory)
			tt.prepare(&tt)

			got, err := kc.GetBasicCatalogInfo(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.GetBasicCatalogInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validator(t, &tt, got)
			}
		})
	}
}

func TestKeeperCatalogImpl_UpdateCatalogStatus(t *testing.T) {
	t.Parallel()
	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts        *schema.UpdateCatalogStatusOpts
		createOpts  *schema.CreateCatalogOpts
		createResp  *schema.CreateCatalogResp
		errorFields []string
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       error
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)

		// validator func(*testing.T, *schema.KeeperSearchCatalogOpts, []schema.KeeperSearchCatalogResp)
	}
	tests := []TC{
		{
			name: "[OK] Draft To Publish",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{"featured_image": featuredImage}})
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    "published",
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[OK] Draft To Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{"featured_image": featuredImage}})
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    "archive",
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[OK] Publish To Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Publish),
						Value:     model.Publish,
						CreatedAt: time.Now(),
					},
				}})
				resp.Status.Value = model.Publish
				resp.Status.Name = strings.Title(model.Publish)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Archive,
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[OK] Publish To Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Publish),
						Value:     model.Publish,
						CreatedAt: time.Now(),
					},
				}})
				resp.Status.Value = model.Publish
				resp.Status.Name = strings.Title(model.Publish)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Unlist,
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[OK] Unlist To Publish",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Unlist),
						Value:     model.Unlist,
						CreatedAt: time.Now(),
					},
				}})
				resp.Status.Value = model.Unlist
				resp.Status.Name = strings.Title(model.Unlist)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Publish,
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[OK] Unlist To Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Unlist),
						Value:     model.Unlist,
						CreatedAt: time.Now(),
					},
				}})
				resp.Status.Value = model.Unlist
				resp.Status.Name = strings.Title(model.Unlist)
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Archive,
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[Error] Draft To Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{"featured_image": featuredImage}})
				tt.args.createOpts = createOpts
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Unlist,
				}
				tt.args.opts = &opts
			},
			wantErr: true,
			err:     errors.Errorf("Status change not allowed from %s to %s", model.Draft, model.Unlist),
		},
		{
			name: "[Error] Archive To Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType

				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Archive),
						Value:     model.Archive,
						CreatedAt: time.Now(),
					},
				}})
				tt.args.createOpts = createOpts
				resp.Status.Value = model.Archive
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Unlist,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Status change not allowed from %s to %s", resp.Status.Value, opts.Status)

			},
			wantErr: true,
		},
		{
			name: "[Error] Archive To Publish",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType

				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Archive),
						Value:     model.Archive,
						CreatedAt: time.Now(),
					},
				}})
				tt.args.createOpts = createOpts
				resp.Status.Value = model.Archive
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Publish,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Status change not allowed from %s to %s", resp.Status.Value, opts.Status)

			},
			wantErr: true,
		},
		{
			name: "[Error] Archive To Draft",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType

				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Archive),
						Value:     model.Archive,
						CreatedAt: time.Now(),
					},
				}})
				tt.args.createOpts = createOpts
				resp.Status.Value = model.Archive
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Draft,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Status change not allowed from %s to %s", resp.Status.Value, opts.Status)

			},
			wantErr: true,
		},
		{
			name: "[Error] Publish To Draft",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType

				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Publish),
						Value:     model.Publish,
						CreatedAt: time.Now(),
					},
				}})
				tt.args.createOpts = createOpts
				resp.Status.Value = model.Publish
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Draft,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Status change not allowed from %s to %s", resp.Status.Value, opts.Status)

			},
			wantErr: true,
		},
		{
			name: "[Error] Unlist To Draft",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.VariantType = model.ColorType

				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
					"status": model.Status{
						Name:      strings.Title(model.Unlist),
						Value:     model.Unlist,
						CreatedAt: time.Now(),
					},
				}})
				tt.args.createOpts = createOpts
				resp.Status.Value = model.Unlist
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Draft,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Status change not allowed from %s to %s", resp.Status.Value, opts.Status)

			},
			wantErr: true,
		},
		{
			name: "[Error] Draft To Publish - Missing Name",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.Name = ""
				createOpts.VariantType = model.ColorType
				createOpts.FilterAttribute = []schema.FilterAttribute{
					{
						Name:  faker.Company().Name(),
						Value: faker.RandomString(5),
					},
				}
				for i := 0; i < 3; i++ {
					createVariantOpts := schema.GetRandomCreateVariantOpts()
					createVariantOpts.Attribute = colors[i]
					createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				}
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				featuredImage := &model.CatalogFeaturedImage{
					ID: primitive.NewObjectID(),
					IMG: model.IMG{
						SRC:    faker.Internet().Url(),
						Height: 20,
						Width:  20,
					},
				}
				//Adding Featured Image to DB
				tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
					"featured_image": featuredImage,
				}})
				tt.args.createOpts = createOpts
				tt.args.errorFields = []string{"Name"}
				tt.args.createResp = resp

			},
			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Publish,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Catalog Data not Complete")

			},
			wantErr: true,
		},
		{
			name: "[Error] Draft To Publish - Missing All Fields",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				createOpts := schema.GetRandomCreateCatalogOpts()
				createOpts.Name = ""
				createOpts.Description = ""
				createOpts.Keywords = nil
				createOpts.ETA = nil
				createOpts.HSNCode = ""
				createOpts.BasePrice = 0
				createOpts.RetailPrice = 0
				createOpts.CategoryID = nil
				// createOpts.VariantType = model.ColorType
				// createOpts.FilterAttribute = []schema.FilterAttribute{
				// 	{
				// 		Name:  faker.Company().Name(),
				// 		Value: faker.RandomString(5),
				// 	},
				// }

				// for i := 0; i < 3; i++ {
				// 	createVariantOpts := schema.GetRandomCreateVariantOpts()
				// 	createVariantOpts.Attribute = colors[i]
				// 	createOpts.Variants = append(createOpts.Variants, *createVariantOpts)
				// }
				var categoryCalls []*gomock.Call
				for _, id := range createOpts.CategoryID {
					path := schema.GetRandomGetCategoryPath(id)
					call := ct.EXPECT().GetCategoryPath(id).Return(path, nil)
					categoryCalls = append(categoryCalls, call)
				}
				gomock.InOrder(categoryCalls...)
				resp, _ := tt.fields.App.KeeperCatalog.CreateCatalog(createOpts)

				// featuredImage := &model.CatalogFeaturedImage{
				// 	ID: primitive.NewObjectID(),
				// 	IMG: model.IMG{
				// 		SRC:    faker.Internet().Url(),
				// 		Height: 20,
				// 		Width:  20,
				// 	},
				// }
				//Adding Featured Image to DB
				// tt.fields.DB.Collection(model.CatalogColl).UpdateOne(context.TODO(), bson.M{"_id": resp.ID}, bson.M{"$set": bson.M{
				// 	"featured_image": featuredImage,
				// }})

				tt.args.createOpts = createOpts
				tt.args.errorFields = []string{
					"Name",
					"Description",
					"Category",
					"Keywords",
					"Featured Image",
					"Filter Attribute",
					"Variants",
					"Variant Type",
					"ETA",
					"HSN Code",
					"Base Price",
					"Retail Price",
				}
				tt.args.createResp = resp

			},

			prepare: func(tt *TC) {
				resp := tt.args.createResp
				opts := schema.UpdateCatalogStatusOpts{
					CatalogID: resp.ID,
					Status:    model.Publish,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("Catalog Data not Complete")

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
			tt.fields.App.KeeperCatalog = kc

			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			resp, err := kc.UpdateCatalogStatus(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Nil(t, resp)
				//Checking if Status Change was Successful
				var catalogResp model.Catalog
				err := tt.fields.DB.Collection(model.CatalogColl).FindOne(context.TODO(), bson.M{"_id": tt.args.opts.CatalogID}).Decode(&catalogResp)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.opts.Status, catalogResp.Status.Value)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				// assert.NotNil(t, resp)
				assert.Equal(t, tt.err.Error(), err.Error())
				if err.Error() == "Catalog Data not Complete" {
					for i := 0; i < len(resp); i++ {
						assert.Equal(t, resp[i].Field, tt.args.errorFields[i])
					}

				}
			}
		})
	}
}

func TestKeeperCatalogImpl_GetCatalogFilter(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		wantErr    bool
		want       *schema.GetCatalogFilterResp
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_1"),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {

			},
			buildStubs: func(tt *TC, c *mock.MockCategory) {
				var resp []schema.GetCategoriesBasicResp
				for i := 0; i < faker.RandomInt(6, 20); i++ {
					resp = append(resp, schema.GetCategoriesBasicResp{
						ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
						Name: faker.Commerce().Department(),
					})
				}
				c.EXPECT().GetCategoriesBasic().Times(1).Return(resp, nil)
				tt.want = &schema.GetCatalogFilterResp{
					Category: resp,
				}
			},
		},
		{
			name: "[Ok] No Categories",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName + "_1"),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {

			},
			buildStubs: func(tt *TC, c *mock.MockCategory) {
				var resp []schema.GetCategoriesBasicResp
				c.EXPECT().GetCategoriesBasic().Times(1).Return(resp, nil)
				tt.want = &schema.GetCatalogFilterResp{
					Category: resp,
				}
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

			tt.fields.App.KeeperCatalog = kc
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory)

			got, err := kc.GetCatalogFilter()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeeperCatalogImpl.GetCatalogFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestKeeperCatalogImpl_AddCatalogContent(t *testing.T) {
	t.Parallel()
	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.AddCatalogContentOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       error
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)

		// validator func(*testing.T, *schema.KeeperSearchCatalogOpts, []schema.KeeperSearchCatalogResp)
	}
	tests := []TC{
		{
			name: "[Error] CatalogID not found",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.AddCatalogContentOpts{
					CatalogID: primitive.NewObjectID(),
				}
				tt.err = errors.Errorf("unable to find the catalog with id: %s", tt.args.opts.CatalogID.Hex())
			},
			wantErr: true,
		},
		{
			name: "[Error] BrandID not found",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// colors := []string{"red", "blue", "green"}
				b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.AddCatalogContentOpts{
					BrandID: primitive.NewObjectID(),
				}
				tt.err = errors.Errorf("unable to find the brand with id: %s", tt.args.opts.BrandID.Hex())
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
			tt.fields.App.KeeperCatalog = kc

			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			resp, err := kc.AddCatalogContent(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				// assert.NotNil(t, resp)
				assert.Equal(t, tt.err.Error(), err.Error())

			}
		})
	}
}

func TestKeeperCatalogImpl_GetCatalogsByFilter(t *testing.T) {
	t.Parallel()
	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)
	ctx := context.TODO()
	bIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID()}
	status := []string{model.Archive, model.Publish, model.Draft, model.Unlist}
	var catalogs []model.Catalog
	for i := 0; i < 10; i++ {
		cat := model.Catalog{
			ID:      primitive.NewObjectID(),
			BrandID: bIDs[faker.RandomInt(0, 2)],
			Status: &model.Status{
				Value: status[faker.RandomInt(0, 3)],
			},
		}
		app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName).Collection(model.CatalogColl).InsertOne(ctx, cat)
		catalogs = append(catalogs, cat)
	}

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.GetCatalogsByFilterOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       []schema.GetCatalogResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)

		// validator func(*testing.T, *schema.KeeperSearchCatalogOpts, []schema.KeeperSearchCatalogResp)
	}

	tests := []TC{
		{
			name: "[OK]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.GetCatalogsByFilterOpts{
					Status:   []string{model.Publish, model.Archive},
					BrandIDs: []primitive.ObjectID{bIDs[1], bIDs[2]},
				}
				var catalogResp []schema.GetCatalogResp
				for i := 0; i < 10; i++ {
					if (catalogs[i].Status.Value == model.Publish || catalogs[i].Status.Value == model.Archive) && (catalogs[i].BrandID == bIDs[1] || catalogs[i].BrandID == bIDs[2]) {
						catalogResp = append(catalogResp, schema.CreateCatalogResp{
							ID:      catalogs[i].ID,
							BrandID: catalogs[i].BrandID,
							Status:  catalogs[i].Status,
						})
					}
				}
				tt.want = catalogResp
			},
			wantErr: false,
		},
		{
			name: "[OK] Status Missing",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.GetCatalogsByFilterOpts{
					BrandIDs: []primitive.ObjectID{bIDs[1]},
				}
				var catalogResp []schema.GetCatalogResp
				for i := 0; i < 10; i++ {
					if catalogs[i].BrandID == bIDs[1] {
						catalogResp = append(catalogResp, schema.CreateCatalogResp{
							ID:      catalogs[i].ID,
							BrandID: catalogs[i].BrandID,
							Status:  catalogs[i].Status,
						})
					}
				}
				tt.want = catalogResp
			},
			wantErr: false,
		},
		{
			name: "[OK] Brand Missing",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.GetCatalogsByFilterOpts{
					Status: []string{model.Publish},
				}
				var catalogResp []schema.GetCatalogResp
				for i := 0; i < 10; i++ {
					if catalogs[i].Status.Value == model.Publish {
						catalogResp = append(catalogResp, schema.CreateCatalogResp{
							ID:      catalogs[i].ID,
							BrandID: catalogs[i].BrandID,
							Status:  catalogs[i].Status,
						})
					}
				}
				tt.want = catalogResp
			},
			wantErr: false,
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

			resp, err := kc.GetCatalogsByFilter(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				// assert.NotNil(t, resp)
				assert.Equal(t, tt.err.Error(), err.Error())

			}
		})
	}

}

func TestKeeperCatalogImpl_GetCatalogBySlug(t *testing.T) {
	t.Parallel()
	// opts := schema.GetRandomKeeperSearchCatalog()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)
	ctx := context.TODO()
	status := []string{model.Archive, model.Publish, model.Draft, model.Unlist}
	var catalogs []model.Catalog
	for i := 0; i < 4; i++ {
		cat := model.Catalog{
			ID: primitive.NewObjectID(),
			Status: &model.Status{
				Value: status[faker.RandomInt(0, 3)],
			},
			Slug: faker.Commerce().ProductName(),
		}
		app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName).Collection(model.CatalogColl).InsertOne(ctx, cat)
		catalogs = append(catalogs, cat)
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		slug string
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.GetCatalogResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockCategory, *mock.MockBrand)
	}
	tests := []TC{
		{
			name: "[OK]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			},
			prepare: func(tt *TC) {
				tt.args.slug = catalogs[0].Slug

				tt.want = &schema.GetCatalogResp{
					ID:     catalogs[0].ID,
					Status: catalogs[0].Status,
					Slug:   catalogs[0].Slug,
				}
			},
			wantErr: false,
		},
		{
			name: "[Error] Catalog with given Slug not found",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand) {
				// b.EXPECT().CheckBrandIDExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			},
			prepare: func(tt *TC) {
				tt.args.slug = faker.Commerce().ProductName()

				tt.err = errors.Errorf("unable to find the catalog with slug: %s", tt.args.slug)
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
			tt.fields.App.KeeperCatalog = kc

			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory

			tt.buildStubs(&tt, mockCategory, mockBrand)
			tt.prepare(&tt)

			resp, err := kc.GetCatalogBySlug(tt.args.slug)
			fmt.Println(resp)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				// assert.NotNil(t, resp)
				assert.Equal(t, tt.err.Error(), err.Error())

			}
		})
	}
}
