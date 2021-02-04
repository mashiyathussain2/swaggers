package app

import (
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func validateCreateCatalogResp(t *testing.T, opts *schema.CreateCatalogOpts, resp *schema.CreateCatalogResp) {
	assert.False(t, resp.ID.IsZero())
	assert.Equal(t, opts.Name, resp.Name)
	assert.Equal(t, opts.BrandID, resp.BrandID)
	// assert.Equal(t, opts.CategoryID, resp.CategoryID)
	assert.Equal(t, opts.Description, resp.Description)
	assert.Equal(t, opts.Keywords, resp.Keywords)
	assert.Equal(t, opts.HSNCode, resp.HSNCode)
	if opts.ETA != nil {
		assert.Equal(t, opts.ETA, resp.ETA)
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
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		buildStubs func(b *mock.MockBrand)
		validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}{
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
			buildStubs: func(b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), opts.BrandID).Times(1).Return(true, nil)
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
			buildStubs: func(b *mock.MockBrand) {
				b.EXPECT().CheckBrandIDExists(gomock.Any(), optsWithVariants.BrandID).Times(1).Return(true, nil)
			},
		},
		{
			name: "[Ok] When brandID does not exists",
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
			buildStubs: func(b *mock.MockBrand) {
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
			kc.App.Brand = mockBrand
			tt.buildStubs(mockBrand)

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
