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

func TestGroupImpl_CreateCatalogGroup(t *testing.T) {

	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type args struct {
		opts       *schema.CreateCatalogGroupOpts
		catalogIDs []string
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
		err        []error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
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
				groupOpts := schema.GetRandomCreateGroupOpts()
				var catalogs []schema.GetCatalogResp
				var catalogIDs []primitive.ObjectID
				for i := 0; i < len(groupOpts.IDs); i++ {
					catalogs = append(catalogs, schema.GetCatalogResp{
						ID: groupOpts.IDs[i],
					})
					catalogIDs = append(catalogIDs, groupOpts.IDs[i])
				}
				tt.args.opts = groupOpts

				kc.EXPECT().GetCatalogByIDs(gomock.Any(), catalogIDs).Times(1).Return(catalogs, nil)

			},
			prepare: func(tt *TC) {
				// opts := schema.CreateCatalogGroupOpts{
				// 	Name:     "Test-Group",
				// 	Catalogs: tt.args.catalogIDs,
				// }
				// tt.args.opts = &opts
			},
		},
		{
			name: "[Error] Few Catalog IDs doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				groupOpts := schema.GetRandomCreateGroupOpts()
				var catalogs []schema.GetCatalogResp
				var catalogIDs []primitive.ObjectID
				for i := 0; i < len(groupOpts.IDs); i++ {
					catalogs = append(catalogs, schema.GetCatalogResp{
						ID: groupOpts.IDs[i],
					})
					catalogIDs = append(catalogIDs, groupOpts.IDs[i])
				}
				tt.args.opts = groupOpts
				//Skipping 0th Index
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), catalogIDs).Times(1).Return(catalogs[1:], nil)
				tt.err = []error{errors.Errorf("catalog with id: %s not found", catalogs[0].ID.Hex())}

			},
			prepare: func(tt *TC) {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = g
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			fmt.Println(tt.args.opts)
			resp, errArray := g.CreateCatalogGroup(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, errArray)
				assert.NotNil(t, resp)
			}

			if tt.wantErr {
				assert.NotNil(t, errArray)
				assert.Equal(t, tt.err[0].Error(), errArray[0].Error())
			}
		})
	}
}

func TestGroupImpl_GetGroups(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	var catalogs []model.Catalog
	var groups []model.Group
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
		})
		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	for i := 0; i < itemLen-1; i++ {

		group := model.Group{
			ID:    primitive.NewObjectID(),
			Basis: faker.Company().Name(),
			Status: &model.GroupStatus{
				Value: faker.RandomChoice([]string{model.Archive, model.Unlist, model.Publish}),
			},
			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}

	type args struct {
		opts *schema.GetGroupsOpts
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
		want       []schema.GroupResp
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	tests := []TC{

		{
			name: "[Ok] All",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GetGroupsOpts{
					Page: 0,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < 10; i++ {
					max := catalogs[i].RetailPrice.Value
					min := catalogs[(i+1)%itemLen].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[(i+1)%itemLen].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[(i+1)%itemLen].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[(i+1)%itemLen]},
					})
				}
				tt.want = groupResps
			},
		},
		{
			name: "[OK] All - Page 1",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GetGroupsOpts{
					Page: 1,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := tt.args.opts.Page * 10; i < (tt.args.opts.Page*10)+10; i++ {
					max := catalogs[i].RetailPrice.Value
					min := catalogs[(i + 1)].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[(i+1)].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[(i + 1)].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[(i + 1)]},
					})
				}
				tt.want = groupResps
			},
		},
		{
			name: "[Ok] Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GetGroupsOpts{
					Page:   0,
					Status: model.Unlist,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].Status.Value != model.Unlist {
						continue
					}
					max := catalogs[i].RetailPrice.Value
					min := catalogs[(i+1)%itemLen].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[(i+1)%itemLen].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[(i+1)%itemLen].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[(i+1)%itemLen]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
		{
			name: "[Ok] Publish",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GetGroupsOpts{
					Page:   0,
					Status: model.Publish,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].Status.Value != model.Publish {
						continue
					}
					max := catalogs[i].RetailPrice.Value
					min := catalogs[(i+1)%itemLen].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[(i+1)%itemLen].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[(i+1)%itemLen].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[(i+1)%itemLen]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
		{
			name: "[Ok] Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GetGroupsOpts{
					Page:   0,
					Status: model.Archive,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].Status.Value != model.Archive {
						continue
					}
					max := catalogs[i].RetailPrice.Value
					min := catalogs[(i+1)%itemLen].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[(i+1)%itemLen].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[(i+1)%itemLen].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[(i+1)%itemLen]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gi := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = gi
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)
			fmt.Println(tt.args.opts.Page)

			resp, err := gi.GetGroups(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)

				for i := 0; i < len(resp); i++ {
					assert.Equal(t, tt.want[i], resp[i])
				}
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestGroupImpl_GetGroupsByCatalogID(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)
	var catalogs []model.Catalog
	var groups []model.Group
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
			FeaturedImage: &model.IMG{
				SRC: faker.Avatar().Url("png", 100, 100),
			},
		})
		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	for i := 0; i < 3; i++ {

		group := model.Group{
			ID:     primitive.NewObjectID(),
			Basis:  faker.Company().Name(),
			Status: &model.GroupStatus{Value: model.Publish},

			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}

	type args struct {
		opts *schema.GetGroupsByCatalogIDOpts
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
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		want       []schema.GetGroupsByCatalogIDResp
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.GetGroupsByCatalogIDOpts{
					ID:   catalogs[1].ID,
					Page: 0,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GetGroupsByCatalogIDResp

				for i := 0; i < len(groups); i++ {

					if groups[i].CatalogIDs[0] != tt.args.opts.ID && groups[i].CatalogIDs[1] != tt.args.opts.ID {
						continue
					}
					if groups[i].Status.Value != model.Publish {
						continue
					}

					groupResps = append(groupResps, schema.GetGroupsByCatalogIDResp{
						ID:          groups[i].ID,
						Basis:       groups[i].Basis,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[i+1]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gi := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = gi
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			resp, err := gi.GetGroupsByCatalogID(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)
				// assert.NotNil(t, resp)
				assert.Equal(t, len(resp), len(tt.want))
				for i := 0; i < len(tt.want); i++ {
					assert.Equal(t, tt.want[i], resp[i])
				}
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestGroupImpl_KeeperGetGroupsByCatalogID(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)
	var catalogs []model.Catalog
	var groups []model.Group
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
		})
		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	for i := 0; i < itemLen-1; i++ {

		group := model.Group{
			ID:     primitive.NewObjectID(),
			Basis:  faker.Company().Name(),
			Status: &model.GroupStatus{Value: faker.RandomChoice([]string{model.Archive, model.Unlist, model.Publish})},

			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}

	type args struct {
		opts *schema.KeeperGetGroupsByCatalogIDOpts
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
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		want       []schema.GroupResp
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}

	tests := []TC{
		{
			name: "[Ok] All",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.KeeperGetGroupsByCatalogIDOpts{
					ID:   catalogs[1].ID,
					Page: 0,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].CatalogIDs[0] != tt.args.opts.ID && groups[i].CatalogIDs[1] != tt.args.opts.ID {
						continue
					}

					max := catalogs[i].RetailPrice.Value
					min := catalogs[i+1].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[i+1].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[i+1].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[i+1]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
		{
			name: "[Ok] Published",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.KeeperGetGroupsByCatalogIDOpts{
					ID:     catalogs[1].ID,
					Page:   0,
					Status: model.Publish,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].Status.Value != model.Publish {
						continue
					}
					if groups[i].CatalogIDs[0] != tt.args.opts.ID && groups[i].CatalogIDs[1] != tt.args.opts.ID {
						continue
					}

					max := catalogs[i].RetailPrice.Value
					min := catalogs[i+1].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[i+1].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[i+1].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[i+1]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
		{
			name: "[Ok] Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.KeeperGetGroupsByCatalogIDOpts{
					ID:     catalogs[1].ID,
					Page:   0,
					Status: model.Unlist,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].Status.Value != model.Unlist {
						continue
					}
					if groups[i].CatalogIDs[0] != tt.args.opts.ID && groups[i].CatalogIDs[1] != tt.args.opts.ID {
						continue
					}

					max := catalogs[i].RetailPrice.Value
					min := catalogs[i+1].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[i+1].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[i+1].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[i+1]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
		{
			name: "[Ok] Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.KeeperGetGroupsByCatalogIDOpts{
					ID:     catalogs[1].ID,
					Page:   0,
					Status: model.Archive,
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				var groupResps []schema.GroupResp

				for i := 0; i < len(groups); i++ {

					if groups[i].Status.Value != model.Archive {
						continue
					}
					if groups[i].CatalogIDs[0] != tt.args.opts.ID && groups[i].CatalogIDs[1] != tt.args.opts.ID {
						continue
					}

					max := catalogs[i].RetailPrice.Value
					min := catalogs[i+1].RetailPrice.Value

					if catalogs[i].RetailPrice.Value < catalogs[i+1].RetailPrice.Value {
						min = catalogs[i].RetailPrice.Value
						max = catalogs[i+1].RetailPrice.Value
					}
					groupResps = append(groupResps, schema.GroupResp{
						ID:          groups[i].ID,
						Status:      *groups[i].Status,
						Basis:       groups[i].Basis,
						Maximum:     max,
						Minimum:     min,
						CatalogInfo: []model.Catalog{catalogs[i], catalogs[i+1]},
					})
					if len(groupResps) == 10 {
						break
					}
				}
				tt.want = groupResps
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gi := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = gi
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			resp, err := gi.KeeperGetGroupsByCatalogID(tt.args.opts)
			if !tt.wantErr {

				assert.Nil(t, err)
				// assert.NotNil(t, resp)
				assert.Equal(t, len(resp), len(tt.want))
				for i := 0; i < len(tt.want); i++ {
					assert.Equal(t, tt.want[i], resp[i])
				}
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestGroupImpl_AddCatalogsInTheGroup(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	var catalogs []model.Catalog
	var groups []model.Group
	var getCatalogByIDsResp []schema.CreateCatalogResp
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
		})

		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	for i := 0; i < 2; i++ {

		group := model.Group{
			ID:     primitive.NewObjectID(),
			Basis:  faker.Company().Name(),
			Status: &model.GroupStatus{Value: faker.RandomChoice([]string{model.Archive, model.Unlist, model.Publish})},

			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}
	getCatalogByIDsResp = append(getCatalogByIDsResp, schema.CreateCatalogResp{
		ID:          catalogs[5].ID,
		Name:        catalogs[5].Name,
		RetailPrice: *catalogs[5].RetailPrice,
	})
	getCatalogByIDsResp = append(getCatalogByIDsResp, schema.CreateCatalogResp{
		ID:          catalogs[10].ID,
		Name:        catalogs[10].Name,
		RetailPrice: *catalogs[10].RetailPrice,
	})

	type args struct {
		opts *schema.AddCatalogsInTheGroupOpts
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

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddCatalogsInTheGroupOpts{
					ID:         groups[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogs[5].ID, catalogs[10].ID},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return(getCatalogByIDsResp, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
				tt.want = true
			},
		},
		{
			name: "[Error] Catalog ID Does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddCatalogsInTheGroupOpts{
					ID:         groups[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogs[5].ID, catalogs[10].ID},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return(getCatalogByIDsResp[1:], nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("catalog with id: %s not found", getCatalogByIDsResp[0].ID.Hex())}
			},
		},
		{
			name: "[Error] Group ID Does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddCatalogsInTheGroupOpts{
					ID:         primitive.NewObjectID(),
					CatalogIDs: []primitive.ObjectID{catalogs[5].ID, catalogs[10].ID},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {
				kc.EXPECT().GetCatalogByIDs(gomock.Any(), tt.args.opts.CatalogIDs).Times(1).Return(getCatalogByIDsResp, nil)

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("group with id:%s not found", tt.args.opts.ID)}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = g
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			fmt.Println(tt.args.opts)
			resp, errArray := g.AddCatalogsInTheGroup(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, errArray)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)

				var groupResp model.Group
				_ = tt.fields.DB.Collection(model.GroupColl).FindOne(context.TODO(), bson.M{"_id": groups[0].ID}).Decode(&groupResp)
				assert.Equal(t, groupResp.CatalogIDs[2], catalogs[5].ID)
				assert.Equal(t, groupResp.CatalogIDs[3], catalogs[10].ID)

			}

			if tt.wantErr {
				assert.NotNil(t, errArray)
				assert.Equal(t, tt.err[0].Error(), errArray[0].Error())
			}
		})
	}
}

func TestGroupImpl_RemoveCatalogsFromTheGroup(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	var catalogs []model.Catalog
	var groups []model.Group
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
		})

		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	for i := 0; i < 2; i++ {

		group := model.Group{
			ID:     primitive.NewObjectID(),
			Basis:  faker.Company().Name(),
			Status: &model.GroupStatus{Value: faker.RandomChoice([]string{model.Archive, model.Unlist, model.Publish})},

			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}

	type args struct {
		opts *schema.AddCatalogsInTheGroupOpts
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

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddCatalogsInTheGroupOpts{
					ID:         groups[0].ID,
					CatalogIDs: []primitive.ObjectID{catalogs[0].ID},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = false
				tt.want = true
			},
		},

		{
			name: "[Error] Group ID Does not exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.AddCatalogsInTheGroupOpts{
					ID:         primitive.NewObjectID(),
					CatalogIDs: []primitive.ObjectID{catalogs[5].ID, catalogs[10].ID},
				},
			},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.wantErr = true
				tt.err = []error{errors.Errorf("group with id:%s not found", tt.args.opts.ID)}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = g
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			fmt.Println(tt.args.opts)
			resp, errArray := g.RemoveCatalogsFromTheGroup(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, errArray)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)

				var groupResp model.Group
				_ = tt.fields.DB.Collection(model.GroupColl).FindOne(context.TODO(), bson.M{"_id": groups[0].ID}).Decode(&groupResp)
				assert.Equal(t, groupResp.CatalogIDs[0], catalogs[1].ID)
			}

			if tt.wantErr {
				assert.NotNil(t, errArray)
				assert.Equal(t, tt.err[0].Error(), errArray[0].Error())
			}
		})
	}
}

func TestGroupImpl_UpdateGroupStatus(t *testing.T) {
	// t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	var catalogs []model.Catalog
	var groups []model.Group
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
		})

		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	statArr := []string{model.Archive, model.Unlist, model.Publish}
	for i := 0; i < 3; i++ {

		group := model.Group{
			ID:     primitive.NewObjectID(),
			Basis:  faker.Company().Name(),
			Status: &model.GroupStatus{Value: statArr[i]},

			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}

	type args struct {
		opts *schema.UpdateGroupStatusOpts
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
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		validator  func(*testing.T, *schema.GroupResp, string)
	}
	tests := []TC{
		{
			name: "[Ok] Publish to Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateGroupStatusOpts{
					ID:     groups[2].ID,
					Status: model.Unlist,
				}
			},

			wantErr: false,
		},
		{
			name: "[Ok] Publish to Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateGroupStatusOpts{
					ID:     groups[2].ID,
					Status: model.Archive,
				}
			},

			wantErr: false,
		},
		{
			name: "[Ok] Unlist to Publish",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateGroupStatusOpts{
					ID:     groups[1].ID,
					Status: model.Publish,
				}
			},

			wantErr: false,
		},
		{
			name: "[Ok] Unlist to Archive",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateGroupStatusOpts{
					ID:     groups[1].ID,
					Status: model.Archive,
				}
			},

			wantErr: false,
		},
		{
			name: "[Error] Archive to Publish",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateGroupStatusOpts{
					ID:     groups[0].ID,
					Status: model.Publish,
				}
				tt.err = errors.Errorf("cannot change status from archive to %s", tt.args.opts.Status)
			},

			wantErr: true,
		},
		{
			name: "[Error] Archive to Unlist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},

			buildStubs: func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateGroupStatusOpts{
					ID:     groups[0].ID,
					Status: model.Unlist,
				}
				tt.err = errors.Errorf("cannot change status from archive to %s", tt.args.opts.Status)
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = g
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			err := g.UpdateGroupStatus(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)

				var groupResp model.Group
				_ = tt.fields.DB.Collection(model.GroupColl).FindOne(context.TODO(), bson.M{"_id": tt.args.opts.ID}).Decode(&groupResp)
				assert.Equal(t, groupResp.Status.Value, tt.args.opts.Status)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestGroupImpl_GetCatalogsByGroupID(t *testing.T) {
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	var catalogs []model.Catalog
	var groups []model.Group
	catalogDB := app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName)
	groupDB := app.MongoDB.Client.Database(app.Config.GroupConfig.DBName)

	itemLen := faker.RandomInt(22, 30)
	for i := 0; i < itemLen; i++ {
		catalogs = append(catalogs, model.Catalog{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: faker.Commerce().ProductName(),
			RetailPrice: &model.Price{
				Value:       faker.Commerce().Price(),
				CurrencyISO: "inr",
			},
		})

		catalogDB.Collection(model.CatalogColl).InsertOne(context.TODO(), catalogs[i])

	}
	statArr := []string{model.Archive, model.Unlist, model.Publish}
	for i := 0; i < 3; i++ {

		group := model.Group{
			ID:     primitive.NewObjectID(),
			Basis:  faker.Company().Name(),
			Status: &model.GroupStatus{Value: statArr[i]},

			CatalogIDs: []primitive.ObjectID{catalogs[i].ID, catalogs[(i + 1)].ID},
		}

		groups = append(groups, group)
		groupDB.Collection(model.GroupColl).InsertOne(context.TODO(), group)
	}

	type args struct {
		id   primitive.ObjectID
		page int
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
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, ct *mock.MockCategory, b *mock.MockBrand, kc *mock.MockKeeperCatalog)
		validator  func(*testing.T, *schema.GroupResp, string)
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

			},
			prepare: func(tt *TC) {
				tt.args.id = groups[0].ID
				tt.args.page = 0
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := &GroupImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}

			tt.fields.App.Group = g
			mockBrand := mock.NewMockBrand(ctrl)
			tt.fields.App.Brand = mockBrand
			mockCategory := mock.NewMockCategory(ctrl)
			tt.fields.App.Category = mockCategory
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockCategory, mockBrand, mockKeeperCatalog)
			tt.prepare(&tt)

			cat, err := g.GetCatalogsByGroupID(tt.args.id, tt.args.page)
			if !tt.wantErr {
				assert.Nil(t, err)
				fmt.Println(cat)
				assert.Equal(t, cat[0].ID, catalogs[0].ID)
				assert.Equal(t, cat[1].ID, catalogs[1].ID)

			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
