package app

import (
	"context"
	"fmt"
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
	"syreclabs.com/go/faker"
)

func TestCollectionImpl_CreateCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts *schema.CreateCollectionOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        []error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []schema.SubCollectionOpts{
		{
			Name:       faker.Name().Name(),
			Image:      faker.Avatar().Url("png", 100, 100),
			CatalogIDs: catalogIDs[0:5],
		},
		{
			Name:       faker.Name().Name(),
			Image:      faker.Avatar().Url("png", 100, 100),
			CatalogIDs: catalogIDs[5:10],
		},
		{
			Name:       faker.Name().Name(),
			Image:      faker.Avatar().Url("png", 100, 100),
			CatalogIDs: catalogIDs[10:15],
		},
		{
			Name:       faker.Name().Name(),
			Image:      faker.Avatar().Url("png", 100, 100),
			CatalogIDs: catalogIDs[15:20],
		},
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

				var catalogCalls []*gomock.Call
				for i := 0; i < 4; i++ {
					call := kc.EXPECT().GetCatalogByIDs(gomock.Any(), catalogIDs[i*5:i*5+5]).Return(catalogs[i*5:i*5+5], nil)
					catalogCalls = append(catalogCalls, call)
				}
				gomock.InOrder(catalogCalls...)

			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.CreateCollectionOpts{
					Type:          model.ProductCollection,
					Genders:       []string{"Male"},
					Title:         "Collection Test",
					SubCollection: subCollections,
				}
				tt.wantErr = false
			},
		},
		{
			name: "CatalogID Does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

				kc.EXPECT().GetCatalogByIDs(gomock.Any(), catalogIDs[:5]).Times(1).Return(catalogs[1:5], nil)
			},
			prepare: func(tt *TC) {

				tt.args.opts = &schema.CreateCollectionOpts{
					Type:          model.ProductCollection,
					Genders:       []string{"Male"},
					Title:         "Collection Test",
					SubCollection: subCollections,
				}
				tt.wantErr = true
				tt.err = []error{errors.Errorf("catalog with id: %s not found", catalogs[0].ID.Hex())}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			resp, errArray := ci.CreateCollection(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, errArray)
				assert.NotNil(t, resp)
				// assert.Equal(t, tt.want, resp)
				var collection model.Collection
				ci.DB.Collection(model.CollectionColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&collection)
				assert.Equal(t, resp.ID, collection.ID)
				assert.NotNil(t, collection.Name)
				assert.Equal(t, resp.Genders, collection.Genders)
				assert.Equal(t, resp.Title, collection.Title)
				assert.Equal(t, resp.Type, collection.Type)
				assert.Equal(t, len(resp.SubCollections), len(collection.SubCollections))
				for i := 0; i < len(resp.SubCollections); i++ {
					resp.SubCollections[i].CreatedAt = time.Time{}
					collection.SubCollections[i].CreatedAt = time.Time{}
					assert.Equal(t, resp.SubCollections[i], collection.SubCollections[i])
				}
			}

			if tt.wantErr {
				assert.NotNil(t, errArray)
				assert.Equal(t, tt.err[0].Error(), errArray[0].Error())
			}
		})
	}
}
func TestCollectionImpl_DeleteCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		id primitive.ObjectID
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[10:15],
		},
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[15:20],
		},
	}
	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	collection = model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "B",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				id: collection.ID,
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
			},
		},
		{
			name: "[Error] Collection doesn't exist'",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				id: primitive.NewObjectID(),
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = errors.Errorf("unable to delete collection with id: %s", tt.args.id.Hex())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ci.DeleteCollection(tt.args.id)
			if !tt.wantErr {
				assert.Nil(t, err)
				// assert.Equal(t, tt.want, resp)
				var collection model.Collection
				ci.DB.Collection(model.CollectionColl).FindOne(context.TODO(), bson.M{"_id": tt.args.id}).Decode(&collection)
				assert.Equal(t, collection, model.Collection{})

			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), tt.err.Error())
			}
		})
	}
}

func TestCollectionImpl_AddSubCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts *schema.AddSubCollectionOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        []error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Avatar().Url("png", 100, 100)},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Avatar().Url("png", 100, 100)},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			Name:       faker.Name().Name(),
			Image:      &model.IMG{SRC: faker.Avatar().Url("png", 100, 100)},
			CatalogIDs: catalogIDs[10:15],
		},
	}

	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddSubCollectionOpts{
					ID: collection.ID,
					SubCollection: &schema.SubCollectionOpts{
						Name:       "New Sub Collection",
						Image:      faker.Avatar().Url("png", 100, 100),
						CatalogIDs: catalogIDs[15:20],
					},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.SubCollection.CatalogIDs).Times(1).Return(catalogs[15:20], nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
			},
		},
		{
			name: "[Error] collection id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddSubCollectionOpts{
					ID: primitive.NewObjectID(),
					SubCollection: &schema.SubCollectionOpts{
						Name:       "New Sub Collection",
						Image:      faker.Avatar().Url("png", 100, 100),
						CatalogIDs: catalogIDs[15:20],
					},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.SubCollection.CatalogIDs).Times(1).Return(catalogs[15:20], nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("collection with id:%s not found", tt.args.opts.ID.Hex())}
			},
		},
		{
			name: "[Error] catalog id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddSubCollectionOpts{
					ID: primitive.NewObjectID(),
					SubCollection: &schema.SubCollectionOpts{
						Name:       "New Sub Collection",
						Image:      faker.Avatar().Url("png", 100, 100),
						CatalogIDs: []primitive.ObjectID{primitive.NewObjectID(), catalogIDs[12]},
					},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.SubCollection.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[12]}, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("catalog with id: %s not found", tt.args.opts.SubCollection.CatalogIDs[0].Hex())}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			_, err := ci.AddSubCollection(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)
				// assert.Equal(t, tt.want, resp)
				var collection model.Collection
				ci.DB.Collection(model.CollectionColl).FindOne(context.TODO(), bson.M{"_id": tt.args.opts.ID}).Decode(&collection)
				// assert.Equal(t, resp.SubCollection, collection.SubCollection)

			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err[0].Error(), tt.err[0].Error())
			}
		})
	}
}

func TestCollectionImpl_DeleteSubCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	// defer CleanTestApp(app)

	type args struct {
		calID primitive.ObjectID
		subID primitive.ObjectID
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			ID:         primitive.NewObjectID(),
			Name:       "a",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "b",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "c",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[10:15],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "d",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[15:20],
		},
	}
	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				calID: collection.ID,
				subID: collection.SubCollections[0].ID,
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
			},
		},
		{
			name: "[Error] CollectionID doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				calID: primitive.NewObjectID(),
				subID: collection.SubCollections[0].ID,
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = errors.Errorf("unable to find collection with id - %s ", tt.args.calID)
			},
		},
		{
			name: "[Error] Sub Collection doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				calID: collection.ID,
				subID: primitive.NewObjectID(),
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = errors.Errorf("unable to delete sub collection")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ci.DeleteSubCollection(tt.args.calID, tt.args.subID)
			if !tt.wantErr {
				assert.Nil(t, err)
				// assert.Equal(t, tt.want, resp)
				var collection model.Collection
				ci.DB.Collection(model.CollectionColl).FindOne(context.TODO(), bson.M{"_id": tt.args.calID}).Decode(&collection)
				for _, sc := range collection.SubCollections {
					assert.NotEqual(t, sc.ID, tt.args.subID)
				}
				assert.Equal(t, len(collection.SubCollections), len(subCollections)-1)

			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), tt.err.Error())
			}
		})
	}
}

func TestCollectionImpl_EditCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts *schema.EditCollectionOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       *schema.CreateCollectionResp
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			ID:         primitive.NewObjectID(),
			Name:       "a",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "b",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "c",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[10:15],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "d",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[15:20],
		},
	}
	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		Genders:        []string{"M"},
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.EditCollectionOpts{
					ID:      collection.ID,
					Genders: []string{"F", "M"},
					Title:   "New Title",
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
				if tt.args.opts.Genders != nil {
					collection.Genders = tt.args.opts.Genders
				}
				if tt.args.opts.Title != "" {
					collection.Title = tt.args.opts.Title
				}
				want := &schema.CreateCollectionResp{
					ID:             collection.ID,
					Type:           collection.Type,
					Genders:        collection.Genders,
					Title:          collection.Title,
					Name:           collection.Name,
					SubCollections: collection.SubCollections,
				}
				tt.want = want
			},
		},
		{
			name: "[Ok] Only Title change",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.EditCollectionOpts{
					ID:    collection.ID,
					Title: "New Title",
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
				if tt.args.opts.Genders != nil {
					collection.Genders = tt.args.opts.Genders
				}
				if tt.args.opts.Title != "" {
					collection.Title = tt.args.opts.Title
				}
				want := &schema.CreateCollectionResp{
					ID:             collection.ID,
					Type:           collection.Type,
					Genders:        collection.Genders,
					Title:          collection.Title,
					Name:           collection.Name,
					SubCollections: collection.SubCollections,
				}
				tt.want = want
			},
		},
		{
			name: "[Error] Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.EditCollectionOpts{
					ID:      primitive.NewObjectID(),
					Genders: []string{"F"},
					Title:   "New Title",
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = errors.Errorf("catalog with id:%s not found", tt.args.opts.ID.Hex())

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			resp, err := ci.EditCollection(tt.args.opts)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, resp)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), tt.err.Error())
			}
		})
	}
}

func TestCollectionImpl_UpdateSubCollectionImage(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts *schema.UpdateSubCollectionImageOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       *schema.CreateCollectionResp
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			ID:         primitive.NewObjectID(),
			Name:       "a",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "b",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "c",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[10:15],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "d",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[15:20],
		},
	}
	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateSubCollectionImageOpts{
					ColID: collection.ID,
					SubID: collection.SubCollections[0].ID,
					Image: "https://www.agencyreporter.com/wp-content/uploads/2021/02/HYPD-Store-raises-pre-seed-strategic-investment-from-ScoopWhoop.jpg",
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {

			},
			wantErr: false,
		},
		{
			name: "[Error] Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateSubCollectionImageOpts{
					ColID: primitive.NewObjectID(),
					SubID: collection.SubCollections[0].ID,
					Image: "https://www.agencyreporter.com/wp-content/uploads/2021/02/HYPD-Store-raises-pre-seed-strategic-investment-from-ScoopWhoop.jpg",
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = errors.Errorf("unable to find the sub collection with id - %s", tt.args.opts.SubID.Hex())

			},
		},
		{
			name: "[Error] Sub Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateSubCollectionImageOpts{
					ColID: collection.ID,
					SubID: primitive.NewObjectID(),
					Image: "https://www.agencyreporter.com/wp-content/uploads/2021/02/HYPD-Store-raises-pre-seed-strategic-investment-from-ScoopWhoop.jpg",
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = errors.Errorf("unable to find the sub collection with id - %s", tt.args.opts.SubID.Hex())

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ci.UpdateSubCollectionImage(tt.args.opts)
			fmt.Println(tt.want)

			if !tt.wantErr {
				assert.Nil(t, err)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), tt.err.Error())
			}
		})
	}
}

func TestCollectionImpl_AddCatalogsToSubCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts *schema.UpdateCatalogsInSubCollectionOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        []error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			ID:         primitive.NewObjectID(),
			Name:       "a",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "b",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "c",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[10:15],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "d",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[15:20],
		},
	}
	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{

		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[11], catalogIDs[12]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[11], catalogs[12]}, nil)

			},
			prepare: func(tt *TC) {

			},
			wantErr: false,
		},
		{
			name: "[Ok] with one same catalog id",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[1], catalogIDs[14]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[14]}, nil)

			},
			prepare: func(tt *TC) {

			},
			wantErr: false,
		},
		{
			name: "[Error] with already added catalog ids",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[1], catalogIDs[2]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				errArray := []error{
					errors.Errorf("unable to the update the sub collection with id - %s", tt.args.opts.SubID.Hex()),
				}
				tt.err = errArray
			},
			wantErr: true,
		},
		{
			name: "[Error] catalog id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{primitive.NewObjectID(), catalogIDs[2]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				errArray := []error{
					errors.Errorf("catalog with id: %s not found", tt.args.opts.CatalogIDs[0].Hex()),
				}
				tt.err = errArray
			},
			wantErr: true,
		},
		{
			name: "[Error] Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      primitive.NewObjectID(),
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[11], catalogIDs[12]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("unable to find the sub collection with id - %s", tt.args.opts.SubID.Hex())}

			},
		},
		{
			name: "[Error] Sub Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      primitive.NewObjectID(),
					CatalogIDs: []primitive.ObjectID{catalogIDs[11], catalogIDs[12]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("unable to find the sub collection with id - %s", tt.args.opts.SubID.Hex())}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ci.AddCatalogsToSubCollection(tt.args.opts)
			fmt.Println(tt.want)

			if !tt.wantErr {
				assert.Nil(t, err)
				var collection model.Collection
				ci.DB.Collection(model.CollectionColl).FindOne(context.TODO(), bson.M{"_id": tt.args.opts.ColID}).Decode(&collection)
				i := 0
				for _, ci := range collection.SubCollections[0].CatalogIDs {
					if ci == tt.args.opts.CatalogIDs[0] {
						i++
					}
					if ci == tt.args.opts.CatalogIDs[1] {
						i++
					}
				}
				assert.Equal(t, i, 2)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err[0].Error(), tt.err[0].Error())
			}
		})
	}
}
func TestCollectionImpl_RemoveCatalogsFromSubCollection(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts *schema.UpdateCatalogsInSubCollectionOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        []error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	var catalogs []schema.GetCatalogResp
	var catalogIDs []primitive.ObjectID
	for i := 0; i < 20; i++ {
		catalogs = append(catalogs, schema.GetCatalogResp{
			ID: primitive.NewObjectID(),
		})
		catalogIDs = append(catalogIDs, catalogs[i].ID)
	}
	subCollections := []model.SubCollection{
		{
			ID:         primitive.NewObjectID(),
			Name:       "a",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[0:5],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "b",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[5:10],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "c",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[10:15],
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "d",
			Image:      &model.IMG{SRC: faker.Internet().Url()},
			CatalogIDs: catalogIDs[15:20],
		},
	}
	collection := model.Collection{
		ID:             primitive.NewObjectID(),
		Name:           "A",
		Type:           model.ProductCollection,
		SubCollections: subCollections,
	}
	app.MongoDB.Client.Database(app.Config.GroupConfig.DBName).Collection(model.CollectionColl).InsertOne(context.TODO(), collection)

	tests := []TC{

		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[1], catalogIDs[2]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {

			},
			wantErr: false,
		},
		{
			name: "[Ok] with one out of two catalog ids in sub collections ",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[3], catalogIDs[14]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[3], catalogs[14]}, nil)

			},
			prepare: func(tt *TC) {

			},
			wantErr: false,
		},
		{
			name: "[Error] catalog ids not in sub collection",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[11], catalogIDs[12]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[11], catalogs[12]}, nil)

			},
			prepare: func(tt *TC) {
				errArray := []error{
					errors.Errorf("unable to the update the sub collection with id - %s", tt.args.opts.SubID.Hex()),
				}
				tt.err = errArray
			},
			wantErr: true,
		},
		{
			name: "[Error] catalog id doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{primitive.NewObjectID(), catalogIDs[2]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				errArray := []error{
					errors.Errorf("catalog with id: %s not found", tt.args.opts.CatalogIDs[0].Hex()),
				}
				tt.err = errArray
			},
			wantErr: true,
		},
		{
			name: "[Error] Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      primitive.NewObjectID(),
					SubID:      collection.SubCollections[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogIDs[11], catalogIDs[12]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("unable to find the sub collection with id - %s", tt.args.opts.SubID.Hex())}

			},
		},
		{
			name: "[Error] Sub Collection ID does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.UpdateCatalogsInSubCollectionOpts{
					ColID:      collection.ID,
					SubID:      primitive.NewObjectID(),
					CatalogIDs: []primitive.ObjectID{catalogIDs[11], catalogIDs[12]},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return([]schema.GetCatalogResp{catalogs[1], catalogs[2]}, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("unable to find the sub collection with id - %s", tt.args.opts.SubID.Hex())}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ci := &CollectionImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Collection = ci
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ci.RemoveCatalogsFromSubCollection(tt.args.opts)
			fmt.Println(tt.want)

			if !tt.wantErr {
				assert.Nil(t, err)
				var collection model.Collection
				ci.DB.Collection(model.CollectionColl).FindOne(context.TODO(), bson.M{"_id": tt.args.opts.ColID}).Decode(&collection)
				i := 0
				for _, ci := range collection.SubCollections[0].CatalogIDs {
					if ci == catalogIDs[1] {
						i++
					}
					if ci == catalogIDs[2] {
						i++
					}
				}
				assert.Equal(t, i, 0)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, err[0].Error(), tt.err[0].Error())
			}
		})
	}
}
