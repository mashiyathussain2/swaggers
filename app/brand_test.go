package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

func validateCreateBrandResp(t *testing.T, opts *schema.CreateBrandOpts, resp *schema.CreateBrandResp) {
	assert.Equal(t, opts.Name, resp.Name)
	assert.Equal(t, opts.Description, resp.Description)
	assert.Equal(t, opts.WebsiteLink, resp.WebsiteLink)
	assert.Equal(t, opts.FulfillmentEmail, resp.Fulfillment.Email)
	assert.False(t, resp.ID.IsZero())
	assert.NotEmpty(t, resp.Slug)
}

func TestBrandImpl_CreateBrand(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	b1 := schema.GetRandomCreateBrandOpts()

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateBrandOpts
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		err       string
		validator func(*testing.T, *schema.CreateBrandOpts, *schema.CreateBrandResp)
	}{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: b1,
			},
			validator: func(t *testing.T, s1 *schema.CreateBrandOpts, s2 *schema.CreateBrandResp) {
				validateCreateBrandResp(t, s1, s2)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			resp, err := b.CreateBrand(tt.args.opts)
			if !tt.wantErr {
				assert.NotNil(t, resp)
				tt.validator(t, tt.args.opts, resp)
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err, err.Error())
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
		createBrandOpts *schema.CreateBrandOpts
		createBrandResp *schema.CreateBrandResp
		editBrandOpts   *schema.EditBrandOpts
	}
	type TC struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		want      *schema.EditBrandResp
		prepare   func(*TC) error
		validator func(*testing.T, *TC, *schema.EditBrandResp)
		err       error
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				createBrandOpts: schema.GetRandomCreateBrandOpts(),
			},
			wantErr: false,
			prepare: func(tt *TC) error {
				res, err := tt.fields.App.Brand.CreateBrand(tt.args.createBrandOpts)
				if err != nil {
					return err
				}
				tt.args.createBrandResp = res
				tt.args.editBrandOpts = &schema.EditBrandOpts{
					ID:   tt.args.createBrandResp.ID,
					Name: tt.args.createBrandOpts.Name + " Edited",
				}

				editBrandResp := schema.EditBrandResp(*res)
				editBrandResp.Name = tt.args.createBrandOpts.Name + " Edited"
				tt.want = &editBrandResp

				return nil
			},
			validator: func(t *testing.T, tt *TC, s1 *schema.EditBrandResp) {
				assert.Equal(t, tt.want, s1)
			},
		},
		{
			name: "[Error] Invalid BrandID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				createBrandOpts: schema.GetRandomCreateBrandOpts(),
			},
			wantErr: true,
			prepare: func(tt *TC) error {
				res, err := tt.fields.App.Brand.CreateBrand(tt.args.createBrandOpts)
				if err != nil {
					return err
				}
				tt.args.createBrandResp = res
				tt.args.editBrandOpts = &schema.EditBrandOpts{
					ID:   primitive.NewObjectID(),
					Name: tt.args.createBrandOpts.Name + " Edited",
				}

				tt.err = errors.Errorf("brand with id:%s not found", tt.args.editBrandOpts.ID)

				return nil
			},
			validator: func(t *testing.T, tt *TC, s1 *schema.EditBrandResp) {
				assert.Equal(t, tt.want, s1)
			},
		},
		{
			name: "[Ok] Multiple Field",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				createBrandOpts: schema.GetRandomCreateBrandOpts(),
			},
			wantErr: false,
			prepare: func(tt *TC) error {
				res, err := tt.fields.App.Brand.CreateBrand(tt.args.createBrandOpts)
				if err != nil {
					return err
				}
				tt.args.createBrandResp = res
				tt.args.editBrandOpts = &schema.EditBrandOpts{
					ID:               tt.args.createBrandResp.ID,
					Name:             tt.args.createBrandOpts.Name + " Edited",
					Description:      tt.args.createBrandOpts.Description + " Edited",
					WebsiteLink:      faker.Internet().Url(),
					FulfillmentEmail: faker.Internet().Email(),
				}

				editBrandResp := res
				editBrandResp.Name = tt.args.editBrandOpts.Name
				editBrandResp.Description = tt.args.editBrandOpts.Description
				editBrandResp.WebsiteLink = tt.args.editBrandOpts.WebsiteLink
				editBrandResp.Fulfillment = &model.Fulfillment{Email: tt.args.editBrandOpts.FulfillmentEmail}
				tt.want = editBrandResp

				return nil
			},
			validator: func(t *testing.T, tt *TC, s1 *schema.EditBrandResp) {
				assert.Equal(t, tt.want, s1)
			},
		},
		{
			name: "[Error] No fields",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				createBrandOpts: schema.GetRandomCreateBrandOpts(),
			},
			wantErr: true,
			prepare: func(tt *TC) error {
				res, err := tt.fields.App.Brand.CreateBrand(tt.args.createBrandOpts)
				if err != nil {
					return err
				}
				tt.args.createBrandResp = res
				tt.args.editBrandOpts = &schema.EditBrandOpts{}

				tt.err = errors.New("no fields found to update")

				return nil
			},
			validator: func(t *testing.T, tt *TC, s1 *schema.EditBrandResp) {
				assert.Equal(t, tt.want, s1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = b
			err := tt.prepare(&tt)
			assert.Nil(t, err)

			resp, err := b.EditBrand(tt.args.editBrandOpts)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				tt.validator(t, &tt, resp)
			}
			if tt.wantErr {
				assert.Nil(t, resp)
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestBrandImpl_CheckBrandIDExists(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		ctx context.Context
		id  primitive.ObjectID
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		prepare  func(*TC) error
		validate func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[OK]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			want:    true,
			wantErr: false,
			prepare: func(tt *TC) error {
				s1 := schema.GetRandomCreateBrandOpts()
				res, err := tt.fields.App.Brand.CreateBrand(s1)
				if err != nil {
					return err
				}
				tt.args.id = res.ID
				tt.args.ctx = context.Background()
				return nil
			},
			validate: func(t *testing.T, tt *TC, s1 bool) {
				assert.Equal(t, tt.want, s1)
			},
		},
		{
			name: "[OK] When ID does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.BrandConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			want:    false,
			wantErr: false,
			prepare: func(tt *TC) error {
				return nil
			},
			validate: func(t *testing.T, tt *TC, s1 bool) {
				assert.Equal(t, tt.want, s1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BrandImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Brand = b
			tt.prepare(&tt)
			res, err := b.CheckBrandIDExists(tt.args.ctx, tt.args.id)
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, res)
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.False(t, res)
			}
		})
	}
}
