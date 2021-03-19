package app

import (
	"fmt"
	"go-app/model"
	"go-app/schema"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

func validateCreateCategoryResp(t *testing.T, opts *schema.CreateCategoryOpts, resp *schema.CreateCategoryResp) {
	assert.False(t, resp.ID.IsZero())
	assert.Equal(t, opts.Name, resp.Name)
	assert.Greater(t, len(resp.Slug), len(opts.Name))
	assert.Equal(t, opts.ParentID, resp.ParentID)
	assert.Equal(t, opts.IsMain, resp.IsMain)
	assert.Equal(t, opts.FeaturedImage.SRC, resp.FeaturedImage.SRC)
	assert.Equal(t, opts.Thumbnail.SRC, resp.Thumbnail.SRC)
	assert.Equal(t, 100, resp.Thumbnail.Width)
	assert.Equal(t, 100, resp.Thumbnail.Height)
	assert.Equal(t, 300, resp.FeaturedImage.Width)
	assert.Equal(t, 300, resp.FeaturedImage.Height)
}

func TestCategoryImpl_CreateCategory(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	opts := schema.GetRandomCreateCategoryOpts()

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		parentOpts *schema.CreateCategoryResp
		opts       *schema.CreateCategoryOpts
	}

	type TC struct {
		name      string
		fields    fields
		args      args
		want      *schema.CreateCategoryResp
		wantErr   bool
		err       error
		prepare   func(*TC)
		validator func(*testing.T, *TC, *schema.CreateCategoryResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parentCatOpts := schema.GetRandomCreateCategoryOpts()
				parentCatOpts.ParentID = primitive.NilObjectID
				resp, err := tt.fields.App.Category.CreateCategory(parentCatOpts)
				if err != nil {
					log.Fatalf("%s", err)
				}
				tt.args.opts.ParentID = resp.ID
				tt.args.parentOpts = resp
			},
			validator: func(t *testing.T, tt *TC, s2 *schema.CreateCategoryResp) {
				validateCreateCategoryResp(t, tt.args.opts, s2)
				assert.Equal(t, []primitive.ObjectID{tt.args.parentOpts.ID}, s2.AncestorID)
			},
		},
		{
			name: "[Ok] Multiple Ancestors",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parentCatOpts := schema.GetRandomCreateCategoryOpts()
				parentCatOpts.ParentID = primitive.NilObjectID
				resp, err := tt.fields.App.Category.CreateCategory(parentCatOpts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				parentCatOpts2 := schema.GetRandomCreateCategoryOpts()
				parentCatOpts2.ParentID = resp.ID
				resp2, err := tt.fields.App.Category.CreateCategory(parentCatOpts2)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.opts.ParentID = resp2.ID
				tt.args.parentOpts = resp2
			},
			validator: func(t *testing.T, tt *TC, s2 *schema.CreateCategoryResp) {
				validateCreateCategoryResp(t, tt.args.opts, s2)
				assert.Len(t, s2.AncestorID, 2)
				var aID []primitive.ObjectID
				aID = append(tt.args.parentOpts.AncestorID, tt.args.parentOpts.ID)
				assert.Equal(t, aID, s2.AncestorID)
			},
		},
		{
			name: "[Error] Invalid Parent ID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts.ParentID = primitive.NewObjectIDFromTimestamp(time.Now())
				tt.err = errors.Errorf("category with id:%s not found", tt.args.opts.ParentID)
			},
		},
		{
			name: "[Ok] With no parent ID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				tt.args.opts.ParentID = primitive.NilObjectID
			},
			validator: func(t *testing.T, tt *TC, s2 *schema.CreateCategoryResp) {
				validateCreateCategoryResp(t, tt.args.opts, s2)
				assert.True(t, s2.ParentID.IsZero())
				assert.Len(t, s2.AncestorID, 0)
			},
		},
		{
			name: "[Error] Invalid Thumbnail SRC",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts.ParentID = primitive.NilObjectID
				tt.args.opts.Thumbnail.SRC = "http://harberweber.name/golden"

				tt.err = errors.New("invalid thumbnail url: Get \"http://harberweber.name/golden\": dial tcp: lookup harberweber.name: no such host")
			},
		},
		{
			name: "[Error] Invalid FeaturedImage SRC",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: opts,
			},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts.ParentID = primitive.NilObjectID
				tt.args.opts.FeaturedImage.SRC = "http://harberweber.name/golden"

				tt.err = errors.New("invalid featured image url: Get \"http://harberweber.name/golden\": dial tcp: lookup harberweber.name: no such host")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = c
			tt.prepare(&tt)
			got, err := c.CreateCategory(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
			if !tt.wantErr {
				tt.validator(t, &tt, got)
			}
		})
	}
}

func TestCategoryImpl_GetAncestorsByID(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	// defer CleanTestApp(app)

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
		want     []primitive.ObjectID
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, []primitive.ObjectID)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}
				tt.args.id = resp1.ID
				tt.want = resp1.AncestorID
			},
			wantErr: false,
		},
		{
			name: "[Ok] With Ancestors",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts2 := schema.GetRandomCreateCategoryOpts()
				opts2.ParentID = resp1.ID
				resp2, err := tt.fields.App.Category.CreateCategory(opts2)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.id = resp2.ID
				tt.want = resp2.AncestorID
			},
			wantErr: false,
		},
		{
			name: "[Error] Invalid ID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.id = primitive.NewObjectIDFromTimestamp(time.Now())
				tt.err = errors.Errorf("category with id:%s not found", tt.args.id)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = c
			tt.prepare(&tt)
			got, err := c.GetAncestorsByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetAncestorsByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestCategoryImpl_EditCategory(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts *schema.CreateCategoryOpts
		createResp *schema.CreateCategoryResp
		opts       *schema.EditCategoryOpts
	}
	type TC struct {
		name    string
		fields  fields
		args    args
		want    *schema.EditCategoryResp
		wantErr bool
		err     error
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: false,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				tt.args.opts = &schema.EditCategoryOpts{
					ID:   resp1.ID,
					Name: resp1.Name + " Edited",
				}

				want := schema.EditCategoryResp(*resp1)
				tt.want = &want
				tt.want.Name = resp1.Name + " Edited"
			},
		},
		{
			name: "[Ok] IsMain from false to true",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: false,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				opts1.IsMain = false
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				t := true
				tt.args.opts = &schema.EditCategoryOpts{
					ID:     resp1.ID,
					IsMain: &t,
				}

				want := schema.EditCategoryResp(*resp1)
				tt.want = &want
				tt.want.IsMain = true
			},
		},
		{
			name: "[Ok] IsMain from true to false",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: false,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				opts1.IsMain = true
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				t := false
				tt.args.opts = &schema.EditCategoryOpts{
					ID:     resp1.ID,
					IsMain: &t,
				}

				want := schema.EditCategoryResp(*resp1)
				tt.want = &want
				tt.want.IsMain = false
			},
		},
		{
			name: "[Ok] Thumbnail & Featured Image",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: false,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				tt.args.opts = &schema.EditCategoryOpts{
					ID:   resp1.ID,
					Name: resp1.Name + " Edited",
					Thumbnail: &schema.Img{
						SRC: faker.Avatar().Url("png", 100, 100),
					},
					FeaturedImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 300, 300),
					},
				}

				want := schema.EditCategoryResp(*resp1)
				tt.want = &want
				tt.want.Name = resp1.Name + " Edited"
				tt.want.FeaturedImage = &model.IMG{
					SRC:    tt.args.opts.FeaturedImage.SRC,
					Width:  300,
					Height: 300,
				}
				tt.want.Thumbnail = &model.IMG{
					SRC:    tt.args.opts.Thumbnail.SRC,
					Width:  100,
					Height: 100,
				}
			},
		},
		{
			name: "[Ok] With ParentID",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: false,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts2 := schema.GetRandomCreateCategoryOpts()
				opts2.ParentID = resp1.ID
				resp2, err := tt.fields.App.Category.CreateCategory(opts2)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts2
				tt.args.createResp = resp2

				tt.args.opts = &schema.EditCategoryOpts{
					ID:   resp2.ID,
					Name: resp2.Name + " Edited",
				}

				want := schema.EditCategoryResp(*resp2)
				tt.want = &want
				tt.want.Name = resp2.Name + " Edited"
			},
		},
		{
			name: "[Error] No Fields",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				tt.args.opts = &schema.EditCategoryOpts{}

				tt.err = errors.New("no fields found to update")
			},
		},
		{
			name: "[Error] Invalid Thumbnail",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				tt.args.opts = &schema.EditCategoryOpts{
					Thumbnail: &schema.Img{
						SRC: "faker.Avatar().Url(\"png\", 100, 100)",
					},
				}
				tt.err = errors.New("invalid thumbnail url: Get \"faker.Avatar%28%29.Url%28%22png%22,%20100,%20100%29\": unsupported protocol scheme \"\"")
			},
		},
		{
			name: "[Error] Invalid Featured Image",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = primitive.NilObjectID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args.createOpts = opts1
				tt.args.createResp = resp1

				tt.args.opts = &schema.EditCategoryOpts{
					FeaturedImage: &schema.Img{
						SRC: "faker.Avatar().Url(\"png\", 100, 100)",
					},
				}
				tt.err = errors.New("invalid featured image url: Get \"faker.Avatar%28%29.Url%28%22png%22,%20100,%20100%29\": unsupported protocol scheme \"\"")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = c
			tt.prepare(&tt)
			got, err := c.EditCategory(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.EditCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestCategoryImpl_GetCategoryPath(t *testing.T) {
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
		want    string
		wantErr bool
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok] Level 0",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				opts0 := schema.GetRandomCreateCategoryOpts()
				opts0.ParentID = primitive.NilObjectID
				resp0, err := tt.fields.App.Category.CreateCategory(opts0)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args = args{
					id: resp0.ID,
				}
				tt.want = fmt.Sprintf("/%s", resp0.ID)
			},
		},
		{
			name: "[Ok] Level 1",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				opts0 := schema.GetRandomCreateCategoryOpts()
				opts0.ParentID = primitive.NilObjectID
				resp0, err := tt.fields.App.Category.CreateCategory(opts0)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = resp0.ID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				tt.args = args{
					id: resp1.ID,
				}
				tt.want = fmt.Sprintf("/%s/%s", resp0.ID, resp1.ID)
			},
		},
		{
			name: "[Ok] Level 2",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				opts0 := schema.GetRandomCreateCategoryOpts()
				opts0.ParentID = primitive.NilObjectID
				resp0, err := tt.fields.App.Category.CreateCategory(opts0)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.ParentID = resp0.ID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts2 := schema.GetRandomCreateCategoryOpts()
				opts2.ParentID = resp1.ID
				resp2, err := tt.fields.App.Category.CreateCategory(opts2)
				if err != nil {
					log.Fatalf("%s", err)
				}
				tt.args = args{
					id: resp2.ID,
				}
				tt.want = fmt.Sprintf("/%s/%s/%s", resp0.ID, resp1.ID, resp2.ID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = c
			tt.prepare(&tt)
			got, err := c.GetCategoryPath(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetCategoryPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CategoryImpl.GetCategoryPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryImpl_GetMainCategoriesMap(t *testing.T) {

	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}

	type TC struct {
		name    string
		fields  fields
		want    map[string]schema.GetMainCategoriesMapResp
		wantErr bool
		prepare func(*TC)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			prepare: func(tt *TC) {
				opts0 := schema.GetRandomCreateCategoryOpts()
				opts0.IsMain = true
				opts0.ParentID = primitive.NilObjectID
				resp0, err := tt.fields.App.Category.CreateCategory(opts0)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts1 := schema.GetRandomCreateCategoryOpts()
				opts1.IsMain = false
				opts1.ParentID = resp0.ID
				resp1, err := tt.fields.App.Category.CreateCategory(opts1)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts11 := schema.GetRandomCreateCategoryOpts()
				opts11.IsMain = false
				opts11.ParentID = resp1.ID
				resp11, err := tt.fields.App.Category.CreateCategory(opts11)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts01 := schema.GetRandomCreateCategoryOpts()
				opts01.IsMain = true
				opts01.ParentID = resp0.ID
				resp01, err := tt.fields.App.Category.CreateCategory(opts01)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts110 := schema.GetRandomCreateCategoryOpts()
				opts110.IsMain = true
				opts110.ParentID = resp11.ID
				resp110, err := tt.fields.App.Category.CreateCategory(opts110)
				if err != nil {
					log.Fatalf("%s", err)
				}

				opts010 := schema.GetRandomCreateCategoryOpts()
				opts010.IsMain = true
				opts010.ParentID = resp01.ID
				resp010, err := tt.fields.App.Category.CreateCategory(opts010)
				if err != nil {
					log.Fatalf("%s", err)
				}

				want := make(map[string]schema.GetMainCategoriesMapResp)
				want[resp0.ID.Hex()] = schema.GetMainCategoriesMapResp{
					ID:            resp0.ID,
					Name:          resp0.Name,
					ParentID:      resp0.ParentID,
					AncestorID:    resp0.AncestorID,
					Thumbnail:     resp0.Thumbnail,
					FeaturedImage: resp0.FeaturedImage,
				}
				want[resp01.ID.Hex()] = schema.GetMainCategoriesMapResp{
					ID:            resp01.ID,
					Name:          resp01.Name,
					ParentID:      resp01.ParentID,
					AncestorID:    resp01.AncestorID,
					Thumbnail:     resp01.Thumbnail,
					FeaturedImage: resp01.FeaturedImage,
				}
				want[resp010.ID.Hex()] = schema.GetMainCategoriesMapResp{
					ID:            resp010.ID,
					Name:          resp010.Name,
					ParentID:      resp010.ParentID,
					AncestorID:    resp010.AncestorID,
					Thumbnail:     resp010.Thumbnail,
					FeaturedImage: resp010.FeaturedImage,
				}
				want[resp110.ID.Hex()] = schema.GetMainCategoriesMapResp{
					ID:            resp110.ID,
					Name:          resp110.Name,
					ParentID:      resp110.ParentID,
					AncestorID:    resp110.AncestorID,
					Thumbnail:     resp110.Thumbnail,
					FeaturedImage: resp110.FeaturedImage,
				}
				tt.want = want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = c
			tt.prepare(&tt)
			got, err := c.GetMainCategoriesMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetCategoriesMapByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCategoryImpl_GetMainCategoriesMapWhenMainCategoryDoesNotExist(t *testing.T) {

	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}

	type TC struct {
		name    string
		fields  fields
		want    map[string]schema.GetMainCategoriesMapResp
		wantErr bool
		prepare func(*TC)
	}
	tests := []TC{
		{
			name: "[Ok] When no main category exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				tt.want = make(map[string]schema.GetMainCategoriesMapResp)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = c
			tt.prepare(&tt)
			got, err := c.GetMainCategoriesMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetCategoriesMapByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCategoryImpl_GetMainParentCategories(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}

	type TC struct {
		name    string
		fields  fields
		want    []schema.GetParentCategoriesResp
		wantErr bool
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parent1Opts := schema.GetRandomCreateCategoryOpts()
				parent1Opts.ParentID = primitive.NilObjectID
				parent1Opts.IsMain = true
				parent1Resp, err := tt.fields.App.Category.CreateCategory(parent1Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				parent2Opts := schema.GetRandomCreateCategoryOpts()
				parent2Opts.ParentID = primitive.NilObjectID
				parent2Opts.IsMain = false
				parent2Resp, err := tt.fields.App.Category.CreateCategory(parent2Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				parent3Opts := schema.GetRandomCreateCategoryOpts()
				parent3Opts.ParentID = primitive.NilObjectID
				parent3Opts.IsMain = true
				parent3Resp, err := tt.fields.App.Category.CreateCategory(parent3Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children1Opts := schema.GetRandomCreateCategoryOpts()
				children1Opts.ParentID = parent1Resp.ID
				children1Opts.IsMain = true
				_, _ = tt.fields.App.Category.CreateCategory(children1Opts)

				children2Opts := schema.GetRandomCreateCategoryOpts()
				children2Opts.ParentID = parent2Resp.ID
				children1Opts.IsMain = true
				_, _ = tt.fields.App.Category.CreateCategory(children2Opts)

				want := []schema.GetParentCategoriesResp{
					{
						ID:        parent1Resp.ID,
						Name:      parent1Resp.Name,
						Thumbnail: parent1Resp.Thumbnail,
					},
					{
						ID:        parent3Resp.ID,
						Name:      parent3Resp.Name,
						Thumbnail: parent3Resp.Thumbnail,
					},
				}
				tt.want = want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = ci
			tt.prepare(&tt)
			got, err := ci.GetMainParentCategories()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetParentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryImpl.GetParentCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryImpl_GetMainCategoriesByParentID(t *testing.T) {
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
		want    []schema.GetMainCategoriesByParentIDResp
		wantErr bool
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parent1Opts := schema.GetRandomCreateCategoryOpts()
				parent1Opts.ParentID = primitive.NilObjectID
				parent1Opts.IsMain = true
				parent1Resp, err := tt.fields.App.Category.CreateCategory(parent1Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children10Opts := schema.GetRandomCreateCategoryOpts()
				children10Opts.ParentID = parent1Resp.ID
				children10Opts.IsMain = true
				children10Resp, err := tt.fields.App.Category.CreateCategory(children10Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children12Opts := schema.GetRandomCreateCategoryOpts()
				children12Opts.ParentID = parent1Resp.ID
				children12Opts.IsMain = false
				_, _ = tt.fields.App.Category.CreateCategory(children12Opts)

				children13Opts := schema.GetRandomCreateCategoryOpts()
				children13Opts.ParentID = parent1Resp.ID
				children13Opts.IsMain = true
				children13Resp, err := tt.fields.App.Category.CreateCategory(children13Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				parent2Opts := schema.GetRandomCreateCategoryOpts()
				parent2Opts.ParentID = primitive.NilObjectID
				parent2Opts.IsMain = true
				parent2Resp, err := tt.fields.App.Category.CreateCategory(parent2Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children21Opts := schema.GetRandomCreateCategoryOpts()
				children21Opts.ParentID = parent2Resp.ID
				children21Opts.IsMain = true
				_, _ = tt.fields.App.Category.CreateCategory(children21Opts)

				want := []schema.GetMainCategoriesByParentIDResp{
					{
						ID:            children10Resp.ID,
						Name:          children10Resp.Name,
						FeaturedImage: children10Resp.FeaturedImage,
					},
					{
						ID:            children13Resp.ID,
						Name:          children13Resp.Name,
						FeaturedImage: children13Resp.FeaturedImage,
					},
				}
				tt.args = args{
					id: parent1Resp.ID,
				}
				tt.want = want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = ci
			tt.prepare(&tt)
			got, err := ci.GetMainCategoriesByParentID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetMainCategoriesByParentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCategoryImpl_GetSubCategoriesByParentID(t *testing.T) {
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
		want    []schema.GetSubCategoriesByParentIDResp
		wantErr bool
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parent1Opts := schema.GetRandomCreateCategoryOpts()
				parent1Opts.ParentID = primitive.NilObjectID
				parent1Opts.IsMain = true
				parent1Resp, err := tt.fields.App.Category.CreateCategory(parent1Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children11Opts := schema.GetRandomCreateCategoryOpts()
				children11Opts.ParentID = parent1Resp.ID
				children11Opts.IsMain = true
				children11Resp, err := tt.fields.App.Category.CreateCategory(children11Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children111Opts := schema.GetRandomCreateCategoryOpts()
				children111Opts.ParentID = children11Resp.ID
				children111Opts.IsMain = true
				children111Resp, err := tt.fields.App.Category.CreateCategory(children111Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children112Opts := schema.GetRandomCreateCategoryOpts()
				children112Opts.ParentID = children11Resp.ID
				children112Opts.IsMain = true
				children112Resp, err := tt.fields.App.Category.CreateCategory(children112Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children113Opts := schema.GetRandomCreateCategoryOpts()
				children113Opts.ParentID = children11Resp.ID
				children113Opts.IsMain = false
				_, _ = tt.fields.App.Category.CreateCategory(children113Opts)

				children12Opts := schema.GetRandomCreateCategoryOpts()
				children12Opts.ParentID = parent1Resp.ID
				children12Opts.IsMain = true
				_, _ = tt.fields.App.Category.CreateCategory(children12Opts)

				want := []schema.GetSubCategoriesByParentIDResp{
					{
						ID:   children111Resp.ID,
						Name: children111Resp.Name,
					},
					{
						ID:   children112Resp.ID,
						Name: children112Resp.Name,
					},
				}
				tt.args = args{
					id: children11Resp.ID,
				}
				tt.want = want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = ci
			tt.prepare(&tt)
			got, err := ci.GetSubCategoriesByParentID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetSubCategoriesByParentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCategoryImpl_GetCategories(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}

	type TC struct {
		name    string
		fields  fields
		want    []schema.GetCategoriesResp
		wantErr bool
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parent1Opts := schema.GetRandomCreateCategoryOpts()
				parent1Opts.ParentID = primitive.NilObjectID
				parent1Opts.IsMain = true
				parent1Resp, err := tt.fields.App.Category.CreateCategory(parent1Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children11Opts := schema.GetRandomCreateCategoryOpts()
				children11Opts.ParentID = parent1Resp.ID
				children11Opts.IsMain = true
				children11Resp, err := tt.fields.App.Category.CreateCategory(children11Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children111Opts := schema.GetRandomCreateCategoryOpts()
				children111Opts.ParentID = children11Resp.ID
				children111Opts.IsMain = true
				children111Resp, err := tt.fields.App.Category.CreateCategory(children111Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children112Opts := schema.GetRandomCreateCategoryOpts()
				children112Opts.ParentID = children11Resp.ID
				children112Opts.IsMain = true
				children112Resp, err := tt.fields.App.Category.CreateCategory(children112Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children113Opts := schema.GetRandomCreateCategoryOpts()
				children113Opts.ParentID = children11Resp.ID
				children113Opts.IsMain = false
				children113Resp, err := tt.fields.App.Category.CreateCategory(children113Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children12Opts := schema.GetRandomCreateCategoryOpts()
				children12Opts.ParentID = parent1Resp.ID
				children12Opts.IsMain = true
				children12Resp, err := tt.fields.App.Category.CreateCategory(children12Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				parent2Opts := schema.GetRandomCreateCategoryOpts()
				parent2Opts.ParentID = primitive.NilObjectID
				parent2Opts.IsMain = true
				parent2Resp, err := tt.fields.App.Category.CreateCategory(parent2Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children21Opts := schema.GetRandomCreateCategoryOpts()
				children21Opts.ParentID = parent1Resp.ID
				children21Opts.IsMain = true
				children21Resp, err := tt.fields.App.Category.CreateCategory(children21Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				want := []schema.GetCategoriesResp{
					{
						ID:            parent1Resp.ID,
						Name:          parent1Resp.Name,
						ParentID:      parent1Resp.ParentID,
						AncestorID:    parent1Resp.AncestorID,
						Thumbnail:     parent1Resp.Thumbnail,
						FeaturedImage: parent1Resp.FeaturedImage,
						IsMain:        parent1Resp.IsMain,
					},
					{
						ID:            children11Resp.ID,
						Name:          children11Resp.Name,
						ParentID:      children11Resp.ParentID,
						AncestorID:    children11Resp.AncestorID,
						Thumbnail:     children11Resp.Thumbnail,
						FeaturedImage: children11Resp.FeaturedImage,
						IsMain:        children11Resp.IsMain,
					},
					{
						ID:            children111Resp.ID,
						Name:          children111Resp.Name,
						ParentID:      children111Resp.ParentID,
						AncestorID:    children111Resp.AncestorID,
						Thumbnail:     children111Resp.Thumbnail,
						FeaturedImage: children111Resp.FeaturedImage,
						IsMain:        children111Resp.IsMain,
					},
					{
						ID:            children112Resp.ID,
						Name:          children112Resp.Name,
						ParentID:      children112Resp.ParentID,
						AncestorID:    children112Resp.AncestorID,
						Thumbnail:     children112Resp.Thumbnail,
						FeaturedImage: children112Resp.FeaturedImage,
						IsMain:        children112Resp.IsMain,
					},
					{
						ID:            children113Resp.ID,
						Name:          children113Resp.Name,
						ParentID:      children113Resp.ParentID,
						AncestorID:    children113Resp.AncestorID,
						Thumbnail:     children113Resp.Thumbnail,
						FeaturedImage: children113Resp.FeaturedImage,
						IsMain:        children113Resp.IsMain,
					},
					{
						ID:            children12Resp.ID,
						Name:          children12Resp.Name,
						ParentID:      children12Resp.ParentID,
						AncestorID:    children12Resp.AncestorID,
						Thumbnail:     children12Resp.Thumbnail,
						FeaturedImage: children12Resp.FeaturedImage,
						IsMain:        children12Resp.IsMain,
					},
					{
						ID:            parent2Resp.ID,
						Name:          parent2Resp.Name,
						ParentID:      parent2Resp.ParentID,
						AncestorID:    parent2Resp.AncestorID,
						Thumbnail:     parent2Resp.Thumbnail,
						FeaturedImage: parent2Resp.FeaturedImage,
						IsMain:        parent2Resp.IsMain,
					},
					{
						ID:            children21Resp.ID,
						Name:          children21Resp.Name,
						ParentID:      children21Resp.ParentID,
						AncestorID:    children21Resp.AncestorID,
						Thumbnail:     children21Resp.Thumbnail,
						FeaturedImage: children21Resp.FeaturedImage,
						IsMain:        children12Resp.IsMain,
					},
				}
				tt.want = want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = ci
			tt.prepare(&tt)
			got, err := ci.GetCategories()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCategoryImpl_GetCategoriesBasic(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}

	type TC struct {
		name    string
		fields  fields
		want    []schema.GetCategoriesBasicResp
		wantErr bool
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.CategoryConfig.DBName),
				Logger: app.Logger,
			},
			wantErr: false,
			prepare: func(tt *TC) {
				parent1Opts := schema.GetRandomCreateCategoryOpts()
				parent1Opts.ParentID = primitive.NilObjectID
				parent1Opts.IsMain = true
				parent1Resp, err := tt.fields.App.Category.CreateCategory(parent1Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children11Opts := schema.GetRandomCreateCategoryOpts()
				children11Opts.ParentID = parent1Resp.ID
				children11Opts.IsMain = true
				children11Resp, err := tt.fields.App.Category.CreateCategory(children11Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children111Opts := schema.GetRandomCreateCategoryOpts()
				children111Opts.ParentID = children11Resp.ID
				children111Opts.IsMain = true
				children111Resp, err := tt.fields.App.Category.CreateCategory(children111Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children112Opts := schema.GetRandomCreateCategoryOpts()
				children112Opts.ParentID = children11Resp.ID
				children112Opts.IsMain = true
				children112Resp, err := tt.fields.App.Category.CreateCategory(children112Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children113Opts := schema.GetRandomCreateCategoryOpts()
				children113Opts.ParentID = children11Resp.ID
				children113Opts.IsMain = false
				children113Resp, err := tt.fields.App.Category.CreateCategory(children113Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children12Opts := schema.GetRandomCreateCategoryOpts()
				children12Opts.ParentID = parent1Resp.ID
				children12Opts.IsMain = true
				children12Resp, err := tt.fields.App.Category.CreateCategory(children12Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				parent2Opts := schema.GetRandomCreateCategoryOpts()
				parent2Opts.ParentID = primitive.NilObjectID
				parent2Opts.IsMain = true
				parent2Resp, err := tt.fields.App.Category.CreateCategory(parent2Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				children21Opts := schema.GetRandomCreateCategoryOpts()
				children21Opts.ParentID = parent1Resp.ID
				children21Opts.IsMain = true
				children21Resp, err := tt.fields.App.Category.CreateCategory(children21Opts)
				if err != nil {
					log.Fatalf("%s", err)
				}

				want := []schema.GetCategoriesBasicResp{
					{
						ID:     parent1Resp.ID,
						Name:   parent1Resp.Name,
						IsMain: parent1Resp.IsMain,
					},
					{
						ID:     children11Resp.ID,
						Name:   children11Resp.Name,
						IsMain: children11Resp.IsMain,
					},
					{
						ID:     children111Resp.ID,
						Name:   children111Resp.Name,
						IsMain: children111Resp.IsMain,
					},
					{
						ID:     children112Resp.ID,
						Name:   children112Resp.Name,
						IsMain: children112Resp.IsMain,
					},
					{
						ID:     children113Resp.ID,
						Name:   children113Resp.Name,
						IsMain: children113Resp.IsMain,
					},
					{
						ID:     children12Resp.ID,
						Name:   children12Resp.Name,
						IsMain: children12Resp.IsMain,
					},
					{
						ID:     parent2Resp.ID,
						Name:   parent2Resp.Name,
						IsMain: parent2Resp.IsMain,
					},
					{
						ID:     children21Resp.ID,
						Name:   children21Resp.Name,
						IsMain: children12Resp.IsMain,
					},
				}
				tt.want = want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Category = ci
			tt.prepare(&tt)
			got, err := ci.GetCategoriesBasic()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
